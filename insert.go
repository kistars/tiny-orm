package my_orm

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

func (e *OrmEngine) Insert(data interface{}) (int64, error) {
	kind := reflect.ValueOf(data).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		return e.batchInsertData(data, "insert")
	case reflect.Struct:
		return e.insertData(data, "insert")
	default:
		return 0, errors.New("插入的数据格式不正确，单个插入格式为: struct，批量插入格式为: []struct")
	}
}

func (e *OrmEngine) Replace(data interface{}) (int64, error) {
	kind := reflect.ValueOf(data).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		return e.batchInsertData(data, "replace")
	case reflect.Struct:
		return e.insertData(data, "replace")
	default:
		return 0, errors.New("插入的数据格式不正确，单个插入格式为: struct，批量插入格式为: []struct")
	}
}

func (e *OrmEngine) insertData(data interface{}, insertType string) (int64, error) {
	// 反射
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	//字段名
	var fieldName []string

	//问号?占位符
	var placeholder []string

	//循环判断
	for i := 0; i < t.NumField(); i++ {

		//小写开头，无法反射，跳过
		if !v.Field(i).CanInterface() {
			continue
		}

		//解析tag，找出真实的sql字段名
		sqlTag := t.Field(i).Tag.Get("sql")
		if sqlTag != "" {
			//跳过自增字段
			if strings.Contains(strings.ToLower(sqlTag), "auto_increment") {
				continue
			} else {
				fieldName = append(fieldName, strings.Split(sqlTag, ",")[0])
				placeholder = append(placeholder, "?")
			}
		} else {
			fieldName = append(fieldName, t.Field(i).Name)
			placeholder = append(placeholder, "?")
		}

		//字段的值
		e.AllExec = append(e.AllExec, v.Field(i).Interface())
	}

	//拼接表，字段名，占位符
	e.Prepare = insertType + " into " + e.GetTable() + " (" + strings.Join(fieldName, ",") + ") values(" + strings.Join(placeholder, ",") + ")"
	fmt.Println("prepare =", e.Prepare)
	fmt.Println("exec =", e.AllExec)

	//prepare
	stmt, err := e.Db.Prepare(e.Prepare)
	if err != nil {
		return 0, e.setErrorInfo(err)
	}

	//执行exec,注意这是stmt.Exec
	result, err := stmt.Exec(e.AllExec...)
	if err != nil {
		return 0, e.setErrorInfo(err)
	}

	//获取自增ID
	id, _ := result.LastInsertId()
	return id, nil
}

//批量插入
func (e *OrmEngine) batchInsertData(batchData interface{}, insertType string) (int64, error) {

	//反射解析
	getValue := reflect.ValueOf(batchData)

	//切片大小
	l := getValue.Len()

	//字段名
	var fieldName []string

	//占位符
	var placeholderString []string

	//循环判断
	for i := 0; i < l; i++ {
		value := getValue.Index(i) // Value of item
		typed := value.Type()      // Type of item
		if typed.Kind() != reflect.Struct {
			panic("批量插入的子元素必须是结构体类型")
		}

		num := value.NumField()

		//子元素值
		var placeholder []string
		//循环遍历子元素
		for j := 0; j < num; j++ {

			//小写开头，无法反射，跳过
			if !value.Field(j).CanInterface() {
				continue
			}

			//解析tag，找出真实的sql字段名
			sqlTag := typed.Field(j).Tag.Get("sql")
			if sqlTag != "" {
				//跳过自增字段
				if strings.Contains(strings.ToLower(sqlTag), "auto_increment") {
					continue
				} else {
					//字段名只记录第一个的
					if i == 1 {
						fieldName = append(fieldName, strings.Split(sqlTag, ",")[0])
					}
					placeholder = append(placeholder, "?")
				}
			} else {
				//字段名只记录第一个的
				if i == 1 {
					fieldName = append(fieldName, typed.Field(j).Name)
				}
				placeholder = append(placeholder, "?")
			}

			//字段值
			e.AllExec = append(e.AllExec, value.Field(j).Interface())
		}

		//子元素拼接成多个()括号后的值
		placeholderString = append(placeholderString, "("+strings.Join(placeholder, ",")+")")
	}

	//拼接表，字段名，占位符
	e.Prepare = insertType + " into " + e.GetTable() + " (" + strings.Join(fieldName, ",") + ") values " + strings.Join(placeholderString, ",")

	//prepare
	var stmt *sql.Stmt
	var err error
	stmt, err = e.Db.Prepare(e.Prepare)
	if err != nil {
		return 0, e.setErrorInfo(err)
	}

	//执行exec,注意这是stmt.Exec
	result, err := stmt.Exec(e.AllExec...)
	if err != nil {
		return 0, e.setErrorInfo(err)
	}

	//获取自增ID
	id, _ := result.LastInsertId()
	return id, nil
}

func (e *OrmEngine) setErrorInfo(err error) error {
	_, file, line, _ := runtime.Caller(1)
	return errors.New("File: " + file + ":" + strconv.Itoa(line) + ", " + err.Error())
}
