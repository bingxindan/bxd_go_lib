package flagutil

import (
	"flag"
)

var signal = flag.String("s", "", "start or stop")
var confpath = flag.String("c", "", "config path")
var cfgpath = flag.String("cfg", "Conf/Config.json", "json config path")
var foreground = flag.Bool("f", false, "foreground")
var mock = flag.Bool("m", false, "mock")
var version = flag.Bool("v", false, "version")
var mode = flag.Int("mode", 0, "mode")
var extendedopt = flag.String("extended", "", "extended options, You can customize the options to bypass the flag.parse() restrictions")
var usr1 = flag.String("usr1", "", "user defined flag -usr1")
var usr2 = flag.String("usr2", "", "user defined flag -usr2")
var usr3 = flag.String("usr3", "", "user defined flag -usr3")
var usr4 = flag.String("usr4", "", "user defined flag -usr4")
var usr5 = flag.String("usr5", "", "user defined flag -usr5")

func GetSignal() *string {
	return signal
}

func GetVersion() *bool {
	return version
}

func GetMode() *int {
	return mode
}

func GetConfig() *string {
	return confpath
}

func SetConfig(path string) {
	confpath = &path
}

func GetCfg() *string {
	return cfgpath
}

func SetCfg(path string) {
	cfgpath = &path
}

func GetForeground() *bool {
	return foreground
}

func GetMock() *bool {
	return mock
}
func SetMock(mval bool) {
	mock = &mval
}

func GetExtendedopt() *string {
	return extendedopt
}

func GetUsr1() *string {
	return usr1
}
func GetUsr2() *string {
	return usr2
}
func GetUsr3() *string {
	return usr3
}
func GetUsr4() *string {
	return usr4
}
func GetUsr5() *string {
	return usr5
}
