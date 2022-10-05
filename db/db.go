package db

import (
	"context"
	"github.com/bingxindan/bxd_go_lib/config"
	"github.com/bingxindan/bxd_go_lib/tools/confutil"
	"log"
	"strings"
	"sync/atomic"
	"time"

	"github.com/spf13/cast"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type DBDao struct {
	Engine *Engine
	quiter chan struct{}
}

var (
	db_instance  map[string][]*DBDao
	curDbPoses   map[string]*uint64 // 当前选择的数据库
	showSql      bool
	showExecTime bool
	slowDuration time.Duration
	maxConn      = 100
	maxIdle      = 30
)

func newDBDaoWithParams(host string, driver string) (Db *DBDao) {
	Db = new(DBDao)

	eng, err := xorm.NewEngine(driver, host)
	engine := &Engine{Engine: eng}

	Db.Engine = engine
	//TODO: 增加存活检查
	if err != nil {
		log.Fatal(err)
	}
	Db.Engine.SetMaxOpenConns(maxConn)
	Db.Engine.SetMaxIdleConns(maxIdle)
	Db.Engine.SetConnMaxLifetime(time.Second * 3000)
	Db.Engine.ShowSQL(showSql)
	Db.Engine.ShowExecTime(showExecTime)
	Db.Engine.SetLogger(dbLogger)
	return
}

func GetDefault(cluster string) *DBDao {
	return GetDbInstance("default", cluster)
}

func init() {
	db_instance = make(map[string][]*DBDao, 0)
	curDbPoses = make(map[string]*uint64)
	//idc := confdao.GetIDC()
	idc := ""
	showSqlStr := config.String("data.mysql.show_sql")
	showSql = showSqlStr == "true"

	showExecTimeStr := config.String("data.mysql.show_exec_time")
	showExecTime = showExecTimeStr == "true"

	slowStr := config.String("data.mysql.slow")
	slowDuration = time.Duration(cast.ToInt(slowStr)) * time.Millisecond

	maxConnStr := config.String("data.mysql.max_conn")
	maxConnConfig := cast.ToInt(maxConnStr)
	if maxConnConfig > 0 {
		maxConn = maxConnConfig
	}

	maxIdleStr := config.String("data.mysql.max_idle_conn")
	maxIdleConfig := cast.ToInt(maxIdleStr)
	if maxIdleConfig > 0 {
		maxIdle = maxIdleConfig
	}

	if maxIdle > maxConn {
		maxIdle = maxConn
	}

	for cluster, hosts := range confutil.GetConfArrayMap("MysqlCluster") {
		items := strings.Split(cluster, ".")
		//必须包含 writer 和 reader
		if len(items) < 2 {
			continue
		}
		//过滤IDC
		if len(items) > 2 && items[2] != idc {
			continue
		}
		instance := items[0] + "." + items[1]
		dbs := make([]*DBDao, 0)
		for _, host := range hosts {
			dbs = append(dbs, newDBDaoWithParams(host, "mysql"))
		}
		db_instance[instance] = dbs
		curDbPoses[instance] = new(uint64)
	}
}

func GetDbInstance(db, cluster string) *DBDao {
	key := db + "." + cluster
	if instances, ok := db_instance[key]; ok {
		// round-robin选择数据库
		cur := atomic.AddUint64(curDbPoses[key], 1) % uint64(len(instances))
		return instances[cur]
	} else {
		return nil
	}
}

func GetDbInstanceWithCtx(ctx context.Context, db, cluster string) *DBDao {
	bench := ctx.Value("IS_BENCHMARK")
	if cast.ToString(bench) == "1" {
		db = "benchmark_" + db
	}
	key := db + "." + cluster
	if instances, ok := db_instance[key]; ok {
		// round-robin选择数据库
		cur := atomic.AddUint64(curDbPoses[key], 1) % uint64(len(instances))
		return instances[cur]
	} else {
		return nil
	}
}

func (this *DBDao) GetSession() *Session {
	return this.Engine.NewSession()
}

func (this *DBDao) Close() {
	this.Engine.Close()
}
