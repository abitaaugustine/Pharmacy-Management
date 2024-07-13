package prescription

import (
	"fmt"
	"pharmacy_management/internal/db"
	"pharmacy_management/internal/medicine"
	"time"
)

type Prescription struct {
	Id                   int
	Doctor               string
	PrescriptionDate     time.Time
	Medicine             []medicine.Medicine
	Frequency            []string
	TimeOfAdministration []string
}

func AddPrescription(userId int) error {
	var doctor string
	fmt.Print("\n Enter doctor name: ")
	fmt.Scanln(&doctor)

	medicineList, err := medicine.GetAllMedicine()
	for _, medicine := range medicineList {
		fmt.Println(medicine.Id, " \t", medicine.Name)
	}

	fmt.Println("Enter no.of medicines to add to prescription: ")
	var n int
	fmt.Scanln(&n)

	var medicineIdList []int
	var frequencyList []string
	var timeOfAdministrationList []string

	for i := 0; i < n; i++ {
		var medicineId int
		fmt.Print("\n Enter medicine ID: ")
		fmt.Scanln(&medicineId)

		var frequency string
		fmt.Print("\n Enter frequency: ")
		fmt.Scanln(&frequency)

		var timeOfAdministration string
		fmt.Print("\n Enter time of administration (before/after food): ")
		fmt.Scanln(&timeOfAdministration)

		medicineIdList = append(medicineIdList, medicineId)
		frequencyList = append(frequencyList, frequency)
		timeOfAdministrationList = append(timeOfAdministrationList, timeOfAdministration)
	}

	if verifyEnteredMedicine(medicineIdList, medicineList) {
		for i := 0; i < n; i++ {
			err = db.AddPrescription(userId, doctor, medicineIdList[i], frequencyList[i], timeOfAdministrationList[i])
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func verifyEnteredMedicine(medicineIdList []int, medicineList []medicine.Medicine) bool {
	if len(medicineIdList) == 0 {
		return false
	}
	for _, medicineId := range medicineIdList {
		flag := false
		for _, medicine := range medicineList {
			if medicineId == medicine.Id {
				flag = true
				break
			}
		}
		if !flag {
			return false
		}
	}
	return true
}

func ViewPreviousPrescriptions(userId int) ([]Prescription, error) {
	prescriptionDetails, err := db.GetPreviousPrescriptions(userId)
	if err != nil {
		return nil, err
	}

	var prescriptionList []Prescription

	for _, p := range prescriptionDetails {
		prescription := Prescription{
			Id:               p.PrescriptionId,
			Doctor:           p.Doctor,
			PrescriptionDate: p.PrescriptionDate,
		}
		for i, medicineId := range p.MedicineId {
			prescription.Medicine = append(prescription.Medicine, medicine.Medicine{Id: medicineId, Name: p.MedicineName[i], Manufacturer: p.Manufacturer[i]})
			prescription.Frequency = append(prescription.Frequency, p.Frequency[i])
			prescription.TimeOfAdministration = append(prescription.TimeOfAdministration, p.TimeOfAdministration[i])
		}
		prescriptionList = append(prescriptionList, prescription)
	}

	return prescriptionList, nil
}
