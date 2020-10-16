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
	"github.com/hslatman/balena-sdk-go/pkg/models"
	"github.com/m7shapan/njson"
	"github.com/tidwall/gjson"
)

type ServicesResource struct {
	client    *Client
	endpoint  string
	modifiers *ODataModifiers

	// TODO: context, configuration
}

func NewServicesResource(c *Client) *ServicesResource {
	return &ServicesResource{
		client:    c,
		endpoint:  string(servicesEndpoint),
		modifiers: NewODataModifiers(c),
	}
}

func (c *Client) Services() *ServicesResource {
	return NewServicesResource(c)
}

func (r *ServicesResource) Select(s string) *ServicesResource {
	r.modifiers.AddSelect(s)
	return r
}

func (r *ServicesResource) Filter(f string) *ServicesResource {
	r.modifiers.AddFilter(f)
	return r
}

func (r *ServicesResource) Get() (map[int]models.Service, error) {

	services := make(map[int]models.Service)

	resp, err := r.client.get(r.endpoint, r.modifiers)

	if err != nil {
		return services, err
	}

	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results

	for _, s := range data.Array() {
		service := models.Service{}
		if err := njson.Unmarshal([]byte(s.Raw), &service); err != nil {
			return services, err // TODO: don't do early return, but just skip this one and aggregate error?
		}
		services[service.ID] = service
	}

	return services, nil
}

func (r *ServicesResource) FindByID(serviceID int) *ServiceResource {
	return NewServiceResource(
		r.client,
		serviceID,
	)
}
