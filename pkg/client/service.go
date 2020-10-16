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
	"fmt"

	"github.com/hslatman/balena-sdk-go/pkg/models"
	"github.com/m7shapan/njson"
	"github.com/tidwall/gjson"
)

type ServiceResource struct {
	client    *Client
	endpoint  string
	serviceID int
	modifiers *ODataModifiers
}

func NewServiceResource(c *Client, serviceID int) *ServiceResource {
	return &ServiceResource{
		client:    c,
		endpoint:  fmt.Sprintf("%s(%d)", servicesEndpoint, serviceID),
		serviceID: serviceID,
		modifiers: NewODataModifiers(c),
	}
}

func (c *Client) Service(serviceID int) *ServiceResource {
	return NewServiceResource(c, serviceID)
}

func (r *ServiceResource) Get() (models.Service, error) {

	service := models.Service{}

	resp, err := r.client.get(r.endpoint, r.modifiers)

	if err != nil {
		return service, err
	}

	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results
	first := data.Get("0")                   // first (and only) device

	if !first.Exists() {
		return service, fmt.Errorf("%s not found", r.endpoint)
	}

	if err := njson.Unmarshal([]byte(first.Raw), &service); err != nil {
		return service, err
	}

	return service, nil
}

func (r *ServiceResource) Select(s string) *ServiceResource {
	r.modifiers.AddSelect(s) // TODO: add validation that fields to be selected are valid fields for Device?
	return r
}
