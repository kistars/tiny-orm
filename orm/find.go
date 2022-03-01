package orm

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

//查询多条，返回值为struct切片
func (e *OrmEngine) Find(result interface{}) error {

	if reflect.ValueOf(result).Kind() != reflect.Ptr {
		return e.setErrorInfo(errors.New("参数请传指针变量！"))
	}

	if reflect.ValueOf(result).IsNil() {
		return e.setErrorInfo(errors.New("参数不能是空指针！"))
	}

	//拼接sql
	e.Prepare = "select * from " + e.GetTable()

	e.AllExec = e.WhereExec

	//query
	rows, err := e.Db.Query(e.Prepare, e.AllExec...)
	if err != nil {
		return e.setErrorInfo(err)
	}

	//读出查询出的列字段名
	column, err := rows.Columns()
	if err != nil {
		return e.setErrorInfo(err)
	}

	//values是每个列的值，这里获取到byte里
	values := make([][]byte, len(column))

	//因为每次查询出来的列是不定长的，用len(column)定住当次查询的长度
	scans := make([]interface{}, len(column))

	//原始struct的切片值
	destSlice := reflect.ValueOf(result).Elem()

	//原始单个struct的类型
	destType := destSlice.Type().Elem()

	for i := range values {
		scans[i] = &values[i]
	}

	//循环遍历
	for rows.Next() {

		dest := reflect.New(destType).Elem()

		if err := rows.Scan(scans...); err != nil {
			//query.Scan查询出来的不定长值放到scans[i] = &values[i],也就是每行都放在values里
			return e.setErrorInfo(err)
		}

		//遍历一行数据的各个字段
		for k, v := range values {
			//每行数据是放在values里面，现在把它挪到row里
			key := column[k]
			value := string(v)

			//遍历结构体
			for i := 0; i < destType.NumField(); i++ {

				//看下是否有sql别名
				sqlTag := destType.Field(i).Tag.Get("sql")
				var fieldName string
				if sqlTag != "" {
					fieldName = strings.Split(sqlTag, ",")[0]
				} else {
					fieldName = destType.Field(i).Name
				}

				//struct里没这个key
				if key != fieldName {
					continue
				}

				//反射赋值
				if err := e.reflectSet(dest, i, value); err != nil {
					return err
				}
			}
		}
		//赋值
		destSlice.Set(reflect.Append(destSlice, dest))
	}

	return nil
}

func (e *OrmEngine) reflectSet(dest reflect.Value, i int, value string) error {
	switch dest.Field(i).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		res, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return e.setErrorInfo(err)
		}
		dest.Field(i).SetInt(res)
	case reflect.String:
		dest.Field(i).SetString(value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		res, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return e.setErrorInfo(err)
		}
		dest.Field(i).SetUint(res)
	case reflect.Float32:
		res, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return e.setErrorInfo(err)
		}
		dest.Field(i).SetFloat(res)
	case reflect.Float64:
		res, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return e.setErrorInfo(err)
		}
		dest.Field(i).SetFloat(res)
	case reflect.Bool:
		res, err := strconv.ParseBool(value)
		if err != nil {
			return e.setErrorInfo(err)
		}
		dest.Field(i).SetBool(res)
	}
	return nil
}
