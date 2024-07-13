package order

import (
	"fmt"
	"time"

	"pharmacy_management/internal/db"
	"pharmacy_management/internal/medicine"
)

type Order struct {
	OrderId    int
	Order_Date time.Time
	Medicine   []medicine.Medicine
	Quantity   []int
	Price      []float64
}

func PlaceOrder(userId int) error {
	var medicineIdList []int
	var quantityList []int

	medicineList, err := medicine.GetAllMedicine()
	if err != nil {
		return err
	}
	for _, medicine := range medicineList {
		fmt.Println(medicine.Id, " \t", medicine.Name, " \t", medicine.Price, " \t", medicine.Quantity, " \t", medicine.Manufacturer)
	}
	fmt.Println("Enter no.of medicines to add to cart: ")
	var n int
	fmt.Scanln(&n)
	for i := 0; i < n; i++ {
		var medicineId int
		fmt.Print("\n Enter medicine ID: ")
		fmt.Scanln(&medicineId)
		medicineIdList = append(medicineIdList, medicineId)

		var quantity int
		fmt.Print("\n Enter quantity: ")
		fmt.Scanln(&quantity)
		quantityList = append(quantityList, quantity)
	}

	if verifyEnteredMedicine(medicineIdList, quantityList, medicineList) {
		priceList := calculatePrice(medicineIdList, quantityList, medicineList)
		err := db.AddOrder(userId, medicineIdList, quantityList, priceList)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("invalid medicine id or quantity")
	}
	return nil
}

func calculatePrice(medicineIdList []int, quantityList []int, medicineList []medicine.Medicine) []float64 {
	var priceList []float64
	for i, medicineId := range medicineIdList {
		for _, medicine := range medicineList {
			if medicineId == medicine.Id {
				priceList = append(priceList, medicine.Price*float64(quantityList[i]))
			}
		}
	}
	return priceList
}

func verifyEnteredMedicine(medicineIdList []int, quantityList []int, medicineList []medicine.Medicine) bool {
	if len(medicineIdList) == 0 {
		return false
	}
	for i, medicineId := range medicineIdList {
		flag := false
		for _, medicine := range medicineList {
			if medicineId == medicine.Id {
				flag = true
				if quantityList[i] > medicine.Quantity {
					return false
				}
			}
		}
		if !flag {
			return false
		}
	}
	return true
}

func ViewPreviousOrders(userId int) ([]Order, error) {
	oderDetails, err := db.GetPreviousOrders(userId)
	if err != nil {
		return nil, err
	}

	var orderList []Order
	for _, o := range oderDetails {
		order := Order{
			OrderId:    o.OrderId,
			Order_Date: o.OrderDate,
		}
		for i, medicineId := range o.MedicineId {
			order.Medicine = append(order.Medicine, medicine.Medicine{Id: medicineId, Name: o.MedicineName[i], Manufacturer: o.Manufacturer[i]})
			order.Quantity = append(order.Quantity, o.Quantity[i])
			order.Price = append(order.Price, o.Price[i])
		}
		orderList = append(orderList, order)
	}

	return orderList, nil
}
