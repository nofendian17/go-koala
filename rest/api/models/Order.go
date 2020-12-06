package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/nofendian17/rest/api/helpers"
	uuid "github.com/satori/go.uuid"
)

type OrderRequest struct {
	Orders        []*OrdersItem `json:"orders" validate:"required,dive,required"`
	PaymentMethod string        `json:"payment_method" validate:"required"`
	CustomerId    string
}

type OrdersItem struct {
	ProductId string `json:"product_id" validate:"required"`
	Qty       int    `json:"qty" validate:"required,numeric"`
}

type Order struct {
	OrderId         uuid.UUID `gorm:"type:varchar(64);primary_key;unique" json:"order_id" validate:"-"`
	CustomerId      string    `gorm:"type:varchar(64);size:64;not null" json:"customer_id"`
	OrderNumber     string    `gorm:"type:varchar(40);size:40;not null" json:"order_number"`
	OrderDate       time.Time `gorm:"default:CURRENT_TIMESTAMP" time_format:"sql_datetime" time_location:"UTC" json:"order_date"`
	PaymentMethodId string    `gorm:"type:varchar(64);size:40;not null" json:"payment_method_id"`
}

type OrderResponse struct {
	OrderId         string                    `json:"order_id,omitempty"`
	OrderNumber     string                    `json:"order_number,omitempty"`
	OrderDate       string                    `json:"order_date,omitempty"`
	PaymentMethodId string                    `json:"payment_method_id,omitempty"`
	Order           []OrderDetailItemResponse `json:"orders,omitempty"`
}

type OrderDetailItemResponse struct {
	ProductId string `json:"product_id"`
	Qty       int    `json:"qty"`
}

type CountOrderMonth struct {
	Counter int
}

func (Order) TableName() string {
	return "orders"
}

func (order *Order) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("order_id", interface{}(uuid.NewV4()))
}

func (orderRequest *OrderRequest) SaveOrder(db *gorm.DB) (*OrderResponse, error) {
	var err error
	orderResponse := &OrderResponse{}

	// save to table orders
	order := &Order{}
	currentTime := time.Now()

	// generate letter code by created_at per month
	var counterMonth CountOrderMonth

	currentMonth := helpers.CurrentMonth()

	err = db.Debug().Raw("SELECT COUNT(order_date) AS counter FROM orders WHERE MONTH(order_date) = ?", currentMonth).Scan(&counterMonth).Error

	if err != nil {
		return &OrderResponse{}, err
	}

	counter := counterMonth.Counter + 1

	generateLetters := helpers.GenerateLetter(counter)

	fmt.Println("debug at SaveOrder generateLetters -> order_number :", generateLetters)

	// order.OrderId = ""
	order.CustomerId = orderRequest.CustomerId
	// order.OrderNumber = "" // should be replace by format PO-123/IX/2020 (IX is current month)(2020 is current year),(123 reset per month).
	order.OrderNumber = generateLetters
	order.OrderDate = currentTime
	order.PaymentMethodId = orderRequest.PaymentMethod

	err = db.Debug().Create(&order).Error
	if err != nil {
		return &OrderResponse{}, err
	}

	// save to table order details
	orderDetail := OrderDetail{}
	for _, v := range orderRequest.Orders {
		// orderDetail.OrderDetailId = ""
		orderDetail.OrderId = order.OrderId.String()
		orderDetail.ProductId = v.ProductId
		orderDetail.Qty = v.Qty
		orderDetail.CreatedAt = currentTime
		err = db.Debug().Create(&orderDetail).Error
		if err != nil {
			return &OrderResponse{}, err
		}
	}

	// bug if insert slice bulk -> field value not valid :)
	od := []OrderDetail{}
	err = db.Debug().Model(OrderDetail{}).Where("order_id = ?", order.OrderId.String()).Find(&od).Error
	fmt.Println(&od)
	if err != nil {
		return &OrderResponse{}, err
	}

	var odir = []OrderDetailItemResponse{}
	for _, rod := range od {
		var odi OrderDetailItemResponse
		odi.ProductId = rod.ProductId
		odi.Qty = rod.Qty
		odir = append(odir, odi)
	}

	orderResponse.OrderId = order.OrderId.String()
	orderResponse.OrderNumber = order.OrderNumber
	orderResponse.OrderDate = order.OrderDate.String()
	orderResponse.PaymentMethodId = order.PaymentMethodId
	orderResponse.Order = odir

	return orderResponse, nil
}
