package repository

import (
	"errors"
	"fmt"
	"nrs16/cme/repository/entities"

	"github.com/gocql/gocql"
)

type Authentication interface {
	Register() error
	GetUserbyUsername(username string) (entities.User, error)
}

type AuthenticationDatabase struct {
	Database *gocql.Session
}

func (r *AuthenticationDatabase) Register(user entities.User) error {
	if user.PasswordHash == nil || user.PasswordSalt == nil {
		return errors.New("missing_password_info")
	}
	cql := `
		insert into user (username, first_name, last_name, password_salt, password_hash, email, ts_created)
		values (
			?,?,?,?,?,?, toTimestamp(now())
		)
	`
	err := r.Database.Query(cql, user.Username, user.FirstName, user.LastName, *user.PasswordSalt, *user.PasswordHash, user.Email).Exec()
	if err != nil {
		fmt.Printf("error inserting: %s", err.Error())
		return err
	}
	return nil
}

func (r *AuthenticationDatabase) GetUserbyUsername(username string) (entities.User, error) {
	var user entities.User
	user.Username = username
	err := r.Database.Query(`select first_name, last_name, email, password_hash, password_salt from user where username = ?`, username).Scan(&user.FirstName, &user.LastName, &user.Email, &user.PasswordHash, &user.PasswordSalt)

	return user, err
}
