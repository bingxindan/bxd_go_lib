package db

import "database/sql"

func (d *DbBaseDao) Create(bean interface{}) (int64, error) {
	if d.Session == nil {
		return d.Engine.Insert(bean)
	} else {
		return d.Session.Insert(bean)
	}
}

func (d *DbBaseDao) Update(bean interface{}) (int64, error) {
	if d.Session == nil {
		return d.Engine.Id(d.Engine.IdOf(bean)).AllCols().Update(bean)
	} else {
		return d.Session.Id(d.Engine.IdOf(bean)).AllCols().Update(bean)
	}
}

func (d *DbBaseDao) UpdateCols(bean interface{}, cols ...string) (int64, error) {
	if d.Session == nil {
		return d.Engine.Id(d.Engine.IdOf(bean)).Cols(cols...).Update(bean)
	} else {
		return d.Session.Id(d.Engine.IdOf(bean)).Cols(cols...).Update(bean)
	}
}

func (d *DbBaseDao) Delete(bean interface{}) (int64, error) {
	if d.Session == nil {
		return d.Engine.Id(d.Engine.IdOf(bean)).Delete(bean)
	} else {
		return d.Session.Id(d.Engine.IdOf(bean)).Delete(bean)
	}
}

func (d *DbBaseDao) Exec(sqlOrArgs ...interface{}) (sql.Result, error) {
	session := d.Engine.NewSession()
	defer session.Close()
	return session.Exec(sqlOrArgs...)
}

// Query a raw sql and return records as []map[string][]byte
func (d *DbBaseDao) Query(sqlorArgs ...interface{}) (resultsSlice []map[string][]byte, err error) {
	session := d.Engine.NewSession()
	defer session.Close()
	return session.Query(sqlorArgs...)
}

// QueryString runs a raw sql and return records as []map[string]string
func (d *DbBaseDao) QueryString(sqlorArgs ...interface{}) ([]map[string]string, error) {
	session := d.Engine.NewSession()
	defer session.Close()
	return session.QueryString(sqlorArgs...)
}

// QueryInterface runs a raw sql and return records as []map[string]interface{}
func (d *DbBaseDao) QueryInterface(sqlorArgs ...interface{}) ([]map[string]interface{}, error) {
	session := d.Engine.NewSession()
	defer session.Close()
	return session.QueryInterface(sqlorArgs...)
}

// Get retrieve one record from table, bean's non-empty fields
// are conditions
func (d *DbBaseDao) Get(bean interface{}) (bool, error) {
	session := d.Engine.NewSession()
	defer session.Close()
	return session.Get(bean)
}

// Exist returns true if the record exist otherwise return false
func (d *DbBaseDao) Exist(bean ...interface{}) (bool, error) {
	session := d.Engine.NewSession()
	defer session.Close()
	return session.Exist(bean...)
}

// Find retrieve records from table, condiBeans's non-empty fields
// are conditions. beans could be []Struct, []*Struct, map[int64]Struct
// map[int64]*Struct
func (d *DbBaseDao) Find(beans interface{}, condiBeans ...interface{}) error {
	session := d.Engine.NewSession()
	defer session.Close()
	return session.Find(beans, condiBeans...)
}

// FindAndCount find the results and also return the counts
func (d *DbBaseDao) FindAndCount(rowsSlicePtr interface{}, condiBean ...interface{}) (int64, error) {
	session := d.Engine.NewSession()
	defer session.Close()
	return session.FindAndCount(rowsSlicePtr, condiBean...)
}

// Count counts the records. bean's non-empty fields are conditions.
func (d *DbBaseDao) Count(bean ...interface{}) (int64, error) {
	session := d.Engine.NewSession()
	defer session.Close()
	return session.Count(bean...)
}

// Sum sum the records by some column. bean's non-empty fields are conditions.
func (d *DbBaseDao) Sum(bean interface{}, colName string) (float64, error) {
	session := d.Engine.NewSession()
	defer session.Close()
	return session.Sum(bean, colName)
}

// SumInt sum the records by some column. bean's non-empty fields are conditions.
func (d *DbBaseDao) SumInt(bean interface{}, colName string) (int64, error) {
	session := d.Engine.NewSession()
	defer session.Close()
	return session.SumInt(bean, colName)
}

// Sums sum the records by some columns. bean's non-empty fields are conditions.
func (d *DbBaseDao) Sums(bean interface{}, colNames ...string) ([]float64, error) {
	session := d.Engine.NewSession()
	defer session.Close()
	return session.Sums(bean, colNames...)
}

// SumsInt like Sums but return slice of int64 instead of float64.
func (d *DbBaseDao) SumsInt(bean interface{}, colNames ...string) ([]int64, error) {
	session := d.Engine.NewSession()
	defer session.Close()
	return session.SumsInt(bean, colNames...)
}
