package db

type DbBaseDao struct {
	Engine  *Engine
	Session *Session
}

type Param interface{}
type ParamNil struct{}
type ParamDesc bool
type ParamIn []interface{}
type ParamRange struct {
	Min interface{}
	Max interface{}
}
type ParamInDesc ParamIn
type ParamRangeDesc ParamRange

func CastToParamIn(input interface{}) ParamIn {
	params := make(ParamIn, 0)
	switch v := input.(type) {
	case []interface{}:
		for _, param := range v {
			params = append(params, param)
		}
	case []int64:
		for _, param := range v {
			params = append(params, param)
		}
	case []int:
		for _, param := range v {
			params = append(params, param)
		}
	case []int32:
		for _, param := range v {
			params = append(params, param)
		}
	case []int8:
		for _, param := range v {
			params = append(params, param)
		}
	case []uint64:
		for _, param := range v {
			params = append(params, param)
		}
	case []uint:
		for _, param := range v {
			params = append(params, param)
		}
	case []uint32:
		for _, param := range v {
			params = append(params, param)
		}
	case []uint8:
		for _, param := range v {
			params = append(params, param)
		}
	case []string:
		for _, param := range v {
			params = append(params, param)
		}
	default:
		params = append(params, 0)
	}
	return params
}

func CastToParamInDesc(input interface{}) ParamInDesc {
	return ParamInDesc(CastToParamIn(input))
}

func (d *DbBaseDao) InitSession() {
	if d.Session == nil {
		d.Session = d.Engine.Where("")
	}
}

//通过此方法，可以指定表名称
func (d *DbBaseDao) SetTable(tableName string) {
	if d.Session == nil {
		d.InitSession()
	}

	d.Session.Table(tableName)
}

func (d *DbBaseDao) BuildQuery(input Param, name string) {
	name = d.Engine.Quote(name)

	switch val := input.(type) {
	case ParamDesc:
		if val {
			d.Session = d.Session.Desc(name)
		}
	case ParamIn:
		if len(val) == 1 {
			d.Session = d.Session.And(name+"=?", val[0])
		} else {
			d.Session = d.Session.In(name, val)
		}
	case ParamInDesc:
		if len(val) == 1 {
			d.Session = d.Session.And(name+"=?", val[0])
		} else {
			d.Session = d.Session.In(name, val)
		}
		d.Session = d.Session.Desc(name)
	case ParamRange:
		if val.Min != nil {
			d.Session = d.Session.And(name+">=?", val.Min)
		}
		if val.Max != nil {
			d.Session = d.Session.And(name+"<?", val.Max)
		}
	case ParamRangeDesc:
		if val.Min != nil {
			d.Session = d.Session.And(name+">=?", val.Min)
		}
		if val.Max != nil {
			d.Session = d.Session.And(name+"<?", val.Max)
		}
		d.Session = d.Session.Desc(name)
	case ParamNil:
	case nil:
	default:
		d.Session = d.Session.And(name+"=?", val)
	}
}

func (d *DbBaseDao) UpdateEngine(v ...interface{}) {
	if len(v) == 0 {
		d.Engine = GetDefault("reader").Engine
		d.Session = nil
	} else if len(v) == 1 {
		param := v[0]
		if engine, ok := param.(*Engine); ok {
			d.Engine = engine
			d.Session = nil
		} else if session, ok := param.(*Session); ok {
			d.Session = session
		} else if tpe, ok := param.(bool); ok {
			cluster := "reader"
			if tpe == true {
				cluster = "writer"
			}
			d.Engine = GetDefault(cluster).Engine
			d.Session = nil
		}
	}
}
