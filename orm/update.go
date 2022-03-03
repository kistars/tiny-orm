package orm

import (
	"database/sql"
	"errors"
	"reflect"
	"strings"
)

//更新
func (e *OrmEngine) Update(data ...interface{}) (int64, error) {
	//判断是结构体还是多个字符串
	var dataType int
	if len(data) == 1 {
		dataType = 1
	} else if len(data) == 2 {
		dataType = 2
	} else {
		return 0, errors.New("参数个数错误")
	}

	//如果是结构体
	if dataType == 1 {
		t := reflect.TypeOf(data[0])
		v := reflect.ValueOf(data[0])

		var fieldNameArray []string
		for i := 0; i < t.NumField(); i++ {

			//首字母小写，不可反射
			if !v.Field(i).CanInterface() {
				continue
			}

			//解析tag，找出真实的sql字段名
			sqlTag := t.Field(i).Tag.Get("sql")
			if sqlTag != "" {
				fieldNameArray = append(fieldNameArray, strings.Split(sqlTag, ",")[0]+"=?")
			} else {
				fieldNameArray = append(fieldNameArray, t.Field(i).Name+"=?")
			}

			e.UpdateExec = append(e.UpdateExec, v.Field(i).Interface())
		}
		e.UpdateParam += strings.Join(fieldNameArray, ",")

	} else if dataType == 2 {
		//直接=的情况
		e.UpdateParam += data[0].(string) + "=?"
		e.UpdateExec = append(e.UpdateExec, data[1])
	}

	//拼接sql
	e.Prepare = "update " + e.GetTable() + " set " + e.UpdateParam

	//如果where不为空
	if e.WhereParam != "" || e.OrWhereParam != "" {
		e.Prepare += " where " + e.WhereParam + e.OrWhereParam
	}

	//limit不为空
	if e.LimitParam != "" {
		e.Prepare += "limit " + e.LimitParam
	}

	//prepare
	var stmt *sql.Stmt
	var err error
	if e.TransStatus == 1 {
		stmt, err = e.Tx.Prepare(e.Prepare)
	} else {
		stmt, err = e.Db.Prepare(e.Prepare)
	}
	if err != nil {
		return 0, e.setErrorInfo(err)
	}

	//合并UpdateExec和WhereExec
	if e.WhereExec != nil {
		e.AllExec = append(e.UpdateExec, e.WhereExec...)
	}

	//执行exec,注意这是stmt.Exec
	result, err := stmt.Exec(e.AllExec...)
	if err != nil {
		return 0, e.setErrorInfo(err)
	}

	//影响的行数
	id, _ := result.RowsAffected()
	return id, nil
}
