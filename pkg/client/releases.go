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
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hslatman/balena-sdk-go/pkg/models"
	"github.com/tidwall/gjson"
)

func (c *Client) Releases() (map[int]models.Release, error) {

	releases := make(map[int]models.Release)

	resp, err := c.send(resty.MethodGet, string(releasesEndpoint))

	if err != nil {
		return releases, err
	}

	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results

	for _, release := range data.Array() {
		r := models.Release{}
		if err := json.Unmarshal([]byte(release.Raw), &r); err != nil {
			return releases, err // TODO: don't do early return, but just skip this one and aggregate error?
		}
		releases[r.ID] = r
	}

	return releases, nil
}

func (c *Client) Release(id int) (models.Release, error) {

	// TODO: also lookup by other identifiers, like UUID, device name, device type, etc

	release := models.Release{}

	resp, err := c.send(resty.MethodGet, fmt.Sprintf("%s(%d)", releasesEndpoint, id))

	if err != nil {
		return release, err
	}

	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results
	first := data.Get("0")                   // first (and only) device

	if !first.Exists() {
		return release, fmt.Errorf("device (@id=%d) not found", id)
	}

	if err := json.Unmarshal([]byte(first.Raw), &release); err != nil {
		return release, err
	}

	return release, nil
}

func (c *Client) ReleasesForApplication(id int) (map[int]models.Release, error) {

	releases := make(map[int]models.Release)

	// TODO: make this a bit nicer to work with? Essentially, it's how OData does filtering and such
	params := make(map[paramOption]string)
	params[filterOption] = "belongs_to__application%20eq%20" + strconv.Itoa(id)
	resp, err := c.request(resty.MethodGet, string(releasesEndpoint), params)

	if err != nil {
		return releases, err
	}

	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results

	for _, release := range data.Array() {
		r := models.Release{}
		if err := json.Unmarshal([]byte(release.Raw), &r); err != nil {
			return releases, err // TODO: don't do early return, but just skip this one and aggregate error?
		}
		releases[r.ID] = r
	}

	return releases, nil
}
