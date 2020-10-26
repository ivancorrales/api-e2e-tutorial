package internal

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

const (
	paramID = "id"
	qStatus = "status"
)

type TodoHandler struct {
	database *MongoHandler
}

func NewTodoHandler(db *MongoHandler) *TodoHandler{
	return &TodoHandler{db}
}

func (h *TodoHandler) AddTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	todo.Status = "PENDING"
	res, err := h.database.AddOne(&todo)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	todo.ID = res.InsertedID.(string)
	respondWithJSON(w, http.StatusCreated, todo)
}

func (h *TodoHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, paramID)
	if id == "" {
		respondWithError(w, http.StatusNotFound, "id is required")
		return
	}
	todo := &Todo{}
	if err := h.database.GetOne(todo, bson.M{paramID: id}); err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("todo with id: %s not found", id))
		return
	}
	respondWithJSON(w, http.StatusOK, todo)
}

func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	existingTodo := &Todo{}
	id := chi.URLParam(r, paramID)
	if id == "" {
		respondWithError(w, http.StatusNotFound, "id is required")
		return
	}
	err := h.database.GetOne(existingTodo, bson.M{paramID: id})
	if err != nil {
		http.Error(w, fmt.Sprintf("Todo with id: %s does not exist", id), 400)
		return
	}
	if _, err = h.database.RemoveOne(bson.M{paramID: id}); err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("todo with id: %s not found", id))
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)
}

func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, paramID)
	if id == "" {
		respondWithError(w, http.StatusNotFound, "id is required")
		return
	}
	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if _, err := h.database.Update(bson.M{"_id": id},todo);err!=nil{
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, todo)
}

func (h *TodoHandler) ListTodosByStatus(w http.ResponseWriter, r *http.Request) {
	status:=r.URL.Query().Get(qStatus)
	todos := h.database.Get(bson.M{"status": status})
	respondWithJSON(w, http.StatusOK, todos)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"message": msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	if payload != nil {
		response, _ := json.Marshal(payload)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
	w.WriteHeader(code)
}
