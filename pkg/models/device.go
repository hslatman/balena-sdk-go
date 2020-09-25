package models

type Device struct {
	ID   int    `json:"id"`
	Name string `json:"device_name"`
	Type string `json:"device_type"`
	UUID string `json:"uuid"`

	// TODO: other fields
}
