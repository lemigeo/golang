package ctrl

import (
	"fmt"
	"strconv"
	"net/http"
	"encoding/json"
	"database/sql"

	"github.com/go-chi/chi"

	"../srv"
	"../data"
)

type UserControl struct {
	service srv.UserService
	db *sql.DB
}

func NewUserControl	(db *sql.DB, service srv.UserService) *UserControl {
	return &UserControl{service, db}
}

func (u *UserControl) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user data.User
	json.NewDecoder(r.Body).Decode(&user)
	if(user.UserId == "") {
		respondWithError(w, -100, "User id not found")
		return
	}
	info, err := u.service.CreateUser(u.db, user)
	if err != nil {
		respondWithError(w, -100, err.Error())
		return
	}
	respondwithJSON(w, info);
}

func (u *UserControl) GetUser(w http.ResponseWriter, r *http.Request) {
	idx := chi.URLParam(r, "idx")
	if idx != "" {
		userIdx, err := strconv.ParseInt(idx, 10, 64)
		if err != nil {
			respondWithError(w, -100, "Invalid user index")
			return
		}
		info, err := u.service.GetUser(u.db, userIdx)
		if err != nil {
			respondWithError(w, -100, err.Error())
			return
		}
		fmt.Println(info)
		respondwithJSON(w, info)
	} else {
		fmt.Println("idx not found")
		respondWithError(w, -100, "User index not found")
	}
}
