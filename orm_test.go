package main

import (
	"fmt"
	. "my-orm/orm"
	"testing"
)

func TestOrmEngine_Insert(t *testing.T) {
	n, err := Engine.Table("user").Insert(User{
		Name: "Tom",
		Age:  12,
	})
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(n)
}

func TestOrmEngine_Where(t *testing.T) {
	Engine.Table("user").Where("username", "EE").OrWhere("departname", "22").OrWhere("status", 1)
	fmt.Println(Engine.WhereParam + Engine.OrWhereParam)
}

func TestOrmEngine_Select(t *testing.T) {
	i, err := Engine.Table("user").Select()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(i)
}

func TestOrmEngine_Update(t *testing.T) {
	i, err := Engine.Table("user").Where("name", "Tom").Update("name", "Jack")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(i)
}

func TestOrmEngine_Delete(t *testing.T) {
	i, err := Engine.Table("user").Where("name", "Tom").Delete()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(i)
}

func TestOrmEngine_Find(t *testing.T) {
	u := []User{}
	err := Engine.Table("user").Where("name", "Jack").Find(&u)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(u)
}
