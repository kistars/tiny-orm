package main

import (
	"fmt"
	. "my-orm/orm"
	"testing"
)

func TestOrmEngine_Insert(t *testing.T) {
	_, _ = DB.Table("user").Insert(User{
		Name: "Tom",
		Age:  12,
	})
}

func TestOrmEngine_Where(t *testing.T) {
	DB.Table("userinfo").Where("username", "EE").OrWhere("departname", "22").OrWhere("status", 1)
	fmt.Println(DB.WhereParam + DB.OrWhereParam)
}

func TestOrmEngine_Select(t *testing.T) {

}
