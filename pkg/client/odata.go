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
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/hslatman/balena-sdk-go/pkg/logger"
)

type ODataModifiers struct {
	modifiers map[modifierOption]string
	logger    logger.Logger
	errors    []error
}

func NewODataModifiers(c *Client) *ODataModifiers {
	return &ODataModifiers{
		modifiers: map[modifierOption]string{},
		logger:    c.logger,
		errors:    []error{},
	}
}

func (o *ODataModifiers) Get() map[modifierOption]string {
	return o.modifiers
}

func (o *ODataModifiers) AddSelect(value string) *ODataModifiers {

	// TODO: vararg approach?

	if _, ok := o.modifiers[selectOption]; ok {
		o.errors = append(o.errors, errors.New("multiple $select parts currently not supported"))
		return o
	}

	o.modifiers[selectOption] = value

	return o
}

func (o *ODataModifiers) AddFilter(value string) *ODataModifiers {

	// TODO: vararg approach?

	if _, ok := o.modifiers[filterOption]; ok {
		o.errors = append(o.errors, errors.New("multiple $select parts currently not supported"))
		return o
	}

	o.modifiers[filterOption] = value

	return o
}

func (o *ODataModifiers) AddExpand(value string) *ODataModifiers {

	// TODO: vararg approach?

	if _, ok := o.modifiers[expandOption]; ok {
		o.errors = append(o.errors, errors.New("multiple $select parts currently not supported"))
		return o
	}

	o.modifiers[expandOption] = value

	return o
}

func (o *ODataModifiers) modifyURL(endpoint string) (string, error) {

	if len(o.errors) > 0 {
		err := "Compound error:"
		for k, v := range o.errors {
			err = strings.Join([]string{err, fmt.Sprintf("(%d) %s", k, v.Error())}, "\n")
		}
		return "", errors.New(err)
	}

	apiURL, _ := url.Parse(endpoint)

	mods := o.Get()

	if len(mods) == 0 {
		return apiURL.String(), nil
	}

	queryParts := []string{}

	for k, v := range mods {
		if k == filterOption {
			part := "$filter=" + trimMultiline(v)
			queryParts = append(queryParts, part)
		}
		if k == selectOption {
			part := "$select=" + trimMultiline(v)
			queryParts = append(queryParts, part)
		}
		if k == expandOption {
			part := "$expand=" + trimMultiline(v)
			queryParts = append(queryParts, part)
		}
	}

	query := "?" + strings.Join(queryParts, "&")

	return apiURL.String() + query, nil

}

func trimMultiline(multi string) string {
	res := ""
	for _, line := range strings.Split(multi, "\n") {
		res += strings.Trim(line, "\t")
	}
	return res
}
