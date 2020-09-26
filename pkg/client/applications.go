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

	"github.com/go-resty/resty/v2"
	"github.com/hslatman/balena-sdk-go/pkg/models"
	"github.com/tidwall/gjson"
)

func (c *Client) Applications() (map[int]models.Application, error) {

	apps := make(map[int]models.Application)

	resp, err := c.send(resty.MethodGet, string(applicationsEndpoint))

	if err != nil {
		return apps, err
	}

	// TODO: decide whether we need gjson dependency, or can do it easily, with a bit more wrangling, ourselves
	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results

	for _, app := range data.Array() {
		a := models.Application{}
		if err := json.Unmarshal([]byte(app.Raw), &a); err != nil {
			return apps, err // TODO: don't do early return, but just skip this one and aggregate error?
		}
		apps[a.ID] = a
	}

	return apps, nil
}

func (c *Client) Application(id int) (models.Application, error) {

	// TODO: also lookup by other identifiers, like UUID, device name, device type, etc

	app := models.Application{}

	resp, err := c.send(resty.MethodGet, fmt.Sprintf("%s(%d)", applicationsEndpoint, id))

	if err != nil {
		return app, err
	}

	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results
	first := data.Get("0")                   // first (and only) application, if found

	if !first.Exists() {
		return app, fmt.Errorf("application (@id=%d) not found", id)
	}

	if err := json.Unmarshal([]byte(first.Raw), &app); err != nil {
		return app, err
	}

	return app, nil
}

func (c *Client) AllApplications() (map[int]models.Application, error) {

	apps := make(map[int]models.Application)

	resp, err := c.send(resty.MethodGet, string(allApplicationsEndpoint))

	if err != nil {
		return apps, err
	}

	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results

	for _, app := range data.Array() {
		a := models.Application{}
		if err := json.Unmarshal([]byte(app.Raw), &a); err != nil {
			return apps, err // TODO: don't do early return, but just skip this one and aggregate error?
		}
		apps[a.ID] = a
	}

	return apps, nil
}
