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

type UsersResource struct {
	client    *Client
	endpoint  string
	modifiers *ODataModifiers

	// TODO: context, configuration
}

func NewUsersResource(c *Client) *UsersResource {
	return &UsersResource{
		client:    c,
		endpoint:  string(usersEndpoint),
		modifiers: NewODataModifiers(c),
	}
}

func (c *Client) Users() *UsersResource {
	return NewUsersResource(c)
}

func (c *Client) WhoAmI() (models.WhoAmI, error) {

	w := models.WhoAmI{}

	url := "https://api.balena-cloud.com/user/v1" // This endpoint uses a different format than all other v5 APIs
	resp, err := c.get(url+string(whoamiEndpoint), nil)

	if err != nil {
		return w, err
	}

	if err := njson.Unmarshal(resp.Body(), &w); err != nil {
		return w, err
	}

	return w, nil
}

func (r *UsersResource) Select(s string) *UsersResource {
	r.modifiers.AddSelect(s)
	return r
}

func (r *UsersResource) Filter(f string) *UsersResource {
	r.modifiers.AddFilter(f)
	return r
}

func (r *UsersResource) Get() (map[int]models.User, error) {

	users := make(map[int]models.User)

	resp, err := r.client.get(r.endpoint, r.modifiers)

	if err != nil {
		return users, err
	}

	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results

	for _, u := range data.Array() {
		user := models.User{}
		if err := njson.Unmarshal([]byte(u.Raw), &user); err != nil {
			return users, err // TODO: don't do early return, but just skip this one and aggregate error?
		}
		users[user.ID] = user
	}

	return users, nil
}

func (r *UsersResource) FindByID(userID int) *UserResource {
	return NewUserResource(
		r.client,
		userID,
	)
}
