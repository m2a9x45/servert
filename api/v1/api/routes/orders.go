package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"../database"
	"../models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"
)

func generateOrderID() string {

	id := ksuid.New()
	orderID := "order_" + id.String()

	println(orderID)

	result, err := database.DBCon.Query("SELECT order_id from orders WHERE order_id=(?)", orderID)
	if err != nil {
		println(err)
	}

	defer result.Close()

	for result.Next() {
		var orderid string
		err := result.Scan(&orderid)
		if err != nil {
			panic(err)
		}
		if orderid != "" {
			// make new user id
			println(orderid, "already exists")
			generateOrderID()
		}
	}

	return orderID
}

func CreatePaymentIntent(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["prodID"]

	if id != "" {
		var allproducts []models.Product

		result, err := database.DBCon.Query("SELECT * from products WHERE prod_id=(?)", id)
		if err != nil {
			println(err)
			res := models.ResObj{Success: false, Message: "Couldn't find product"}
			json.NewEncoder(w).Encode(res)
			return
		}

		defer result.Close()

		for result.Next() {
			var product models.Product
			err := result.Scan(&product.ID, &product.UIDD, &product.Name, &product.Des, &product.CPU, &product.RAM, &product.Disk, &product.Price, &product.Instock, &product.Setupfee, &product.Discount)
			if err != nil {
				panic(err.Error())
			}
			allproducts = append(allproducts, product)
		}

		if len(allproducts) == 0 {
			w.WriteHeader(http.StatusNotFound)
			res := models.ResObj{Success: false, Message: "product not found"}
			json.NewEncoder(w).Encode(res)
			return
		}

		stripe.Key = "sk_test_OGXIlmLXL1Gvhpa9jqBdxutN00YB96uOjP"

		price := allproducts[0].Price
		pricePennies := int64(price * 100)
		println(price, " in £")
		println(price, " in p")

		params := &stripe.PaymentIntentParams{
			Amount:   stripe.Int64(pricePennies),
			Currency: stripe.String(string(stripe.CurrencyGBP)),
		}

		intent, _ := paymentintent.New(params)

		data := models.CheckoutData{
			ClientSecret: intent.ClientSecret,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)

	}
}

func Order(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodOptions {
		return
	}

	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tknStr := c.Value

	claims := &models.Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	order := models.OrderData{}

	jsn, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Error wilst reading r body", err)
	}

	err = json.Unmarshal(jsn, &order)
	if err != nil {
		log.Fatal("Error wilst unmarshaling json", err)
	}

	OrderID := generateOrderID()

	println(OrderID, "seletced")

	result, err := database.DBCon.Query("INSERT INTO orders (order_id, user_id, payment_id, prod_id) VALUES (?,?,?,?)", OrderID, claims.UserID, order.PaymentID, order.ProductID)
	if err != nil {
		log.Fatal("Error wilst inserting into DB", err)
	}

	defer result.Close()

	fmt.Println("Inserted Into DB")

	res := models.ResObj{Success: true, Message: "Details inserted into DB"}

	json.NewEncoder(w).Encode(res)

}

// GetOrders will return a list of orders for an account
func GetOrders(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tknStr := c.Value

	claims := &models.Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	result, err := database.DBCon.Query("SELECT order_id, prod_id from orders WHERE user_id=(?)", claims.UserID)
	if err != nil {
		println(err)
	}

	var allorders []models.OrderObj

	for result.Next() {
		var order models.OrderObj
		err := result.Scan(&order.OrderID, &order.ProdID)
		if err != nil {
			panic(err.Error())
		}
		allorders = append(allorders, order)
	}

	println(len(allorders))

	if len(allorders) == 0 {
		w.WriteHeader(http.StatusNotFound)
		res := models.ResObj{Success: false, Message: "No orders found if you think this is wrong please contact us"}
		json.NewEncoder(w).Encode(res)
		return
	}

	fmt.Println(allorders)
	json.NewEncoder(w).Encode(allorders)
}
