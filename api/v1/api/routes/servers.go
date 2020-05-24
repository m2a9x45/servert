package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"../database"
	"../models"
)

func GetServer(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value("user")
	claims := uid.(*models.Claims)

	var allservers []models.Server

	fmt.Println(claims.UserID)

	result, err := database.DBCon.Query("SELECT server_id, virtiual_id, user_id, active_order from servers WHERE user_id=(?)", claims.UserID)
	if err != nil {
		println(err.Error())
		res := models.ResObj{Success: false, Message: "Couldn't find servers"}
		json.NewEncoder(w).Encode(res)
		return
	}

	defer result.Close()

	for result.Next() {
		var server models.Server
		err := result.Scan(&server.UUID, &server.VirtiualID, &server.UserID, &server.ActiveOrder)
		if err != nil {
			panic(err.Error())
		}
		allservers = append(allservers, server)
	}

	if len(allservers) == 0 {
		w.WriteHeader(http.StatusNotFound)
		res := models.ResObj{Success: false, Message: "product not found"}
		json.NewEncoder(w).Encode(res)
		return
	}

	json.NewEncoder(w).Encode(allservers)

}
