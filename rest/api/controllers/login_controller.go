package controllers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator"
	"github.com/nofendian17/rest/api/auth"
	"github.com/nofendian17/rest/api/helpers"
	"github.com/nofendian17/rest/api/models"
	"github.com/nofendian17/rest/api/responses"
	formaterror "github.com/nofendian17/rest/api/utils"
)

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	customer := models.CustomerAuth{}
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

	token, err := server.SignIn(customer.Email, customer.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}

	customerToken := models.CustomerToken{}
	customerToken.AccessToken = token.AccessToken
	customerToken.RefreshToken = token.RefreshToken
	responses.JSON(w, http.StatusOK, customerToken)
}

func (server *Server) SignIn(email, password string) (*models.TokenDetails, error) {

	var err error

	customer := models.CustomerAuth{}

	err = server.DB.Debug().Model(models.CustomerAuth{}).Where("email = ?", email).Take(&customer).Error
	if err != nil {
		return &models.TokenDetails{}, err
	}

	hashedPassword := customer.Password

	salt, err := base64.URLEncoding.DecodeString(customer.Salt)
	if err != nil {
		return &models.TokenDetails{}, err
	}

	authPassword := helpers.DoPasswordsMatch(hashedPassword, password, salt)
	if !authPassword {
		return &models.TokenDetails{}, err
	}

	ts, err := auth.CreateToken(customer.CustomerID.String())
	if err != nil {
		return &models.TokenDetails{}, err
	}

	saveErr := auth.CreateAuth(server.Redis, customer.CustomerID.String(), ts)
	if saveErr != nil {
		return &models.TokenDetails{}, err
	}

	td := &models.TokenDetails{
		AccessToken:  ts.AccessToken,
		RefreshToken: ts.RefreshToken,
	}

	return td, nil
}

func (server *Server) Refresh(w http.ResponseWriter, r *http.Request) {
	var err error
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	refreshTokenRequest := models.RefreshTokenRequest{}
	err = json.Unmarshal(body, &refreshTokenRequest)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	validate := validator.New()
	err = validate.Struct(refreshTokenRequest)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	os.Setenv("REFRESH_SECRET", "zzzwwwasd") //this should be in an env file
	token, err := jwt.Parse(refreshTokenRequest.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	//if there is an error, the token must have expired
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}
		// customerId, err := fmt.Sprintf("%.f", claims["customer_id"])
		customerId, ok := claims["customer_id"].(string)
		if !ok {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}
		//Delete the previous Refresh Token
		deleted, delErr := auth.DeleteAuth(server.Redis, refreshUuid)
		if delErr != nil || deleted == 0 { //if any goes wrong
			responses.ERROR(w, http.StatusUnauthorized, err)
			return
		}
		//Create new pairs of refresh and access tokens
		ts, createErr := auth.CreateToken(customerId)
		if createErr != nil {
			responses.ERROR(w, http.StatusForbidden, err)
			return
		}
		//save the tokens metadata to redis
		saveErr := auth.CreateAuth(server.Redis, customerId, ts)
		if saveErr != nil {
			responses.ERROR(w, http.StatusForbidden, err)
			return
		}
		refreshTokenResponse := models.RefreshTokenResponse{}
		refreshTokenResponse.AccesToken = ts.AccessToken
		refreshTokenResponse.RefreshToken = ts.RefreshToken
		responses.JSON(w, http.StatusOK, refreshTokenResponse)
	} else {
		responses.ERROR(w, http.StatusUnauthorized, err)
	}
}

func (server *Server) LogOut(w http.ResponseWriter, r *http.Request) {
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
	fmt.Println("debug at LogOut -> customer_id : ", customerId)

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
	fmt.Println(deleted)
	responses.JSON(w, http.StatusOK, map[string]string{"message": "Successfully logged out"})
}
