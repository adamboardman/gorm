package gorm

import "sync"

type Transaction struct {
	once sync.Once
	rollback bool
	tx *DB
}

func (t *Transaction) Close () {
	t.once.Do(func() {
		if t.rollback {
			t.tx.Rollback()
		} else {
			t.tx.Commit()
		}
	})
}

func (t *Transaction) Fail() {
	t.rollback = true
}

func NewTransaction(db *DB) (*DB, *Transaction) {
	return db, &Transaction{tx: db}
}