package test

import (
	"fmt"
	"my-orm"
	"testing"
)

func TestOrmEngine_Insert(t *testing.T) {
	orm, err := my_orm.NewMysql("root", "451244939", "localhost:3306", "mydb")
	if err != nil {
		t.Error(err)
		return
	}

	_, _ = orm.Table("user").Insert(my_orm.User{
		Name: "Tom",
		Age:  12,
	})
}

func TestOrmEngine_Where(t *testing.T) {
	orm, err := my_orm.NewMysql("root", "451244939", "localhost:3306", "mydb")
	if err != nil {
		t.Error(err)
		return
	}

	orm.Table("userinfo").Where("username", "EE").OrWhere("departname", "22").OrWhere("status", 1)
	fmt.Println(orm.WhereParam + orm.OrWhereParam)
}
