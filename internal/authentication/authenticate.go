package authentication

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/pkg/errors"
	"os"
	"regexp"
	"strings"

	"pharmacy_management/internal/db"
	"pharmacy_management/internal/logformated"
	user1 "pharmacy_management/internal/user"
)

var log = logformated.GetLogger(logformated.ComponentAuthentication)

func Register() error {
	var user user1.User
	fmt.Print("\n Register as : 1. Customer \n2. Pharmacist \n3. Administrator \n Enter your choice : ")
	_, err := fmt.Scanln(&user.Role.RoleId)
	if err != nil {
		return errors.Wrap(err, "Error while reading input from user : ")
	}
	fmt.Print("\n Enter user name : ")
	_, err = fmt.Scanln(&user.UserName)
	if err != nil {
		return errors.Wrap(err, "Error while reading username from user : ")
	}
	fmt.Print("\n Enter password : ")
	passwordByte, err := gopass.GetPasswdMasked()
	if err != nil {
		return errors.Wrap(err, "Error while reading password from user : ")
	}
	fmt.Print("\n Enter name : ")
	_, err = fmt.Scanln(&user.Name)
	if err != nil {
		return errors.Wrap(err, "Error while reading name from user : ")
	}
	fmt.Print("\n Enter phone(10 digit mobile number) : ")
	_, err = fmt.Scanln(&user.Phone)
	if err != nil {
		return errors.Wrap(err, "Error while reading phone from user : ")
	}
	fmt.Print("\n Enter email : ")
	_, err = fmt.Scanln(&user.Email)
	if err != nil {
		return errors.Wrap(err, "Error while reading email from user : ")
	}
	fmt.Print("\n Enter address (as comma separated values) : ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	if err != nil {
		return errors.Wrap(err, "Error while reading address from user : ")
	}
	user.Address = strings.TrimSpace(input)
	if !validUserRegister(user, passwordByte) {
		return errors.New("Invalid input. Please try again")
	}
	passwordHash := sha256.Sum256(passwordByte)
	user.PasswordHash = string(passwordHash[:])
	if db.InsertUser(user.UserName, user.PasswordHash, user.Name, user.Phone, user.Email, user.Address, user.Role.RoleId) != nil {
		return errors.New("Error while registering user. Please try again")
	}
	return nil
}

func Login() (error, user1.User) {
	var user user1.User
	var userName string
	fmt.Print("\n Enter user name : ")
	_, err := fmt.Scanln(&userName)
	if err != nil {
		return errors.Wrap(err, "Error while reading username from user : "), user
	}
	fmt.Print("\n Enter password : ")
	passwordByte, err := gopass.GetPasswdMasked()
	if err != nil {
		return errors.Wrap(err, "Error while reading password from user : "), user
	}
	passwordHash := sha256.Sum256(passwordByte)
	PasswordHashString := string(passwordHash[:])
	if userId, name, phone, email, address, roleId, success := db.VerifyAndGetUser(userName, PasswordHashString); success {
		user.UserId = userId
		user.UserName = userName
		user.PasswordHash = PasswordHashString
		user.Name = name
		user.Phone = phone
		user.Email = email
		user.Address = address
		user.Role = user1.DefaultRole[roleId]
		return nil, user
	}
	return errors.New("Invalid user name or password. Please try again"), user
}

func validUserRegister(user user1.User, passwordByte []byte) bool {
	err := ""
	if user.Role.RoleId < 1 || user.Role.RoleId > 3 {
		err += "Invalid role id. "
	}
	if len(user.UserName) < 1 {
		err += "Invalid user name. "
	}
	if !db.IsUnusedUserName(user.UserName) {
		err += "User name already exists. "
	}
	password := string(passwordByte)
	if len(password) < 4 && len(password) > 16 {
		err += "Password must be 4-16 characters long"
	}
	if len(user.Name) < 1 {
		err += "Invalid name. "
	}
	match, _ := regexp.MatchString("^[0-9]{10}$", user.Phone)
	if !match {
		err += "Invalid phone number. Phone Number should be 10 digits long. "
	}
	match, _ = regexp.MatchString("^[a-zA-Z0-9+_.-]+@[a-zA-Z0-9.-]+$", user.Email)
	if !match {
		err += "Invalid email. "
	}
	if len(user.Address) < 1 {
		err += "Invalid address. "
	}
	if len(err) > 0 {
		log.Info(err)
		return false
	}
	return true
}

//mettlexams@wilp.bits-pilani.ac.in
