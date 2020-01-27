package form

import (
	"fmt"
	"net/url"
	"regexp"
	"time"
	"unicode/utf8"
	. "xCut/constants"
)

// PhoneRX represents phone number maching pattern
var PhoneRX = regexp.MustCompile("(^\\+[0-9]{2}|^\\+[0-9]{2}\\(0\\)|^\\(\\+[0-9]{2}\\)\\(0\\)|^00[0-9]{2}|^0)([0-9]{9}$|[0-9\\-\\s]{10}$)")

// EmailRX represents email address maching pattern
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")


var WebsiteRX = regexp.MustCompile("(http(s)?:\\/\\/.)?(www\\.)?[-a-zA-Z0-9@:%._\\+~#=]{2,256}\\.[a-z]{2,6}\\b([-a-zA-Z0-9@:%_\\+.~#?&//=]*)")


var LngLatRX = regexp.MustCompile("\\-?\\d+(\\.\\d+)?")





// Input represents form input values and validations
type Input struct {
	Values  url.Values
	VErrors ValidationErrors
	CSRF    string
}


// MinLength checks if a given minium length is satisfied
func (inVal *Input) MinLength(field string, d int) {
	value := inVal.Values.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) < d {
		inVal.VErrors.Add(field, fmt.Sprintf("This field is too short (minimum is %d characters)", d))
	}
}

// Required checks if list of provided form input fields have values
func (inVal *Input) ValidateRequiredFields(fields ...string) {
	for _, f := range fields {
		value := inVal.Values.Get(f)
		if value == "" {
			inVal.VErrors.Add(f, "This field is required field")
		}
	}
}

// MatchesPattern checks if a given input form field matchs a given pattern
func (inVal *Input) MatchesPattern(field string, pattern *regexp.Regexp) {
	value := inVal.Values.Get(field)
	if value == "" {
		return
	}
	if !pattern.MatchString(value) {
		inVal.VErrors.Add(field, "The value entered is invalid")
	}
}

// PasswordMatches checks if Password and Confirm Password fields match
func (inVal *Input) PasswordMatches(password string, confPassword string) {
	pwd := inVal.Values.Get(password)
	confPwd := inVal.Values.Get(confPassword)

	if pwd == "" || confPwd == "" {
		return
	}

	if pwd != confPwd {
		inVal.VErrors.Add(password, "The Password and Confirm Password values did not match")
		inVal.VErrors.Add(confPassword, "The Password and Confirm Password values did not match")
	}
}

func (inVal *Input) ValidateStartAndEnd(start string, end string) {
	start_time,err := time.Parse("15:04",inVal.Values.Get(start))
	if(err != nil){
		inVal.VErrors.Add(OpenHours, "Wrong time format")
		return
	}
	end_time,err := time.Parse("15:04",inVal.Values.Get(end))
	if(err != nil){
		inVal.VErrors.Add(OpenHours, "Wrong time format")
		return
	}



	if end_time.Before(start_time) {
		inVal.VErrors.Add(OpenHours, "Closing time can't be before opening hours")
	}
}


// Valid checks if any form input validation has failed or not
func (inVal *Input) IsValid() bool {
	return len(inVal.VErrors) == 0
}
