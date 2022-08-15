package user

import "gorm.io/gorm"

type Repository interface {
	Save(user User) (User, error)
	FindByEmail(email string) (User, error)
	FindByID(ID int) (User, error)
	Update(user User) (User, error)
}

type repository struct {
	db *gorm.DB
}

// this code call in main.go
func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

// input to table user
func (r *repository) Save(user User) (User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

// get user email
func (r *repository) FindByEmail(email string) (User, error) {

	var user User

	err := r.db.Where("email = ?", email).Find(&user).Error // ? is a value from email
	if err != nil {
		return user, err
	}

	return user, nil
}

// get user id
func (r *repository) FindByID(ID int) (User, error) {
	var user User

	err := r.db.Where("id = ?", ID).Find(&user).Error // ? is a valu from ID
	if err != nil {
		return user, err
	}

	return user, nil
}

// update avatar user
func (r *repository) Update(user User) (User, error) {
	err := r.db.Save(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}
