package models

import (
	"errors"
	"example.com/rest-api/db"
	"example.com/rest-api/utils"
)

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u User) Save() error {
	query := `INSERT INTO users (email, password) VALUES (?, ?)`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()
	// encrypt password
	hashedPassword, err := utils.HashPassword(u.Password)
	result, err := stmt.Exec(u.Email, hashedPassword)

	if err != nil {
		return err
	}

	userId, err := result.LastInsertId()
	u.ID = userId
	return err
}

func (u *User) ValidateCredentials() error {
	query := `SELECT id, password FROM users WHERE email = ?`
	row := db.DB.QueryRow(query, u.Email)

	var retrievedPassword string
	err := row.Scan(&u.ID, &retrievedPassword)
	if err != nil {
		return errors.New("System is error")
	}
	isCorrect := utils.IsPasswordCorrect(u.Password, retrievedPassword)

	if !isCorrect {
		return errors.New("Invalid username or password")
	}
	return nil
}

func GetAllUsers() ([]User, error) {
	query := `SELECT * FROM users`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	var users []User
	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
