package orm

import (
	"reflect"
	"strings"
)

func (e *OrmEngine) Group(group ...string) *OrmEngine {
	if len(group) == 1 {
		e.GroupParam = group[0]
	} else if len(group) > 1 {
		e.GroupParam = strings.Join(group, ",")
	}
	return e
}

//having过滤
func (e *OrmEngine) Having(having ...interface{}) *OrmEngine {

	//判断是结构体还是多个字符串
	var dataType int
	if len(having) == 1 {
		dataType = 1
	} else if len(having) == 2 {
		dataType = 2
	} else if len(having) == 3 {
		dataType = 3
	} else {
		panic("having个数错误")
	}

	//多次调用判断
	if e.HavingParam != "" {
		e.HavingParam += "and ("
	} else {
		e.HavingParam += "("
	}

	//如果是结构体
	if dataType == 1 {
		t := reflect.TypeOf(having[0])
		v := reflect.ValueOf(having[0])

		var fieldNameArray []string
		for i := 0; i < t.NumField(); i++ {

			//小写开头，无法反射，跳过
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

			e.WhereExec = append(e.WhereExec, v.Field(i).Interface())
		}
		e.HavingParam += strings.Join(fieldNameArray, " and ") + ") "

	} else if dataType == 2 {
		//直接=的情况
		e.HavingParam += having[0].(string) + "=?) "
		e.WhereExec = append(e.WhereExec, having[1])
	} else if dataType == 3 {
		//3个参数的情况
		e.HavingParam += having[0].(string) + " " + having[1].(string) + " ?) "
		e.WhereExec = append(e.WhereExec, having[2])
	}

	return e
}
