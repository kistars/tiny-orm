package main

import (
	"fmt"
	. "my-orm/orm"
	"testing"
)

func TestOrmEngine_Insert(t *testing.T) {
	_, _ = Engine.Table("user").Insert(User{
		Name: "Tom",
		Age:  12,
	})
}

func TestOrmEngine_Where(t *testing.T) {
	Engine.Table("userinfo").Where("username", "EE").OrWhere("departname", "22").OrWhere("status", 1)
	fmt.Println(Engine.WhereParam + Engine.OrWhereParam)
}

func TestOrmEngine_Select(t *testing.T) {

}
