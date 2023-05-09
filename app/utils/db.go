package utils

import "gorm.io/gorm"

func DBExists[T any](db *gorm.DB, query any, where ...any) (bool, error) {
	var count int64
	if err := db.Model(new(T)).Where(query, where...).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func DBExistsByRaw(db *gorm.DB, query string, where ...interface{}) (bool, error) {
	var count int64
	if err := db.Raw(query, where...).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
