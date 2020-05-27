package models

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/1414C/lw"
	"github.com/1414C/sqac"
	"golang.org/x/crypto/bcrypt"
)

// Usr represents the signed-in usr
type Usr struct {
	ID           uint64     `db:"id" sqac:"primary_key:inc"`
	Href         string     `json:"href" db:"href" sqac:"-"`
	Name         string     `db:"name" sqac:"nullable:false"`
	Email        string     `db:"email" sqac:"nullable:false;index:unique"`
	Password     string     `db:"password" sqac:"-"` // "-" indicates that the field is not to be stored in the db
	PasswordHash string     `db:"password_hash" sqac:"nullable:false"`
	CreatedOn    *time.Time `json:"created_on,omitempty" db:"created_on" sqac:"nullable:false;default:now()"`
	UpdatedOn    *time.Time `json:"updated_on,omitempty" db:"updated_on" sqac:"nullable:false;default:now()"`
	Active       bool       `db:"active" sqac:"nullable:false;default:false"`
	Groups       *string    `json:"groups,omitempty" db:"groups" sqac:"nullable:true"` // try with inline{group1;group2;group3} for now
}

// UsrDB is an interface that outlines the methods that can be
// used to manipulate usr records.
// for single usr queries, any error but ErrNotFound will
// result in a http status codd of 500.
type UsrDB interface {

	// methods for altering single Usr entities
	Create(usr *Usr) error
	Update(usr *Usr) error
	Delete(usr *Usr) error
	Get(usr *Usr) error
	GetUsrs() []Usr

	// methods for querying single Usr entities
	ByEmail(email string) (*Usr, error)
	// ByID(id uint) (*Usr, error) // testing-only
}

// UsrService exposes the set of methods that are made available
// to manipulate and work with the usr model. (called from controller)
type UsrService interface {

	// Authenticate will verify that the usr credentials are
	// correct.  Errors will be:
	// ErrNotFound, ErrInvalidPassword, or other(!)
	Authenticate(email string, password string) (*Usr, error)
	UsrDB
}

// usrService supports internal interaction with Usr
type usrService struct {
	UsrDB
	pepper string
}

// NewUsrService needs some work:
func NewUsrService(handle sqac.PublicDB, pepper string) UsrService {

	us := &usrSqac{handle}

	// create a new usrValidator
	uv := newUsrValidator(us, pepper)

	return &usrService{
		UsrDB:  uv,
		pepper: pepper,
	}
}

// Authenticate can be used to authenticate a usr with the provided email address and password
func (us *usrService) Authenticate(email string, password string) (*Usr, error) {

	// lookup usr record
	foundUsr, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}

	// user is active?
	if !foundUsr.Active {
		return nil, ErrUserIsNotActive
	}

	// compare password + pepper with stored hash
	err = bcrypt.CompareHashAndPassword([]byte(foundUsr.PasswordHash), []byte(password+us.pepper))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrInvalidPassword
		default:
			return nil, err
		}
	}
	return foundUsr, nil
}

// usrValidator is a layer that validates things before
// they go to the db to perform queries.  normalization
// and validation...
type usrValidator struct {
	UsrDB
	emailRegex *regexp.Regexp
	pepper     string
}

// ensure consistency
var _ UsrDB = &usrValidator{}

func newUsrValidator(udb UsrDB, pepper string) *usrValidator {
	return &usrValidator{
		UsrDB:      udb,
		emailRegex: regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
		pepper:     pepper,
	}
}

// usrValFunc is the function prototype for discrete usr validation
// functions and methods.
type usrValFunc func(*Usr) error

// runUsrValFuncs is a function that accepts a usr and
// then runs a list of discrete validation functions against
// it.
func runUsrValFuncs(usr *Usr, fns ...usrValFunc) error {

	// iterate over the slice of function names
	for _, fn := range fns {
		err := fn(usr)
		if err != nil {
			return err
		}
	}
	return nil
}

// Create backfills usr data like the ID, CreatedAt and
// UpdatedAt fields prior to the calling of the UsrDB.Create
// method which will create the usr in the storage-layer.
func (uv *usrValidator) Create(usr *Usr) error {

	// call discrete usrValFuncs(...)
	err := runUsrValFuncs(usr,
		uv.passwordRequired,     // check that a password has been provided
		uv.passwordMinLength,    // check that the password meets min length criteria
		uv.bcryptPassword,       // bcrypt usr.Password -> usr.PasswordHash
		uv.passwordHashRequired, // check that a passwordHash was computed
		uv.requireEmail,         // check email address format
		uv.normalizeEmail,       // normalize content of usr.Email
		uv.emailFormat,          // check the email format via regexp
		uv.emailIsAvailable,     // check that the email address is available
		uv.requireGroups,        // check that at least one group has been assigned to the user
		uv.normalizeGroups,      // normalize the format of the Groups string
	)

	if err != nil {
		return err
	}
	return uv.UsrDB.Create(usr)
}

