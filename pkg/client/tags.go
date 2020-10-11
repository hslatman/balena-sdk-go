package client

import (
	"fmt"

	"github.com/m7shapan/njson"
	"github.com/tidwall/gjson"
)

type DeviceTagsResource struct {
	client    *Client
	endpoint  string
	modifiers *ODataModifiers
}

type DeviceTagResource struct {
	client      *Client
	endpoint    string
	deviceTagID int
	modifiers   *ODataModifiers
}

type DeviceTag struct {
	// Tag
	ID       int    `njson:"id"`
	TagKey   string `njson:"tag_key"`
	Value    string `njson:"value"`
	DeviceID int    `njson:"device.__id"`
}

func NewDeviceTagsResource(c *Client, e string) *DeviceTagsResource {
	return &DeviceTagsResource{
		client:    c,
		endpoint:  e,
		modifiers: NewODataModifiers(),
	}
}

func NewDeviceTagResource(c *Client, e string, deviceTagID int) *DeviceTagResource {
	return &DeviceTagResource{
		client:      c,
		endpoint:    e,
		deviceTagID: deviceTagID,
		modifiers:   NewODataModifiers(),
	}
}

func (c *Client) DeviceTags() *DeviceTagsResource {
	return NewDeviceTagsResource(c, string(deviceTagsEndpoint))
}

func (c *Client) DeviceTag(deviceTagID int) *DeviceTagResource {
	endpoint := fmt.Sprintf("%s(%d)", deviceTagsEndpoint, deviceTagID)
	return NewDeviceTagResource(c, endpoint, deviceTagID)
}

func (d *DeviceTagsResource) Get() (map[int]DeviceTag, error) {

	tags := make(map[int]DeviceTag)

	//params[filterOption] = "device/id%20eq%20'" + strconv.Itoa(deviceID) + "'"

	url := createURL(d.endpoint, d.modifiers)

	params := make(map[paramOption]string)

	resp, err := d.client.get(url, params)

	if err != nil {
		return tags, err
	}

	// TODO: decide whether we need gjson dependency, or can do it easily, with a bit more wrangling, ourselves
	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results

	for _, tag := range data.Array() {
		t := DeviceTag{}
		if err := njson.Unmarshal([]byte(tag.Raw), &t); err != nil {
			return tags, err // TODO: don't do early return, but just skip this one and aggregate error?
		}
		tags[t.ID] = t
	}

	return tags, nil

}

func (d *DeviceTagResource) Get() (DeviceTag, error) {

	tag := DeviceTag{}

	url := createURL(d.endpoint, d.modifiers)

	params := make(map[paramOption]string)

	resp, err := d.client.get(url, params)

	if err != nil {
		return tag, err
	}

	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results
	first := data.Get("0")                   // first (and only) tag, if found

	if !first.Exists() {
		return tag, fmt.Errorf("%s not found", d.endpoint)
	}

	if err := njson.Unmarshal([]byte(first.Raw), &tag); err != nil {
		return tag, err
	}

	return tag, nil
}
