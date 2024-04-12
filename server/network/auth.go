package network

import (
	"mintalk/server/db"

	"golang.org/x/crypto/bcrypt"
)

func ValidateAuthRequest(database *db.Connection, data NetworkData) (bool, error) {
	var user db.User
	user.Name = data["username"].(string)
	err := database.Where(&user).First(&user).Error
	if err != nil {
		return false, err
	}
	if user.Password == "" {
		CreatePassword(&user, data["password"].(string))
		err = database.Save(user).Error
		if err != nil {
			return false, err
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"].(string)))
	return err == nil, nil
}

func CreatePassword(user *db.User, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	return nil
}
