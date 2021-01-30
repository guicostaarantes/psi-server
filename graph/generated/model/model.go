// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type AuthenticateUserInput struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	IPAddress string `json:"ipAddress"`
}

type CreatePatientInput struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type CreateUserInput struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Role      Role   `json:"role"`
}

type ResetPasswordInput struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

type Token struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiresAt"`
}

type UpdatePatientInput struct {
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
}

type UpdateUserInput struct {
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Role      *Role   `json:"role"`
}

type User struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Role      Role   `json:"role"`
}

type Role string

const (
	RoleCoordinator  Role = "COORDINATOR"
	RolePsychologist Role = "PSYCHOLOGIST"
	RolePatient      Role = "PATIENT"
)

var AllRole = []Role{
	RoleCoordinator,
	RolePsychologist,
	RolePatient,
}

func (e Role) IsValid() bool {
	switch e {
	case RoleCoordinator, RolePsychologist, RolePatient:
		return true
	}
	return false
}

func (e Role) String() string {
	return string(e)
}

func (e *Role) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Role(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Role", str)
	}
	return nil
}

func (e Role) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
