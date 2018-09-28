package main

import "net/http"


var db map[string]Code

func init{
	db = make(map[string]Code)
}

//Save some shortened code
func SaveCode(dto CodeDto) (bool, error) {

	if code.Url == "" {
		//TODO this should be a specific error, instead of soemthing so generic as error type
		return false, error.New("Url provided is missing.")
	}
	
	if code.Shortcode == "" {
		code.Shortcode = genCode()
	}

	db[code.Shortcode] = code

	return true, nil
}

//Get some code
func RetrieveCode(code string) Code {
	return nil
}

//Check if some code is present
func IsCodePresent(code string) bool {
	return false
}
