# Servert API V1

## Routes

```go
	// auth.go
	r.HandleFunc("/isloggedin", routes.IsLoggedIn).Methods("GET")
	r.HandleFunc("/signup", routes.Signup).Methods("POST", "OPTIONS")
	r.HandleFunc("/signin", routes.Signin).Methods("POST", "OPTIONS")
	r.HandleFunc("/refresh", routes.Refresh).Methods("GET")

	//account.go
	r.HandleFunc("/intrest", routes.Intrest).Methods("POST", "OPTIONS")
	r.HandleFunc("/account", routes.Account).Methods("GET")
	r.HandleFunc("/acountinfo", routes.AccountInfo).Methods("GET")

	//products.go
	r.HandleFunc("/products/{prodID}", routes.GetProducts).Methods("GET")
	r.HandleFunc("/products", routes.GetProducts).Methods("GET")

	//orders.go
	r.HandleFunc("/create-payment-intent/{prodID}", routes.CreatePaymentIntent).Methods("GET")
	r.HandleFunc("/makeorder", routes.MakeOrder).Methods("POST", "OPTIONS")
	r.HandleFunc("/getorders", routes.GetOrders).Methods("GET")
```


