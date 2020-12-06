package controllers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/nofendian17/rest/api/auth"
	"github.com/nofendian17/rest/api/models"
	"github.com/nofendian17/rest/api/responses"
	formaterror "github.com/nofendian17/rest/api/utils"
)

func (server *Server) CreateOrder(w http.ResponseWriter, r *http.Request) {
	tokenAuth, err := auth.ExtractTokenMetadata(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	customerId, err := auth.FetchAuth(server.Redis, tokenAuth)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	fmt.Println("debug at CreateOrder -> customer_id : ", customerId)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	orderRequest := models.OrderRequest{}
	err = json.Unmarshal(body, &orderRequest)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	validate := validator.New()
	err = validate.Struct(orderRequest)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	for _, v := range orderRequest.Orders {
		fmt.Println("product_id : ", v.ProductId)
		fmt.Println("qty : ", v.Qty)
	}

	fmt.Println("payment_method :", orderRequest.PaymentMethod)

	// bind customer id
	orderRequest.CustomerId = customerId

	orderCreated, err := orderRequest.SaveOrder(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	// destroy token for 1 time usage
	au, err := auth.ExtractTokenMetadata(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	//fmt.Println(au)
	deleted, delErr := auth.DeleteAuth(server.Redis, au.AccessUuid)
	if delErr != nil || deleted == 0 { //if any goes wrong
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, orderCreated))
	responses.JSON(w, http.StatusCreated, orderCreated)
}
