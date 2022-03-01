package orm

//查询多条，返回值为map切片
func (e *OrmEngine) Select() ([]map[string]string, error) {

	//拼接sql
	e.Prepare = "select * from " + e.GetTable()

	//如果where不为空
	if e.WhereParam != "" || e.OrWhereParam != "" {
		e.Prepare += " where " + e.WhereParam + e.OrWhereParam
	}

	e.AllExec = e.WhereExec

	//query
	rows, err := e.Db.Query(e.Prepare, e.AllExec...)
	if err != nil {
		return nil, e.setErrorInfo(err)
	}

	//读出查询出的列字段名
	column, err := rows.Columns()
	if err != nil {
		return nil, e.setErrorInfo(err)
	}

	//values是每个列的值，这里获取到byte里
	values := make([][]byte, len(column))

	//因为每次查询出来的列是不定长的，用len(column)定住当次查询的长度
	scans := make([]interface{}, len(column))

	for i := range values {
		scans[i] = &values[i]
	}

	results := make([]map[string]string, 0)
	for rows.Next() {
		if err := rows.Scan(scans...); err != nil {
			//query.Scan查询出来的不定长值放到scans[i] = &values[i],也就是每行都放在values里
			return nil, e.setErrorInfo(err)
		}

		//每行数据
		row := make(map[string]string)

		//循环values数据，通过相同的下标，取column里面对应的列名，生成1个新的map
		for k, v := range values {
			key := column[k]
			row[key] = string(v)
		}

		//添加到map切片中
		results = append(results, row)
	}

	return results, nil
}

//查询1条
func (e *OrmEngine) SelectOne() (map[string]string, error) {
	//limit 1 单个查询
	results, err := e.Limit("1").Select()
	if err != nil {
		return nil, e.setErrorInfo(err)
	}

	//判断是否为空
	if len(results) == 0 {
		return nil, nil
	} else {
		return results[0], nil
	}
}
