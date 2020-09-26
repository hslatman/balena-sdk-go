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

func (c *Client) ApplicationTagsByApplicationName(applicationName string) (map[int]models.ApplicationTag, error) {

	tags := make(map[int]models.ApplicationTag)

	//"https://api.balena-cloud.com/v5/application_tag?\$filter=application/app_name%20eq%20'<NAME>'" \

	url := string(applicationTagsEndpoint) + "?" + "$filter=application/app_name%20eq%20'" + applicationName + "'"
	resp, err := c.send(resty.MethodGet, url)

	if err != nil {
		return tags, err
	}

	// TODO: decide whether we need gjson dependency, or can do it easily, with a bit more wrangling, ourselves
	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results

	for _, tag := range data.Array() {
		t := models.ApplicationTag{}
		if err := json.Unmarshal([]byte(tag.Raw), &t); err != nil {
			return tags, err // TODO: don't do early return, but just skip this one and aggregate error?
		}
		tags[t.ID] = t
	}

	return tags, nil
}

func (c *Client) DeviceTagsByDeviceUUID(deviceUUID string) (map[int]models.DeviceTag, error) {

	tags := make(map[int]models.DeviceTag)

	url := string(deviceTagsEndpoint) + "?" + "$filter=device/uuid%20eq%20'" + deviceUUID + "'"
	resp, err := c.send(resty.MethodGet, url)

	if err != nil {
		return tags, err
	}

	// TODO: decide whether we need gjson dependency, or can do it easily, with a bit more wrangling, ourselves
	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results

	for _, tag := range data.Array() {
		t := models.DeviceTag{}
		if err := json.Unmarshal([]byte(tag.Raw), &t); err != nil {
			return tags, err // TODO: don't do early return, but just skip this one and aggregate error?
		}
		tags[t.ID] = t
	}

	return tags, nil
}

func (c *Client) ReleaseTagsByReleaseCommit(commit string) (map[int]models.ReleaseTag, error) {

	tags := make(map[int]models.ReleaseTag)

	url := string(releaseTagsEndpoint) + "?" + "$filter=release/commit%20eq%20'" + commit + "'"
	resp, err := c.send(resty.MethodGet, url)

	if err != nil {
		return tags, err
	}

	// TODO: decide whether we need gjson dependency, or can do it easily, with a bit more wrangling, ourselves
	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results

	for _, tag := range data.Array() {
		t := models.ReleaseTag{}
		if err := json.Unmarshal([]byte(tag.Raw), &t); err != nil {
			return tags, err // TODO: don't do early return, but just skip this one and aggregate error?
		}
		tags[t.ID] = t
	}

	return tags, nil
}

func (c *Client) ReleaseTagsByReleaseID(id int) (map[int]models.ReleaseTag, error) {

	tags := make(map[int]models.ReleaseTag)

	url := string(releaseTagsEndpoint) + "?" + "$filter=release/id%20eq%20'" + strconv.Itoa(id) + "'"
	resp, err := c.send(resty.MethodGet, url)

	if err != nil {
		return tags, err
	}

	// TODO: decide whether we need gjson dependency, or can do it easily, with a bit more wrangling, ourselves
	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results

	for _, tag := range data.Array() {
		t := models.ReleaseTag{}
		if err := json.Unmarshal([]byte(tag.Raw), &t); err != nil {
			return tags, err // TODO: don't do early return, but just skip this one and aggregate error?
		}
		tags[t.ID] = t
	}

	return tags, nil
}
