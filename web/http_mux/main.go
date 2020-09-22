package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	fmt.Println("\n--- Starting server ---------------------")
	port := os.Getenv("API_PORT")
	db = configureDB()

	fmt.Println("Starting Restful services...")
	r := mux.NewRouter()
	r.HandleFunc("/tasks", handleGetCreateTasks).
		Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/tasks/{id}", handleGetEditDeleteTasks).
		Methods(http.MethodGet, http.MethodPut, http.MethodDelete)
	r.HandleFunc("/tasks/{id}/toggle", handleTaskToggle).
		Methods(http.MethodPut)

	http.Handle("/", r)

	err := http.ListenAndServe(":"+port, nil)
	log.Fatal(err)
	fmt.Println("Listening on port", port)
}

func respond(w http.ResponseWriter, status int, response string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

func handleGetCreateTasks(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%v - %v - %v\n", r.Method, r.Host, r.RequestURI)

	switch r.Method {
	case http.MethodGet: /* get all tasks */
		tc, _ := dbGetAllTasks()
		jsonTasks, _ := toJSON(tc)
		respond(w, http.StatusOK, jsonTasks)

	case http.MethodPost: /* create task */
		var t Task
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			respond(w, http.StatusBadRequest, makeMessage("Task has invalid fields"))
			return
		}
		valErrs, jsonString := validateTask(t)
		if valErrs != nil {
			respond(w, http.StatusUnprocessableEntity, fmt.Sprintf(`{"message":%v}`, jsonString))
			return
		}

		id, err := dbCreateTask(t)
		if err != nil {
			respond(w, http.StatusBadRequest, makeMessage("Couldn't create new task"))
			return
		}
		respond(w, http.StatusCreated, makeMessage("Created new task with id %v", id))
	}
}

func handleGetEditDeleteTasks(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%v - %v - %v\n", r.Method, r.Host, r.RequestURI)

	strID, hasID := mux.Vars(r)["id"]
	id, err := strconv.Atoi(strID)

	if !hasID || (hasID && err != nil /* id isn't an int */) {
		respond(w, http.StatusUnprocessableEntity,
			makeMessage("Invalid task id, please specify an existing integer id"))
		return
	}

	switch r.Method {
	case http.MethodGet: /* get task by id */
		t, err := dbGetTask(id)
		if err != nil {
			respond(w, http.StatusNotFound,
				makeMessage("No task with the id %v was found", id))
			return
		}

		jsonTask, _ := toJSON(t)
		respond(w, http.StatusOK, jsonTask)

	case http.MethodPut: /* edit existing task */
		var t Task
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			respond(w, http.StatusBadRequest, makeMessage("Task has invalid fields"))
			return
		}

		valErrs, jsonString := validateTask(t)
		if valErrs != nil {
			respond(w, http.StatusUnprocessableEntity,
				fmt.Sprintf(`{"message":%v}`, jsonString))
			return
		}

		updatedTask, err := dbUpdateTask(t)
		if err != nil {
			respond(w, http.StatusBadRequest, makeMessage("Couldn't update task"))
			return
		}

		response, _ := toJSON(updatedTask)
		respond(w, http.StatusOK, response)

	case http.MethodDelete: /* delete task */
		if err := dbDeleteTask(id); err != nil {
			respond(w, http.StatusNotFound,
				makeMessage("Couldn't delete task with id %v", id))
			return
		}
		respond(w, http.StatusNoContent, "")
	}
}

func handleTaskToggle(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%v - %v - %v\n", r.Method, r.Host, r.RequestURI)

	if r.Method == http.MethodPut {
		strID, hasID := mux.Vars(r)["id"]
		id, err := strconv.Atoi(strID)

		if !hasID || (hasID && err != nil /* id isn't an int */) {
			respond(w, http.StatusUnprocessableEntity,
				makeMessage("Invalid task id, please specify an existing integer id"))
			return
		}

		t, err := dbGetTask(id)
		if err != nil {
			respond(w, http.StatusNotFound,
				makeMessage("No task with the id %v was found", id))
			return
		}

		t.Completed = !t.Completed
		updatedTask, err := dbUpdateTask(t)
		if err != nil {
			respond(w, http.StatusBadRequest,
				makeMessage("Couldn't toggle task completed"))
			return
		}
		response, _ := toJSON(updatedTask)
		respond(w, http.StatusOK, response)
	}
}

func toJSON(i interface{}) (string, error) {
	tj, err := json.Marshal(i)
	if err != nil {
		return "", err
	}
	return string(tj), nil
}

func makeMessage(template string, args ...interface{}) string {
	str := fmt.Sprintf(`{"message": "%v"}`, template)
	if len(args) > 0 {
		return fmt.Sprintf(str, args)
	}
	return str
}
