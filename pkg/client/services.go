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
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hslatman/balena-sdk-go/pkg/models"
	"github.com/tidwall/gjson"
)

func (c *Client) Services() (map[int]models.Service, error) {

	services := make(map[int]models.Service)

	resp, err := c.send(resty.MethodGet, string(servicesEndpoint))

	if err != nil {
		return services, err
	}

	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results

	for _, service := range data.Array() {
		s := models.Service{}
		if err := json.Unmarshal([]byte(service.Raw), &s); err != nil {
			return services, err // TODO: don't do early return, but just skip this one and aggregate error?
		}
		services[s.ID] = s
	}

	return services, nil
}

func (c *Client) ServicesForApplication(id int) (map[int]models.Service, error) {

	services := make(map[int]models.Service)

	// TODO: make this a bit nicer to work with? Essentially, it's how OData does filtering and such
	params := make(map[paramOption]string)
	params[filterOption] = "application/id%20eq%20" + strconv.Itoa(id)
	resp, err := c.request(resty.MethodGet, string(servicesEndpoint), params)

	if err != nil {
		return services, err
	}

	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results

	for _, service := range data.Array() {
		s := models.Service{}
		if err := json.Unmarshal([]byte(service.Raw), &s); err != nil {
			return services, err // TODO: don't do early return, but just skip this one and aggregate error?
		}
		services[s.ID] = s
	}

	return services, nil
}
