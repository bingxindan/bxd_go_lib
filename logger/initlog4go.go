package logger

import (
	"encoding/xml"
	"fmt"
	"github.com/bingxindan/bxd_go_lib/logger/log4go"
	"github.com/bingxindan/bxd_go_lib/logger/logutils"
	"io/ioutil"
	"os"
	"strings"

)

type xmlProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

type xmlFilter struct {
	Enabled  string        `xml:"enabled,attr"`
	Tag      string        `xml:"tag"`
	Level    string        `xml:"level"`
	Type     string        `xml:"type"`
	Property []xmlProperty `xml:"property"`
}

type xmlLoggerConfig struct {
	Filter []xmlFilter `xml:"filter"`
}

//main函数调用 使用example.xml配置
func InitLogger(logpath string) {
	if logpath == "" {
		logpath = "conf/log.xml"
	}
	checkLogPath(logpath)
	log4go.LoadConfiguration(logpath)
	logutils.Inited = true
}

//使用自定义option
func NewLogConfig() *log4go.LogConfig {
	config := new(log4go.LogConfig)
	//存储路径
	config.LogPath = "/home/logs/bxdlog/default/default.log"
	//日志级别
	config.Level = "INFO"
	//日志标签 多日志时使用
	config.Tag = "default"
	//日志格式
	config.Format = "%G %L %S %M"
	//最大行数切割
	config.RotateLines = "0K"
	//最大容量切割
	config.RotateSize = "0M"
	//按日期切割
	config.RotateHourly = true
	//是否启用切割
	config.Rotate = true
	//日志保留时间，day
	config.Retention = "0"
	return config
}

//自定义config Init
func InitLogWithConfig(config *log4go.LogConfig) {
	if config.LogPath == "" {
		fmt.Fprintf(os.Stderr, "InitLoggerConfig: Error: config could not found logpath %s\n", config.LogPath)
		os.Exit(1)
	}
	checkLogConfig(config)
	log4go.LoadLogConfig(config)
	logutils.Inited = true
}

func Close() {
	log4go.Close()
}
func checkLogConfig(config *log4go.LogConfig) {
	if _, ok := logutils.LevelMap[config.Level]; ok {
		if logutils.LevelMap[config.Level] < logutils.SortLevel {
			logutils.SortLevel = logutils.LevelMap[config.Level]
			logutils.Level = config.Level
		}
	}
	paths := strings.Split(config.LogPath, "/")
	if len(paths) > 1 {
		dir := strings.Join(paths[0:len(paths)-1], "/")
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Could not create directory %s, err:%s\n", dir, err)
			os.Exit(1)
		}
	} else {
		fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: log directory invalid %s\n", config.LogPath)
		os.Exit(1)
	}
}

func checkLogPath(filename string) {
	fd, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Could not open %q for reading: %s\n", filename, err)
		os.Exit(1)
	}

	contents, err := ioutil.ReadAll(fd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Could not read %q: %s\n", filename, err)
		os.Exit(1)
	}

	xc := new(xmlLoggerConfig)
	if err := xml.Unmarshal(contents, xc); err != nil {
		fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Could not parse XML configuration in %q: %s\n", filename, err)
		os.Exit(1)
	}
	//stdout kworker trace--new
	for _, xmlfilt := range xc.Filter {
		if xmlfilt.Enabled == "true" {
			//获取 log.xml中的level
			if _, ok := logutils.LevelMap[xmlfilt.Level]; ok {
				//取出配置文件中最小的级别
				if logutils.LevelMap[xmlfilt.Level] < logutils.SortLevel {
					logutils.SortLevel = logutils.LevelMap[xmlfilt.Level]
					logutils.Level = xmlfilt.Level
				}
			}
			if xmlfilt.Tag != "stdout" {
				for _, prop := range xmlfilt.Property {
					if prop.Name == "filename" {
						paths := strings.Split(prop.Value, "/")
						if len(paths) > 1 {
							dir := strings.Join(paths[0:len(paths)-1], "/")
							err := os.MkdirAll(dir, 0755)
							if err != nil {
								fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Could not create directory %s, err:%s\n", dir, err)
								os.Exit(1)
							}
						} else {
							fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: log directory invalid %s, err:%s\n", prop.Value, err)
							os.Exit(1)
						}

					}
				}
			}
		}
	}

}
