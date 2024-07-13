package user

import (
	"pharmacy_management/internal/logformated"
	"pharmacy_management/internal/db"
)

var log = logformated.GetLogger(logformated.ComponentUser)

type IUser interface {
	OpenUserMenu(user User)
	Logout()
}

type User struct {
	UserId       int
	UserName     string
	PasswordHash string
	Name         string
	Phone        string
	Email        string
	Address      string
	Role         Role
}

type Role struct {
	RoleId   int
	RoleType string
}

var DefaultRole map[int]Role = make(map[int]Role)

func LoadDefaultRole() {
	roleList := db.LoadDefaultRoleList()
	for roleId, role := range roleList {
		DefaultRole[roleId] = Role{RoleId: roleId, RoleType: role}
	}
}
