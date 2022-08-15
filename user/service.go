package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

/**
1. mapping input struct to struct user
2. save user struct through repository
*/

// interface for connected another file or code to service.go
type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
	Login(input LoginInput) (User, error)
	IsEmailAvailable(input CheckEmailInput) (bool, error)
	SaveAvatar(ID int, fileLocation string) (User, error)
}

type service struct {
	repository Repository
}

// this code call in main.go
func NewService(repository Repository) *service {
	return &service{repository}
}

// register user
func (s *service) RegisterUser(input RegisterUserInput) (User, error) {

	user := User{}
	user.Name = input.Name
	user.Email = input.Email
	user.Occupation = input.Occupation
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}
	user.PasswordHash = string(passwordHash)
	user.Role = "user"

	newUser, err := s.repository.Save(user)
	if err != nil {
		return newUser, err
	}

	return newUser, nil
}

// login user/ session
func (s *service) Login(input LoginInput) (User, error) {
	// get email and password from user.input.go
	email := input.Email
	password := input.Password

	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return user, err
	}

	// check available user email
	if user.ID == 0 {
		return user, errors.New("User not found!")
	}

	// check match password when user available
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return user, err
	}

	return user, nil
}

// check availability email
func (s *service) IsEmailAvailable(input CheckEmailInput) (bool, error) {

	email := input.Email
	user, err := s.repository.FindByEmail(email)

	if err != nil {
		return false, err
	}

	if user.ID == 0 {
		return true, nil
	}

	return false, nil //looking for safety valu return will be 'false', instead of user have double or same email
}

// upload avatar user
func (s *service) SaveAvatar(ID int, fileLocation string) (User, error) {

	user, err := s.repository.FindByID(ID) //1.find user appropriate ID
	if err != nil {
		return user, err
	}

	user.AvatarFileName = fileLocation //2.update attribute avatar filename

	updatedUser, err := s.repository.Update(user) //3.save avatar to database
	if err != nil {
		return updatedUser, err
	}
	return updatedUser, nil
}
