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
