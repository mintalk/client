package network

import (
	"mintalk/server/db"

	"golang.org/x/crypto/bcrypt"
)

func ValidateAuthRequest(database *db.Connection, data NetworkData) (bool, error) {
	var user *db.User
	err := database.Where(&db.User{Name: data["username"].(string)}).First(&user).Error
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, nil
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"].(string)))
	return err == nil, nil
}
