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

type ReleaseTagResource struct {
	client       *Client
	endpoint     string
	releaseTagID int
	modifiers    *ODataModifiers
}

func NewReleaseTagResource(c *Client, releaseTagID int) *ReleaseTagResource {
	return &ReleaseTagResource{
		client:       c,
		endpoint:     fmt.Sprintf("%s(%d)", releaseTagsEndpoint, releaseTagID),
		releaseTagID: releaseTagID,
		modifiers:    NewODataModifiers(c),
	}
}

func (c *Client) ReleaseTag(releaseTagID int) *ReleaseTagResource {
	return NewReleaseTagResource(c, releaseTagID)
}

func (r *ReleaseTagResource) Get() (models.ReleaseTag, error) {

	tag := models.ReleaseTag{}

	resp, err := r.client.get(r.endpoint, r.modifiers)

	if err != nil {
		return tag, err
	}

	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results
	first := data.Get("0")                   // first (and only) tag, if found

	if !first.Exists() {
		return tag, fmt.Errorf("%s not found", r.endpoint)
	}

	if err := njson.Unmarshal([]byte(first.Raw), &tag); err != nil {
		return tag, err
	}

	return tag, nil
}

func (r *ReleaseTagResource) Select(s string) *ReleaseTagResource {
	r.modifiers.AddSelect(s)
	return r
}

func (r *ReleaseTagResource) Expand(s string) *ReleaseTagResource {
	r.modifiers.AddExpand(s)
	return r
}
