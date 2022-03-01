package orm

import (
	"reflect"
	"strings"
)

func (e *OrmEngine) Where(data ...interface{}) *OrmEngine {
	//判断是结构体还是多个字符串
	var dataType int
	if len(data) == 1 {
		dataType = 1
	} else if len(data) == 2 {
		dataType = 2
	} else if len(data) == 3 {
		dataType = 3
	} else {
		panic("参数个数错误")
	}

	//多次调用判断
	if e.WhereParam != "" {
		e.WhereParam += "and ("
	} else {
		e.WhereParam += "("
	}

	//如果是结构体
	if dataType == 1 {
		t := reflect.TypeOf(data[0])
		v := reflect.ValueOf(data[0])

		//字段名
		var fieldNameArray []string

		//循环解析
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

			e.WhereExec = append(e.WhereExec, v.Field(i).Interface())
		}

		//拼接
		e.WhereParam += strings.Join(fieldNameArray, " and ") + ") "

	} else if dataType == 2 {
		//直接=的情况
		e.WhereParam += data[0].(string) + " = ?) "
		e.WhereExec = append(e.WhereExec, data[1])
	} else if dataType == 3 {
		//3个参数的情况

		//区分是操作符in的情况
		data2 := strings.Trim(strings.ToLower(data[1].(string)), " ")
		if data2 == "in" || data2 == "not in" {
			//判断传入的是切片
			reType := reflect.TypeOf(data[2]).Kind()
			if reType != reflect.Slice && reType != reflect.Array {
				panic("in/not in 操作传入的数据必须是切片或者数组")
			}

			//反射值
			v := reflect.ValueOf(data[2])
			//数组/切片长度
			dataNum := v.Len()
			//占位符
			ps := make([]string, dataNum)
			for i := 0; i < dataNum; i++ {
				ps[i] = "?"
				e.WhereExec = append(e.WhereExec, v.Index(i).Interface())
			}

			//拼接
			e.WhereParam += data[0].(string) + " " + data2 + " (" + strings.Join(ps, ",") + "))"

		} else {
			e.WhereParam += data[0].(string) + " " + data[1].(string) + " ?)"
			e.WhereExec = append(e.WhereExec, data[2])
		}
	}

	return e
}

func (e *OrmEngine) OrWhere(data ...interface{}) *OrmEngine {
	//判断是结构体还是多个字符串
	var dataType int
	if len(data) == 1 {
		dataType = 1
	} else if len(data) == 2 {
		dataType = 2
	} else if len(data) == 3 {
		dataType = 3
	} else {
		panic("参数个数错误")
	}

	//多次调用判断
	if e.WhereParam == "" {
		panic("or必须在where后面")
	}
	e.OrWhereParam += "or ("

	//如果是结构体
	if dataType == 1 {
		t := reflect.TypeOf(data[0])
		v := reflect.ValueOf(data[0])

		//字段名
		var fieldNameArray []string

		//循环解析
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

			e.WhereExec = append(e.WhereExec, v.Field(i).Interface())
		}

		//拼接
		e.OrWhereParam += strings.Join(fieldNameArray, " and ") + ") "

	} else if dataType == 2 {
		//直接=的情况
		e.OrWhereParam += data[0].(string) + " = ?) "
		e.WhereExec = append(e.WhereExec, data[1])
	} else if dataType == 3 {
		//3个参数的情况

		//区分是操作符in的情况
		data2 := strings.Trim(strings.ToLower(data[1].(string)), " ")
		if data2 == "in" || data2 == "not in" {
			//判断传入的是切片
			reType := reflect.TypeOf(data[2]).Kind()
			if reType != reflect.Slice && reType != reflect.Array {
				panic("in/not in 操作传入的数据必须是切片或者数组")
			}

			//反射值
			v := reflect.ValueOf(data[2])
			//数组/切片长度
			dataNum := v.Len()
			//占位符
			ps := make([]string, dataNum)
			for i := 0; i < dataNum; i++ {
				ps[i] = "?"
				e.WhereExec = append(e.WhereExec, v.Index(i).Interface())
			}

			//拼接
			e.OrWhereParam += data[0].(string) + " " + data2 + " (" + strings.Join(ps, ",") + "))"

		} else {
			e.OrWhereParam += data[0].(string) + " " + data[1].(string) + " ?)"
			e.WhereExec = append(e.WhereExec, data[2])
		}
	}

	return e
}
