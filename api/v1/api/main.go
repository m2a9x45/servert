package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/handlers"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type product struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Des      string  `json:"des"`
	Price    float64 `json:"price"`
	Instock  bool    `json:"instock"`
	Setupfee float64 `json:"setupfee"`
	Discount float64 `json:"discount"`
}

type details struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type resObj struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

var db *sql.DB
var err error

func main() {

	db, err = sql.Open("mysql", "root:99dZ%dtw&gE@tcp(127.0.0.1:4000)/servert")
	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()

	r := mux.NewRouter()

	header := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origin := handlers.AllowedOrigins([]string{"*"})

	r.HandleFunc("/products", GetProducts).Methods("GET")
	r.HandleFunc("/intrest", RegIntrest).Methods("POST", "OPTIONS")

	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(header, methods, origin)(r)))
}

// RegIntrest will add name and email to DB
func RegIntrest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		return
	}

	destails := details{}

	jsn, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Error wilst reading r body", err)
	}

	err = json.Unmarshal(jsn, &destails)
	if err != nil {
		log.Fatal("Error wilst unmarshaling json", err)
	}

	log.Println(destails.Email, destails.Name)

	result, err := db.Query("INSERT INTO reg (name, email) VALUES (?,?)", destails.Name, destails.Email)
	if err != nil {
		log.Fatal("Error wilst inserting into DB", err)
	}

	defer result.Close()

	log.Println(result)

	res := resObj{true, "Details inserted into DB"}

	json.NewEncoder(w).Encode(res)
}

// GetProducts will return a list of products
func GetProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	type Category product

	var allproducts []product

	result, err := db.Query("SELECT * from products")
	if err != nil {
		println(err)
	}

	defer result.Close()

	for result.Next() {
		var product product
		err := result.Scan(&product.ID, &product.Name, &product.Des, &product.Price, &product.Instock, &product.Setupfee, &product.Discount)
		if err != nil {
			panic(err.Error())
		}
		allproducts = append(allproducts, product)
	}

	fmt.Println(allproducts)
	json.NewEncoder(w).Encode(allproducts)
}
