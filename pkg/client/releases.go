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

type ReleasesResource struct {
	client    *Client
	endpoint  string
	modifiers *ODataModifiers

	// TODO: context, configuration
}

func NewReleasesResource(c *Client) *ReleasesResource {
	return &ReleasesResource{
		client:    c,
		endpoint:  string(releasesEndpoint),
		modifiers: NewODataModifiers(c),
	}
}

func (c *Client) Releases() *ReleasesResource {
	return NewReleasesResource(c)
}

func (r *ReleasesResource) Select(s string) *ReleasesResource {
	r.modifiers.AddSelect(s)
	return r
}

func (r *ReleasesResource) Filter(f string) *ReleasesResource {
	r.modifiers.AddFilter(f)
	return r
}

func (r *ReleasesResource) Get() (map[int]models.Release, error) {

	releases := make(map[int]models.Release)

	resp, err := r.client.get(r.endpoint, r.modifiers)

	if err != nil {
		return releases, err
	}

	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results

	for _, r := range data.Array() {
		release := models.Release{}
		if err := njson.Unmarshal([]byte(r.Raw), &release); err != nil {
			return releases, err // TODO: don't do early return, but just skip this one and aggregate error?
		}
		releases[release.ID] = release
	}

	return releases, nil
}

func (r *ReleasesResource) FindByID(releaseID int) *ReleaseResource {
	return NewReleaseResource(
		r.client,
		releaseID,
	)
}
