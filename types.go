package main

// Code represents a data layer type. This kind of types
// should not be exposed outside of the system.
// This type represents a shortened code. The url field
// is the url in its raw version. The shortcode field
// holds the shortened version of the url.
type Code struct {
	Url       string `json:"url"`
	Shortcode string `json:"shortcode"`
}

// CodeDto type should be used when answering some request.
type CodeDto struct {
	Shortcode string `json:"shortcode"`
}

// APIError should be used to represent erros that might happen
// during request/responses cycles.
type APIError struct {
	Error int    `json:"error"`
	Desc  string `json:"description"`
}
