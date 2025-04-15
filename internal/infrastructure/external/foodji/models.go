package foodji

import "github.com/google/uuid"

type Response struct {
	Data Machine `json:"data"`
}

type Machine struct {
	MachineProducts []Product `json:"machineProducts"`
}

// Product represents a product in the machine
type Product struct {
	ID uuid.UUID `json:"id"`
	// Including only ID as per the sample, but can be expanded with more fields as needed
}