// Update computes a rememberHash value for the usr
func (uv *usrValidator) Update(usr *Usr) error {

	// call discrete usrValFuncs(...)
	err := runUsrValFuncs(usr) // uv.passwordMinLength,    // check that the password meets min length criteria
	// uv.bcryptPassword,       // bcrypt usr.Password -> usr.PasswordHash
	// uv.passwordHashRequired, // check that a passwordHash was computed
	// uv.normalizeEmail,       // normalize content of usr.Email
	// uv.requireEmail,         // check email address format
	// uv.emailFormat,          // check the email format via regexp
	// uv.emailIsAvailable,     // check that the email address is available
	// uv.requireGroups,        // check that at least one group has been assigned to the user
	// uv.normalizeGroups,      // normalize the format of the Groups string

	if err != nil {
		return err
	}
	return uv.UsrDB.Update(usr)
}

// Get is passed through to the ORM with no real
// validations.  id is checked in the controller.
func (uv *usrValidator) Get(usr *Usr) error {

	return uv.UsrDB.Get(usr)
}

// Delete validates the usr related to the specified ID
func (uv *usrValidator) Delete(usr *Usr) error {

	return uv.UsrDB.Delete(usr)
}

// GetUsrs is passed through to the ORM with no validation
func (uv *usrValidator) GetUsrs() []Usr {

	return uv.UsrDB.GetUsrs()
}

// ByEmail calls the normalization function(s) for the
// usr.Email address, then calls the storage-layer
// if the normalization was successful.
func (uv *usrValidator) ByEmail(email string) (*Usr, error) {

	usr := Usr{
		Email: email,
	}
	err := runUsrValFuncs(&usr, uv.normalizeEmail)
	if err != nil {
		return nil, err
	}
	return uv.UsrDB.ByEmail(usr.Email)
}

// private usrValidator methods
// passwordRequired checks that a password has been provided
func (uv usrValidator) passwordRequired(usr *Usr) error {

	if usr.Password == "" {
		return ErrPasswordTooShort
	}
	return nil
}

// passwordMinLength checks that the provided password contains
// at least 6 characters
func (uv *usrValidator) passwordMinLength(usr *Usr) error {

	if usr.Password == "" {
		return nil
	}
	if len(usr.Password) < 6 {
		return ErrPasswordTooShort
	}
	return nil
}

// bcryptPassword conforms to type usrValFunc func(*Usr) error.
// the method will hash the usr password with a predefined
// pepper and bcrypt if the Password field is not an empty string.
func (uv *usrValidator) bcryptPassword(usr *Usr) error {

	// if no password was provided, leave
	if usr.Password == "" {
		return nil
	}

	// add pepper to passwd
	pwBytes := []byte(usr.Password + uv.pepper)

	// create hash value from password text and set value on usr
	hashedByes, err := bcrypt.GenerateFromPassword(pwBytes, 14)
	if err != nil {
		return err
	}
	usr.PasswordHash = string(hashedByes)
	usr.Password = ""
	return nil
}

// passwordHashRequired checks that a password hash has been computed
func (uv usrValidator) passwordHashRequired(usr *Usr) error {

	if usr.PasswordHash == "" {
		return ErrPasswordHashRequired
	}
	return nil
}

// requireEmail checks the post-normalized value of the email
// to ensure that it contains a value.
func (uv *usrValidator) requireEmail(usr *Usr) error {

	if usr.Email == "" {
		return ErrEmailRequired
	}
	return nil
}

// normalizeEmail checks the email address for spaces and
// ensures that the address contains only lower-case.
func (uv *usrValidator) normalizeEmail(usr *Usr) error {

	usr.Email = strings.ToLower(usr.Email)
	usr.Email = strings.TrimSpace(usr.Email)
	return nil
}

// emailFormat checks the format of the email via a simple
// regex-based pattern match.
func (uv *usrValidator) emailFormat(usr *Usr) error {

	if usr.Email != "admin" {
		if !uv.emailRegex.MatchString(usr.Email) {
			return ErrEmailInvalid
		}
	}
	return nil
}

// emaiIsAvailable checks the backend db to ensure that the email address
// provided by the usr is available for use.
func (uv *usrValidator) emailIsAvailable(usr *Usr) error {

	// check the db using the email address - beware of cycles here
	existing, err := uv.ByEmail(usr.Email)
	if err == ErrNotFound {
		// email address is not taken
		return nil
	}
	if err != nil {
		return nil
	}

	// found a Usr with this email address..
	// if the found usr has the same ID as this usr, it
	// is an update.
	if usr.ID != existing.ID {
		return ErrEmailTaken
	}
	return nil
}

