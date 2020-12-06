package models

import (
	"encoding/base64"
	"html"
	mathRand "math/rand"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/nofendian17/rest/api/helpers"
	uuid "github.com/satori/go.uuid"
)

// Define salt size
const saltSize = 16

type Customer struct {
	CustomerID   uuid.UUID  `gorm:"type:varchar(64);primary_key;unique" json:"customer_id" validate:"-"`
	CustomerName string     `gorm:"size:80;not null" json:"customer_name" validate:"required"`
	Email        string     `gorm:"size:50;not null;unique" json:"email" validate:"required,email"`
	PhoneNumber  string     `gorm:"size:20;not null;unique" json:"phone_number" validate:"required"`
	DOB          *time.Time `gorm:"not null" time_format:"sql_date" time_utc:"true" json:"dob" validate:"required"`
	Sex          *bool      `gorm:"not null" json:"sex" validate:"required"`
	Salt         string     `gorm:"size:80;not null" json:"salt"`
	Password     string     `gorm:"not null" json:"password"`
	CreatedAt    time.Time  `gorm:"default:CURRENT_TIMESTAMP" time_format:"sql_datetime" time_location:"UTC" json:"created_at" validate:"required"`
}

type GenerateCredential struct {
	CustomerID  string `json:"customer_id,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Password    string `json:"password,omitempty"`
}

type CustomerAuth struct {
	CustomerID uuid.UUID
	Email      string `json:"email" validate:"required,email"`
	Salt       string
	Password   string `json:"password" validate:"required"`
}

type CustomerToken struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func (Customer) TableName() string {
	return "customers"
}

func (CustomerAuth) TableName() string {
	return "customers"
}

func (customer *Customer) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("customer_id", interface{}(uuid.NewV4()))
}

func (customer *Customer) Prepare() {
	customer.CustomerName = html.EscapeString(strings.TrimSpace(customer.CustomerName))
	customer.Email = html.EscapeString(strings.TrimSpace(customer.Email))
	customer.CreatedAt = time.Now()
}

// func VerifyPassword(salt string, hashPassword string, plainPassword string) error {
// 	combination := salt + plainPassword
// 	hash := sha256.New()
// 	io.WriteString(hash, combination)
// 	fmt.Println(plainPassword)
// 	fmt.Println(hashPassword)
// 	fmt.Printf("%x \n", hash.Sum(nil))
// 	match := bytes.Equal(hash.Sum(nil), []byte(hashPassword))
// 	if !match {
// 		return errors.New("password doesnt match.")
// 	}
// 	return nil
// }

func (customer *Customer) FindAllCustomers(db *gorm.DB) (*[]Customer, error) {
	var err error
	customers := []Customer{}
	err = db.Debug().Model(&Customer{}).Limit(100).Find(&customers).Error
	if err != nil {
		return &[]Customer{}, err
	}
	return &customers, err
}

func (customer *Customer) SaveCustomer(db *gorm.DB) (*GenerateCredential, error) {
	var err error
	mathRand.Seed(time.Now().UnixNano())
	password := helpers.RandomString(12)

	// First generate random 16 byte salt
	var salt = helpers.GenerateRandomSalt(saltSize)

	// Hash password using the salt
	var hashedPassword = helpers.HashPassword(password, salt)

	// fmt.Println("Password Hash:", hashedPassword)
	// fmt.Println("Salt:", salt)
	// fmt.Println("Salt string as base64 :", base64.URLEncoding.EncodeToString(salt))
	encodeSalt := base64.URLEncoding.EncodeToString(salt)
	// decodeSalt, err := base64.URLEncoding.DecodeString(encodeSalt)
	// if err != nil {
	// 	fmt.Println("error:", err)
	// }
	// fmt.Println("Salt string base64 decode to byte :", decodeSalt)

	// Check if passed password matches the original password by hashing it
	// with the original password's salt and check if the hashes match
	// fmt.Println("Password Match:", helpers.DoPasswordsMatch(hashedPassword, password, salt))

	customer.Salt = encodeSalt
	customer.Password = hashedPassword

	credential := &GenerateCredential{
		Email:       customer.Email,
		PhoneNumber: customer.PhoneNumber,
		Password:    password,
	}
	err = db.Debug().Create(&customer).Error
	if err != nil {
		return &GenerateCredential{}, err
	}

	return credential, nil
}
