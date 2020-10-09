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

import "github.com/hslatman/balena-sdk-go/pkg/models"

type Device struct {
	fluent
	ID   int    `json:"id"`
	Name string `json:"device_name"`
	Type string `json:"device_type"`
	UUID string `json:"uuid"`

	// TODO: other fields
}

func (d Device) GetTags() (map[int]models.DeviceTag, error) {
	return d.client.DeviceTagsByDeviceID(d.ID)
}
