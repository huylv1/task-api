package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Task struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
}

var tasks []Task
var currentID int

type App struct {
	Router *mux.Router
}

func (app *App) handleRoutes() {
	app.Router.HandleFunc("/tasks", app.getTasks).Methods("GET")
	app.Router.HandleFunc("/task/{id}", app.readTask).Methods("GET")
	app.Router.HandleFunc("/task", app.createTask).Methods("POST")
	app.Router.HandleFunc("/task/{id}", app.updateTask).Methods("PUT")
	app.Router.HandleFunc("/task/{id}", app.deleteTask).Methods("DELETE")
}

func (app *App) Initialise(initialTasks []Task, id int) {
	tasks = initialTasks
	currentID = id
	app.Router = mux.NewRouter().StrictSlash(true)
	app.handleRoutes()
}
func main() {
	app := App{}
	tasks, id := CreateInitialTasks()
	app.Initialise(tasks, id)
	app.Run("localhost:10000")
}

func (app *App) Run(address string) {
	log.Fatal(http.ListenAndServe(address, app.Router))
}

func sendResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

func sendError(w http.ResponseWriter, statusCode int, err string) {
	message_error := map[string]string{"error": err}
	sendResponse(w, statusCode, message_error)
}

func (app *App) getTasks(writer http.ResponseWriter, request *http.Request) {
	tasks, err := getTasks()
	if err != nil {
		sendError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponse(writer, http.StatusOK, tasks)
}

func (app *App) createTask(writer http.ResponseWriter, r *http.Request) {
	var t Task

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		sendError(writer, http.StatusBadRequest, "Invalid request payload")
		return
	}
	err = t.createTask()
	if err != nil {
		sendError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponse(writer, http.StatusCreated, t)
}

func (app *App) readTask(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(writer, http.StatusBadRequest, "invalid task ID")
		return
	}

	var task Task = Task{ID: id}
	err = task.getTask()

	if err != nil {
		sendError(writer, http.StatusNotFound, "task not found")
	} else {
		sendResponse(writer, http.StatusOK, task)
	}

}

func (app *App) updateTask(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(writer, http.StatusBadRequest, "invalid task ID")
		return
	}

	var newTask Task
	err = json.NewDecoder(request.Body).Decode(&newTask)
	if err != nil {
		sendError(writer, http.StatusBadRequest, "Invalid request payload")
		return
	}

	newTask.ID = id
	err = newTask.updateTask()

	if err != nil {
		sendError(writer, http.StatusInternalServerError, err.Error())
	} else {
		sendResponse(writer, http.StatusOK, newTask)
	}

}

func (app *App) deleteTask(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(writer, http.StatusBadRequest, "invalid task ID")
		return
	}

	var task Task = Task{ID: id}
	err = task.deleteTask()

	if err != nil {
		sendError(writer, http.StatusNotFound, "task not found")
	} else {
		sendResponse(writer, http.StatusOK, map[string]string{"result": "successful deletion"})
	}

}
