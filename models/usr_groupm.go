package models

//=============================================================================================
// base UsrGroup entity model code generated on 27 May 20 17:57 CDT
//=============================================================================================

import (
	"fmt"

	"github.com/1414C/lw"
	"github.com/1414C/sqac"
)

// UsrGroup structure
type UsrGroup struct {
	ID          uint64 `json:"id" db:"id" sqac:"primary_key:inc"`
	Href        string `json:"href" db:"href" sqac:"-"`
	GroupName   string `json:"group_name" db:"group_name" sqac:"nullable:false;index:non-unique"`
	Description string `json:"description" db:"description" sqac:"nullable:false"`
}

// UsrGroupDB is a CRUD-type interface specifically for dealing with UsrGroups.
type UsrGroupDB interface {
	Create(usrgroup *UsrGroup) error
	Update(usrgroup *UsrGroup) error
	Delete(usrgroup *UsrGroup) error
	Get(usrgroup *UsrGroup) error
	GetUsrGroups() []UsrGroup
	GetUsrGroupsByGroupName(op string, GroupName string) []UsrGroup
	GetUsrGroupsByDescription(op string, Description string) []UsrGroup
	// GetRelUsrGroupToUGResources(librarys *[]Library, mapParams map[string]interface{}) error
}

// usrgroupValidator checks and normalizes data prior to
// db access.
type usrgroupValidator struct {
	UsrGroupDB
}

// usrgroupValFunc type is the prototype for discrete UsrGroup normalization
// and validation functions that will be executed by func runUsrGroupValidationFuncs(...)
type usrgroupValFunc func(*UsrGroup) error

// UsrGroupService is the public interface to the Usr entity
type UsrGroupService interface {
	UsrGroupDB
}

// private service for usrgroup
type usrgroupService struct {
	UsrGroupDB
}

// usrgroupSqac is a sqac-based implementation of the UsrGroupDB interface.
type usrgroupSqac struct {
	handle sqac.PublicDB
}

var _ UsrGroupDB = &usrgroupSqac{}

// newUsrGroupValidator returns a new usrgroupValidator
func newUsrGroupValidator(udb UsrGroupDB) *usrgroupValidator {
	return &usrgroupValidator{
		UsrGroupDB: udb,
	}
}

