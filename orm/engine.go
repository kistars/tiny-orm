package orm

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
	"os"
	"strconv"
	"strings"
)

type OrmEngine struct {
	Db           *sql.DB // 数据库操作
	TableName    string
	Prepare      string
	AllExec      []interface{}
	Sql          string
	WhereParam   string
	LimitParam   string
	OrderParam   string
	OrWhereParam string
	WhereExec    []interface{}
	UpdateParam  string
	UpdateExec   []interface{}
	FieldParam   string
	TransStatus  int
	Tx           *sql.Tx // 事务
	GroupParam   string
	HavingParam  string
}

var Engine *OrmEngine
var DB *sql.DB

func init() {
	content, err := os.ReadFile("local.yaml")
	if err != nil {
		fmt.Println(err)
		panic("config file not found")
	}
	kv := map[string]string{}
	err = yaml.Unmarshal(content, &kv)
	if err != nil {
		panic("unmarshal failed")
	}
	Engine, err = NewMysql(kv["username"], kv["password"], kv["address"], kv["dbname"])
	if err != nil {
		panic("open db failed")
	}
}

//新建Mysql连接
func NewMysql(Username string, Password string, Address string, Dbname string) (*OrmEngine, error) {
	dsn := Username + ":" + Password + "@tcp(" + Address + ")/" + Dbname + "?charset=utf8&timeout=5s&readTimeout=6s"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	DB = db

	//最大连接数等配置，先占个位
	//db.SetMaxOpenConns(3)
	//db.SetMaxIdleConns(3)

	return &OrmEngine{
		Db:         DB,
		FieldParam: "*",
	}, nil
}

func (e *OrmEngine) Limit(param ...int) *OrmEngine {
	if len(param) == 1 {
		e.LimitParam = strconv.Itoa(param[0])
	} else if len(param) == 2 {
		e.LimitParam = strconv.Itoa(param[0]) + "," + strconv.Itoa(param[1])
	} else {
		panic("invalid parameters")
	}
	return e
}

// 查询指定字段
func (e *OrmEngine) Field(param string) *OrmEngine {
	e.FieldParam = param
	return e
}

//设置表名
func (e *OrmEngine) Table(name string) *OrmEngine {
	e.TableName = name

	//重置引擎
	e.resetOrmEngine()
	return e
}

//获取表名
func (e *OrmEngine) GetTable() string {
	return e.TableName
}

//重置引擎
func (e *OrmEngine) resetOrmEngine() {
	e = &OrmEngine{
		Db:         DB,
		FieldParam: "*",
	}
}

// 查询生成sql
func (e *OrmEngine) generateSql() {
	e.Sql = e.Prepare
	for _, v := range e.AllExec {
		switch v.(type) {
		case int:
			e.Sql = strings.Replace(e.Sql, "?", strconv.Itoa(v.(int)), 1)
		case int64:
			e.Sql = strings.Replace(e.Sql, "?", strconv.FormatInt(v.(int64), 10), 1)
		case bool:
			e.Sql = strings.Replace(e.Sql, "?", strconv.FormatBool(v.(bool)), 1)
		default:
			e.Sql = strings.Replace(e.Sql, "?", "'"+v.(string)+"'", 1)
		}
	}
}

func (e *OrmEngine) GetLastSql() string {
	e.generateSql()
	return e.Sql
}

// 执行原生sql
func (e *OrmEngine) Exec(sql string) (int64, error) {
	result, err := e.Db.Exec(sql)
	e.Sql = sql
	if err != nil {
		return 0, e.setErrorInfo(err)
	}

	//区分是insert还是其他(update,delete)
	if strings.Contains(sql, "insert") {
		lastInsertId, _ := result.LastInsertId()
		return lastInsertId, nil
	} else {
		rowsAffected, _ := result.RowsAffected()
		return rowsAffected, nil
	}
}

func (e *OrmEngine) Query(sql string) ([]map[string]string, error) {
	rows, err := e.Db.Query(sql)
	e.Sql = sql
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

	//最后得到的map
	results := make([]map[string]string, 0)
	for rows.Next() {
		if err := rows.Scan(scans...); err != nil {
			//query.Scan查询出来的不定长值放到scans[i] = &values[i],也就是每行都放在values里
			return nil, e.setErrorInfo(err)
		}

		row := make(map[string]string) //每行数据
		for k, v := range values {
			//每行数据是放在values里面，现在把它挪到row里
			key := column[k]
			row[key] = string(v)
		}
		results = append(results, row)
	}

	return results, nil
}
