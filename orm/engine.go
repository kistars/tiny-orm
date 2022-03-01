package orm

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
	"os"
	"strconv"
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
