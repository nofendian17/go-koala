package models

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type OrderDetail struct {
	OrderDetailId uuid.UUID `gorm:"type:varchar(64);primary_key;unique"`
	OrderId       string    `gorm:"type:varchar(64);not null;foreignkey:OrderId;association_foreignkey:OrderId" `
	ProductId     string    `gorm:"type:varchar(64);not null"`
	Qty           int       `gorm:"type:integer(64);not null"`
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP;not null" time_format:"sql_datetime" time_location:"UTC"`
}

func (OrderDetail) TableName() string {
	return "order_details"
}

func (orderDetail *OrderDetail) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("order_detail_id", interface{}(uuid.NewV4()))
}
