package my_orm

import "database/sql"

//删除
func (e *OrmEngine) Delete() (int64, error) {

	//拼接delete sql
	e.Prepare = "delete from " + e.GetTable()

	//如果where不为空
	if e.WhereParam != "" || e.OrWhereParam != "" {
		e.Prepare += " where " + e.WhereParam + e.OrWhereParam
	}

	//limit不为空
	if e.LimitParam != "" {
		e.Prepare += "limit " + e.LimitParam
	}

	//第一步：Prepare
	var stmt *sql.Stmt
	var err error
	stmt, err = e.Db.Prepare(e.Prepare)
	if err != nil {
		return 0, err
	}

	e.AllExec = e.WhereExec

	//第二步：执行exec,注意这是stmt.Exec
	result, err := stmt.Exec(e.AllExec...)
	if err != nil {
		return 0, e.setErrorInfo(err)
	}

	//影响的行数
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, e.setErrorInfo(err)
	}

	return rowsAffected, nil
}