// runUsrGroupValFuncs executes a list of discrete validation
// functions against a usrgroup.
func runUsrGroupValFuncs(usrgroup *UsrGroup, fns ...usrgroupValFunc) error {

	// iterate over the slice of function names and execute
	// each in-turn.  the order in which the lists are made
	// can matter...
	for _, fn := range fns {
		err := fn(usrgroup)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewUsrGroupService ...
func NewUsrGroupService(handle sqac.PublicDB) UsrGroupService {

	us := &usrgroupSqac{handle}

	uv := newUsrGroupValidator(us) // *db
	return &usrgroupService{
		UsrGroupDB: uv,
	}
}

// ensure consistency (build error if delta exists)
var _ UsrGroupDB = &usrgroupValidator{}

//-------------------------------------------------------------------------------------------------------
// CRUD-type model methods for UsrGroup
//-------------------------------------------------------------------------------------------------------
//
// Create validates and normalizes data used in the usrgroup creation.
// Create then calls the creation code contained in UsrGroupService.
func (uv *usrgroupValidator) Create(usrgroup *UsrGroup) error {

	// perform normalization and validation -- comment out checks that are not required
	// note that the check calls are generated as a straight enumeration of the entity
	// structure.  It may be neccessary to adjust the calling order depending on the
	// relationships between the fields in the entity structure.
	err := runUsrGroupValFuncs(usrgroup,
		uv.normvalGroupName,
	)

	if err != nil {
		return err
	}
	return uv.UsrGroupDB.Create(usrgroup)
}

// Update validates and normalizes the content of the UsrGroup
// being updated by way of executing a list of predefined discrete
// checks.  if the checks are successful, the entity is updated
// on the db via the ORM.
func (uv *usrgroupValidator) Update(usrgroup *UsrGroup) error {

	// perform normalization and validation -- comment out checks that are not required
	// note that the check calls are generated as a straight enumeration of the entity
	// structure.  It may be neccessary to adjust the calling order depending on the
	// relationships between the fields in the entity structure.
	err := runUsrGroupValFuncs(usrgroup,
		uv.normvalGroupName,
	)

	if err != nil {
		return err
	}
	return uv.UsrGroupDB.Update(usrgroup)
}

// Delete is passed through to the ORM with no real
// validations.  id is checked in the controller.
func (uv *usrgroupValidator) Delete(usrgroup *UsrGroup) error {

	return uv.UsrGroupDB.Delete(usrgroup)
}

// Get is passed through to the ORM with no real
// validations.  id is checked in the controller.
func (uv *usrgroupValidator) Get(usrgroup *UsrGroup) error {

	return uv.UsrGroupDB.Get(usrgroup)
}

// GetUsrGroups is passed through to the ORM with no validation
func (uv *usrgroupValidator) GetUsrGroups() []UsrGroup {

	return uv.UsrGroupDB.GetUsrGroups()
}

//-------------------------------------------------------------------------------------------------------
// internal usrgroupValidator funcs
//-------------------------------------------------------------------------------------------------------
// These discrete functions are used to normalize and validate the Entity fields
// from with in the Create and Update methods.  See the comments in the model's
// Create and Update methods for details regarding use.

// normvalGroupName normalizes and validates field GroupName
func (uv *usrgroupValidator) normvalGroupName(usrgroup *UsrGroup) error {

	// TODO: implement normalization and validation for UsrGroup
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
// internal usrgroup Simple Query Validator funcs
//-------------------------------------------------------------------------------------------------------
// Simple query normalization and validation occurs in the controller to an
// extent, as the URL has to be examined closely in order to determine what to call
// in the model.  This section may be blank if no model fields were marked as
// selectable in the <models>.json file.
// GetUsrGroupsByGroupName is passed through to the ORM with no validation.
func (uv *usrgroupValidator) GetUsrGroupsByGroupName(op string, GroupName string) []UsrGroup {

	// TODO: implement normalization and validation for the GetUsrGroupsByGroupName call.
	// TODO: typically no modifications are required here.
	return uv.UsrGroupDB.GetUsrGroupsByGroupName(op, GroupName)
}

// GetUsrGroupsByDesription is passed through to the ORM with no validation.
func (uv *usrgroupValidator) GetUsrGroupsByDescription(op string, Description string) []UsrGroup {

	// TODO: implement normalization and validation for the GetUsrGroupsByDescription call.
	// TODO: typically no modifications are required here.
	return uv.UsrGroupDB.GetUsrGroupsByDescription(op, Description)
}

//-------------------------------------------------------------------------------------------------------
// ORM db CRUD access methods
//-------------------------------------------------------------------------------------------------------
//
// Create a new UsrGroup in the database via the ORM
func (us *usrgroupSqac) Create(usrgroup *UsrGroup) error {
	return us.handle.Create(usrgroup)
}

// Update an existng UsrGroup in the database via the ORM
func (us *usrgroupSqac) Update(usrgroup *UsrGroup) error {
	return us.handle.Update(usrgroup)
}

// Delete an existing UsrGroup in the database via the ORM
func (us *usrgroupSqac) Delete(usrgroup *UsrGroup) error {
	return us.handle.Delete(usrgroup)
}

// Get an existing UsrGroup from the database via the ORM
func (us *usrgroupSqac) Get(usrgroup *UsrGroup) error {
	return us.handle.GetEntity(usrgroup)
}

// Get all existing UsrGroups from the db via the ORM
func (us *usrgroupSqac) GetUsrGroups() []UsrGroup {

	getEnts := UsrGroupGetEntitiesT{}

	err := us.handle.GetEntities2(&getEnts)
	if err != nil {
		lw.Warning("GetUsrGroups got: %s", err.Error())
		return getEnts.ents
	}
	return getEnts.ents
}

//- - - - - - - - - - - - - - - - - - -- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
//  Implement support for the sqac.GetEnt{} interface to support GetEntities2 calls
//- - - - - - - - - - - - - - - - - - -- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -

// UsrGroupGetEntitiesT is used to support the implmentation of the sqac.GetEnt{} interface.
type UsrGroupGetEntitiesT struct {
	ents []UsrGroup
	sqac.GetEnt
}

// Exec implments the one and only method in the  sqac.GetEnt{} interface for UsrGroupGetEntitiesT.
// An implemented interface is passed to sqac.GetEntities2 and the code contained withing the Exec()
// method is executed to retrieve the entities from the DB.  Retrieved entities are stored in the
// UsrGroupGetEntitiesT struct's ents []UsrGroup field.
func (ge *UsrGroupGetEntitiesT) Exec(sqh sqac.PublicDB) error {

	selQuery := "SELECT * FROM usrgroup;"

	// read the table rows
	rows, err := sqh.ExecuteQueryx(selQuery)
	if err != nil {
		lw.Warning("GetEntities2 for table usrgroup returned error: %s", err.Error())
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var ent UsrGroup
		err = rows.StructScan(&ent)
		if err != nil {
			lw.ErrorWithPrefixString("error reading UsrGroup rows via StructScan:", err)
			return err
		}
		ge.ents = append(ge.ents, ent)
	}
	if sqh.IsLog() {
		lw.Info("UsrGroupGetEntitiesT.Exec() got:", ge.ents)
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
// Get all existing UsrGroupsByGroupName from the db via the ORM
func (us *usrgroupSqac) GetUsrGroupsByGroupName(op string, GroupName string) []UsrGroup {

	var usrgroups []UsrGroup
	var c string

	switch op {
	case "EQ":
		c = "group_name = ?"
	case "LIKE":
		c = "group_name like ?"
	default:
		return nil
	}
	qs := fmt.Sprintf("SELECT * FROM usrgroup WHERE %s;", c)
	err := us.handle.Select(&usrgroups, qs, GroupName)
	if err != nil {
		lw.Warning("GetUsrGroupsByGroupName got: %s", err.Error())
		return nil
	}

	if us.handle.IsLog() {
		lw.Info("GetUsrGroupsByGroupName found: %v based on (%s %v)", usrgroups, op, GroupName)
	}
	return usrgroups
}

// GetUsrGroupsByDescription gets all existing UsrGroupsByDescription from the db via the ORM
func (us *usrgroupSqac) GetUsrGroupsByDescription(op string, Description string) []UsrGroup {

	var usrgroups []UsrGroup
	var c string

	switch op {
	case "EQ":
		c = "description = ?"
	case "LIKE":
		c = "description like ?"
	default:
		return nil
	}
	qs := fmt.Sprintf("SELECT * FROM usrgroup WHERE %s;", c)
	err := us.handle.Select(&usrgroups, qs, Description)
	if err != nil {
		lw.Warning("GetUsrGroupsByDescription got: %s", err.Error())
		return nil
	}

	if us.handle.IsLog() {
		lw.Info("GetUsrGroupsByDescription found: %v based on (%s %v)", usrgroups, op, Description)
	}
	return usrgroups
}
