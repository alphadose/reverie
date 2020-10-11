package types

// M is a shorthand notation for map[string]interface{}
type M map[string]interface{}

// Inventory stores the items in a vendor's inventory
type Inventory struct {
	Truck         int `json:"Truck,omitempty" bson:"Truck,omitempty"`
	Crane         int `json:"Crane,omitempty" bson:"Crane,omitempty"`
	Tanker        int `json:"Tanker,omitempty" bson:"Tanker,omitempty"`
	RoadRoller    int `json:"RoadRoller,omitempty" bson:"RoadRoller,omitempty"`
	ForkLift      int `json:"ForkLift,omitempty" bson:"ForkLift,omitempty"`
	BoomLifter    int `json:"BoomLifter,omitempty" bson:"BoomLifter,omitempty"`
	ManLifter     int `json:"ManLifter,omitempty" bson:"ManLifter,omitempty"`
	HydraulicJack int `json:"HydraulicJack,omitempty" bson:"HydraulicJack,omitempty"`
	Manpower      int `json:"Manpower,omitempty" bson:"Manpower,omitempty"`
}
