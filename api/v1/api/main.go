package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"./database"
	"./routes"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

var jwtKey []byte

func main() {

	jwtSecret, exists := os.LookupEnv("JWT_SECRET")

	if exists {
		println(jwtSecret)
		jwtKey = []byte(jwtSecret)
	}

	var err error
	database.DBCon, err = sql.Open("mysql", "root:99dZ%dtw&gE@tcp(127.0.0.1:4000)/servert")
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	header := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "Access-Control-Allow-Credentials", "Access-Control-Allow-Origin", "Access-Control-Request-Headers"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	origin := handlers.AllowedOrigins([]string{"http://127.0.0.1:8080"})
	creds := handlers.AllowCredentials()

	r.HandleFunc("/intrest", routes.RegIntrest).Methods("POST", "OPTIONS")
	r.HandleFunc("/signup", routes.Signup).Methods("POST", "OPTIONS")
	r.HandleFunc("/signin", routes.Signin).Methods("POST", "OPTIONS")
	r.HandleFunc("/account", routes.Account).Methods("GET")
	r.HandleFunc("/acountInfo", routes.AccountInfo).Methods("GET")
	r.HandleFunc("/loggedIn", routes.IsLoggedIn).Methods("GET")
	r.HandleFunc("/refresh", routes.Refresh).Methods("GET")

	r.HandleFunc("/products/{prodID}", routes.GetProducts).Methods("GET")
	r.HandleFunc("/products", routes.GetProducts).Methods("GET")

	r.HandleFunc("/create-payment-intent/{prodID}", routes.CreatePaymentIntent).Methods("GET")

	r.HandleFunc("/order", routes.Order).Methods("POST", "OPTIONS")
	r.HandleFunc("/order", routes.GetOrders).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(header, methods, origin, creds)(r)))

}
