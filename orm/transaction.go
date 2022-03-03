package orm

func (e *OrmEngine) Begin() error {
	tx, err := e.Db.Begin()
	if err != nil {
		return e.setErrorInfo(err)
	}
	e.Tx = tx
	e.TransStatus = 1
	return nil
}

func (e *OrmEngine) Rollback() error {
	e.TransStatus = 0
	return e.Tx.Rollback()
}

/*
  确认提交表示我们所有的执行都是 OK 的，这个时候我们需要向 mysql 服务器发出确认提交指令，
  它才会真正意义上将 sql 给执行。如果不执行这个指令，实际上数据并不会执行。
  所以，我们最后一定不要忘记执行确认提交操作
*/
func (e *OrmEngine) Commit() error {
	e.TransStatus = 0
	return e.Tx.Commit()
}
