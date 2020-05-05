package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"../database"
	"../models"
)

func GetTasks(w http.ResponseWriter, r *http.Request) {

	// Add auth for internal product

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
