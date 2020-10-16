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
	"strconv"

	"github.com/hslatman/balena-sdk-go/pkg/models"
	"github.com/m7shapan/njson"
	"github.com/tidwall/gjson"
)

type ApplicationResource struct {
	client        *Client
	endpoint      string
	applicationID int
	modifiers     *ODataModifiers
}

func NewApplicationResource(c *Client, applicationID int) *ApplicationResource {
	return &ApplicationResource{
		client:        c,
		endpoint:      fmt.Sprintf("%s(%d)", applicationsEndpoint, applicationID),
		applicationID: applicationID,
		modifiers:     NewODataModifiers(c),
	}
}

func (c *Client) Application(applicationID int) *ApplicationResource {
	return NewApplicationResource(c, applicationID)
}

func (r *ApplicationResource) Get() (models.Application, error) {

	application := models.Application{}

	resp, err := r.client.get(r.endpoint, r.modifiers)

	if err != nil {
		return application, err
	}

	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results
	first := data.Get("0")                   // first (and only) device

	if !first.Exists() {
		return application, fmt.Errorf("%s not found", r.endpoint)
	}

	if err := njson.Unmarshal([]byte(first.Raw), &application); err != nil {
		return application, err
	}

	return application, nil
}

func (r *ApplicationResource) Select(s string) *ApplicationResource {
	r.modifiers.AddSelect(s) // TODO: add validation that fields to be selected are valid fields for Device?
	return r
}

func (r *ApplicationResource) Tags() *ApplicationTagsResource {
	tr := NewApplicationTagsResource(
		r.client,
	)
	// TODO: improve the formatting of this filter; Balena API seems to require this like this?
	tr.modifiers.AddFilter("application/id%20eq%20'" + strconv.Itoa(r.applicationID) + "'")
	return tr
}

func (r *ApplicationResource) Releases() *ReleasesResource {
	rr := NewReleasesResource(
		r.client,
	)
	filter := "belongs_to__application%20eq%20'" + strconv.Itoa(r.applicationID) + "'"
	rr.modifiers.AddFilter(filter)
	return rr
}

func (r *ApplicationResource) Services() *ServicesResource {
	sr := NewServicesResource(
		r.client,
	)
	filter := "application/id%20eq%20'" + strconv.Itoa(r.applicationID) + "'"
	sr.modifiers.AddFilter(filter)
	return sr
}

func (r *ApplicationResource) AddTag(key string, value string) (models.ApplicationTag, error) {

	ar := NewApplicationTagsResource(
		r.client,
	)

	return ar.Create(r.applicationID, key, value)
}

func (r *ApplicationResource) UpdateTag(key string, value string) error {

	ar := NewApplicationTagsResource(
		r.client,
	)

	ar.modifiers.AddFilter("application/id%20eq%20'" + strconv.Itoa(r.applicationID) + "'%20and%20tag_key%20eq%20'" + key + "'")

	return ar.Update(value)
}

func (r *ApplicationResource) DeleteTag(key string) error {

	ar := NewApplicationTagsResource(
		r.client,
	)

	ar.modifiers.AddFilter("application/id%20eq%20'" + strconv.Itoa(r.applicationID) + "'%20and%20tag_key%20eq%20'" + key + "'")

	return ar.Delete()
}
