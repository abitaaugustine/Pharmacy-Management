package db

import (
	"database/sql"
	"pharmacy_management/internal/logformated"
	"time"

	"github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

var log = logformated.GetLogger(logformated.ComponentDB)

type UserDetails struct {
	UserId  int
	Name    string
	Phone   string
	Email   string
	Address string
}

type MedicineDetails struct {
	MedicineId   int
	Name         string
	Manufacturer string
	Quantity     int
	Price        float64
}

type OrderDetails struct {
	OrderId      int
	OrderDate    time.Time
	MedicineId   []int
	MedicineName []string
	Manufacturer []string
	Price        []float64
	Quantity     []int
}

type PrescriptionDetails struct {
	PrescriptionId       int
	Doctor               string
	PrescriptionDate     time.Time
	MedicineId           []int
	MedicineName         []string
	Manufacturer         []string
	Frequency            []string
	TimeOfAdministration []string
}

const DbFileName = "./pharmacy.db"

func RegisterDriver() {
	sql.Register("sqlite3Driver", &sqlite3.SQLiteDriver{})
}

func DbConnection() (*sql.DB, func()) {
	db, err := sql.Open("sqlite3", DbFileName)
	if err != nil {
		log.Errorf("Failed to open database : %s", err.Error())
	}
	return db, func() {
		err := db.Close()
		if err != nil {
			log.Errorf("Error while trying to close database : %s", err.Error())
		}
	}
}

func ExecuteQuery(dbConnection *sql.DB, query string) {
	if _, err := dbConnection.Exec(query); err != nil {
		log.Errorf("Error while executing query : %s", err.Error())
	}
}

func CreateTables() {
	dbConnection, closeDb := DbConnection()
	defer closeDb()
	log.Debug("Creating tables if not exists already")
	ExecuteQuery(dbConnection, "CREATE TABLE IF NOT EXISTS Role( Role_Id INTEGER PRIMARY KEY AUTOINCREMENT, Role_Type VARCHAR)")
	ExecuteQuery(dbConnection, "CREATE TABLE IF NOT EXISTS User( User_Id INTEGER PRIMARY KEY AUTOINCREMENT, User_Name VARCHAR NOT NULL UNIQUE, Password_Hash VARCHAR, Name VARCHAR, Phone VARCHAR, Email VARCHAR NOT NULL UNIQUE, Address VARCHAR, Role_Id INTEGER, FOREIGN KEY(Role_Id) REFERENCES Role(Role_Id))")
	ExecuteQuery(dbConnection, "CREATE TABLE IF NOT EXISTS Prescription( Prescription_Id INTEGER PRIMARY KEY AUTOINCREMENT, Doctor VARCHAR, Prescription_date DATE, User_Id INTEGER, FOREIGN KEY(User_Id) REFERENCES User(User_Id))")
	ExecuteQuery(dbConnection, "CREATE TABLE IF NOT EXISTS Order_Placed( Order_Id INTEGER PRIMARY KEY AUTOINCREMENT, Order_Date DATETIME, User_Id INTEGER, FOREIGN KEY(User_Id) REFERENCES User(User_Id))")
	ExecuteQuery(dbConnection, "CREATE TABLE IF NOT EXISTS Medicine( Medicine_Id INTEGER PRIMARY KEY AUTOINCREMENT, Name VARCHAR, Manufacturer VARCHAR, Stock_Available INT, Price FLOAT)")
	ExecuteQuery(dbConnection, "CREATE TABLE IF NOT EXISTS Prescription_Medicine( Prescription_Id INTEGER, Medicine_Id INTEGER, Frequency VARCHAR, Time_of_Administration VARCHAR, FOREIGN KEY(Prescription_Id) REFERENCES Prescription(Prescription_Id), FOREIGN KEY(Medicine_Id) REFERENCES Medicine(Medicine_Id), PRIMARY KEY(Prescription_Id, Medicine_Id))")
	ExecuteQuery(dbConnection, "CREATE TABLE IF NOT EXISTS Order_Medicine( Order_Id INTEGER, Medicine_Id INTEGER, Quantity INT, Price FLOAT, FOREIGN KEY(Order_Id) REFERENCES Order_Placed(Order_Id), FOREIGN KEY(Medicine_Id) REFERENCES Medicine(Medicine_Id),  PRIMARY KEY(Order_Id, Medicine_Id))")
	log.Debug("Tables created if not exists already")
}

func InsertDefaultValuesIntoRole() {
	dbConnection, closeDb := DbConnection()
	defer closeDb()
	log.Debug("Inserting default values to role if not done already")
	ExecuteQuery(dbConnection, "INSERT OR REPLACE INTO Role VALUES(1, 'Customer'), (2, 'Pharmacist'), (3, 'Administrator')")
}

func LoadDefaultRoleList() map[int]string {
	dbConnection, closeDb := DbConnection()
	defer closeDb()

	rows, err := dbConnection.Query("SELECT * FROM Role")
	if err != nil {
		log.Errorf("Error while executing query : %s", err.Error())
		return nil
	}
	defer rows.Close()

	roleList := make(map[int]string)
	var roleId int
	var roleType string

	for rows.Next() {
		err := rows.Scan(&roleId, &roleType)
		if err != nil {
			log.Errorf("Error while scanning rows : %s", err.Error())
			return nil
		}
		roleList[roleId] = roleType
	}
	return roleList
}

func IsUnusedUserName(userName string) bool {
	dbConnection, closeDb := DbConnection()
	defer closeDb()

	row := dbConnection.QueryRow("SELECT User_Name FROM User WHERE User_Name = ?", userName)

	var userNameFromDb string
	err := row.Scan(&userNameFromDb)
	if err == sql.ErrNoRows {
		return true
	} else if err != nil {
		log.Errorf("Error while scanning rows : %s", err.Error())
		return false
	}
	log.Error("User name already exists. Please try again with a different user name. ")
	return false
}

func InsertUser(userName, passwordHash, name, phone, email, address string, roleId int) error {
	dbConnection, closeDb := DbConnection()
	defer closeDb()

	query := "INSERT INTO User(User_Name, Password_Hash, Name, Phone, Email, Address, Role_Id) VALUES(?, ?, ?, ?, ?, ?, ?)"

	if _, err := dbConnection.Exec(query, userName, passwordHash, name, phone, email, address, roleId); err != nil {
		return errors.Wrap(err, "Error while inserting user details to db : ")
	}
	return nil
}

func VerifyAndGetUser(userName, passwordHash string) (userId int, name string, phone string, email string, address string, roleId int, isSuccess bool) {
	dbConnection, closeDb := DbConnection()
	defer closeDb()

	row := dbConnection.QueryRow("SELECT * FROM User WHERE User_Name = ? AND Password_Hash = ?", userName, passwordHash)

	err := row.Scan(&userId, &userName, &passwordHash, &name, &phone, &email, &address, &roleId)
	if err == sql.ErrNoRows {
		return userId, name, phone, email, address, roleId, false
	} else if err != nil {
		log.Errorf("Error while scanning rows : %s", err.Error())
		return userId, name, phone, email, address, roleId, false
	}
	return userId, name, phone, email, address, roleId, true
}

func GetAllCustomers() (map[int]UserDetails, error) {
	dbConnection, closeDb := DbConnection()
	defer closeDb()

	rows, err := dbConnection.Query("SELECT u.User_Id, u.Name, u.Phone, u.Email, u.Address FROM User u, Role r WHERE r.Role_Type = 'Customer' AND u.Role_Id = r.Role_Id")
	if err != nil {
		log.Errorf("Error while executing query : %s", err.Error())
		return nil, err
	}
	defer rows.Close()

	return processUserDetails(rows)
}

func GetAllPharmacists() (map[int]UserDetails, error) {
	dbConnection, closeDb := DbConnection()
	defer closeDb()

	rows, err := dbConnection.Query("SELECT u.User_Id, u.Name, u.Phone, u.Email, u.Address FROM User u, Role r WHERE r.Role_Type = 'Pharmacist' AND u.Role_Id = r.Role_Id")
	if err != nil {
		log.Errorf("Error while executing query : %s", err.Error())
		return nil, err
	}
	defer rows.Close()
	return processUserDetails(rows)
}

func DeletePharmacist(pharmacistId int) error {
	dbConnection, closeDb := DbConnection()
	defer closeDb()

	query := "DELETE FROM User WHERE User_Id = ? and Role_Id = 2"

	if _, err := dbConnection.Exec(query, pharmacistId); err != nil {
		return errors.Wrap(err, "Error while deleting pharmacist details from db : ")
	}
	return nil
}

func processUserDetails(rows *sql.Rows) (map[int]UserDetails, error) {
	userDetails := make(map[int]UserDetails)
	var userId int
	var name string
	var phone string
	var email string
	var address string

	for rows.Next() {
		err := rows.Scan(&userId, &name, &phone, &email, &address)
		if err != nil {
			log.Errorf("Error while scanning rows : %s", err.Error())
			return userDetails, err
		}
		userDetails[userId] = UserDetails{userId, name, phone, email, address}
	}
	return userDetails, nil
}

func AddPrescription(userId int, doctor string, medicineId int, frequency string, timeOfAdministration string) error {
	dbConnection, closeDb := DbConnection()
	defer closeDb()

	query := "INSERT INTO Prescription(Doctor, Prescription_date, User_Id) VALUES(?, ?, ?)"
	prescriptionDate := time.Now()
	if _, err := dbConnection.Exec(query, doctor, prescriptionDate, userId); err != nil {
		return errors.Wrap(err, "Error while inserting prescription details to db : ")
	}

	prescriptionId, err := LastInsertedPrescriptionId(userId, prescriptionDate)
	if err != nil {
		return errors.Wrap(err, "Error while getting last inserted prescription id : ")
	}

	query = "INSERT INTO Prescription_Medicine(Prescription_Id, Medicine_Id, Frequency, Time_of_Administration) VALUES(?, ?, ?, ?)"
	if _, err := dbConnection.Exec(query, prescriptionId, medicineId, frequency, timeOfAdministration); err != nil {
		return errors.Wrap(err, "Error while inserting prescription medicine details to db : ")
	}
	return nil
}

func LastInsertedPrescriptionId(userId int, prescriptionDate time.Time) (int, error) {
	dbConnection, closeDb := DbConnection()
	defer closeDb()

	row := dbConnection.QueryRow("SELECT Prescription_Id FROM Prescription WHERE User_Id = ? AND Prescription_date = ?", userId, prescriptionDate)

	var prescriptionId int
	err := row.Scan(&prescriptionId)
	if err == sql.ErrNoRows {
		return prescriptionId, nil
	} else if err != nil {
		log.Errorf("Error while scanning rows : %s", err.Error())
		return prescriptionId, err
	}
	return prescriptionId, nil
}

func AddOrder(userId int, medicineIdList []int, quantityList []int, priceList []float64) error {
	dbConnection, closeDb := DbConnection()
	defer closeDb()

	query := "INSERT INTO Order_Placed(Order_Date, User_Id) VALUES(?, ?)"
	t := time.Now()
	if _, err := dbConnection.Exec(query, t, userId); err != nil {
		return errors.Wrap(err, "Error while inserting order details to db : ")
	}

	orderId, err := LastInsertedOrderId(userId, t)
	if err != nil {
		return errors.Wrap(err, "Error while getting last inserted order id : ")
	}

	query = "INSERT INTO Order_Medicine(Order_Id, Medicine_Id, Quantity, Price) VALUES(?, ?, ?, ?)"
	for i, medicineId := range medicineIdList {
		if _, err := dbConnection.Exec(query, orderId, medicineId, quantityList[i], priceList[i]); err != nil {
			return errors.Wrap(err, "Error while inserting order medicine details to db : ")
		}
	}

	query = "UPDATE Medicine SET Stock_Available = Stock_Available - ? WHERE Medicine_Id = ?"
	for i, medicineId := range medicineIdList {
		if _, err := dbConnection.Exec(query, quantityList[i], medicineId); err != nil {
			return errors.Wrap(err, "Error while updating medicine stock details after placing order : ")
		}
	}
	return nil
}

func LastInsertedOrderId(userId int, t time.Time) (int, error) {
	dbConnection, closeDb := DbConnection()
	defer closeDb()

	row := dbConnection.QueryRow("SELECT Order_Id FROM Order_Placed WHERE User_Id = ? AND Order_Date = ?", userId, t)

	var orderId int
	err := row.Scan(&orderId)
	if err != nil {
		log.Errorf("Error while scanning rows : %s", err.Error())
		return orderId, err
	}
	return orderId, nil
}

func GetAllMedicine() (map[int]MedicineDetails, error) {
	dbConnection, closeDb := DbConnection()
	defer closeDb()

	rows, err := dbConnection.Query("SELECT * FROM Medicine")
	if err != nil {
		log.Errorf("Error while executing query : %s", err.Error())
		return nil, err
	}
	defer rows.Close()

	medicineList := make(map[int]MedicineDetails)
	var medicineId int
	var name string
	var manufacturer string
	var stockAvailable int
	var price float64

	for rows.Next() {
		err := rows.Scan(&medicineId, &name, &manufacturer, &stockAvailable, &price)
		if err != nil {
			log.Errorf("Error while scanning rows : %s", err.Error())
			return nil, err
		}
		medicineList[medicineId] = MedicineDetails{
			MedicineId:   medicineId,
			Name:         name,
			Manufacturer: manufacturer,
			Quantity:     stockAvailable,
			Price:        price,
		}
	}
	return medicineList, nil
}

func AddNewMedicine(name string, price float64, quantity int, manufacturer string) error {
	dbConnection, closeDb := DbConnection()
	defer closeDb()

	query := "INSERT OR REPLACE INTO Medicine(Name, Price, Stock_Available, Manufacturer) VALUES(?, ?, ?, ?)"

	if _, err := dbConnection.Exec(query, name, price, quantity, manufacturer); err != nil {
		return errors.Wrap(err, "Error while inserting medicine details to db : ")
	}
	return nil
}

func UpdateMedicine(medicineId int, name string, price float64, quantity int, manufacturer string) error {
	dbConnection, closeDb := DbConnection()
	defer closeDb()

	query := "UPDATE Medicine SET Name = ?, Price = ?, Stock_Available = ?, Manufacturer = ? WHERE Medicine_Id = ?"

	if _, err := dbConnection.Exec(query, name, price, quantity, manufacturer, medicineId); err != nil {
		return errors.Wrap(err, "Error while updating medicine details to db : ")
	}
	return nil
}

func IsMedicinePresent(medicineId int) bool {
	dbConnection, closeDb := DbConnection()
	defer closeDb()

	row := dbConnection.QueryRow("SELECT Medicine_Id FROM Medicine WHERE Medicine_Id = ?", medicineId)

	var medicineIdFromDb int
	err := row.Scan(&medicineIdFromDb)
	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		log.Errorf("Error while scanning rows : %s", err.Error())
		return false
	}
	return true
}

