package client

import (
	"net/url"
	"strings"
)

type ODataModifiers struct {
	modifiers map[string]string
}

func NewODataModifiers() *ODataModifiers {
	return &ODataModifiers{
		// TODO: different type?
		modifiers: map[string]string{},
	}
}

func (o *ODataModifiers) Get() map[string]string {
	return o.modifiers
}

// func (o *ODataModifiers) AddSelect(values string) *ODataModifiers {
// 	// if o.mods == nil {
// 	// 	o.mods = map[string]string{}
// 	// }
// 	// TODO: do vararg approach?
// 	o.modifiers["$select"] = values
// 	return o
// }

func (o *ODataModifiers) AddFilter(values string) *ODataModifiers {
	// if o.mods == nil {
	// 	o.mods = map[string]string{}
	// }
	o.modifiers["$filter"] = values
	return o
}

func createURL(endpoint string, modifiers *ODataModifiers) string {

	// apiURL, _ := url.Parse(endpoint)
	// query := "" //apiURL.Query() // url.Values{}
	// for k, v := range modifiers.Get() {
	// 	//query.Set(k, trimMultiline(v))

	// }
	// //apiURL.RawQuery = query.Encode()
	// return apiURL.String()

	apiURL, _ := url.Parse(endpoint)

	mods := modifiers.Get()

	if len(mods) == 0 {
		return apiURL.String()
	}

	queryParts := []string{}

	for k, v := range mods {
		if k == "$filter" {
			part := "$filter=" + trimMultiline(v)
			queryParts = append(queryParts, part)
		}
	}

	// _, foundFilters := params[filtersOption]
	// if foundFilters {
	// 	// TODO: process multiple filters;
	// } else {
	// 	expand, foundExpand := params[expandOption]
	// 	if foundExpand {
	// 		part := "$expand=" + expand
	// 		queryParts = append(queryParts, part)
	// 	}
	// 	filter, foundFilter := params[filterOption]
	// 	if foundFilter {
	// 		part := "$filter=" + filter
	// 		queryParts = append(queryParts, part)
	// 	}
	// }

	query := "?" + strings.Join(queryParts, "&")

	return apiURL.String() + query

}

func trimMultiline(multi string) string {
	res := ""
	for _, line := range strings.Split(multi, "\n") {
		res += strings.Trim(line, "\t")
	}
	return res
}
