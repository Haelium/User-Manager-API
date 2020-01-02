package validation

import (
	"errors"
	"testing"
)

var valid_users = map[string]string{
	"bobman12": `{"username": "bobman12", "email": "bob@bobmail.com"}`,
	"jcdenton": `{"username": "jcdenton", "email": "jc@unatco.org"}`,
	"herpderp": `{"username": "herpderp", "email": "herp@derp.io"}`,
}

var invalid_users = map[string]string{
	"missingusername": `{"email": "missing@username.com"}`,
	"missingemail":    `{"username": "missingemail"}`,
}

var invalid_json = []string{
	`{"Json}`,
	`{"key": val}`,
	`{key: "val"}`,
	`{"key:" "val"}`,
	`{"key": {"subkey": "subval}`,
	``,
}

func Test_ValidateUser_invalid_json(t *testing.T) {
	for _, input := range invalid_json {
		_, err := ValidateUser(input)

		if err != nil {
			t.Logf("Invalid json %s was not rejected", input)
		}
	}
}

func Test_ValidateUser_UserMissingField(t *testing.T) {
	// Missing fullname
	user_missing_fullname := `{"username": "billy", "email": "Bob@bobmail.bob", "address": {"name": "Bob", "Line 1": "44 Bobstreet", "region": "Bobville", "country": "Bobland"}}`

	expected_err := errors.New("Fullname is a required field")
	returnval, actual_err := ValidateUser(user_missing_fullname)

	if actual_err != expected_err {
		t.Logf("Expected: %s\n Got: %s\n", expected_err, actual_err)
	}
	if returnval != "" {
		t.Logf("Return value should be empty string, Got: %s\n", returnval)
	}

	// Missing fullname
	user_missing_username := `{"fullname": "billy bobson", "email": "Bob@bobmail.bob", "address": {"name": "Bob", "Line 1": "44 Bobstreet", "region": "Bobville", "country": "Bobland"}}`

	expected_err = errors.New("Username is a required field")
	returnval, actual_err = ValidateUser(user_missing_username)

	if actual_err != expected_err {
		t.Logf("Expected: %s\n Got: %s\n", expected_err, actual_err)
	}
	if returnval != "" {
		t.Logf("Return value should be empty string, Got: %s\n", returnval)
	}

	// Missing email
	user_missing_email := `{"fullname": "billy bobson", "username": "bobman2000", "address": {"name": "Bob", "Line 1": "44 Bobstreet", "region": "Bobville", "country": "Bobland"}}`

	expected_err = errors.New("Email is a required field")
	returnval, actual_err = ValidateUser(user_missing_email)

	if actual_err != expected_err {
		t.Logf("Expected: %s\n Got: %s\n", expected_err, actual_err)
	}
	if returnval != "" {
		t.Logf("Return value should be empty string, Got: %s\n", returnval)
	}

}

func Test_validateAddress(t *testing.T) {
	valid_addresses := []address{
		address{
			Name:     "Mr One",
			Line1:    "22 Woodroad",
			Line2:    "Ballyville",
			Region:   "Cork",
			Postcode: "11111",
			Country:  "Ireland",
		},
		address{
			Name:    "Mr Two",
			Line1:   "16 Fallsroad",
			Line2:   "Derry",
			Region:  "Belfast",
			Country: "Northern Ireland",
		},
		address{
			Name:    "The Occupant",
			Line1:   "22 Slumville",
			Region:  "Limerick",
			Country: "Something",
		},
		address{
			Name:     "Lord O'Fancy",
			Line1:    "Fancyland Manor",
			Line2:    "Use all lines",
			Line3:    "allfields",
			Region:   "Someregion",
			Postcode: "12123123",
			Country:  "Wakanda",
		},
		address{
			Name:    "名称",
			Line1:   "Somewhere",
			Region:  "عنوان",
			Country: "薛대한민국",
		},
	}

	for _, address := range valid_addresses {
		err := validateAddress(address)
		if err != nil {
			t.Logf("Expected nil, got: %s", err)
		}
	}

	missing_name := address{
		Line1:   "44 address blah blah",
		Region:  "Dublin 15",
		Country: "Ireland",
	}
	expected_err := errors.New("Address Name is a required field")
	err := validateAddress(missing_name)
	if err != expected_err {
		t.Logf("Expected: %s\nGot: %s\n", expected_err, err)
	}

}
