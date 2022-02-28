package my_orm

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
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

//新建Mysql连接
func NewMysql(Username string, Password string, Address string, Dbname string) (*OrmEngine, error) {
	dsn := Username + ":" + Password + "@tcp(" + Address + ")/" + Dbname + "?charset=utf8&timeout=5s&readTimeout=6s"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	//最大连接数等配置，先占个位
	//db.SetMaxOpenConns(3)
	//db.SetMaxIdleConns(3)

	return &OrmEngine{
		Db:         db,
		FieldParam: "*",
	}, nil
}

func (e *OrmEngine) Limit(param string) *OrmEngine {
	e.LimitParam = param
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

}