func GetPreviousOrders(userId int) (map[int]OrderDetails, error) {
	dbConnection, closeDb := DbConnection()
	defer closeDb()

	rows, err := dbConnection.Query("SELECT o.Order_Id, o.Order_Date, m.Medicine_Id, m.Name, m.Manufacturer, om.Quantity, om.Price FROM Order_Placed o JOIN Order_Medicine om ON o.Order_Id = om.Order_Id JOIN Medicine m ON om.Medicine_Id = m.Medicine_Id WHERE o.User_Id = ?", userId)
	if err != nil {
		log.Errorf("Error while executing query : %s", err.Error())
		return nil, err
	}
	defer rows.Close()

	var orderList = make(map[int]OrderDetails)
	var orderId int
	var orderDate time.Time
	var medicineId int
	var medicineName string
	var manufacturer string
	var price float64
	var quantity int

	for rows.Next() {
		err := rows.Scan(&orderId, &orderDate, &medicineId, &medicineName, &manufacturer, &quantity, &price)
		if err != nil {
			log.Errorf("Error while scanning rows : %s", err.Error())
			return nil, err
		}
		if _, ok := orderList[orderId]; ok {
			order := orderList[orderId]
			order.MedicineId = append(order.MedicineId, medicineId)
			order.MedicineName = append(order.MedicineName, medicineName)
			order.Manufacturer = append(order.Manufacturer, manufacturer)
			order.Price = append(order.Price, price)
			order.Quantity = append(order.Quantity, quantity)
			orderList[orderId] = order
		} else {
			orderList[orderId] = OrderDetails{
				OrderId:      orderId,
				OrderDate:    orderDate,
				MedicineId:   []int{medicineId},
				MedicineName: []string{medicineName},
				Manufacturer: []string{manufacturer},
				Price:        []float64{price},
				Quantity:     []int{quantity},
			}
		}
	}
	return orderList, nil
}

