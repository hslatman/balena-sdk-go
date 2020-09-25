package models

type Application struct {
	ID         int    `json:"id"`
	Name       string `json:"app_name"`
	Slug       string `json:"slug"`
	DeviceType string `json:"device_type"`
	IsPublic   bool   `json:"is_public"`
	IsHost     bool   `json:"is_host"`
	IsArchived bool   `json:"is_archived"`
}
