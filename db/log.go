package db

import (
	"fmt"
	"time"

	logger ""
	"xorm.io/core"
)

var dbLogger = &dbLog{}

type dbLog struct{}

func (d *dbLog) Debug(v ...interface{}) {
	logger.D("[dbdao]", fmt.Sprint(v...))
}

func (d *dbLog) Debugf(format string, v ...interface{}) {
	logger.D("[dbdao]", format, v...)
}

func (d *dbLog) Info(v ...interface{}) {
	logger.I("[dbdao]", fmt.Sprint(v...))
}

func (d *dbLog) Infof(format string, v ...interface{}) {
	if len(v) > 0 {
		duration, ok := v[len(v)-1].(time.Duration)
		if ok && duration < slowDuration {
			return
		}
	}
	logger.I("[dbdao]", format, v...)
}

func (d *dbLog) Warn(v ...interface{}) {
	logger.W("[dbdao]", fmt.Sprint(v...))
}

func (d *dbLog) Warnf(format string, v ...interface{}) {
	logger.W("[dbdao]", format, v...)
}

func (d *dbLog) Error(v ...interface{}) {
	logger.E("[dbdao]", fmt.Sprint(v...))
}

func (d *dbLog) Errorf(format string, v ...interface{}) {
	logger.E("[dbdao]", format, v...)
}

func (d *dbLog) Level() core.LogLevel {
	return core.LOG_INFO
}

func (d *dbLog) SetLevel(l core.LogLevel) {
	return
}

func (d *dbLog) ShowSQL(show ...bool) {
	return
}

func (d *dbLog) IsShowSQL() bool {
	return showSql
}
