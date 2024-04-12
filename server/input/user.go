package input

import (
	"fmt"
	"mintalk/server/db"
)

func (console *Console) user(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("user requires an argument")
	}
	switch args[0] {
	case "add":
		return console.useradd(args[1:])
	case "del":
		return console.userdel(args[1:])
	case "list":
		return console.userlist(args[1:])
	}
	return fmt.Errorf("user subcommand not found: %s", args[0])
}

func (console *Console) useradd(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("user add requires 1 argument")
	}
	user := db.User{Name: args[0]}
	err := console.database.Create(&user).Error
	if err != nil {
		err = fmt.Errorf("failed to create user: %v", err)
	}
	return err
}

func (console *Console) userdel(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("user del requires 1 argument")
	}
	user := db.User{Name: args[0]}
	err := console.database.Where(user).First(&user).Error
	if err != nil {
		return fmt.Errorf("failed to find user: %v", err)
	}
	err = console.database.Delete(&user).Error
	if err != nil {
		err = fmt.Errorf("failed to delete user: %v", err)
	}
	return err
}

func (console *Console) userlist(args []string) error {
	var users []db.User
	err := console.database.Find(&users).Error
	if err != nil {
		return fmt.Errorf("failed to list users: %v", err)
	}
	fmt.Printf("Name\tRole\n")
	for _, user := range users {
		role := "member"
		if user.Operator {
			role = "operator"
		}
		fmt.Printf("%v\t%v\n", user.Name, role)
	}
	return nil
}
