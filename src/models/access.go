package models

import (
	"log"
	"strings"
	"encoding/base64"
)

type (
	Authorize struct {}
	User struct {
		Auid string `json:"auid"`
	}
)


// Access exposes all of the access methods
func Access() *Authorize {
	return &Authorize{}
}


func (uc Authorize) VerifyHeader(header map[string][]string) bool {

	// make sure content type has been set
	if val, ok := header["Content-Type"]; ok {

		// make sure content type is correct
		if val[0] == "application/json" {

			// make sure Authorization has been set
			if _, ok := header["Authorization"]; ok {

				// make sure Key has been set
				if _, ok := header["Key"]; ok {

					return true // everything was set correctly

				} else {
					log.Println("Key is missing")
					return false
				}

			} else {
				log.Println("Authorization is missing")
				return false
			}

		} else {
			log.Println("Incorrect Content Type")
			return false
		}

	} else {
		log.Println("Missing Content Type")
		return false
	}
}

func (uc Authorize) VerifyUploadHeader(header map[string][]string) bool {

	// make sure content type has been set
	if val, ok := header["Content-Type"]; ok {

		// make sure content type is correct
		if strings.Contains(val[0], "multipart/form-data") {

			// make sure Authorization has been set
			if _, ok := header["Authorization"]; ok {

				// make sure Key has been set
				if _, ok := header["Key"]; ok {

					return true // everything was set correctly

				} else {
					log.Println("Key is missing")
					return false
				}

			} else {
				log.Println("Authorization is missing")
				return false
			}

		} else {
			log.Println("Incorrect Content Type")
			return false
		}

	} else {
		log.Println("Missing Content Type")
		return false
	}
}

func (uc Authorize) VerifyToken(header map[string][]string) bool {


	if val, ok := header["Authorization"]; ok {

		// get token
		token64 := val[0]
		token64 = strings.TrimPrefix(token64, "Bearer ")

		// decode token and return raw token
		rawtoken, err := base64.StdEncoding.DecodeString(token64)
		if err != nil {
			log.Println(err)
			return false
		}

		// return token
		token := strings.TrimSuffix(string(rawtoken), ":")

		//fmt.Println(token)
		if token == "eyJpYXQiOjE0OTcyMzA5MjYsImFsZyI6IkhTMjU2IiwiZXhwIjoxNDk5ODIyOTI2fQ.eyJpZCI6bnVsbH0.oHgJxBDcXwgw_v92IELbdUfIc8e2f-wDHIV-MTTKa7E" {
			return true
		}

		log.Println("Bad Token")
		return false

	} else {
		log.Println("Authorization is missing")
		return false
	}

}

func (uc Authorize) VerifyKey(header map[string][]string) (*User, bool) {

	if val, ok := header["Key"]; ok {

		// get key
		key := val[0]

		user := &User {
			Auid: "123456",
		}

		//fmt.Println(token)
		if key == "z4JoA3e5XOETpt1ymxhHWs2ltM0EE4RAV8siCElT" {
			return user, true
		}

		log.Println("Bad Key")
		return user, false

	} else {
		log.Println("Key is missing")
		return &User{}, false
	}

}
