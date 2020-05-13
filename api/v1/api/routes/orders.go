package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"../database"
	"../models"
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

func generateTaskID() string {

	id := ksuid.New()
	taskID := "task_" + id.String()

	println(taskID)

	result, err := database.DBCon.Query("SELECT task_id from tasks WHERE task_id=(?)", taskID)
	if err != nil {
		println(err)
	}

	defer result.Close()

	for result.Next() {
		var taskid string
		err := result.Scan(&taskid)
		if err != nil {
			panic(err)
		}
		if taskid != "" {
			// make new user id
			println(taskid, "already exists")
			generateTaskID()
		}
	}

	return taskID
}

func createTask(user string, taskType string) {
	fmt.Println(user)
	taskID := generateTaskID()
	status := "open"

	now := time.Now()
	secs := now.Unix()

	result, err := database.DBCon.Query("INSERT INTO tasks (task_id, user_id, status, link_id, created_at) VALUES (?,?,?,?,?)", taskID, user, status, taskType, secs)
	if err != nil {
		log.Fatal("Error wilst inserting into DB", err)
	}

	defer result.Close()

	fmt.Println("Inserted Into DB")
}

func createServer(userID string, orderID string) {

	// this is where we make calls to promox to setup the server

	result, err := database.DBCon.Query("INSERT INTO servers (server_id, virtiual_id, user_id, active_order) VALUES (?,?,?,?)", "server_1", "placeholder", userID, orderID)
	if err != nil {
		log.Fatal("Error wilst inserting into DB", err)
	}

	defer result.Close()

	fmt.Println("Inserted Into DB")
}

func CreatePaymentIntent(w http.ResponseWriter, r *http.Request) {

	uid := r.Context().Value("user")
	claims := uid.(*models.Claims)

	vars := mux.Vars(r)
	id := vars["prodID"]
	dur := vars["dur"]
	cardID := vars["cardID"]

	durI, err := strconv.Atoi(dur)
	if err != nil {
		panic(err)
	}

	println("Hi", dur)

	if durI != 1 && durI != 3 && durI != 6 && durI != 12 {
		res := models.ResObj{Success: false, Message: "You cannot rent the sever for that period of time"}
		json.NewEncoder(w).Encode(res)
		return
	}

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

		// get customerID from userID

		var customerID string

		result, err = database.DBCon.Query("SELECT stripe_id FROM users WHERE user_id=(?)", claims.UserID)
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

		// stripe.Key = "sk_test_OGXIlmLXL1Gvhpa9jqBdxutN00YB96uOjP"

		stripeKey, exists := os.LookupEnv("STRIPE_KEY")

		if exists {
			stripe.Key = stripeKey
		}

		price := allproducts[0].Price
		pricePennies := int64(price * 100)
		chargePrice := pricePennies * int64(durI)
		println(price, " in £")
		println(pricePennies, " in p")
		println(chargePrice, "for", durI, "months")

		var params *stripe.PaymentIntentParams

		if cardID != "" {
			// using saved card
			params = &stripe.PaymentIntentParams{
				Amount:        stripe.Int64(chargePrice),
				Currency:      stripe.String(string(stripe.CurrencyGBP)),
				Customer:      stripe.String(customerID),
				PaymentMethod: stripe.String(cardID),
			}
		} else {
			// using a new card
			params = &stripe.PaymentIntentParams{
				Amount:   stripe.Int64(chargePrice),
				Currency: stripe.String(string(stripe.CurrencyGBP)),
				Customer: stripe.String(customerID),
			}
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

func MakeOrder(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodOptions {
		return
	}

	uid := r.Context().Value("user")
	claims := uid.(*models.Claims)

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
	now := time.Now()
	secs := now.Unix()

	println(OrderID, "seletced")
	// converts duration which is a string to an int
	dur, err := strconv.Atoi(order.Dur)
	if err != nil {
		panic(err)
	}
	// checks that int is 1,3,6,12 for the months we offer for products
	if dur != 1 && dur != 3 && dur != 6 && dur != 12 {
		res := models.ResObj{Success: false, Message: "You cannot rent the sever for that period of time"}
		json.NewEncoder(w).Encode(res)
		return
	}

	factor := dur * 2629743 // 2629743 number of seconds in a month
	expires := secs + int64(factor)

	result, err := database.DBCon.Query("INSERT INTO orders (order_id, user_id, payment_id, prod_id, created_at, duration, expires_at ) VALUES (?,?,?,?,?,?,?)", OrderID, claims.UserID, order.PaymentID, order.ProductID, secs, order.Dur, expires)
	if err != nil {
		log.Fatal("Error wilst inserting into DB", err)
	}

	defer result.Close()

	fmt.Println("Inserted Into DB")

	createServer(claims.UserID, OrderID)
	createTask(claims.UserID, OrderID)

	res := models.ResObj{Success: true, Message: OrderID}

	json.NewEncoder(w).Encode(res)

}

// GetOrders will return a list of orders for an account
func GetOrders(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value("user")
	claims := uid.(*models.Claims)

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
