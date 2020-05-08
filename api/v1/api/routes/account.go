package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"../database"
	"../models"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/gorilla/mux"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentmethod"
)

var jwtKey []byte

func checkifemailexistsreg(email string) bool {

	result, err := database.DBCon.Query("SELECT email from reg WHERE email=(?)", email)
	if err != nil {
		println(err)
	}

	defer result.Close()

	for result.Next() {
		var emailfound string
		err := result.Scan(&emailfound)
		if err != nil {
			panic(err)
		}
		if emailfound != "" {
			// make new user id
			println(emailfound, "already exists")
			return true
		}
	}

	return false

}

func Intrest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		return
	}

	details := models.Details{}

	jsn, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Error wilst reading r body", err)
	}

	err = json.Unmarshal(jsn, &details)
	if err != nil {
		log.Fatal("Error wilst unmarshaling json", err)
	}

	log.Println(details.Email, details.Name)

	data := data{Name: details.Name, Email: details.Email}

	err = data.Validate()

	if err != nil {
		fmt.Println(err)
		res := models.ResObj{Success: false, Message: "something went wrong"}
		json.NewEncoder(w).Encode(res)
		return
	}

	emailInDB := checkifemailexistsreg(details.Email)
	if emailInDB != false {
		res := models.ResObj{Success: false, Message: "Email already in use"}
		json.NewEncoder(w).Encode(res)
		return
	}

	result, err := database.DBCon.Query("INSERT INTO reg (name, email) VALUES (?,?)", details.Name, details.Email)
	if err != nil {
		log.Fatal("Error wilst inserting into DB", err)
	}

	defer result.Close()

	log.Println(result)

	res := models.ResObj{Success: true, Message: "Details inserted into DB"}

	json.NewEncoder(w).Encode(res)
}

func Account(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value("user")
	claims := uid.(*models.Claims)

	res := models.ResObj{Success: true, Message: claims.UserID}

	json.NewEncoder(w).Encode(res)
}

func AccountInfo(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value("user")
	claims := uid.(*models.Claims)

	result, err := database.DBCon.Query("SELECT name, email from users WHERE user_id=(?)", claims.UserID)
	if err != nil {
		println(err)
	}

	var userDetailsList []models.UserDetails

	for result.Next() {
		var user models.UserDetails
		err := result.Scan(&user.Name, &user.Email)
		if err != nil {
			panic(err.Error())
		}
		userDetailsList = append(userDetailsList, user)
	}

	println(len(userDetailsList))

	if len(userDetailsList) == 0 {
		w.WriteHeader(http.StatusNotFound)
		res := models.ResObj{Success: false, Message: "No orders found if you think this is wrong please contact us"}
		json.NewEncoder(w).Encode(res)
		return
	}

	fmt.Println(userDetailsList)
	json.NewEncoder(w).Encode(userDetailsList)

}

func Getreceipt(w http.ResponseWriter, r *http.Request) {

	// Look into further but I think an authed user can access a recipt that isn't theirs

	if r.Method == http.MethodOptions {
		return
	}

	// uid := r.Context().Value("user")
	// claims := uid.(*models.Claims)

	// fmt.Println(claims.UserID)

	vars := mux.Vars(r)
	id := vars["receiptID"]

	if id != "" {

		var allorders []models.Order

		result, err := database.DBCon.Query("SELECT * from orders WHERE order_id=(?)", id)
		if err != nil {
			println(err)
		}

		defer result.Close()

		for result.Next() {
			var order models.Order
			err := result.Scan(&order.ID, &order.OrderID, &order.UserID, &order.PaymentID, &order.ProdID, &order.Time, &order.Dur, &order.Expires)
			if err != nil {
				panic(err.Error())
			}
			allorders = append(allorders, order)
		}

		println(len(allorders))

		if len(allorders) == 0 {
			w.WriteHeader(http.StatusNotFound)
			res := models.ResObj{Success: false, Message: "receipt not found"}
			json.NewEncoder(w).Encode(res)
			return
		}

		fmt.Println(allorders)

		// get product

		var allproducts []models.Product

		result, err = database.DBCon.Query("SELECT * from products WHERE prod_id=(?)", allorders[0].ProdID)
		if err != nil {
			println(err)
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

		println(len(allproducts))

		if len(allproducts) == 0 {
			w.WriteHeader(http.StatusNotFound)
			res := models.ResObj{Success: false, Message: "product not found"}
			json.NewEncoder(w).Encode(res)
		}

		fmt.Println(allproducts)

		receipt := models.Receipt{Success: true, Order: allorders[0], Product: allproducts[0]}
		json.NewEncoder(w).Encode(receipt)
		return
	}

	res := models.ResObj{Success: false, Message: "Sorry we couldn't find your receipts"}
	json.NewEncoder(w).Encode(res)

}

func Getcustomercards(w http.ResponseWriter, r *http.Request) {

	uid := r.Context().Value("user")
	claims := uid.(*models.Claims)
	// get customer ID from DB using JWT userID

	var customerID string

	result, err := database.DBCon.Query("SELECT stripe_id FROM users WHERE user_id=(?)", claims.UserID)
	if err != nil {
		println(err)
	}

	for result.Next() {
		var custID string
		err := result.Scan(&custID)
		if err != nil {
			panic(err.Error())
		}
		customerID = custID
	}

	stripeKey, exists := os.LookupEnv("STRIPE_KEY")

	if exists {
		stripe.Key = stripeKey
	}

	// params := &stripe.CardListParams{
	// 	Customer: stripe.String(customerID),
	// }
	// params.Filters.AddFilter("limit", "", "3")
	// i := card.List(params)

	params := &stripe.PaymentMethodListParams{
		Customer: stripe.String(customerID),
		Type:     stripe.String("card"),
	}

	cards := []models.Card{}

	i := paymentmethod.List(params)
	for i.Next() {
		pm := i.PaymentMethod()
		fmt.Println(pm.Card.Fingerprint)
		card := models.Card{ID: pm.ID, Fingerprint: pm.Card.Fingerprint, Last4: pm.Card.Last4, Brand: pm.Card.Brand, Exp_month: pm.Card.ExpMonth, Exp_year: pm.Card.ExpYear}
		cards = append(cards, card)
	}

	// for i.Next() {
	// 	c := i.Card()
	// 	fmt.Println(c.ID)
	// 	card := models.Card{ID: c.ID, Last4: c.Last4, Brand: c.Brand, Exp_month: c.ExpMonth, Exp_year: c.ExpYear}
	// 	cards = append(cards, card)
	// }

	// defult card shows as index 0 the first

	json.NewEncoder(w).Encode(cards)
}

type data struct {
	Name  string
	Email string
}

func (d data) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Name, validation.Required, validation.Length(2, 45)),
		validation.Field(&d.Email, validation.Required, is.Email),
	)
}
