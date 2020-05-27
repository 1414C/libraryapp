package models

//=============================================================================================
// GroupAuth entity model code generated on 27 May 20 17:32 CDT
//=============================================================================================

import (
	"fmt"

	"github.com/1414C/lw"
	"github.com/1414C/sqac"
)

// GroupAuth structure
type GroupAuth struct {
	ID          uint64 `json:"id" db:"id" sqac:"primary_key:inc"`
	Href        string `json:"href" db:"href" sqac:"-"`
	GroupID     uint64 `json:"group_id" db:"group_id" sqac:"nullable:false"`
	GroupName   string `json:"group_name" db:"group_name" sqac:"-"`
	AuthID      uint64 `json:"auth_id" db:"auth_id" sqac:"nullable:false"`
	AuthName    string `json:"auth_name" db:"auth_name" sqac:"-"`
	AuthType    string `json:"auth_type" db:"auth_type" sqac:"-"`
	Description string `json:"description" db:"description" sqac:"-"`
}

// GroupAuthDB is a CRUD-type interface specifically for dealing with GroupAuths.
type GroupAuthDB interface {
	Create(groupauth *GroupAuth) error
	Update(groupauth *GroupAuth) error
	Delete(groupauth *GroupAuth) error
	Get(groupauth *GroupAuth) error
	GetGroupAuths() []GroupAuth
	GetGroupAuthsByAuthName(op string, AuthName string) []GroupAuth
	GetGroupAuthsByDescription(op string, Description string) []GroupAuth
	GetGroupAuthsByGroupID(op string, GroupID string) []GroupAuth
	CreateGroupAuthDirect(GroupID, AuthID uint64) error
	DeleteGroupAuthsByGroupID(GroupID string) error
	// GetRelGroupAuthToUGResources(librarys *[]Library, mapParams map[string]interface{}) error
}

// groupauthValidator checks and normalizes data prior to
// db access.
type groupauthValidator struct {
	GroupAuthDB
}

// groupauthValFunc type is the prototype for discrete GroupAuth normalization
// and validation functions that will be executed by func runGroupAuthValidationFuncs(...)
type groupauthValFunc func(*GroupAuth) error

// GroupAuthService is the public interface to the Usr entity
type GroupAuthService interface {
	GroupAuthDB
}

// private service for groupauth
type groupauthService struct {
	GroupAuthDB
}

// groupauthSqac is a sqac-based implementation of the GroupAuthDB interface.
type groupauthSqac struct {
	handle sqac.PublicDB
}

var _ GroupAuthDB = &groupauthSqac{}

// newGroupAuthValidator returns a new groupauthValidator
func newGroupAuthValidator(gdb GroupAuthDB) *groupauthValidator {
	return &groupauthValidator{
		GroupAuthDB: gdb,
	}
}

