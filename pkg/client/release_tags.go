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

type ReleaseTagsResource struct {
	client    *Client
	endpoint  string
	modifiers *ODataModifiers
}

func NewReleaseTagsResource(c *Client) *ReleaseTagsResource {
	return &ReleaseTagsResource{
		client:    c,
		endpoint:  string(releaseTagsEndpoint),
		modifiers: NewODataModifiers(c),
	}
}

func (c *Client) ReleaseTags() *ReleaseTagsResource {
	return NewReleaseTagsResource(c)
}

func (r *ReleaseTagsResource) FindByID(releaseTagID int) *ReleaseTagResource {
	return NewReleaseTagResource(
		r.client,
		releaseTagID,
	)
}

func (r *ReleaseTagsResource) Get() (map[int]models.ReleaseTag, error) {

	tags := make(map[int]models.ReleaseTag)

	resp, err := r.client.get(r.endpoint, r.modifiers)

	if err != nil {
		return tags, err
	}

	// TODO: decide whether we need gjson dependency, or can do it easily, with a bit more wrangling, ourselves
	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results

	for _, t := range data.Array() {
		tag := models.ReleaseTag{}
		if err := njson.Unmarshal([]byte(t.Raw), &tag); err != nil {
			return tags, err // TODO: don't do early return, but just skip this one and aggregate error?
		}
		tags[tag.ID] = tag
	}

	return tags, nil

}

func (r *ReleaseTagsResource) Create(releaseID int, key string, value string) (models.ReleaseTag, error) {

	tag := models.ReleaseTag{}

	body := map[string]interface{}{
		"release": releaseID,
		"tag_key": key,
		"value":   value,
	}

	resp, err := r.client.post(r.endpoint, r.modifiers, body)
	if err != nil {
		return tag, err
	}

	if err := njson.Unmarshal(resp.Body(), &tag); err != nil {
		return tag, err
	}

	return tag, nil
}

func (r *ReleaseTagsResource) Update(value string) error {

	// TODO: should modifiers (also) be handled here?

	body := map[string]interface{}{
		"value": value,
	}

	_, err := r.client.patch(r.endpoint, r.modifiers, body)
	if err != nil {
		return err
	}

	return nil
}

func (r *ReleaseTagsResource) Delete() error {

	// TODO: should modifiers (also) be handled here?

	_, err := r.client.delete(r.endpoint, r.modifiers)
	if err != nil {
		return err
	}

	return nil
}

func (r *ReleaseTagsResource) Select(s string) *ReleaseTagsResource {
	r.modifiers.AddSelect(s)
	return r
}

func (r *ReleaseTagsResource) Filter(s string) *ReleaseTagsResource {
	r.modifiers.AddFilter(s)
	return r
}
