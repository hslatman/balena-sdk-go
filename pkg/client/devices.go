// Copyright 2020 Herman Slatman
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
)

func (c *Client) Devices() (map[int]Device, error) {

	devices := make(map[int]Device)

	params := make(map[paramOption]string)
	resp, err := c.get(string(devicesEndpoint), params)

	if err != nil {
		return devices, err
	}

	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results

	for _, device := range data.Array() {
		d := Device{}
		if err := json.Unmarshal([]byte(device.Raw), &d); err != nil {
			return devices, err // TODO: don't do early return, but just skip this one and aggregate error?
		}
		d.client = c
		devices[d.ID] = d
	}

	return devices, nil
}

func (c *Client) Device(id int) (Device, error) {

	// TODO: also lookup by other identifiers, like UUID, device name, device type, etc

	device := Device{}

	params := make(map[paramOption]string)
	resp, err := c.get(fmt.Sprintf("%s(%d)", devicesEndpoint, id), params)

	if err != nil {
		return device, err
	}

	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results
	first := data.Get("0")                   // first (and only) device

	if !first.Exists() {
		return device, fmt.Errorf("device (@id=%d) not found", id)
	}

	if err := json.Unmarshal([]byte(first.Raw), &device); err != nil {
		return device, err
	}

	device.client = c

	return device, nil
}
