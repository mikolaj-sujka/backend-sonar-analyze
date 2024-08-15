package models

import (
	"gorm.io/gorm"
)

func PricedAbove(price float64) func(*gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("price > ?", price)
    }
}

func InCategory(categoryID uint) func(*gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("category_id = ?", categoryID)
    }
}