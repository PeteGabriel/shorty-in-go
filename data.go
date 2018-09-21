package main

/*
Code ...
*/
type Code struct {
	Url       string `json:"url"`
	Shortcode string `json:"shortcode"`
}

/*
ApiError ..
*/
type ApiError struct {
	Error int    `json:"error"`
	Desc  string `json:"description"`
}
