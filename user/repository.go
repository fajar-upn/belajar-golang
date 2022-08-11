package user

import "gorm.io/gorm"

type Repository interface {
	Save(user User) (User, error)
	FindByEmail(email string) (User, error)
}

type repository struct {
	db *gorm.DB
}

// this code call in main.go
func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

// input to database
func (r *repository) Save(user User) (User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *repository) FindByEmail(email string) (User, error) {

	var user User

	err := r.db.Where("email = ?", email).Find(&user).Error // ? is value from email
	if err != nil {
		return user, err
	}

	return user, nil
}
