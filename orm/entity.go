package orm

type User struct {
	Name string `sql:"name"`
	Age  int    `sql:"age"`
}
