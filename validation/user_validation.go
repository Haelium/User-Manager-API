package validation

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
)

type address struct {
	Name     string `json:"name"`
	Line1    string `json:"line 1"`
	Line2    string `json:"line 2,omitempty"`
	Line3    string `json:"line 3,omitempty"`
	Region   string `json:"region"`
	Postcode string `json:"postcode,omitempty"`
	Country  string `json:"country"`
}

type user struct {
	Username string  `json:"username"`
	FullName string  `json:"fullname"`
	Email    string  `json:"email"`
	Address  address `json:"address"`
}

// email_regexp

func validateEmail(input string) error {

	return nil
}

func validateFullname(input string) error {
	// Single unicode character names exist, assuming 2 names seperated by space, 3 character is minimum
	if len(input) < 3 {
		return errors.New("Fullname is less than 3 characters")
	} else if len(input) > 128 {
		return errors.New("Fullname is greater than 128 characters")
	}

	return nil
}

func validateUsername(input string) error {
	// Only accepting 8-64 alphanumeric characters. First character must be alphabetic
	if len(input) < 8 {
		return errors.New("Username is less than 8 characters")
	}

	if len(input) > 64 {
		return errors.New("Username is greater than 64 characters")
	}

	if input[0] < 'A' || input[0] > 'z' {
		return errors.New("Username does not begin with a roman alphabetic character")
	}

	if !isAlphaNumeric.MatchString(input) {
		return errors.New("Username is not alphanumeric")
	}

	return nil
}

func validateAddress(input address) error {
	if input.Name == "" {
		return errors.New("Address Name is a required field")
	} else if input.Line1 == "" {
		return errors.New("Address Line1 is a required field")
	} else if input.Region == "" {
		return errors.New("Address Region is a required field")
	} else if input.Country == "" {
		return errors.New("Address Country is a required field")
	} else {
		return nil
	}
}

func ValidateUser(input string) (string, error) {
	var newuser user

	err := json.Unmarshal([]byte(input), &newuser)
	if err != nil {
		return "", err
	}

	if newuser.Username == "" {
		return "", errors.New("Username is a required field")
	} else if newuser.FullName == "" {
		return "", errors.New("Fullname is a required field")
	} else if newuser.Email == "" {
		return "", errors.New("Email is a required field")
	}

	if err = validateUsername(newuser.Username); err != nil {
		return "", err
	} else if err = validateFullname(newuser.FullName); err != nil {
		return "", err
	} else if err = validateEmail(newuser.Email); err != nil {
		return "", err
	}

	// Address only requires Name, Address Line 1, region, and Country
	if err = validateAddress(newuser.Address); err != nil {
		return "", err
	}

	// Convert all new usernames to lowercase, as their input should be case insensitive
	return strings.ToLower(newuser.Username), nil
}

var isAlphaNumeric regexp.Regexp

func init() {
	isAlphaNumeric = *regexp.MustCompile(`^[A-Za-z0-9]+$`)
}
