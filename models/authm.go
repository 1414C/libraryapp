package models

//=============================================================================================
// Auth entity model code generated on 27 May 20 17:57 CDT
//=============================================================================================

import (
	"fmt"

	"github.com/1414C/lw"
	"github.com/1414C/sqac"
)

// Auth structure
type Auth struct {
	ID          uint64 `json:"id" db:"id" sqac:"primary_key:inc"`
	Href        string `json:"href" db:"href" sqac:"-"`
	AuthName    string `json:"auth_name" db:"auth_name" sqac:"nullable:false;index:non-unique"`
	AuthType    string `json:"auth_type" db:"auth_type" sqac:"nullable:false"`
	Description string `json:"description" db:"description" sqac:"nullable:false"`
}

// AuthDB is a CRUD-type interface specifically for dealing with Auths.
type AuthDB interface {
	Create(auth *Auth) error
	Update(auth *Auth) error
	Delete(auth *Auth) error
	Get(auth *Auth) error
	GetAuths() []Auth
	GetAuthsByAuthName(op string, AuthName string) []Auth
	GetAuthsByDescription(op string, Description string) []Auth
}

// authValidator checks and normalizes data prior to
// db access.
type authValidator struct {
	AuthDB
}

// authValFunc type is the prototype for discrete Auth normalization
// and validation functions that will be executed by func runAuthValidationFuncs(...)
type authValFunc func(*Auth) error

// AuthService is the public interface to the Usr entity
type AuthService interface {
	AuthDB
}

// private service for auth
type authService struct {
	AuthDB
}

// authSqac is a sqac-based implementation of the AuthDB interface.
type authSqac struct {
	handle sqac.PublicDB
}

var _ AuthDB = &authSqac{}

// newAuthValidator returns a new authValidator
func newAuthValidator(gdb AuthDB) *authValidator {
	return &authValidator{
		AuthDB: gdb,
	}
}

