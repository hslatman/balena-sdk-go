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

type ReleaseResource struct {
	client    *Client
	endpoint  string
	releaseID int
	modifiers *ODataModifiers
}

func NewReleaseResource(c *Client, releaseID int) *ReleaseResource {
	return &ReleaseResource{
		client:    c,
		endpoint:  fmt.Sprintf("%s(%d)", releasesEndpoint, releaseID),
		releaseID: releaseID,
		modifiers: NewODataModifiers(c),
	}
}

func (c *Client) Release(releaseID int) *ReleaseResource {
	return NewReleaseResource(c, releaseID)
}

func (r *ReleaseResource) Get() (models.Release, error) {

	release := models.Release{}

	resp, err := r.client.get(r.endpoint, r.modifiers)

	if err != nil {
		return release, err
	}

	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results
	first := data.Get("0")                   // first (and only) device

	if !first.Exists() {
		return release, fmt.Errorf("%s not found", r.endpoint)
	}

	if err := njson.Unmarshal([]byte(first.Raw), &release); err != nil {
		return release, err
	}

	return release, nil
}

func (r *ReleaseResource) Select(s string) *ReleaseResource {
	r.modifiers.AddSelect(s) // TODO: add validation that fields to be selected are valid fields for Device?
	return r
}

func (r *ReleaseResource) Tags() *ReleaseTagsResource {
	tr := NewReleaseTagsResource(
		r.client,
	)
	// TODO: improve the formatting of this filter; Balena API seems to require this like this?
	tr.modifiers.AddFilter("release/id%20eq%20'" + strconv.Itoa(r.releaseID) + "'")
	return tr
}
