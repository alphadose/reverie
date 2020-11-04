package types

// M is a shorthand notation for map[string]interface{}
type M map[string]interface{}

const (
	// Success holds the response key for "success" field
	Success = "success"

	// Error holds the response key for "error" field
	Error = "error"

	// EmailKey is the key holding the email of a user
	EmailKey = "email"

	// UsernameKey is the key holding the username of a user
	UsernameKey = "username"

	// RoleKey is the key denoting the role of a user
	RoleKey = "role"

	// EMPTY denotes the empty string
	EMPTY = ""
)

// Inventory stores the items in a vendor's inventory
type Inventory struct {
	Truck         int64 `json:"Truck,omitempty" bson:"Truck,omitempty"`
	Crane         int64 `json:"Crane,omitempty" bson:"Crane,omitempty"`
	Tanker        int64 `json:"Tanker,omitempty" bson:"Tanker,omitempty"`
	RoadRoller    int64 `json:"RoadRoller,omitempty" bson:"RoadRoller,omitempty"`
	ForkLift      int64 `json:"ForkLift,omitempty" bson:"ForkLift,omitempty"`
	BoomLifter    int64 `json:"BoomLifter,omitempty" bson:"BoomLifter,omitempty"`
	ManLifter     int64 `json:"ManLifter,omitempty" bson:"ManLifter,omitempty"`
	HydraulicJack int64 `json:"HydraulicJack,omitempty" bson:"HydraulicJack,omitempty"`
	Manpower      int64 `json:"Manpower,omitempty" bson:"Manpower,omitempty"`
}
