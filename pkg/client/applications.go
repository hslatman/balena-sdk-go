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

type ApplicationsResource struct {
	client    *Client
	endpoint  string
	modifiers *ODataModifiers

	// TODO: context, configuration
}

func NewApplicationsResource(c *Client) *ApplicationsResource {
	return &ApplicationsResource{
		client:    c,
		endpoint:  string(applicationsEndpoint),
		modifiers: NewODataModifiers(c),
	}
}

func (c *Client) Applications() *ApplicationsResource {
	return NewApplicationsResource(c)
}

func (r *ApplicationsResource) Select(s string) *ApplicationsResource {
	r.modifiers.AddSelect(s)
	return r
}

func (r *ApplicationsResource) Filter(f string) *ApplicationsResource {
	r.modifiers.AddFilter(f)
	return r
}

func (r *ApplicationsResource) All() *ApplicationsResource {
	r.endpoint = string(allApplicationsEndpoint)
	return r
}

func (r *ApplicationsResource) Get() (map[int]models.Application, error) {

	applications := make(map[int]models.Application)

	resp, err := r.client.get(r.endpoint, r.modifiers)

	if err != nil {
		return applications, err
	}

	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results

	for _, a := range data.Array() {
		application := models.Application{}
		if err := njson.Unmarshal([]byte(a.Raw), &application); err != nil {
			return applications, err // TODO: don't do early return, but just skip this one and aggregate error?
		}
		applications[application.ID] = application
	}

	return applications, nil
}

func (r *ApplicationsResource) FindByID(applicationID int) *ApplicationResource {
	return NewApplicationResource(
		r.client,
		applicationID,
	)
}