// runAuthValFuncs executes a list of discrete validation
// functions against a auth.
func runAuthValFuncs(auth *Auth, fns ...authValFunc) error {

	// iterate over the slice of function names and execute
	// each in-turn.  the order in which the lists are made
	// can matter...
	for _, fn := range fns {
		err := fn(auth)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewAuthService ...
func NewAuthService(handle sqac.PublicDB) AuthService {

	gs := &authSqac{handle}

	gv := newAuthValidator(gs) // *db
	return &authService{
		AuthDB: gv,
	}
}

// ensure consistency (build error if delta exists)
var _ AuthDB = &authValidator{}

//-------------------------------------------------------------------------------------------------------
// CRUD-type model methods for Auth
//-------------------------------------------------------------------------------------------------------
//
// Create validates and normalizes data used in the auth creation.
// Create then calls the creation code contained in AuthService.
func (gv *authValidator) Create(auth *Auth) error {

	// perform normalization and validation -- comment out checks that are not required
	// note that the check calls are generated as a straight enumeration of the entity
	// structure.  It may be neccessary to adjust the calling order depending on the
	// relationships between the fields in the entity structure.
	err := runAuthValFuncs(auth,
		gv.normvalAuthName,
	)

	if err != nil {
		return err
	}
	return gv.AuthDB.Create(auth)
}

// Update validates and normalizes the content of the Auth
// being updated by way of executing a list of predefined discrete
// checks.  if the checks are successful, the entity is updated
// on the db via the ORM.
func (gv *authValidator) Update(auth *Auth) error {

	// perform normalization and validation -- comment out checks that are not required
	// note that the check calls are generated as a straight enumeration of the entity
	// structure.  It may be neccessary to adjust the calling order depending on the
	// relationships between the fields in the entity structure.
	err := runAuthValFuncs(auth,
		gv.normvalAuthName,
	)

	if err != nil {
		return err
	}
	return gv.AuthDB.Update(auth)
}

// Delete is passed through to the ORM with no real
// validations.  id is checked in the controller.
func (gv *authValidator) Delete(auth *Auth) error {

	return gv.AuthDB.Delete(auth)
}

// Get is passed through to the ORM with no real
// validations.  id is checked in the controller.
func (gv *authValidator) Get(auth *Auth) error {

	return gv.AuthDB.Get(auth)
}

// GetAuths is passed through to the ORM with no validation
func (gv *authValidator) GetAuths() []Auth {

	return gv.AuthDB.GetAuths()
}

//-------------------------------------------------------------------------------------------------------
// internal authValidator funcs
//-------------------------------------------------------------------------------------------------------
// These discrete functions are used to normalize and validate the Entity fields
// from with in the Create and Update methods.  See the comments in the model's
// Create and Update methods for details regarding use.

// normvalAuthName normalizes and validates field GroupName
func (gv *authValidator) normvalAuthName(auth *Auth) error {

	// TODO: implement normalization and validation for Auth
	return nil
}

//-------------------------------------------------------------------------------------------------------
// internal book relations Validator funcs
//-------------------------------------------------------------------------------------------------------

// // GetRelBookToLibrary is passed through to the ORM with no validation.
// // belongsTo relationship
// func (bv *bookValidator) GetRelBookToLibrary(librarys *[]Library, mapParams map[string]interface{}) error {

// 	// TODO: implement normalization and validation for the GetRelBookToLibrary call.
// 	// TODO: typically no modifications are required here.
// 	return bv.BookDB.GetRelBookToLibrary(librarys, mapParams)
// }

//-------------------------------------------------------------------------------------------------------
// internal auth Simple Query Validator funcs
//-------------------------------------------------------------------------------------------------------
// Simple query normalization and validation occurs in the controller to an
// extent, as the URL has to be examined closely in order to determine what to call
// in the model.  This section may be blank if no model fields were marked as
// selectable in the <models>.json file.
// GetAuthsByAuthName is passed through to the ORM with no validation.
func (gv *authValidator) GetAuthsByAuthName(op string, AuthName string) []Auth {

	// TODO: implement normalization and validation for the GetAuthsByAuthName call.
	// TODO: typically no modifications are required here.
	return gv.AuthDB.GetAuthsByAuthName(op, AuthName)
}

// GetAuthsByDescription is passed through to the ORM with no validation.
func (gv *authValidator) GetAuthsByDescription(op string, Description string) []Auth {

	// TODO: implement normalization and validation for the GetAuthsByDescription call.
	// TODO: typically no modifications are required here.
	return gv.AuthDB.GetAuthsByDescription(op, Description)
}

//-------------------------------------------------------------------------------------------------------
// ORM db CRUD access methods
//-------------------------------------------------------------------------------------------------------
//
// Create a new Auth in the database via the ORM
func (gs *authSqac) Create(auth *Auth) error {
	return gs.handle.Create(auth)
}

// Update an existng Auth in the database via the ORM
func (gs *authSqac) Update(auth *Auth) error {
	return gs.handle.Update(auth)
}

// Delete an existing Auth in the database via the ORM
func (gs *authSqac) Delete(auth *Auth) error {
	return gs.handle.Delete(auth)
}

// Get an existing Auth from the database via the ORM
func (gs *authSqac) Get(auth *Auth) error {
	return gs.handle.GetEntity(auth)
}

// Get all existing Auths from the db via the ORM
func (gs *authSqac) GetAuths() []Auth {

	getEnts := AuthGetEntitiesT{}

	err := gs.handle.GetEntities2(&getEnts)
	if err != nil {
		lw.Error(err)
		return getEnts.ents
	}
	return getEnts.ents
}

//- - - - - - - - - - - - - - - - - - -- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
//  Implement support for the sqac.GetEnt{} interface to support GetEntities2 calls
//- - - - - - - - - - - - - - - - - - -- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -

// AuthGetEntitiesT is used to support the implmentation of the sqac.GetEnt{} interface.
type AuthGetEntitiesT struct {
	ents []Auth
	sqac.GetEnt
}

// Exec implments the one and only method in the  sqac.GetEnt{} interface for AuthGetEntitiesT.
// An implemented interface is passed to sqac.GetEntities2 and the code contained withing the Exec()
// method is executed to retrieve the entities from the DB.  Retrieved entities are stored in the
// AuthGetEntitiesT struct's ents []Auth field.
func (ge *AuthGetEntitiesT) Exec(sqh sqac.PublicDB) error {

	selQuery := "SELECT * FROM auth;"

	// read the table rows
	rows, err := sqh.ExecuteQueryx(selQuery)
	if err != nil {
		lw.Warning("GetEntities2 for table auth returned error: %s", err.Error())
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var ent Auth
		err = rows.StructScan(&ent)
		if err != nil {
			lw.ErrorWithPrefixString("error reading Auth rows via StructScan:", err)
			return err
		}
		ge.ents = append(ge.ents, ent)
	}
	if sqh.IsLog() {
		lw.Debug("AuthGetEntitiesT.Exec() got: %v", ge.ents)
	}
	return nil
}

//- - - - - - - - - - - - - - - - - - -- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -

//-------------------------------------------------------------------------------------------------------
// ORM db relations selector access methods
//-------------------------------------------------------------------------------------------------------
//
// func (bs *bookSqac) GetRelBookToLibrary(librarys *[]Library, mapParams map[string]interface{}) error {

// 	var c string
// 	var s string

// 	for k, v := range mapParams {
// 		switch reflect.TypeOf(v).String() {
// 		case "string":
// 			s = fmt.Sprintf("%s%s = '%v' AND ", c, common.CamelToSnake(k), v)
// 		default:
// 			s = fmt.Sprintf("%s%s = %v AND ", c, common.CamelToSnake(k), v)
// 		}
// 		c = s
// 	}
// 	c = strings.TrimSuffix(c, " AND ")

// 	qs := fmt.Sprintf("SELECT * FROM library WHERE %s;", c)
// 	fmt.Println("qs:", qs)
// 	err := bs.handle.Select(librarys, qs)
// 	if err != nil {
// 		log.Println("GetRelBookToLibrary:", err)
// 		return nil
// 	}

// 	if bs.handle.IsLog() {
// 		log.Printf("GetRelLibraryBooks found: %v \n", librarys)
// 	}
// 	return nil
// }

//-------------------------------------------------------------------------------------------------------
// ORM db simple selector access methods
//-------------------------------------------------------------------------------------------------------
//
// Get all existing AuthsByAuthName from the db via the ORM
func (gs *authSqac) GetAuthsByAuthName(op string, AuthName string) []Auth {

	var auths []Auth
	var c string

	switch op {
	case "EQ":
		c = "auth_name = ?"
	case "LIKE":
		c = "auth_name like ?"
	default:
		return nil
	}
	qs := fmt.Sprintf("SELECT * FROM auth WHERE %s;", c)
	err := gs.handle.Select(&auths, qs, AuthName)
	if err != nil {
		lw.Warning("GetAuthsByAuthName: %s", err.Error())
		return nil
	}

	if gs.handle.IsLog() {
		lw.Debug("GetAuthsByAuthName found: %v based on (%s %v)", auths, op, AuthName)
	}
	return auths
}

// GetAuthsByDescription gets all existing AuthsByDescription from the db via the ORM
func (gs *authSqac) GetAuthsByDescription(op string, Description string) []Auth {

	var auths []Auth
	var c string

	switch op {
	case "EQ":
		c = "description = ?"
	case "LIKE":
		c = "description like ?"
	default:
		return nil
	}
	qs := fmt.Sprintf("SELECT * FROM auth WHERE %s;", c)
	err := gs.handle.Select(&auths, qs, Description)
	if err != nil {
		lw.Warning("GetAuthsByDescription: %s", err.Error())
		return nil
	}

	if gs.handle.IsLog() {
		lw.Debug("GetAuthsByDescription found: %v based on (%s %v)", auths, op, Description)
	}
	return auths
}
