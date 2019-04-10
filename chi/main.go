package main

import (
	"fmt"
	"sync"
	"time"
	"net/http"
	"database/sql"
	"./db"
	"./ctrl"
	"./srv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

var once sync.Once
var instance *sql.DB
var connString = "nodejs:nodejs@tcp(localhost:3306)/nodejs?charset=utf8&parseTime=true"

func GetConnection() *sql.DB {
	once.Do(func() {
		instance = db.NewMysql(connString)
	})
	return instance
}

func main() {
	r := Routes()
	r.Mount("/v1", v1Router())
	r.Mount("/v2", v2Router())
	err := http.ListenAndServe(":3000", r)
	if err != nil {
		fmt.Println("error:", err)
	}
}

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Timeout(3 * time.Second))
	router.Use(render.SetContentType(render.ContentTypeJSON))
	// router.Use(middleware.DefaultCompress)
	// router.Use(middleware.RedirectSlashes)
	// router.Use(middleware.Recoverer)
	return router
}

func v1Router() chi.Router {
	user := ctrl.NewUserControl(GetConnection(), srv.NewUserService());
	router := Routes()
	router.Route("/", func(r chi.Router) {
		r.Get("/", greeting1)
		r.Post("/user", user.CreateUser)
		r.Get("/user/{idx}", user.GetUser)
	})
	return router
}

func v2Router() chi.Router {
	router := Routes()
	router.Route("/", func(r chi.Router) {
		r.Get("/", greeting2)
	})
	return router
}

func greeting1(w http.ResponseWriter, r *http.Request) {
	val := chi.URLParam(r, "name")
	if val != "" {
		w.Write([]byte(fmt.Sprintf("Hello: %s", val)))
	} else {
		w.Write([]byte(fmt.Sprintf("Hello: None")))
	}
}

func greeting2(w http.ResponseWriter, r *http.Request) {
	val := chi.URLParam(r, "name")
	if val != "" {
		w.Write([]byte(fmt.Sprintf("Hi: %s", val)))
	} else {
		w.Write([]byte(fmt.Sprintf("Hi: None")))
	}
}