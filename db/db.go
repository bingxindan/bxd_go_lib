package db

import (
	"context"
	"github.com/bingxindan/bxd_go_lib/config"
	"log"
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
	dbInstance   map[string][]*DBDao
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
	dbInstance = make(map[string][]*DBDao, 0)
	curDbPoses = make(map[string]*uint64)

	showSql = config.Bool("data.mysql.show_sql")
	showExecTime = config.Bool("data.mysql.show_exec_time")
	slowDuration = config.Duration("data.mysql.slow")
	maxConnConfig := config.Int("data.mysql.max_conn")
	if maxConnConfig > 0 {
		maxConn = maxConnConfig
	}
	maxIdleConfig := config.Int("data.mysql.max_idle_conn")
	if maxIdleConfig > 0 {
		maxIdle = maxIdleConfig
	}
	if maxIdle > maxConn {
		maxIdle = maxConn
	}

	// 主实例
	keyWriter := config.String("data.mysql.source.key_writer")
	hostWriters := config.StringList("data.mysql.source.host_writer")
	dbWriters := make([]*DBDao, 0)
	for _, hostWriter := range hostWriters {
		dbWriters = append(dbWriters, newDBDaoWithParams(hostWriter, "mysql"))
	}
	dbInstance[keyWriter] = dbWriters
	curDbPoses[keyWriter] = new(uint64)

	// 从实例
	keyReader := config.String("data.mysql.source.key_reader")
	hostReaders := config.StringList("data.mysql.source.host_reader")
	dbReaders := make([]*DBDao, 0)
	for _, hostReader := range hostReaders {
		dbReaders = append(dbReaders, newDBDaoWithParams(hostReader, "mysql"))
	}
	dbInstance[keyReader] = dbReaders
	curDbPoses[keyReader] = new(uint64)
}

func GetDbInstance(dbCluster string) *DBDao {
	if instances, ok := dbInstance[dbCluster]; ok {
		// round-robin选择数据库
		cur := atomic.AddUint64(curDbPoses[dbCluster], 1) % uint64(len(instances))
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
	if instances, ok := dbInstance[key]; ok {
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
