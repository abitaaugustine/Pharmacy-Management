package medicine

import (
	"fmt"
	"pharmacy_management/internal/db"
)

type Medicine struct {
	Id           int
	Name         string
	Price        float64
	Quantity     int
	Manufacturer string
}

func GetAllMedicine() ([]Medicine, error) {
	medicineDetails, err := db.GetAllMedicine()
	if err != nil {
		return nil, err
	}
	var medicineList []Medicine
	for _, medicine := range medicineDetails {
		medicineList = append(medicineList, Medicine{
			Id:           medicine.MedicineId,
			Name:         medicine.Name,
			Price:        medicine.Price,
			Quantity:     medicine.Quantity,
			Manufacturer: medicine.Manufacturer,
		})
	}
	return medicineList, nil

}

func AddOrUpdateMedicine() error {
	fmt.Println(GetAllMedicine())
	fmt.Println("Enter medicine id to update or 0 to add new medicine: ")
	var medicineId int
	fmt.Scanln(&medicineId)
	if medicineId != 0 {
		err := updateMedicine(medicineId)
		if err != nil {
			return err
		}
	} else {
		err := addMedicine()
		if err != nil {
			return err
		}
	}
	return nil
}

func updateMedicine(medicineId int) error {
	if !db.IsMedicinePresent(medicineId) {
		return fmt.Errorf("medicine with id %d not found", medicineId)
	}
	fmt.Println("Enter new details for medicine: ")
	var medicine Medicine
	fmt.Print("\n Enter medicine name: ")
	fmt.Scanln(&medicine.Name)
	fmt.Print("\n Enter medicine price: ")
	fmt.Scanln(&medicine.Price)
	fmt.Print("\n Enter medicine quantity: ")
	fmt.Scanln(&medicine.Quantity)
	fmt.Print("\n Enter medicine manufacturer: ")
	fmt.Scanln(&medicine.Manufacturer)

	err := db.UpdateMedicine(medicineId, medicine.Name, medicine.Price, medicine.Quantity, medicine.Manufacturer)
	if err != nil {
		return err
	}
	return nil
}

func addMedicine() error {
	var medicine Medicine
	fmt.Print("\n Enter medicine name: ")
	fmt.Scanln(&medicine.Name)
	fmt.Print("\n Enter medicine price: ")
	fmt.Scanln(&medicine.Price)
	fmt.Print("\n Enter medicine quantity: ")
	fmt.Scanln(&medicine.Quantity)
	fmt.Print("\n Enter medicine manufacturer: ")
	fmt.Scanln(&medicine.Manufacturer)

	err := db.AddNewMedicine(medicine.Name, medicine.Price, medicine.Quantity, medicine.Manufacturer)
	if err != nil {
		return err
	}
	return nil
}