// runGroupAuthValFuncs executes a list of discrete validation
// functions against a groupauth.
func runGroupAuthValFuncs(groupauth *GroupAuth, fns ...groupauthValFunc) error {

	// iterate over the slice of function names and execute
	// each in-turn.  the order in which the lists are made
	// can matter...
	for _, fn := range fns {
		err := fn(groupauth)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewGroupAuthService ...
func NewGroupAuthService(handle sqac.PublicDB) GroupAuthService {

	gs := &groupauthSqac{handle}

	gv := newGroupAuthValidator(gs) // *db
	return &groupauthService{
		GroupAuthDB: gv,
	}
}

// ensure consistency (build error if delta exists)
var _ GroupAuthDB = &groupauthValidator{}

//-------------------------------------------------------------------------------------------------------
// CRUD-type model methods for GroupAuth
//-------------------------------------------------------------------------------------------------------
//
// Create validates and normalizes data used in the groupauth creation.
// Create then calls the creation code contained in GroupAuthService.
func (gv *groupauthValidator) Create(groupauth *GroupAuth) error {

	// perform normalization and validation -- comment out checks that are not required
	// note that the check calls are generated as a straight enumeration of the entity
	// structure.  It may be neccessary to adjust the calling order depending on the
	// relationships between the fields in the entity structure.
	err := runGroupAuthValFuncs(groupauth,
		gv.normvalGroupName,
	)

	if err != nil {
		return err
	}
	return gv.GroupAuthDB.Create(groupauth)
}

// Update validates and normalizes the content of the GroupAuth
// being updated by way of executing a list of predefined discrete
// checks.  if the checks are successful, the entity is updated
// on the db via the ORM.
func (gv *groupauthValidator) Update(groupauth *GroupAuth) error {

	// perform normalization and validation -- comment out checks that are not required
	// note that the check calls are generated as a straight enumeration of the entity
	// structure.  It may be neccessary to adjust the calling order depending on the
	// relationships between the fields in the entity structure.
	err := runGroupAuthValFuncs(groupauth,
		gv.normvalGroupName,
	)

	if err != nil {
		return err
	}
	return gv.GroupAuthDB.Update(groupauth)
}

// Delete is passed through to the ORM with no real
// validations.  id is checked in the controller.
func (gv *groupauthValidator) Delete(groupauth *GroupAuth) error {

	return gv.GroupAuthDB.Delete(groupauth)
}

// Get is passed through to the ORM with no real
// validations.  id is checked in the controller.
func (gv *groupauthValidator) Get(groupauth *GroupAuth) error {

	return gv.GroupAuthDB.Get(groupauth)
}

// GetGroupAuths is passed through to the ORM with no validation
func (gv *groupauthValidator) GetGroupAuths() []GroupAuth {

	return gv.GroupAuthDB.GetGroupAuths()
}

//-------------------------------------------------------------------------------------------------------
// internal groupauthValidator funcs
//-------------------------------------------------------------------------------------------------------
// These discrete functions are used to normalize and validate the Entity fields
// from with in the Create and Update methods.  See the comments in the model's
// Create and Update methods for details regarding use.

// normvalGroupName normalizes and validates field GroupName
func (gv *groupauthValidator) normvalGroupName(groupauth *GroupAuth) error {

	// TODO: implement normalization and validation for GroupAuth
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
// internal groupauth Simple Query Validator funcs
//-------------------------------------------------------------------------------------------------------
// Simple query normalization and validation occurs in the controller to an
// extent, as the URL has to be examined closely in order to determine what to call
// in the model.  This section may be blank if no model fields were marked as
// selectable in the <models>.json file.
// GetGroupAuthsByAuthName is passed through to the ORM with no validation.
func (gv *groupauthValidator) GetGroupAuthsByAuthName(op string, AuthName string) []GroupAuth {

	// TODO: implement normalization and validation for the GetGroupAuthsByAuthName call.
	// TODO: typically no modifications are required here.
	return gv.GroupAuthDB.GetGroupAuthsByAuthName(op, AuthName)
}

// GetGroupAuthsByDescription is passed through to the ORM with no validation.
func (gv *groupauthValidator) GetGroupAuthsByDescription(op string, Description string) []GroupAuth {

	// TODO: implement normalization and validation for the GetGroupAuthsByDescription call.
	// TODO: typically no modifications are required here.
	return gv.GroupAuthDB.GetGroupAuthsByDescription(op, Description)
}

// GetGroupAuthsByGroupID is passed through to the ORM with no validation.
func (gv *groupauthValidator) GetGroupAuthsByGroupID(op string, GroupID string) []GroupAuth {

	// TODO: implement normalization and validation for the GetGroupAuthsByGroupID call.
	// TODO: typically no modifications are required here.
	return gv.GroupAuthDB.GetGroupAuthsByGroupID(op, GroupID)
}

//-------------------------------------------------------------------------------------------------------
// ORM db CRUD access methods
//-------------------------------------------------------------------------------------------------------
//
// Create a new GroupAuth in the database via the ORM
func (gs *groupauthSqac) Create(groupauth *GroupAuth) error {
	return gs.handle.Create(groupauth)
}

// Update an existng GroupAuth in the database via the ORM
func (gs *groupauthSqac) Update(groupauth *GroupAuth) error {
	return gs.handle.Update(groupauth)
}

// Delete an existing GroupAuth in the database via the ORM
func (gs *groupauthSqac) Delete(groupauth *GroupAuth) error {
	return gs.handle.Delete(groupauth)
}

// Get an existing GroupAuth from the database via the ORM
func (gs *groupauthSqac) Get(groupauth *GroupAuth) error {

	// do not call the standard sqac CRUD Get(), as it would be
	// nice to provide the caller with some text from the resource
	c := fmt.Sprintf("groupauth.id = ?")

	qs := fmt.Sprintf(`SELECT groupauth.id, groupauth.group_id, usrgroup.group_name, groupauth.auth_id, auth.auth_name, auth.auth_type, auth.description
 		FROM groupauth INNER JOIN auth ON (groupauth.auth_id = auth.id) 
 		               INNER JOIN usrgroup ON (usrgroup.id = groupauth.group_id) WHERE %s;`, c)

	// read the table row - can only be one due to db key contraint
	rows, err := gs.handle.ExecuteQueryx(qs, groupauth.ID)
	if err != nil {
		lw.Warning("Get for table groupauth returned error: %v", err.Error())
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.StructScan(groupauth)
		if err != nil {
			lw.ErrorWithPrefixString("error reading GroupAuth rows via StructScan:", err)
			return err
		}
	}

	if gs.handle.IsLog() {
		lw.Debug("Get for table groupauth got: %v", groupauth)
	}

	if groupauth.GroupID == 0 || groupauth.AuthID == 0 {
		return fmt.Errorf("groupauth %d not found", groupauth.ID)
	}
	return nil
}

// Get all existing GroupAuths from the db via the ORM
func (gs *groupauthSqac) GetGroupAuths() []GroupAuth {

	groupauth := GroupAuth{}
	groupauths := []GroupAuth{}

	qs := fmt.Sprintf(`SELECT groupauth.id, groupauth.group_id, usrgroup.group_name, groupauth.auth_id, auth.auth_name, auth.auth_type, auth.description
 		FROM groupauth INNER JOIN auth ON (groupauth.auth_id = auth.id) 
 		               INNER JOIN usrgroup ON (usrgroup.id = groupauth.group_id);`)

	// read the table row - can only be one due to db key contraint
	rows, err := gs.handle.ExecuteQueryx(qs)
	if err != nil {
		lw.ErrorWithPrefixString("Get for table groupauth returned error: ", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.StructScan(&groupauth)
		if err != nil {
			lw.ErrorWithPrefixString("error reading GroupAuth rows via StructScan: ", err)
			return nil
		}
		groupauths = append(groupauths, groupauth)
	}

	if gs.handle.IsLog() {
		lw.Debug("Get for table groupauth got: %v", groupauths)
	}
	return groupauths
}

//- - - - - - - - - - - - - - - - - - -- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
//  Implement support for the sqac.GetEnt{} interface to support GetEntities2 calls
//- - - - - - - - - - - - - - - - - - -- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -

// GroupAuthGetEntitiesT is used to support the implmentation of the sqac.GetEnt{} interface.
type GroupAuthGetEntitiesT struct {
	ents []GroupAuth
	sqac.GetEnt
}

// Exec implments the one and only method in the  sqac.GetEnt{} interface for GroupAuthGetEntitiesT.
// An implemented interface is passed to sqac.GetEntities2 and the code contained withing the Exec()
// method is executed to retrieve the entities from the DB.  Retrieved entities are stored in the
// GroupAuthGetEntitiesT struct's ents []GroupAuth field.
func (ge *GroupAuthGetEntitiesT) Exec(sqh sqac.PublicDB) error {

	selQuery := `SELECT groupauth.id, groupauth.group_id, usrgroup.group_name, groupauth.auth_id, auth.auth_name, auth.auth_type, auth.description
 	FROM groupauth INNER JOIN auth ON (groupauth.auth_id = auth.id) 
 				   INNER JOIN usrgroup ON (usrgroup.id = groupauth.group_id);`

	// read the table rows
	rows, err := sqh.ExecuteQueryx(selQuery)
	if err != nil {
		lw.Warning("GetAuthEntitiesT for table groupauth returned error: %s", err.Error())
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var ent GroupAuth
		err = rows.StructScan(&ent)
		if err != nil {
			lw.ErrorWithPrefixString("error reading GroupAuth rows via StructScan: ", err)
			return err
		}
		ge.ents = append(ge.ents, ent)
	}
	if sqh.IsLog() {
		lw.Debug("GroupAuthGetEntitiesT.Exec() got: %v", ge.ents)
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
// Get all existing GroupAuthsByAuthName from the db via the ORM
func (gs *groupauthSqac) GetGroupAuthsByAuthName(op string, AuthName string) []GroupAuth {

	var groupauths []GroupAuth
	var c string

	switch op {
	case "EQ":
		c = "auth_name = ?"
	case "LIKE":
		c = "auth_name like ?"
	default:
		return nil
	}
	qs := fmt.Sprintf(`SELECT groupauth.id, groupauth.group_id, usrgroup.group_name, groupauth.auth_id, auth.auth_name, auth.auth_type, auth.description
 		FROM groupauth INNER JOIN auth ON (groupauth.auth_id = auth.id) 
 		               INNER JOIN usrgroup ON (usrgroup.id = groupauth.group_id) WHERE %s;`, c)

	err := gs.handle.Select(&groupauths, qs, AuthName)
	if err != nil {
		lw.Warning("GetGroupAuthsByAuthName got: %s", err.Error())
		return nil
	}

	if gs.handle.IsLog() {
		lw.Debug("GetGroupAuthsByAuthName found: %v based on (%s %v)", groupauths, op, AuthName)
	}
	return groupauths
}

// GetGroupAuthsByDescription gets all existing GroupAuthsByDescription from the db via the ORM
func (gs *groupauthSqac) GetGroupAuthsByDescription(op string, Description string) []GroupAuth {

	var groupauths []GroupAuth
	var c string

	switch op {
	case "EQ":
		c = "description = ?"
	case "LIKE":
		c = "description like ?"
	default:
		return nil
	}

	qs := fmt.Sprintf(`SELECT groupauth.id, groupauth.group_id, usrgroup.group_name, groupauth.auth_id, auth.auth_name, auth.auth_type, auth.description
 		FROM groupauth INNER JOIN auth ON (groupauth.auth_id = auth.id) 
 		               INNER JOIN usrgroup ON (usrgroup.id = groupauth.group_id) WHERE %s;`, c)

	err := gs.handle.Select(&groupauths, qs, Description)
	if err != nil {
		lw.Warning("GetGroupAuthsByDescription got: %s", err.Error())
		return nil
	}

	if gs.handle.IsLog() {
		lw.Debug("GetGroupAuthsByDescription found: %v based on (%s %v)", groupauths, op, Description)
	}
	return groupauths
}

// GetGroupAuthsByGroupID gets all existing GroupAuthsByGroupID from the db via the ORM
func (gs *groupauthSqac) GetGroupAuthsByGroupID(op string, GroupID string) []GroupAuth {

	var groupauths []GroupAuth
	var c string

	switch op {
	case "EQ":
		c = "group_id = ?"
	case "LT":
		c = "group_id < ?"
	case "GT":
		c = "group_id > ?"
	default:
		return nil
	}

	qs := fmt.Sprintf(`SELECT groupauth.id, groupauth.group_id, usrgroup.group_name, groupauth.auth_id, auth.auth_name, auth.auth_type, auth.description
 		FROM groupauth INNER JOIN auth ON (groupauth.auth_id = auth.id) 
 		               INNER JOIN usrgroup ON (usrgroup.id = groupauth.group_id) WHERE %s;`, c)

	err := gs.handle.Select(&groupauths, qs, GroupID)
	if err != nil {
		lw.Warning("GetGroupAuthsByGroupID got: %s", err.Error())
		return nil
	}

	if gs.handle.IsLog() {
		lw.Debug("GetGroupAuthsByGroupID found: %v based on (%s %v)", groupauths, op, GroupID)
	}
	return groupauths
}

// CreateGroupAuthDirect facilitates the creation of an Auth allocation to a UsrGroup
// from a source other than an incoming http request.  For example, this method is called
// as a part of the Super Group initialization on application startup.
func (gs *groupauthSqac) CreateGroupAuthDirect(GroupID, AuthID uint64) error {

	if GroupID == 0 || AuthID == 0 {
		return fmt.Errorf("GroupID and AuthID must be specfied when allocating and Auth to a UsrGroup")
	}

	// fill the model
	groupauth := GroupAuth{
		GroupID: GroupID,
		AuthID:  AuthID,
	}

	// call the Create method on the groupauth model
	err := gs.Create(&groupauth)
	if err != nil {
		lw.ErrorWithPrefixString("GroupAuth Create got:", err)
		return err
	}
	return nil
}

// DeleteGroupAuthsByGroupID performs the deletion of all assigned Auths
// to the specified UsrGroup.
func (gs *groupauthSqac) DeleteGroupAuthsByGroupID(GroupID string) error {

	ds := fmt.Sprintf("DELETE FROM groupauth WHERE group_id = ?;")

	// delete the requested rows
	_, err := gs.handle.ExecuteQueryx(ds, GroupID)
	if err != nil {
		return err
	}
	return nil
}
