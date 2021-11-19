package db

import "reflect"

var DbDaos = make(map[string]reflect.Type, 0)
var DbStructs = make(map[string]reflect.Type, 0)

type dbDao struct {
	dao       reflect.Value
	tableName string
}

func NewDbDao(table string) (this *dbDao) {
	this = new(dbDao)
	this.tableName = table
	if t, o := DbDaos[table]; o {
		this.dao = reflect.New(t)
	}
	this.UpdateEngine()
	return
}

func (this *dbDao) UpdateEngine(v ...interface{}) {
	if !this.dao.IsValid() {
		return
	}
	method := this.dao.MethodByName("UpdateEngine")
	if method.IsValid() {
		method.Call([]reflect.Value{})
	}
}

func (this *dbDao) GetCountByIndex(index string, params ...Param) (ret int64, err error) {
	if r, e := this.GetByIndex(index+"Count", params...); e == nil {
		if len(r) > 0 {
			ret, _ = r[0].(int64)
		}
	} else {
		err = e
	}
	return
}

func (this *dbDao) GetByIndex(index string, params ...Param) (ret []interface{}, err error) {
	ret, err = this.GetLimitByIndex(0, 0, index, params...)
	return
}

func (this *dbDao) GetLimitByIndex(offset, limit int64, index string, params ...Param) (ret []interface{}, err error) {
	ret = make([]interface{}, 0)
	if !this.dao.IsValid() {
		return
	}
	methodName := "Get" + index
	if limit > 0 {
		methodName += "Limit"
	}
	method := this.dao.MethodByName(methodName)
	if !method.IsValid() {
		return
	}
	innum := method.Type().NumIn()
	ps := this.CastParam(innum, offset, limit, params...)
	vs := method.Call(ps)
	length := len(vs)
	for i := 0; i < length; i++ {
		ret = append(ret, vs[i].Interface())
	}
	err, _ = vs[length-1].Interface().(error)
	return
}

func (this *dbDao) CastParam(innum int, offset, limit int64, params ...Param) (ret []reflect.Value) {
	ret = make([]reflect.Value, 0, innum)
	for _, p := range params {
		if p != nil {
			ret = append(ret, reflect.ValueOf(p))
		} else {
			var in ParamNil
			ret = append(ret, reflect.ValueOf(in))
		}
	}
	if limit > 0 {
		innum -= 2
	}
	for i := len(ret); i < innum; i++ {
		var in ParamNil
		ret = append(ret, reflect.ValueOf(in))
	}
	if limit > 0 {
		ret = append(ret, reflect.ValueOf(int(offset)), reflect.ValueOf(int(limit)))
	}
	return
}
