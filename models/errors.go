package models

//=============================================================================================
// start of generated code: please do not modify this section
// code generated on 27 May 20 17:32 CDT
//=============================================================================================

import (
	"strings"
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {

	// replace 'models: ' with ""
	s := strings.Replace(string(e), "models: ", "", 1)
	return strings.Title(s)
}

type privateError string

func (e privateError) Error() string {
	return string(e)
}

// ErrPasswordTooShort - used in the User model validations.  Passwords must contain at least six characters.
const ErrPasswordTooShort modelError = "models: password must be at least 6 characters"

// ErrPasswordHashRequired - used in the User model validations.
const ErrPasswordHashRequired modelError = "models: a password hash is required"

// ErrEmailRequired  - used in the User model validations.
const ErrEmailRequired modelError = "models: an email address is required"

// ErrEmailInvalid - used in the User model validations.  Checks for
// regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`).
const ErrEmailInvalid modelError = "models: the provided email address is not valid"

// ErrUserIsNotActive - indicates that the user is currently locked
const ErrUserIsNotActive modelError = "models: user is currently locked"

// ErrNotFound - the requested resource was not found in the db
const ErrNotFound modelError = "models: resource not found"

// ErrEmailTaken - the email address has already been taken by another user
const ErrEmailTaken modelError = "models: email address has been taken by another user"

// ErrGroupRequired - at least one Group must be specified when creating or updating a user record
const ErrGroupRequired modelError = "models: at least one Group must be specified when creating or updating a user record"

// ErrInvalidPassword - used in the User model.
const ErrInvalidPassword modelError = "models: the provided password is not valid"

//=============================================================================================
// end of generated code
//=============================================================================================

//=============================================================================================
// implement additional error messages below:
//=============================================================================================

// const ErrNewCustomerErr modelError = "models: your message here"
