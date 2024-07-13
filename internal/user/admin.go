package user

import (
	"fmt"
	"pharmacy_management/internal/db"
	"pharmacy_management/internal/medicine"
)

type Admin struct {
	user User
}

func (admin *Admin) OpenUserMenu(user User) {
	admin.user = user
	log.Info("Welcome to the Pharmacy Management System! You are logged in as ", user.Name)
	for {
		fmt.Print("1. View all customers \n" +
			"2. View all pharmacists \n" +
			"3. View all medicine \n" +
			"4. Delete pharmacist \n" +
			"5. Logout \n" +
			"Enter your choice : ")
		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			log.Info("Viewing all customers...")
			admin.ViewAllCustomers()
		case 2:
			log.Info("Viewing all pharmacists...")
			admin.ViewAllPharmacists()
		case 3:
			log.Info("Viewing all medicine...")
			admin.ViewAllMedicine()
		case 4:
			log.Info("Deleting pharmacist...")
			admin.DeletePharmacist()
		case 5:
			admin.Logout()
			log.Info("Logged out successfully.")
			return
		default:
			log.Info("Invalid option. Please try again.")
		}
	}
}

func (admin *Admin) GenerateReport() {
	panic("Not implemented yet")
}

func (admin Admin) ViewAllCustomers() {
	customerList, err := db.GetAllCustomers()
	if err != nil {
		log.Error("Error while viewing all customers : ", err.Error())
		return
	}
	for _, customer := range customerList {
		fmt.Println(customer.UserId, "\t", customer.Name, " \t", customer.Email, " \t", customer.Phone, " \t", customer.Address)
	}
}

func (admin *Admin) ViewAllPharmacists() {
	pharmacists, err := db.GetAllPharmacists()
	if err != nil {
		log.Error("Error while viewing all pharmacists : ", err.Error())
		return
	}
	for _, pharmacist := range pharmacists {
		fmt.Println(pharmacist.UserId, "\t", pharmacist.Name, " \t", pharmacist.Email, " \t", pharmacist.Phone, " \t", pharmacist.Address)
	}
}

func (admin *Admin) DeletePharmacist() {
	admin.ViewAllPharmacists()
	fmt.Print("Enter pharmacist id to delete : ")
	var pharmacistId int
	fmt.Scanln(&pharmacistId)
	err := db.DeletePharmacist(pharmacistId)
	if err != nil {
		log.Error("Error while deleting pharmacist : ", err.Error())
		return
	}
	log.Info("Pharmacist deleted successfully.")
}

func (admin *Admin) ViewAllMedicine() {
	medicineList, err := medicine.GetAllMedicine()
	if err != nil {
		log.Error("Error while viewing all medicine : ", err.Error())
		return
	}
	for _, medicine := range medicineList {
		fmt.Println(medicine.Id, " \t", medicine.Name, " \t", medicine.Price, " \t", medicine.Quantity, " \t", medicine.Manufacturer)
	}
}

func (admin *Admin) Logout() {
	log.Info("Logging out...")
}
