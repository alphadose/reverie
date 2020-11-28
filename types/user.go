package types

const (
	// Client holds the name  for the client role
	Client = "client"

	// Vendor holds the name  for the vendor role
	Vendor = "vendor"

	// Admin holds the name  for the admin role
	Admin = "admin"
)

// Claims store the JWT claims of a user
// Its a minified version of the User struct to reduce the memory footprint
type Claims struct {
	Email    string
	Username string
	Role     string
}

// GetName returns the user's username
func (claims *Claims) GetName() string {
	return claims.Username
}

// GetEmail returns the user's email
func (claims *Claims) GetEmail() string {
	return claims.Email
}

// IsClient checks of the user is a client or not
func (claims *Claims) IsClient() bool {
	return claims.Role == Client
}

// IsVendor checks of the user is a vendor or not
func (claims *Claims) IsVendor() bool {
	return claims.Role == Vendor
}

// IsAdmin checks of the user is an admin or not
func (claims *Claims) IsAdmin() bool {
	return claims.Role == Admin
}

// User stores user related information
type User struct {
	Email     string     `json:"email" binding:"required" bson:"email" valid:"required~Field 'email' is required but was not provided,email"`
	Password  string     `json:"password,omitempty" bson:"password" binding:"required" valid:"required~Field 'password' is required but was not provided"`
	Username  string     `json:"username" bson:"username" binding:"required" valid:"required~Field 'username' is required but was not provided,matches(^[\\p{L} .'-]+$)"`
	Role      string     `json:"-" bson:"role"`
	Inventory *Inventory `json:"inventory,omitempty" bson:"inventory,omitempty"`
	Success   bool       `json:"success,omitempty" bson:"-"`
}

// GetName returns the user's username
func (user *User) GetName() string {
	return user.Username
}

// GetEmail returns the user's email
func (user *User) GetEmail() string {
	return user.Email
}

// SetEmail sets the user's email in its context
func (user *User) SetEmail(email string) {
	user.Email = email
}

// GetPassword returns the user's password
// The password will be hashed if retrieving from database
func (user *User) GetPassword() string {
	return user.Password
}

// SetPassword sets a password in the user's context
func (user *User) SetPassword(password string) {
	user.Password = password
}

// SetRole sets the role in the user's context
func (user *User) SetRole(role string) {
	user.Role = role
}

// GetRole returns the user's current role
func (user *User) GetRole() string {
	return user.Role
}

// IsClient checks of the user is a client or not
func (user *User) IsClient() bool {
	return user.Role == Client
}

// IsVendor checks of the user is a vendor or not
func (user *User) IsVendor() bool {
	return user.Role == Vendor
}

// IsAdmin checks of the user is an admin or not
func (user *User) IsAdmin() bool {
	return user.Role == Admin
}

// SetSuccess defines the success of user creation
func (user *User) SetSuccess(success bool) {
	user.Success = success
}
