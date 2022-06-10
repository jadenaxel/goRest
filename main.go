package main

import (
	"fmt"
	"strconv"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/gorilla/mux"
)

type Task struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Content string `json:"content"`
}

type allTasks []Task

var tasks = allTasks{
	{
		ID: 1,
		Name: "Task 1",
		Content: "This is task 1",
	},
}

func indexRoute(res http.ResponseWriter, req *http.Request){
	fmt.Fprintf(res, "Welcome to the Homepage!")
}

func getTasks(res http.ResponseWriter, req *http.Request){
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(tasks)
}

func createTask(res http.ResponseWriter, req *http.Request){
	var newTask Task
	request, err := ioutil.ReadAll(req.Body)

	if err != nil {
		fmt.Fprintf(res, "Error: %v", err)
	}

	json.Unmarshal(request, &newTask)

	newTask.ID = len(tasks) + 1

	tasks = append(tasks, newTask)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(newTask)
}

func getTask(res http.ResponseWriter, req *http.Request){
	params := mux.Vars(req)
	taskID, err := strconv.Atoi(params["id"])

	if err != nil {
		fmt.Fprintf(res, "Error: %v", err)
		return;
	}

	for _, task := range tasks {
		if task.ID == taskID {
			res.Header().Set("Content-Type", "application/json")
			json.NewEncoder(res).Encode(task)
		} else{
			res.WriteHeader(http.StatusNotFound)
		}
	}
}

func updateTask(res http.ResponseWriter, req *http.Request){
	params := mux.Vars(req);
	taskID, errParams := strconv.Atoi(params["id"]);

	if errParams != nil{
		fmt.Fprintf(res, "Error: %v", errParams);
		return;
	}

	data, errData := ioutil.ReadAll(req.Body);

	if errData != nil {
		fmt.Fprintf(res, "Error: %v", errData);
		return;
	}

	for _, task := range tasks {
		if task.ID == taskID {
			res.Header().Set("Content-Type", "application/json");
			fmt.Fprintf(res, "Task updated: %v", string(data));
		} else {
			res.WriteHeader(http.StatusNotFound);
		}
	}
}

func main(){
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", indexRoute).Methods("GET")
	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/tasks", createTask).Methods("POST")
	router.HandleFunc("/tasks/{id}", getTask).Methods("GET")
	router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")

	http.ListenAndServe(":8080", router)
}