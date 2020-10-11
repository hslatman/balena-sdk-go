package client

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/tidwall/gjson"
)

type DevicesResource struct {
	client    *Client
	endpoint  string
	modifiers *ODataModifiers

	// TODO: context, configuration
}

type DeviceResource struct {
	client    *Client
	endpoint  string
	deviceID  int
	modifiers *ODataModifiers
}

type Device struct {
	ID   int    `json:"id"`
	Name string `json:"device_name"`
	Type string `json:"device_type"`
	UUID string `json:"uuid"`
}

func NewDevicesResource(c *Client, e string) *DevicesResource {
	return &DevicesResource{
		client:    c,
		endpoint:  e,
		modifiers: NewODataModifiers(),
	}
}

func NewDeviceResource(c *Client, e string, deviceID int) *DeviceResource {
	return &DeviceResource{
		client:    c,
		endpoint:  e,
		deviceID:  deviceID,
		modifiers: NewODataModifiers(),
	}
}

func (c *Client) Devices() *DevicesResource {
	return NewDevicesResource(c, string(devicesEndpoint))
}

func (c *Client) Device(deviceID int) *DeviceResource {
	endpoint := fmt.Sprintf("%s(%d)", devicesEndpoint, deviceID)
	return NewDeviceResource(c, endpoint, deviceID)
}

// func (d *Devices) Select(s string) *Devices {
// 	d.modifiers.AddSelect(s)
// 	return d
// }

func (d *DevicesResource) Filter(f string) *DevicesResource {
	d.modifiers.AddFilter(f)
	return d
}

// func (d *DevicesResource) FilterByID(id int) *DevicesResource {
// 	d.Filter("id%20eq%20'" + strconv.Itoa(id) + "'")
// 	return d
// }

func (d *DevicesResource) Get() (map[int]Device, error) {

	devices := make(map[int]Device)

	fmt.Println(d.modifiers)

	url := createURL(d.endpoint, d.modifiers)

	params := make(map[paramOption]string)

	resp, err := d.client.get(url, params)

	if err != nil {
		return devices, err
	}

	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results

	for _, device := range data.Array() {
		d := Device{}
		if err := json.Unmarshal([]byte(device.Raw), &d); err != nil {
			return devices, err // TODO: don't do early return, but just skip this one and aggregate error?
		}
		devices[d.ID] = d
	}

	return devices, nil
}

func (d *DevicesResource) GetByID(deviceID int) *DeviceResource {
	return NewDeviceResource(
		d.client,
		fmt.Sprintf("%s(%d)", devicesEndpoint, deviceID),
		deviceID,
	)
}

func (d *DeviceResource) Get() (Device, error) {

	device := Device{}

	fmt.Println(d.modifiers)

	url := createURL(d.endpoint, d.modifiers)

	params := make(map[paramOption]string)

	resp, err := d.client.get(url, params)

	if err != nil {
		return device, err
	}

	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results
	first := data.Get("0")                   // first (and only) device

	if !first.Exists() {
		return device, fmt.Errorf("%s not found", url)
	}

	if err := json.Unmarshal([]byte(first.Raw), &device); err != nil {
		return device, err
	}

	return device, nil
}

func (d *DeviceResource) Tags() *DeviceTagsResource {
	r := NewDeviceTagsResource(
		d.client,
		string(deviceTagsEndpoint),
	)
	// TODO: improve the formatting of this filter
	r.modifiers.AddFilter("device/id%20eq%20'" + strconv.Itoa(d.deviceID) + "'")
	return r
}

// TODO: split devices vs device resource
// TODO: add Tags resource; how to build it on the Device resource?
// TODO: add more elegant handling of createURL

// GetByID gets a user by his/her ID (numeric ID from User Information List)
// func (users *Users) GetByID(userID int) *User {
// 	return NewUser(
// 		users.client,
// 		fmt.Sprintf("%s/GetById(%d)", users.endpoint, userID),
// 		users.config,
// 	)
// }

// func (d Device) GetTags() (map[int]models.DeviceTag, error) {
// 	return d.client.DeviceTagsByDeviceID(d.ID)
// }
