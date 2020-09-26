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

package models

type Tag struct {
	ID     int    `json:"id"`
	TagKey string `json:"tag_key"`
	Value  string `json:"value"`
}

type ApplicationTag struct {
	Tag
	//ApplicationID int `json:"application.__id"`// TODO: can we get this easily?
}

type DeviceTag struct {
	Tag
	// TODO: device?
}

type ReleaseTag struct {
	Tag
	// TODO: release?
}
