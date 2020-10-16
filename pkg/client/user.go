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

type UserResource struct {
	client    *Client
	endpoint  string
	userID    int
	modifiers *ODataModifiers
}

func NewUserResource(c *Client, userID int) *UserResource {
	return &UserResource{
		client:    c,
		endpoint:  fmt.Sprintf("%s(%d)", usersEndpoint, userID),
		userID:    userID,
		modifiers: NewODataModifiers(c),
	}
}

func (c *Client) User(userID int) *UserResource {
	return NewUserResource(c, userID)
}

func (r *UserResource) Get() (models.User, error) {

	user := models.User{}

	resp, err := r.client.get(r.endpoint, r.modifiers)

	if err != nil {
		return user, err
	}

	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results
	first := data.Get("0")                   // first (and only) device

	if !first.Exists() {
		return user, fmt.Errorf("%s not found", r.endpoint)
	}

	if err := njson.Unmarshal([]byte(first.Raw), &user); err != nil {
		return user, err
	}

	return user, nil
}

func (r *UserResource) Select(s string) *UserResource {
	r.modifiers.AddSelect(s) // TODO: add validation that fields to be selected are valid fields for Device?
	return r
}
