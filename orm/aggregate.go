package orm

func (e *OrmEngine) aggregateQuery(name, param string) (interface{}, error) {

	//拼接sql
	e.Prepare = "select " + name + "(" + param + ") as cnt from " + e.GetTable()

	//如果where不为空
	if e.WhereParam != "" || e.OrWhereParam != "" {
		e.Prepare += " where " + e.WhereParam + e.OrWhereParam
	}

	//limit不为空
	if e.LimitParam != "" {
		e.Prepare += " limit " + e.LimitParam
	}

	e.AllExec = e.WhereExec

	//生成sql
	//e.generateSql()

	//执行绑定
	var cnt interface{}

	//queryRows
	err := e.Db.QueryRow(e.Prepare, e.AllExec...).Scan(&cnt)
	if err != nil {
		return nil, e.setErrorInfo(err)
	}

	return cnt, err
}

func (e *OrmEngine) Count() (int64, error) {
	res, err := e.aggregateQuery("count", "1")
	if err != nil {
		return 0, e.setErrorInfo(err)
	}
	return res.(int64), nil
}

func (e *OrmEngine) Max(field string) (string, error) {
	res, err := e.aggregateQuery("max", field)
	if err != nil {
		return "", e.setErrorInfo(err)
	}
	return string(res.([]byte)), nil
}

func (e *OrmEngine) Min(field string) (string, error) {
	res, err := e.aggregateQuery("min", field)
	if err != nil {
		return "", e.setErrorInfo(err)
	}
	return string(res.([]byte)), nil
}

func (e *OrmEngine) Avg(field string) (string, error) {
	res, err := e.aggregateQuery("avg", field)
	if err != nil {
		return "", e.setErrorInfo(err)
	}
	return string(res.([]byte)), nil
}

func (e *OrmEngine) Sum(field string) (string, error) {
	res, err := e.aggregateQuery("sum", field)
	if err != nil {
		return "", e.setErrorInfo(err)
	}
	return string(res.([]byte)), nil
}
