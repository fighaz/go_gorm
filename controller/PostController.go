package controller

import (
	"blog/config"
	"blog/model"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Response struct {
	Message string
	Data    interface{}
}

func HandlerIndex(w http.ResponseWriter, r *http.Request) {
	var response Response
	var post []model.Post
	err := config.DB.Find(&post).Error
	if err != nil {
		response.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
	}
	response.Message = "Get Data Success"
	response.Data = post
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(post)
}
func HandlerDetail(w http.ResponseWriter, r *http.Request) {
	var response Response
	vars := mux.Vars(r)
	id := vars["id"]
	var post model.Post
	err := config.DB.First(&post, &id).Error
	if err != nil {
		response.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
	}
	response.Message = "Get Data Success"
	response.Data = post
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}
func HandlerInsert(w http.ResponseWriter, r *http.Request) {
	var post model.Post
	var response Response
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		response.Message = err.Error()
	}
	err = config.DB.Create(&post).Error
	if err != nil {
		response.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
	}
	fmt.Println("Insert Succes")

	response.Message = "Insert Succes"
	response.Data = post
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}

func HandlerUpdate(w http.ResponseWriter, r *http.Request) {
	var post model.Post
	var response Response
	vars := mux.Vars(r)
	id := vars["id"]
	err := config.DB.First(&post, id).Error

	if err != nil {
		response.Message = err.Error()

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
	}

	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		response.Message = err.Error()
	}
	err = config.DB.Save(&post).Error
	if err != nil {
		response.Message = err.Error()

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
	}

	fmt.Println("Update Succes")

	response.Message = "Update Succes"
	response.Data = post

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}
func HandlerDelete(w http.ResponseWriter, r *http.Request) {
	var response Response
	vars := mux.Vars(r)
	id := vars["id"]
	var post model.Post
	err := config.DB.Delete(&post, id).Error
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
	}
	fmt.Println("Delete Succes")

	response.Message = "Delete Succes"
	response.Data = post
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}
