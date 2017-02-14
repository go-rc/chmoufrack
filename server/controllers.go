package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// POSTExercise ...
func POSTExercise(writer http.ResponseWriter, reader *http.Request) {
	var exercise Exercise
	var err error

	if reader.Body == nil {
		http.Error(writer, "Please send a request body", http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(reader.Body).Decode(&exercise)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = addExercise(exercise)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	writer.WriteHeader(http.StatusCreated)
	return
}

func GETExercise(writer http.ResponseWriter, reader *http.Request) {
	var exercise Exercise
	var err error

	vars := mux.Vars(reader)
	exerciseID := vars["id"]

	i, err := strconv.ParseInt(exerciseID, 10, 64)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}

	exercise, err = getExercise(i)
	if err != nil {
		if _, ok := err.(*error404); ok {
			http.Error(writer, err.Error(), http.StatusNotFound)
		} else {
			http.Error(writer, err.Error(), http.StatusBadRequest)
		}
		return
	}

	if err = json.NewEncoder(writer).Encode(exercise); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
}
