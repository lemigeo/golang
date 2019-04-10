package srv

import (
	"fmt"
	"time"
	"database/sql"
	"../data"
)

type UserService struct {
}

func NewUserService() UserService {
	service := UserService{}
	return service
}

func (s *UserService) GetAllUser(db *sql.DB, idx int64) ([]data.User, error) {
	var list []data.User
	rs, err := db.Query("SELECT * FROM user")
	if err != nil {
		fmt.Println(err.Error())
		return list, err
	}
	for rs.Next() {
		var user data.User
		err := rs.Scan(&user.IDX, &user.UserId, &user.CreateDt, &user.UpdateDt)
		if err != nil {
			fmt.Println(err.Error())
			return list, err
		}
		list = append(list, user)
	}
	return list, err
}

func (s *UserService) GetUser(db *sql.DB, idx int64) (data.User, error) {
	var user data.User
	err := db.QueryRow("SELECT * FROM user WHERE idx=?", idx).Scan(&user.IDX, &user.UserId, &user.CreateDt, &user.UpdateDt)
	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}
	return user, err
}

func (s *UserService) CreateUser(db *sql.DB, user data.User) (data.User, error) {
	now := time.Now()
	rs, err := db.Exec("INSERT INTO user(user_id, create_dt, update_dt) VALUES(?, ?, ?)", user.UserId, now, now)
	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}
	user.IDX, err = rs.LastInsertId()
	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}
	return user, err
}