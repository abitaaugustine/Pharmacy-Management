package user

import (
	"fmt"
	"pharmacy_management/internal/order"
	"pharmacy_management/internal/prescription"
)

type Customer struct {
	user User
}

func (customer *Customer) OpenUserMenu(user User) {
	customer.user = user
	log.Info("Welcome to the Pharmacy Management System! You are logged in as ", user.Name)
	for {
		fmt.Print("\n1. Place an order\n" +
			"2. Add a prescription\n" +
			"3. View previous orders\n" +
			"4. View previous prescriptions\n" +
			"5. Logout\n" +
			"Enter your choice : ")
		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			log.Info("Placing an order...")
			customer.PlaceOrder()
		case 2:
			log.Info("Adding a prescription...")
			customer.AddPrescription()
		case 3:
			log.Info("Viewing previous orders...")
			customer.ViewPreviousOrders()
		case 4:
			log.Info("Viewing previous prescriptions...")
			customer.ViewPreviousPrescriptions()
		case 5:
			customer.Logout()
			log.Info("Logged out successfully.")
			return
		default:
			log.Info("Invalid option. Please try again.")
		}
	}
}

func (customer *Customer) PlaceOrder() {
	err := order.PlaceOrder(customer.user.UserId)
	if err != nil {
		log.Error("Error while placing order : ", err.Error())
		return
	}
	log.Info("Order placed successfully.")
}

func (customer *Customer) AddPrescription() {
	err := prescription.AddPrescription(customer.user.UserId)
	if err != nil {
		log.Error("Error while adding prescription : ", err.Error())
		return
	}
	log.Info("Prescription added successfully.")
}

func (customer *Customer) ViewPreviousOrders() {
	previousOrders, err := order.ViewPreviousOrders(customer.user.UserId)
	if err != nil {
		log.Error("Error while viewing previous orders : ", err.Error())
		return
	}
	for _, order := range previousOrders {
		fmt.Print("\n\n", order.OrderId, " \t", order.Order_Date)
		for i, medicine := range order.Medicine {
			fmt.Println(" \t", medicine.Name, "\t", medicine.Manufacturer, " \t", order.Quantity[i], " \t", order.Price[i])
		}
	}
}

func (customer *Customer) ViewPreviousPrescriptions() {
	previousPrescriptions, err := prescription.ViewPreviousPrescriptions(customer.user.UserId)
	if err != nil {
		log.Error("Error while viewing previous prescriptions : ", err.Error())
		return
	}
	for _, prescription := range previousPrescriptions {
		fmt.Print("\n\n", prescription.Id, " \t", prescription.Doctor, " \t", prescription.PrescriptionDate)
		for i, medicine := range prescription.Medicine {
			fmt.Println(" \t", medicine.Name, "\t", medicine.Manufacturer, " \t", prescription.Frequency[i], " \t", prescription.TimeOfAdministration[i])
		}
	}
}

func (customer *Customer) Logout() {
	log.Info("Logging out...")
}
