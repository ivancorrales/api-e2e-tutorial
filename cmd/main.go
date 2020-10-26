package main

import (
	"github.com/go-chi/chi"
	"github.com/wesovilabs/api-e2e-tutorial/internal"
	"log"
	"net/http"
	"os"
)

func registerRoutes(mh *internal.MongoHandler) http.Handler {
	h := internal.NewTodoHandler(mh)
	r := chi.NewRouter()
	r.Route("/todos", func(r chi.Router) {
		r.Get("/", h.ListTodosByStatus) //GET /todos
		r.Get("/{id}", h.GetTodo)       //GET /todos/2
		r.Post("/", h.AddTodo)          //POST /todos
		r.Put("/{id}", h.UpdateTodo)    //PUT /todos/2
		r.Delete("/{id}", h.DeleteTodo) //DELETE /todos/2
	})
	return r
}

func buildMongoHandler(uri string, dbName string) *internal.MongoHandler {
	mongoDbConnection := "mongodb://localhost:27017"
	return internal.NewMongoHandler(mongoDbConnection, "")
}

func main() {
	args := os.Args[1:]
	mongoHandler := buildMongoHandler(args[0], args[1])
	r := registerRoutes(mongoHandler)
	log.Fatal(http.ListenAndServe(":9000", r))
}
