package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"../database"
	"../models"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func IsLoggedInStaff(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the JWT string from the cookie
	tknStr := c.Value

	// Initialize a new instance of `Claims`
	claims := &models.ClaimsStaff{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
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

	fmt.Println(claims.Role)

	if claims.Role == "Staff" {
		res := models.ResObj{Success: true, Message: claims.StaffID}
		json.NewEncoder(w).Encode(res)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Finally, return the welcome message to the user, along with their
	// username given in the token

}

func SigninStaff(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		return
	}

	signin := models.SignUpObj{}

	jsn, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Error wilst reading r body", err)
	}

	err = json.Unmarshal(jsn, &signin)
	if err != nil {
		log.Fatal("Error wilst unmarshaling json", err)
	}

	log.Println(signin, "here")

	// basic logging cred checking
	if signin.Email == "" || signin.Password == "" {
		log.Println("No email")
		res := models.ResObj{Success: false, Message: "Problem signing in"}
		json.NewEncoder(w).Encode(res)
		return
	}

	// check password.

	result, err := database.DBCon.Query("SELECT staff_id, password from staff WHERE email=(?)", signin.Email)
	if err != nil {
		println(err)
	}

	log.Println("DB check complete")

	defer result.Close()

	for result.Next() {
		var hashedPassword string
		var userID string
		err := result.Scan(&userID, &hashedPassword)
		if err != nil {
			panic(err)
		}
		// Comparing the password with the hash
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(signin.Password))
		if err == nil {
			expirationTime := time.Now().Add(5 * time.Minute)

			claims := &models.ClaimsStaff{
				StaffID: userID,
				Role:    "Staff",
				StandardClaims: jwt.StandardClaims{
					// In JWT, the expiry time is expressed as unix milliseconds
					ExpiresAt: expirationTime.Unix(),
				},
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			// Create the JWT string
			tokenString, err := token.SignedString(jwtKey)
			if err != nil {
				// If there is an error in creating the JWT return an internal server error
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "token",
				Value:    tokenString,
				Expires:  expirationTime,
				HttpOnly: true,
				Path:     "/",
			})

			res := models.ResObj{Success: true, Message: "Signed In"}
			json.NewEncoder(w).Encode(res)
		}
		fmt.Println(err) // nil means it is a match
		return
	}

	res := models.ResObj{Success: false, Message: "Something went wrong"}
	json.NewEncoder(w).Encode(res)

}

func RefreshStaff(w http.ResponseWriter, r *http.Request) {
	// (BEGIN) The code uptil this point is the same as the first part of the `Welcome` route
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
	claims := &models.ClaimsStaff{}
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
	// (END) The code up-till this point is the same as the first part of the `Welcome` route

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
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

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
