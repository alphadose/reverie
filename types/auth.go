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

// PasswordUpdate is the request body for updating a user's password
type PasswordUpdate struct {
	OldPassword string `form:"old_password" json:"old_password,omitempty" binding:"required"`
	NewPassword string `form:"new_password" json:"new_password,omitempty" binding:"required"`
}

// GetOldPassword returns the user's old password
func (pw *PasswordUpdate) GetOldPassword() string {
	return pw.OldPassword
}

// GetNewPassword returns the user's new password
// which will replace the old password
func (pw *PasswordUpdate) GetNewPassword() string {
	return pw.NewPassword
}
