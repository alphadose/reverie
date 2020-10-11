package types

// M is a shorthand notation for map[string]interface{}
type M map[string]interface{}

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
