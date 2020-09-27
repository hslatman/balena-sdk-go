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
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/hslatman/balena-sdk-go/pkg/logger"
)

const (
	version = "0.1.0"

	defaultBaseURL           = "https://api.balena-cloud.com/v5"
	defaultUserAgent         = "https://github.com/hslatman/balena-sdk-go"
	defaultContentTypeHeader = "application/json"
	defaultAcceptHeader      = "application/json"
	defaultTimeOut           = 30 * time.Second
)

type endpoint string

const (
	allApplicationsEndpoint endpoint = "/application"
	applicationsEndpoint    endpoint = "/my_application"
	devicesEndpoint         endpoint = "/device"
	usersEndpoint           endpoint = "/user"
	whoamiEndpoint          endpoint = "/whoami"
	releasesEndpoint        endpoint = "/release"
	applicationTagsEndpoint endpoint = "/application_tag"
	deviceTagsEndpoint      endpoint = "/device_tag"
	releaseTagsEndpoint     endpoint = "/release_tag"
	servicesEndpoint        endpoint = "/service"
)

type paramOption string

const (
	filterOption  paramOption = "filter"
	filtersOption paramOption = "filters"
	expandOption  paramOption = "expand"
	selectOption  paramOption = "select"
)

type ClientModifier func(c *Client)

type Client struct {
	rc           *resty.Client
	modifiers    []ClientModifier
	logger       logger.Logger
	debugEnabled bool
	traceEnabled bool
}

func New(token string, modifiers ...ClientModifier) (*Client, error) {

	// TODO: add functionality for authenticating with Balena?
	// TODO: add additional configuration options / modifiers?
	// TODO: default retries? support for proxy? TLS settings? other transports?

	// Creating a new Resty client with defaults for all requests
	rc := resty.New()

	rc.SetHostURL(defaultBaseURL)
	rc.SetHeader("User-Agent", defaultUserAgent) // TODO: add version?
	rc.SetHeader("Content-Type", defaultContentTypeHeader)
	rc.SetHeader("Accept", defaultAcceptHeader)

	rc.SetAuthScheme("Bearer")
	rc.SetAuthToken(token)

	rc.SetTimeout(defaultTimeOut)
	rc.SetDebug(false)

	c := &Client{
		rc:           rc,
		modifiers:    []ClientModifier{},
		debugEnabled: false,
		traceEnabled: false,
	}

	c.modifiers = append(c.modifiers, modifiers...)
	for _, modifier := range c.modifiers {
		modifier(c)
	}

	// TODO: this output should be made a little nicer to view
	fmt.Println(rc)
	fmt.Println(c)

	return c, nil
}

func WithLogger(logger logger.Logger) ClientModifier {
	return func(c *Client) {
		c.logger = logger
	}
}

func WithTimeout(timeout time.Duration) ClientModifier {
	return func(c *Client) {
		c.rc.SetTimeout(timeout)
	}
}

func WithDebug() ClientModifier {
	return func(c *Client) {
		c.debugEnabled = true
		c.rc.SetDebug(true)
	}
}

func WithTrace() ClientModifier {
	return func(c *Client) {
		c.traceEnabled = true
		c.rc.EnableTrace()
	}
}

func (c *Client) get(url string, params map[paramOption]string) (*resty.Response, error) {
	return c.request(resty.MethodGet, url, params)
}

func (c *Client) post(url string, params map[paramOption]string, body interface{}) (*resty.Response, error) {
	// TODO: handle body
	return c.request(resty.MethodPost, url, params)
}

func (c *Client) patch(url string, params map[paramOption]string, body interface{}) (*resty.Response, error) {
	// TODO: handle body
	return c.request(resty.MethodPatch, url, params)
}

func (c *Client) request(method string, url string, params map[paramOption]string) (*resty.Response, error) {

	query := c.createQuery(params)

	resp, err := c.rc.R().Execute(method, url+query)
	if err != nil {
		return nil, err
	}

	if c.debugEnabled {
		c.debug(resp.Status())
		c.debug(resp.String())
	}

	if c.traceEnabled {
		// TODO: do more with the TraceInfo struct?
		c.info("trace: " + fmt.Sprint(resp.Request.TraceInfo()))
	}

	statusCode := resp.StatusCode()
	if statusCode < 200 || statusCode >= 300 { // TODO: check that this is OK
		return nil, fmt.Errorf("response failed with status: %d (%s)", statusCode, http.StatusText(statusCode))
	}

	return resp, nil
}

func (c *Client) createQuery(params map[paramOption]string) string {

	queryParts := []string{}

	_, foundFilters := params[filtersOption]
	if foundFilters {
		// TODO: process multiple filters;
	} else {
		expand, foundExpand := params[expandOption]
		if foundExpand {
			part := "$expand=" + expand
			queryParts = append(queryParts, part)
		}
		filter, foundFilter := params[filterOption]
		if foundFilter {
			part := "$filter=" + filter
			queryParts = append(queryParts, part)
		}
	}

	query := "?" + strings.Join(queryParts, "&")

	return query
}

func (c *Client) info(message string) {

	if c.logger == nil {
		return
	}

	c.logger.Info(message)
}

func (c *Client) debug(message string) {

	if c.logger == nil {
		return
	}

	c.logger.Debug(message)
}
