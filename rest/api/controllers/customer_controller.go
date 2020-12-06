package controllers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/liamylian/jsontime"
	"github.com/nofendian17/rest/api/auth"
	"github.com/nofendian17/rest/api/models"
	"github.com/nofendian17/rest/api/responses"
	formaterror "github.com/nofendian17/rest/api/utils"
)

var json = jsontime.ConfigWithCustomTimeFormat

func (server *Server) CreateCustomer(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	customer := models.Customer{}
	err = json.Unmarshal(body, &customer)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	validate := validator.New()
	err = validate.Struct(customer)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	customer.Prepare()

	CustomerCreated, err := customer.SaveCustomer(server.DB)

	if err != nil {

		formattedError := formaterror.FormatError(err.Error())

		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, CustomerCreated.CustomerID))
	responses.JSON(w, http.StatusCreated, CustomerCreated)
}

func (server *Server) GetCustomers(w http.ResponseWriter, r *http.Request) {

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
	fmt.Println("debug at GetCustomers -> customer_id : ", customerId)
	customer := models.Customer{}

	customers, err := customer.FindAllCustomers(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, customers)
}

// func (server *Server) GetCustomer(w http.ResponseWriter, r *http.Request) {

// 	vars := mux.Vars(r)
// 	uid, err := strconv.ParseUint(vars["id"], 10, 32)
// 	if err != nil {
// 		responses.ERROR(w, http.StatusBadRequest, err)
// 		return
// 	}
// 	Customer := models.Customer{}
// 	CustomerGotten, err := Customer.FindCustomerByID(server.DB, uint32(uid))
// 	if err != nil {
// 		responses.ERROR(w, http.StatusBadRequest, err)
// 		return
// 	}
// 	responses.JSON(w, http.StatusOK, CustomerGotten)
// }

// func (server *Server) UpdateCustomer(w http.ResponseWriter, r *http.Request) {

// 	vars := mux.Vars(r)
// 	uid, err := strconv.ParseUint(vars["id"], 10, 32)
// 	if err != nil {
// 		responses.ERROR(w, http.StatusBadRequest, err)
// 		return
// 	}
// 	body, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		responses.ERROR(w, http.StatusUnprocessableEntity, err)
// 		return
// 	}
// 	Customer := models.Customer{}
// 	err = json.Unmarshal(body, &Customer)
// 	if err != nil {
// 		responses.ERROR(w, http.StatusUnprocessableEntity, err)
// 		return
// 	}
// 	tokenID, err := auth.ExtractTokenID(r)
// 	if err != nil {
// 		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
// 		return
// 	}
// 	if tokenID != uint32(uid) {
// 		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
// 		return
// 	}
// 	Customer.Prepare()
// 	err = Customer.Validate("update")
// 	if err != nil {
// 		responses.ERROR(w, http.StatusUnprocessableEntity, err)
// 		return
// 	}
// 	updatedCustomer, err := Customer.UpdateACustomer(server.DB, uint32(uid))
// 	if err != nil {
// 		formattedError := formaterror.FormatError(err.Error())
// 		responses.ERROR(w, http.StatusInternalServerError, formattedError)
// 		return
// 	}
// 	responses.JSON(w, http.StatusOK, updatedCustomer)
// }

// func (server *Server) DeleteCustomer(w http.ResponseWriter, r *http.Request) {

// 	vars := mux.Vars(r)

// 	Customer := models.Customer{}

// 	uid, err := strconv.ParseUint(vars["id"], 10, 32)
// 	if err != nil {
// 		responses.ERROR(w, http.StatusBadRequest, err)
// 		return
// 	}
// 	tokenID, err := auth.ExtractTokenID(r)
// 	if err != nil {
// 		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
// 		return
// 	}
// 	if tokenID != 0 && tokenID != uint32(uid) {
// 		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
// 		return
// 	}
// 	_, err = Customer.DeleteACustomer(server.DB, uint32(uid))
// 	if err != nil {
// 		responses.ERROR(w, http.StatusInternalServerError, err)
// 		return
// 	}
// 	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
// 	responses.JSON(w, http.StatusNoContent, "")
// }
