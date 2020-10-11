package types

// Login is the request body binding for login
type Login struct {
	Email    string `form:"email" json:"email" bson:"email" binding:"required"`
	Password string `form:"password" json:"password,omitempty" bson:"password" binding:"required"`
}

// GetEmail returns the user's email
func (auth *Login) GetEmail() string {
	return auth.Email
}

// GetPassword returns the user's password
// The password will be hashed if retrieving from database
func (auth *Login) GetPassword() string {
	return auth.Password
}
