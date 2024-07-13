package user

import (
	"fmt"
	"pharmacy_management/internal/medicine"
)

type Pharmacist struct {
	user User
}

func (pharmacist *Pharmacist) OpenUserMenu(user User) {
	pharmacist.user = user
	log.Info("Welcome to the Pharmacy Management System! You are logged in as ", user.Name)
	for {
		fmt.Print("\n1. View all medicine available \n" +
			"2. Add/Update medicine \n" +
			"3. Logout \n" +
			"Enter your choice : ")
		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			log.Info("Viewing all medicine available...")
			pharmacist.ViewAllMedicine()
		case 2:
			log.Info("Adding/Updating medicine...")
			pharmacist.AddOrUpdateMedicine()
		case 3:
			pharmacist.Logout()
			log.Info("Logged out successfully.")
			return
		default:
			log.Info("Invalid option. Please try again.")
		}
	}
}

func (pharmacist *Pharmacist) ViewAllMedicine() {
	medicineList, err := medicine.GetAllMedicine()
	if err != nil {
		log.Error("Error while viewing all medicine : ", err.Error())
		return
	}
	for _, medicine := range medicineList {
		fmt.Println(medicine)
	}
}

func (pharmacist *Pharmacist) AddOrUpdateMedicine() {
	err := medicine.AddOrUpdateMedicine()
	if err != nil {
		log.Error("Error while adding/updating medicine : ", err.Error())
	}
}

func (pharmacist *Pharmacist) Logout() {
	log.Info("Logging out...")
}
