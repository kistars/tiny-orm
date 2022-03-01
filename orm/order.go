package orm

import "strings"

//order排序
func (e *OrmEngine) Order(order ...string) *OrmEngine {
	orderLen := len(order)
	if orderLen%2 != 0 {
		panic("order by参数错误，请保证个数为偶数个")
	}

	//排序的个数
	orderNum := orderLen / 2

	//多次调用的情况
	if e.OrderParam != "" {
		e.OrderParam += ","
	}

	for i := 0; i < orderNum; i++ {
		keyString := strings.ToLower(order[i*2+1])
		if keyString != "desc" && keyString != "asc" {
			panic("排序关键字为：desc和asc")
		}
		if i < orderNum-1 {
			e.OrderParam += order[i*2] + " " + order[i*2+1] + ","
		} else {
			e.OrderParam += order[i*2] + " " + order[i*2+1]
		}
	}

	return e
}
