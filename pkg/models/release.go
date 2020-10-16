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

import "time"

type Release struct {
	ID              int       `njson:"id"`
	Commit          string    `njson:"commit"`
	Status          string    `njson:"status"`
	Source          string    `njson:"source"`
	CreatedAt       time.Time `njson:"created_at"`
	StartTimestamp  time.Time `njson:"start_timestamp"`
	EndTimestamp    time.Time `njson:"end_timestamp"`
	UpdateTimestamp time.Time `njson:"update_timestamp"`
}

// id
// created_at
// belongs_to__application
// is_created_by__user
// commit
// composition
// status
// source
// start_timestamp
// end_timestamp
// update_timestamp