func GetPreviousPrescriptions(userId int) (map[int]PrescriptionDetails, error) {
	dbConnection, closeDb := DbConnection()
	defer closeDb()

	rows, err := dbConnection.Query("SELECT p.Prescription_Id, p.Doctor, p.Prescription_Date, m.Medicine_Id, m.Name, m.Manufacturer, pm.Frequency, pm.Time_of_Administration FROM Prescription p JOIN Prescription_Medicine pm ON p.Prescription_Id = pm.Prescription_Id JOIN Medicine m ON pm.Medicine_Id = m.Medicine_Id WHERE p.User_Id = ?", userId) // TODO: verify the correctness of this query to handle m:n prescription medicine
	if err != nil {
		log.Errorf("Error while executing query : %s", err.Error())
		return nil, err
	}
	defer rows.Close()

	var prescriptionList = make(map[int]PrescriptionDetails)
	var prescriptionId int
	var doctor string
	var prescriptionDate time.Time
	var medicineId int
	var medicineName string
	var manufacturer string
	var frequency string
	var timeOfAdministration string

	for rows.Next() {
		err := rows.Scan(&prescriptionId, &doctor, &prescriptionDate, &medicineId, &medicineName, &manufacturer, &frequency, &timeOfAdministration)
		if err != nil {
			log.Errorf("Error while scanning rows : %s", err.Error())
			return nil, err
		}
		if _, ok := prescriptionList[prescriptionId]; ok {
			prescription := prescriptionList[prescriptionId]
			prescription.MedicineId = append(prescription.MedicineId, medicineId)
			prescription.MedicineName = append(prescription.MedicineName, medicineName)
			prescription.Manufacturer = append(prescription.Manufacturer, manufacturer)
			prescription.Frequency = append(prescription.Frequency, frequency)
			prescription.TimeOfAdministration = append(prescription.TimeOfAdministration, timeOfAdministration)
			prescriptionList[prescriptionId] = prescription
		} else {
			prescriptionList[prescriptionId] = PrescriptionDetails{
				PrescriptionId:       prescriptionId,
				Doctor:               doctor,
				PrescriptionDate:     prescriptionDate,
				MedicineId:           []int{medicineId},
				MedicineName:         []string{medicineName},
				Manufacturer:         []string{manufacturer},
				Frequency:            []string{frequency},
				TimeOfAdministration: []string{timeOfAdministration},
			}
		}
	}
	return prescriptionList, nil
}
