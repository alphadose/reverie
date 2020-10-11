package types

const (
	// Client holds the name  for the client role
	Client = "client"

	// Vendor holds the name  for the vendor role
	Vendor = "vendor"

	// Admin holds the name  for the admin role
	Admin = "admin"
)

// Inventory stores the items in a vendor's inventory
type Inventory struct {
	Truck         string `json:"Truck,omitempty" bson:"Truck,omitempty"`
	Crane         string `json:"Crane,omitempty" bson:"Crane,omitempty"`
	Tanker        string `json:"Tanker,omitempty" bson:"Tanker,omitempty"`
	RoadRoller    string `json:"RoadRoller,omitempty" bson:"RoadRoller,omitempty"`
	ForkLift      string `json:"ForkLift,omitempty" bson:"ForkLift,omitempty"`
	BoomLifter    string `json:"BoomLifter,omitempty" bson:"BoomLifter,omitempty"`
	ManLifter     string `json:"ManLifter,omitempty" bson:"ManLifter,omitempty"`
	HydraulicJack string `json:"HydraulicJack,omitempty" bson:"HydraulicJack,omitempty"`
	Manpower      string `json:"Manpower,omitempty" bson:"Manpower,omitempty"`
}

// User stores user related information
type User struct {
	Email     string    `form:"email" json:"email" binding:"required" bson:"email" valid:"required~Field 'email' is required but was not provided,email"`
	Password  string    `form:"password" json:"password,omitempty" bson:"password" binding:"required" valid:"required~Field 'password' is required but was not provided"`
	Username  string    `form:"username" json:"username" bson:"username" binding:"required" valid:"required~Field 'username' is required but was not provided,alphanum~Field 'username' should only have alphanumeric characters,stringlength(5|40)~Field 'username' should have length between 5 to 40 characters"`
	Role      string    `form:"role" json:"role" bson:"role" valid:"required~Field 'role' is required but was not provided,in(client|vendor)~Field 'Role' should be either client or vendor"`
	Inventory Inventory `form:"inventory" json:"inventory,omitempty" bson:"inventory,omitempty"`
	Success   bool      `json:"success,omitempty" bson:"-"`
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
