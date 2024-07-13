package main

import (
	"fmt"
	"os"
	"pharmacy_management/internal/authentication"
	"pharmacy_management/internal/db"
	"pharmacy_management/internal/logformated"
	"pharmacy_management/internal/user"
)

var log = logformated.GetLogger(logformated.ComponentMain)

func CreateAndInsertDefaultValuesToDb() {
	log.Info("Creating and inserting values to DB only if not done already")
	db.RegisterDriver()
	db.CreateTables()
	db.InsertDefaultValuesIntoRole()
}

func main() {
	CreateAndInsertDefaultValuesToDb()
	user.LoadDefaultRole()
	for {
		OpenMainMenu()
	}
}

func OpenMainMenu() {
	fmt.Print("\n1. Login \n" +
		"2. Register \n" +
		"3. Exit \n" +
		"Enter your choice : ")
	var choice int
	_, err := fmt.Scanln(&choice)
	if err != nil {
		log.Info("Error while reading input from user : ", err.Error())
		return
	}
	switch choice {
	case 1:
		err, loggedInUser := authentication.Login()
		if err != nil {
			log.Error("Error while logging in : ", err.Error(), " Please try again..")
			break
		}
		switch loggedInUser.Role.RoleId {
		case 1:
			var customer user.Customer
			customer.OpenUserMenu(loggedInUser)
		case 2:
			var pharmacist user.Pharmacist
			pharmacist.OpenUserMenu(loggedInUser)
		case 3:
			var admin user.Admin
			admin.OpenUserMenu(loggedInUser)
		}
	case 2:
		err := authentication.Register()
		if err != nil {
			log.Error("Error while registering user : ", err.Error())
			break
		}
		log.Info("User registered successfully. Please login to continue. ")
	case 3:
		log.Info("Exiting application..")
		os.Exit(0)
	default:
		log.Info("Invalid choice. Please try again..")
	}
}
