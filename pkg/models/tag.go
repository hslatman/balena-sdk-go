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

// type Tag struct {
// 	ID     int    `njson:"id"`
// 	TagKey string `njson:"tag_key"`
// 	Value  string `njson:"value"`
// }

type ApplicationTag struct {
	//Tag // NOTE: njson doesn't seem to work with embedded structs
	// TODO: njson also doesn't seem to support omitempty and such
	ID            int    `njson:"id"`
	TagKey        string `njson:"tag_key"`
	Value         string `njson:"value"`
	ApplicationID int    `njson:"application.__id"`
}

type DeviceTag struct {
	// Tag
	ID       int    `njson:"id"`
	TagKey   string `njson:"tag_key"`
	Value    string `njson:"value"`
	DeviceID int    `njson:"device.__id"`
}

type ReleaseTag struct {
	// Tag
	ID        int    `njson:"id"`
	TagKey    string `njson:"tag_key"`
	Value     string `njson:"value"`
	ReleaseID int    `njson:"release.__id"`
}
