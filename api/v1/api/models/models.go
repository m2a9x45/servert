package models

import (
	jwt "github.com/dgrijalva/jwt-go"
	stripe "github.com/stripe/stripe-go"
)

type Product struct {
	ID       string  `json:"id"`
	UIDD     string  `json:"uuid"`
	Name     string  `json:"name"`
	Des      string  `json:"des"`
	CPU      string  `json:"cpu"`
	RAM      string  `json:"ram"`
	Disk     string  `json:"disk"`
	Price    float64 `json:"price"`
	Instock  bool    `json:"instock"`
	Setupfee float64 `json:"setupfee"`
	Discount float64 `json:"discount"`
}

type Details struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ResObj struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type SignUpObj struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

type CheckoutData struct {
	ClientSecret string `json:"clientecret"`
}

type OrderData struct {
	PaymentID string `json="PaymentID"`
	ProductID string `json="ProductID"`
	Dur       string `json=""Dur`
}

type OrderObj struct {
	OrderID string `json="order_id"`
	ProdID  string `json="prod_id"`
}

type Order struct {
	ID        string `json:"id"`
	OrderID   string `json="order_id"`
	UserID    string `json="userid"`
	PaymentID string `json="PaymentID"`
	ProdID    string `json="prod_id"`
	Time      int    `json="createdAt"`
	Dur       int    `json=""Dur`
	Expires   int    `josn="expires"`
}

type Receipt struct {
	Success bool    `json:"success"`
	Order   Order   `json:"order"`
	Product Product `json:"product"`
}

type UserDetails struct {
	Name  string `json="name"`
	Email string `json="email"`
}

type Hello struct {
	Name string `json="name"`
}

type Reset struct {
	Email string `json="email"`
}

type ResetPassword struct {
	Password string `json:"password"`
	Token    string `json:"token"`
}

type ResetPasswordDB struct {
	Email   string `json="email"`
	Expires int64  `json="expires"`
}

type Card struct {
	ID          string                        `json="id"`
	Fingerprint string                        `json="fingerprint"`
	Brand       stripe.PaymentMethodCardBrand `json="brand"`
	Exp_month   uint64                        `json="exp_month"`
	Exp_year    uint64                        `json="exp_year"`
	Last4       string                        `json="last4"`
}
