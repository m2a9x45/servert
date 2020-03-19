package routes

import (
	"encoding/json"
	"net/http"

	"../database"
	"../models"
	"github.com/gorilla/mux"
)

func GetProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	id := vars["prodID"]

	if id != "" {
		var allproducts []models.Product

		result, err := database.DBCon.Query("SELECT * from products WHERE prod_id=(?)", id)
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
			return
		}

		// fmt.Println(allproducts)
		json.NewEncoder(w).Encode(allproducts)
		return
	}

	// type Category product

	var allproducts []models.Product

	result, err := database.DBCon.Query("SELECT * from products")
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

	// fmt.Println(allproducts)
	json.NewEncoder(w).Encode(allproducts)
}
