package input

import (
	"fmt"
	"mintalk/server/db"
)

func (console *Console) op(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("op requires 1 argument")
	}
	user := db.User{Name: args[0]}
	err := console.database.Where(user).First(&user).Error
	if err != nil {
		return fmt.Errorf("failed to find user: %v", err)
	}
	user.Operator = true
	err = console.database.Save(&user).Error
	if err != nil {
		err = fmt.Errorf("failed to save user: %v", err)
	}
	return err
}

func (console *Console) deop(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("deop requires 1 argument")
	}
	user := db.User{Name: args[0]}
	err := console.database.Where(user).First(&user).Error
	if err != nil {
		return fmt.Errorf("failed to find user: %v", err)
	}
	user.Operator = false
	err = console.database.Save(&user).Error
	if err != nil {
		err = fmt.Errorf("failed to save user: %v", err)
	}
	return err
}
