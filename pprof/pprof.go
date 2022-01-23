package pprof

import (
	"github.com/bingxindan/bxd_go_lib/logger"
	"github.com/bingxindan/bxd_go_lib/tools/confutil"
	"github.com/spf13/cast"
	"log"
	"net"
	"net/http"
	"runtime/debug"
)

func Pprof() {
	grace := confutil.GetConf("Server", "grace")
	if grace == "true" {
		InitPort()
	} else {
		go pprofstartWithoutGrace()
	}
}

func pprofstartWithoutGrace() {
	enable := confutil.GetConf("Pprof", "enable")
	if enable != "true" {
		return
	}

	port := confutil.GetConf("Pprof", "port")
	if len(port) <= 0 {
		logger.E("Pprof", "pprof port:%s format wrong", port)
		return
	}
	logger.I("Pprof", "open pprof on port:%s", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

var PprofPort string

var Pserver = &http.Server{Addr: PprofPort}

func InitPort() {
	enable := confutil.GetConf("Pprof", "enable")
	if enable != "true" {
		return
	}

	PprofPort = ":" + confutil.GetConf("Pprof", "port")
	if len(PprofPort) <= 0 {
		logger.E("Pprof", "pprof port:%s format wrong", PprofPort)
		return
	}

}

func Start(l net.Listener) {
	err := Pserver.Serve(l)
	if err != nil {
		logger.W("ServerError", "Unhandled error: %v\n stack:%v", err.Error(), cast.ToString(debug.Stack()))
	}
}
