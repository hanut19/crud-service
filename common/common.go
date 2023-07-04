package common

import (
	"net/http"
)

type User struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role" validate:"required"`
}

type Authentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	Role        string `json:"role"`
	Email       string `json:"email"`
	TokenString string `json:"token"`
}

type Error struct {
	IsError bool   `json:"isError"`
	Message string `json:"message"`
}

type APIResponse struct {
	IsError    bool        `json:"isError"`
	Message    string      `json:"message"`
	Result     interface{} `json:"result"`
	StatusCode int         `json:"statusCode"`
}

type Product struct {
	Id          string `bson:"_id,omitempty" json:"id,omitempty"`
	ProductId   string `bson:"productid,omitempty" json:"productid" validate:"required"`
	Name        string `bson:"name,omitempty" json:"name" validate:"required"`
	Description string `bson:"description,omitempty" json:"description"`
	Price       string `bson:"price,omitempty" json:"price" validate:"required"`
	File        []byte `bson:"file,omitempty" json:"file"`
	FileType    string `bson:"filetype,omitempty" json:"filetype"`
}

type RespProduct struct {
	ProductId   string `json:"productid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
	File        string `json:"file"`
}

func GetHost(r *http.Request) string {
	customURL := ""
	if r.TLS == nil {
		customURL = "http"
	} else {
		customURL = "https"
	}
	return customURL + "://" + r.Host
}