// requireGroups checks that at least one Group has been provided
func (uv *usrValidator) requireGroups(usr *Usr) error {

	if *usr.Groups == "" {
		return ErrGroupRequired
	}
	return nil
}

// normalizeGroups adjusts the format of the provided Groups string
// input:  "Group 1 ; Group 2;Group 3 ; Group 4   "
// output: "Group 1;Group 2;Group 3;Group 4"
func (uv *usrValidator) normalizeGroups(usr *Usr) error {

	sa := strings.Split(*usr.Groups, ";")
	*usr.Groups = ""
	for i := range sa {
		*usr.Groups = fmt.Sprintf("%s%s;", *usr.Groups, strings.TrimSpace(sa[i]))
	}
	*usr.Groups = strings.TrimSuffix(*usr.Groups, ";")
	return nil
}

//***************************************************************************
//
//		db-access for the usrValidator->usrSqac interface chain
//
//****************************************************************************
// usrSqac is a sqac-based implementation of the UsrDB interface.
type usrSqac struct {
	handle sqac.PublicDB
}

// inclusion of this line ensures that usrSqac will always adhere
// to the UsrDB interface.  Non-compliance will result in a
// compilation / linter error.
var _ UsrDB = &usrSqac{}

// ByEmail - lookup a Usr using the provided email address
// 1 - usr, nil
// 2 - nil, ErrNotFound
// 3 - nil, otherError
//
func (us *usrSqac) ByEmail(email string) (*Usr, error) {

	var usr Usr
	err := us.handle.Get(&usr, "SELECT * FROM usr WHERE email = ?;", email)
	if err != nil {
		lw.Warning("reading Usr by email got: %s", err.Error())
		return nil, err
	}
	return &usr, nil
}

// ByID - lookup a Usr using the provided id
// 1 - usr, nil
// 2 - nil, ErrNotFound
// 3 - nil, otherError
//
func (us *usrSqac) ByID(id uint64) (*Usr, error) {

	var usr Usr
	err := us.handle.Get(&usr, "SELECT * FROM usr WHERE id = ?;", id)
	if err != nil {
		lw.Warning("reading Usr by id got: %s", err.Error())
		return nil, err
	}
	return &usr, nil
}

// Create a new Usr in the db
func (us *usrSqac) Create(usr *Usr) error {
	return us.handle.Create(usr)
}

// Update an existing usr in the db
func (us *usrSqac) Update(usr *Usr) error {
	return us.handle.Update(usr)
}

// Get an existing Usr from the database via the ORM
func (us *usrSqac) Get(usr *Usr) error {
	return us.handle.GetEntity(usr)
}

// Delete the usr related to the specified ID
func (us *usrSqac) Delete(usr *Usr) error {

	return us.handle.Delete(usr)
}

// Get all existing Usrs from the db via the ORM
func (us *usrSqac) GetUsrs() []Usr {

	getEnts := UsrGetEntitiesT{}

	err := us.handle.GetEntities2(&getEnts)
	if err != nil {
		lw.Warning("GetUsrs got: %s", err.Error())
		return getEnts.ents
	}
	return getEnts.ents
}

//- - - - - - - - - - - - - - - - - - -- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
//  Support for the sqac.GetEnt{} interface to support GetEntities2 calls
//- - - - - - - - - - - - - - - - - - -- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -

// UsrGetEntitiesT is used to support the implmentation of the sqac.GetEnt{} interface.
type UsrGetEntitiesT struct {
	ents []Usr
	sqac.GetEnt
}

// Exec implements the one and only method in the  sqac.GetEnt{} interface for UsrGroupGetEntitiesT.
// An implemented interface is passed to sqac.GetEntities2 and the code contained withing the Exec()
// method is executed to retrieve the entities from the DB.  Retrieved entities are stored in the
// UsrGetEntitiesT struct's ents []Usr field.
func (ge *UsrGetEntitiesT) Exec(sqh sqac.PublicDB) error {

	selQuery := "SELECT * FROM usr;"

	// read the table rows
	rows, err := sqh.ExecuteQueryx(selQuery)
	if err != nil {
		lw.Warning("GetEntities2 for table usr returned error: %v", err.Error())
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var ent Usr
		err = rows.StructScan(&ent)
		if err != nil {
			lw.ErrorWithPrefixString("error reading usr rows via StructScan:", err)
			return err
		}
		ge.ents = append(ge.ents, ent)
	}
	if sqh.IsLog() {
		lw.Debug("UsrGroupGetEntitiesT.Exec() got: %v", ge.ents)
	}
	return nil
}
