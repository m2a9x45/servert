package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"../database"
	"../models"
	jwt "github.com/dgrijalva/jwt-go"
)

func RefreshStaff(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value("user")
	claims := uid.(*models.ClaimsStaff)

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 75*time.Second {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Now, create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if claims.Role == "Staff" {
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    tokenString,
			Expires:  expirationTime,
			HttpOnly: true,
			Path:     "/",
		})
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
	return

	// Set the new token as the users `token` cookie

}

func GetTasks(w http.ResponseWriter, r *http.Request) {

	result, err := database.DBCon.Query("SELECT * from tasks")
	if err != nil {
		println(err)
	}

	var alltasks []models.Task

	for result.Next() {
		var task models.Task
		err := result.Scan(&task.ID, &task.UUID, &task.UserID, &task.QueueID, &task.LinkID, &task.Assigned, &task.Status, &task.CreatedAt)
		if err != nil {
			panic(err.Error())
		}
		alltasks = append(alltasks, task)
	}

	println(len(alltasks))

	if len(alltasks) == 0 {
		w.WriteHeader(http.StatusNotFound)
		res := models.ResObj{Success: false, Message: "No tasks"}
		json.NewEncoder(w).Encode(res)
		return
	}

	fmt.Println(alltasks)
	json.NewEncoder(w).Encode(alltasks)
}
