package repository

import "gorm.io/gorm"

type transactionManager struct {
	db *gorm.DB
}

func (t *transactionManager) begin() (*gorm.DB, error) {
	return t.db.Begin(), nil
}

func (t *transactionManager) commit(tx *gorm.DB) error {
	return tx.Commit().Error
}

func (t *transactionManager) rollback(tx *gorm.DB) error {
	return tx.Rollback().Error
}
