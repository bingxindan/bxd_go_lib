package confutil

import (
	"errors"
	"github.com/bingxindan/bxd_go_lib/tools/flagutil"
	"github.com/kardianos/osext"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
)

//Config interface definition
//developers can implement this interface to combine with other config library
type Config interface {
	//get value by section and key,if not exist,return the defaultVal
	MustValue(section, key string, defaultVal ...string) string
	//get value by sction and key and split value by delim,return array
	MustValueArray(section, key, delim string) []string
	//get section value,return array
	GetKeyList(section string) []string
	//get section value,return map
	GetSection(section string) (map[string]string, error)
	//get all section list
	GetSectionList() []string
	//get object value by section
	GetSectionObject(section string, obj interface{}) error
	//set value by section and key when need
	Set(section, key string, value interface{})
}

var (
	//global cache
	//only load the config file once
	//after config init complete,all config data get from cache
	config_cache = struct {
		sync.RWMutex
		cache map[string]Config
	}{cache: make(map[string]Config, 0)}
	//when you don't want to use ini or yaml file source,you need the "any" pattern
	//the "any" pattern can customize the data source by register a plugin
	anyFileMap = make(map[string]func() (Config, error), 0)
	//global config
	g_cfg Config
	//config path prefix if your config path is not a absolute path
	USER_CONF_PATH string
)

//config init function
//include 3 load module(ini,yaml,any),any is a plugin module,support second develop
func InitConfig() {
	//check if config has inited
	if g_cfg != nil {
		return
	}

	//get the path args
	configPath := flagutil.GetConfig()
	var err error

	//set the default path
	if len(*configPath) == 0 {
		*configPath = "/home/zhangming/jz_api/Conf/Conf.ini"
	}

	log.Printf("CONF INIT,path:%s", *configPath)
	//load config from path
	if g_cfg, err = Load(*configPath); err != nil {
		g_cfg = nil
		log.Printf("Conf,err%v", err)
	}
}

//the default base directory if file path not exist or invalid,default "/home/dev"
func home() string {
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return "/home/dev"
}

//the default conf directory,generated by home()
func Binhome() string {
	if len(USER_CONF_PATH) > 0 {
		return USER_CONF_PATH
	}
	if path, err := osext.ExecutableFolder(); err == nil {
		if strings.HasPrefix(path, "/tmp/go-build") {
			return home() + "/Conf/"
		}
		return path
	} else {
		return "."
	}
}

//path prefix,if the path is not absolute path,you can
//set the prefix manual
func SetConfPathPrefix(fullPathPrefix string) {
	if len(fullPathPrefix) != 0 {
		USER_CONF_PATH = fullPathPrefix
	}
}

//config set function
func Set(section, key string, value interface{}) {
	InitConfig()
	if g_cfg == nil {
		log.Printf("Conf,NOT_FOUND[sec:%s,key:%s]", section, key)
		return
	}
	g_cfg.Set(section, key, value)
}

//file load function
//include 3 load module(ini,yaml,any),any module support plugins
//by use flag -c=xxx,and you need provide a xxx.go which implement the
//Config interface and register the load function in plugin.go
func Load(path string) (cfg Config, err error) {
	//load any module
	//path is not a valid path
	if !strings.Contains(path, "/") {
		fn := anyFileMap[path]
		cfg, err = fn()
		return
	}
	//load ini or yaml
	//path must has more than 3 bytes
	if len([]byte(path)) < 4 {
		return nil, errors.New("path invalid")
	}

	//default load module
	fileType := "ini"

	if path[len(path)-4:] == "yaml" {
		fileType = "yaml"
		//if the path suffix is not ".ini", completed path by append ".ini"
	} else if path[len(path)-4:] != ".ini" {
		path = path + ".ini"
	}
	//config cache
	var ok bool
	config_cache.RLock()
	//read cache first
	cfg, ok = config_cache.cache[path]
	config_cache.RUnlock()
	//no cache
	if !ok {
		//path invalid,path completed
		if !strings.HasPrefix(path, "/") {
			path = Binhome() + "/" + path
			if _, err := os.Stat(path); os.IsNotExist(err) {
				path = home() + "/Conf/" + filepath.Base(path)
			}
		}
		//load file and create cache
		config_cache.Lock()
		if fileType == "yaml" {
			if cfg, err = loadYamlFile(path); err == nil {
				config_cache.cache[path] = cfg
			}
		} else {
			if cfg, err = loadIniFile(path); err == nil {
				config_cache.cache[path] = cfg
			}
		}
		config_cache.Unlock()
	}
	return
}

//cache force clear
//NOTICE:clear the cache only you change the config source
func ClearConfigCache() {
	config_cache.Lock()
	config_cache.cache = make(map[string]Config, 0)
	g_cfg = nil
	config_cache.Unlock()
}

//get config function
//section: first key
//key:second key
func GetConf(sec, key string) string {
	//init
	InitConfig()
	if g_cfg == nil {
		log.Printf("Conf,NOT_FOUND[sec:%s,key:%s] AT:%v", sec, key, flagutil.GetConfig())
		return ""
	}
	//if value not existed return ""
	return g_cfg.MustValue(sec, key, "")
}

//get config with default
//if value not existed,return default value def
func GetConfDefault(sec, key, def string) string {
	//init
	InitConfig()
	if g_cfg == nil {
		log.Printf("Conf,NOT_FOUND[sec:%s,key:%s] AT:%v", sec, key, flagutil.GetConfig())
		return ""
	}
	//if value not existed return def
	return g_cfg.MustValue(sec, key, def)
}

//get configs function,like:
/*
[Redis]
redis = 127.0.0.1:6379 127.0.0.1:7379
*/
//GetConfs("Redis") like "redis = 127.0.0.1:6379 127.0.0.1:7379",return []string{127.0.0.1:6379,127.0.0.1:7379}
func GetConfs(sec, key string) []string {
	//init
	InitConfig()
	if g_cfg == nil {
		log.Printf("Conf,NOT_FOUND[sec:%s,key:%s]", sec, key)
		return []string{}
	}
	//if value not existed return " "
	return g_cfg.MustValueArray(sec, key, " ")
}

//get configmap
//return map[string]string
func GetConfStringMap(sec string) (ret map[string]string) {
	//init
	InitConfig()
	if g_cfg == nil {
		log.Printf("Conf,NOT_FOUND[sec:%s]", sec)
		return nil
	}
	var err error
	//if value not existed return empty map
	if ret, err = g_cfg.GetSection(sec); err != nil {
		log.Printf("Conf,err:%v", err)
		ret = make(map[string]string, 0)
	}
	return
}

//get config map
//return map[string][]string,like:
/*
   [Redis]
   reids = 127.0.0.1:6379 127.0.0.0:7379
*/
//GetConfArrayMap("Redis") return map[string][]string{"redis":[127.0.0.1:6379,127.0.0.1:7379]}
func GetConfArrayMap(sec string) (ret map[string][]string) {
	//init
	InitConfig()
	if g_cfg == nil {
		log.Printf("Conf,NOT_FOUND[sec:%s]", sec)
		return nil
	}
	ret = make(map[string][]string, 0)
	//get all keys
	confs := g_cfg.GetKeyList(sec)
	//get all config by range keys
	for _, k := range confs {
		ret[k] = g_cfg.MustValueArray(sec, k, " ")
	}
	return
}

// get config value with object value return
func ConfMapToStruct(sec string, v interface{}) error {
	//init
	InitConfig()
	if g_cfg == nil {
		log.Printf("Conf,NOT_FOUND[sec:%s]", sec)
		return nil
	}
	return g_cfg.GetSectionObject(sec, v)
}
