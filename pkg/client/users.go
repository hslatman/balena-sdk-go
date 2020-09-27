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

	"github.com/tidwall/gjson"

	"github.com/hslatman/balena-sdk-go/pkg/models"
)

func (c *Client) Users() (map[int]models.User, error) {

	users := make(map[int]models.User)

	params := make(map[paramOption]string)
	resp, err := c.get(string(usersEndpoint), params)

	if err != nil {
		return users, err
	}

	// TODO: decide whether we need gjson dependency, or can do it easily, with a bit more wrangling, ourselves
	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results

	for _, app := range data.Array() {
		u := models.User{}
		if err := json.Unmarshal([]byte(app.Raw), &u); err != nil {
			return users, err // TODO: don't do early return, but just skip this one and aggregate error?
		}
		users[u.ID] = u
	}

	return users, nil
}

func (c *Client) WhoAmI() (models.WhoAmI, error) {

	w := models.WhoAmI{}

	url := "https://api.balena-cloud.com/user/v1" // This endpoint uses a different format than all other v5 APIs
	params := make(map[paramOption]string)
	resp, err := c.get(url+string(whoamiEndpoint), params)

	if err != nil {
		return w, err
	}

	if err := json.Unmarshal(resp.Body(), &w); err != nil {
		return w, err
	}

	return w, nil
}

func (c *Client) UsersAssociatedWithApplication(id int) (map[int]models.User, error) {

	// TODO: doesn't work yet.

	users := make(map[int]models.User)

	// TODO: make this a bit nicer to work with? Essentially, it's how OData does filtering and such
	// TODO: this does return something, but it's an empty list; is that correct?
	params := make(map[paramOption]string)
	params[filterOption] = "is_member_of__application%20eq%20" + strconv.Itoa(id)
	params[expandOption] = "user($select=id,username,actor),application_membership_role($select=id,name,actor)"

	url := string(usersEndpoint) + "__is_member_of__application"
	resp, err := c.get(url, params)

	if err != nil {
		return users, err
	}

	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results

	for _, app := range data.Array() {
		u := models.User{}
		if err := json.Unmarshal([]byte(app.Raw), &u); err != nil {
			return users, err // TODO: don't do early return, but just skip this one and aggregate error?
		}
		users[u.ID] = u
	}

	return users, nil
}
