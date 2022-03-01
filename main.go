package main

import (
	"fmt"
	. "my-orm/orm"
)

func main() {
	DB.Table("userinfo").Where("username", "EE").OrWhere("departname", "22").OrWhere("status", 1)
	fmt.Println(DB.WhereParam + DB.OrWhereParam)
}
