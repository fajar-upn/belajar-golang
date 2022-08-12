package user

/**
this struct represent input from register user form
*/
type RegisterUserInput struct {
	Name       string `json:"name" binding:"required"`
	Occupation string `json:"occupation" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
}

/**
this struct represent input from login user form
*/
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

/**
this struct represent input from register user form, when email available or not
*/
type CheckEmailInput struct {
	Email string `json:"email" binding:"required,email"`
}
