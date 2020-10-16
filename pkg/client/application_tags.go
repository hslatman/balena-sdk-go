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

type ApplicationTagsResource struct {
	client    *Client
	endpoint  string
	modifiers *ODataModifiers
}

func NewApplicationTagsResource(c *Client) *ApplicationTagsResource {
	return &ApplicationTagsResource{
		client:    c,
		endpoint:  string(applicationTagsEndpoint),
		modifiers: NewODataModifiers(c),
	}
}

func (c *Client) ApplicationTags() *ApplicationTagsResource {
	return NewApplicationTagsResource(c)
}

func (r *ApplicationTagsResource) FindByID(applicationTagID int) *ApplicationTagResource {
	return NewApplicationTagResource(
		r.client,
		applicationTagID,
	)
}

func (r *ApplicationTagsResource) Get() (map[int]models.ApplicationTag, error) {

	tags := make(map[int]models.ApplicationTag)

	resp, err := r.client.get(r.endpoint, r.modifiers)

	if err != nil {
		return tags, err
	}

	// TODO: decide whether we need gjson dependency, or can do it easily, with a bit more wrangling, ourselves
	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results

	for _, t := range data.Array() {
		tag := models.ApplicationTag{}
		if err := njson.Unmarshal([]byte(t.Raw), &tag); err != nil {
			return tags, err // TODO: don't do early return, but just skip this one and aggregate error?
		}
		tags[tag.ID] = tag
	}

	return tags, nil

}

func (r *ApplicationTagsResource) Select(s string) *ApplicationTagsResource {
	r.modifiers.AddSelect(s)
	return r
}

func (r *ApplicationTagsResource) Filter(s string) *ApplicationTagsResource {
	r.modifiers.AddFilter(s)
	return r
}
