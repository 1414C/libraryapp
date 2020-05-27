package models

//=============================================================================================
// base Library entity model code generated on 27 May 20 17:32 CDT
//=============================================================================================

import (
	"fmt"

	"github.com/1414C/lw"
	"github.com/1414C/sqac"
)

// Library structure
type Library struct {
	ID   uint64 `json:"id" db:"id" sqac:"primary_key:inc"`
	Href string `json:"href" db:"href" sqac:"-"`
	Name string `json:"name" db:"name" sqac:"nullable:false;index:non-unique;index:idx_library_name_city"`
	City string `json:"city" db:"city" sqac:"nullable:false;index:idx_library_name_city"`
}

// LibraryDB is a CRUD-type interface specifically for dealing with Librarys.
type LibraryDB interface {
	Create(library *Library) error
	Update(library *Library) error
	Delete(library *Library) error
	Get(library *Library) error
	GetLibrarys(params []sqac.GetParam, cmdMap map[string]interface{}) ([]Library, uint64) // uint64 holds $count result
	GetLibrarysByName(op string, Name string) []Library
	GetLibrarysByCity(op string, City string) []Library
}

// libraryValidator checks and normalizes data prior to
// db access.
type libraryValidator struct {
	LibraryDB
}

// libraryValFunc type is the prototype for discrete Library normalization
// and validation functions that will be executed by func runLibraryValidationFuncs(...)
type libraryValFunc func(*Library) error

// LibraryService is the public interface to the Library entity
type LibraryService interface {
	LibraryDB
}

// private service for library
type libraryService struct {
	LibraryDB
}

// librarySqac is a sqac-based implementation of the LibraryDB interface.
type librarySqac struct {
	handle sqac.PublicDB
	ep     LibraryMdlExt
}

var _ LibraryDB = &librarySqac{}

// newLibraryValidator returns a new libraryValidator
func newLibraryValidator(ldb LibraryDB) *libraryValidator {
	return &libraryValidator{
		LibraryDB: ldb,
	}
}

// runLibraryValFuncs executes a list of discrete validation
// functions against a library.
func runLibraryValFuncs(library *Library, fns ...libraryValFunc) error {

	// iterate over the slice of function names and execute
	// each in-turn.  the order in which the lists are made
	// can matter...
	for _, fn := range fns {
		err := fn(library)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewLibraryService needs some work:
func NewLibraryService(handle sqac.PublicDB) LibraryService {

	ls := &librarySqac{
		handle: handle,
		ep:     *InitLibraryMdlExt(),
	}

	lv := newLibraryValidator(ls) // *db
	return &libraryService{
		LibraryDB: lv,
	}
}

// ensure consistency (build error if delta exists)
var _ LibraryDB = &libraryValidator{}

//-------------------------------------------------------------------------------------------------------
// CRUD-type model methods for Library
//-------------------------------------------------------------------------------------------------------
//
// Create validates and normalizes data used in the library creation.
// Create then calls the creation code contained in LibraryService.
func (lv *libraryValidator) Create(library *Library) error {

	// perform normalization and validation -- comment out checks that are not required
	// note that the check calls are generated as a straight enumeration of the entity
	// structure.  It may be neccessary to adjust the calling order depending on the
	// relationships between the fields in the entity structure.
	err := runLibraryValFuncs(library,
		lv.normvalName,
		lv.normvalCity,
	)

	if err != nil {
		return err
	}
	return lv.LibraryDB.Create(library)
}

// Update validates and normalizes the content of the Library
// being updated by way of executing a list of predefined discrete
// checks.  if the checks are successful, the entity is updated
// on the db via the ORM.
func (lv *libraryValidator) Update(library *Library) error {

	// perform normalization and validation -- comment out checks that are not required
	// note that the check calls are generated as a straight enumeration of the entity
	// structure.  It may be neccessary to adjust the calling order depending on the
	// relationships between the fields in the entity structure.
	err := runLibraryValFuncs(library,
		lv.normvalName,
		lv.normvalCity,
	)

	if err != nil {
		return err
	}
	return lv.LibraryDB.Update(library)
}

// Delete is passed through to the ORM with no real
// validations.  id is checked in the controller.
func (lv *libraryValidator) Delete(library *Library) error {

	return lv.LibraryDB.Delete(library)
}

// Get is passed through to the ORM with no real
// validations.  id is checked in the controller.
func (lv *libraryValidator) Get(library *Library) error {

	return lv.LibraryDB.Get(library)
}

// GetLibrarys is passed through to the ORM with no validation
func (lv *libraryValidator) GetLibrarys(params []sqac.GetParam, cmdMap map[string]interface{}) ([]Library, uint64) {

	return lv.LibraryDB.GetLibrarys(params, cmdMap)
}

//-------------------------------------------------------------------------------------------------------
// internal libraryValidator funcs
//-------------------------------------------------------------------------------------------------------
// These discrete functions are used to normalize and validate the Entity fields
// from with in the Create and Update methods.  See the comments in the model's
// Create and Update methods for details regarding use.

// normvalName normalizes and validates field Name
func (lv *libraryValidator) normvalName(library *Library) error {

	// TODO: implement normalization and validation for Name
	return nil
}

// normvalCity normalizes and validates field City
func (lv *libraryValidator) normvalCity(library *Library) error {

	// TODO: implement normalization and validation for City
	return nil
}

//-------------------------------------------------------------------------------------------------------
// internal library Simple Query Validator funcs
//-------------------------------------------------------------------------------------------------------
// Simple query normalization and validation occurs in the controller to an
// extent, as the URL has to be examined closely in order to determine what to call
// in the model.  This section may be blank if no model fields were marked as
// selectable in the <models>.json file.
// GetLibrarysByName is passed through to the ORM with no validation.
func (lv *libraryValidator) GetLibrarysByName(op string, name string) []Library {

	// TODO: implement normalization and validation for the GetLibrarysByName call.
	// TODO: typically no modifications are required here.
	return lv.LibraryDB.GetLibrarysByName(op, name)
}

// GetLibrarysByCity is passed through to the ORM with no validation.
func (lv *libraryValidator) GetLibrarysByCity(op string, city string) []Library {

	// TODO: implement normalization and validation for the GetLibrarysByCity call.
	// TODO: typically no modifications are required here.
	return lv.LibraryDB.GetLibrarysByCity(op, city)
}

//-------------------------------------------------------------------------------------------------------
// ORM db CRUD access methods
//-------------------------------------------------------------------------------------------------------
//
// Create a new Library in the database via the ORM
func (ls *librarySqac) Create(library *Library) error {

	err := ls.ep.CrtEp.BeforeDB(library)
	if err != nil {
		return err
	}
	err = ls.handle.Create(library)
	if err != nil {
		return err
	}
	err = ls.ep.CrtEp.AfterDB(library)
	if err != nil {
		return err
	}
	return err
}

// Update an existng Library in the database via the ORM
func (ls *librarySqac) Update(library *Library) error {

	err := ls.ep.UpdEp.BeforeDB(library)
	if err != nil {
		return err
	}
	err = ls.handle.Update(library)
	if err != nil {
		return err
	}
	err = ls.ep.UpdEp.AfterDB(library)
	if err != nil {
		return err
	}
	return err
}

// Delete an existing Library in the database via the ORM
func (ls *librarySqac) Delete(library *Library) error {
	return ls.handle.Delete(library)
}

// Get an existing Library from the database via the ORM
func (ls *librarySqac) Get(library *Library) error {

	err := ls.ep.GetEp.BeforeDB(library)
	if err != nil {
		return err
	}
	err = ls.handle.GetEntity(library)
	if err != nil {
		return err
	}
	err = ls.ep.GetEp.AfterDB(library)
	if err != nil {
		return err
	}
	return err
}

// Get all existing Librarys from the db via the ORM
func (ls *librarySqac) GetLibrarys(params []sqac.GetParam, cmdMap map[string]interface{}) ([]Library, uint64) {

	var err error

	// create a slice to read into
	librarys := []Library{}

	// call the ORM
	result, err := ls.handle.GetEntitiesWithCommands(librarys, params, cmdMap)
	if err != nil {
		lw.Warning("LibraryModel GetLibrarys() error: %s", err.Error())
		return nil, 0
	}

	// check to see what was returned
	switch result.(type) {
	case []Library:
		librarys = result.([]Library)

		// call the extension-point
		for i := range librarys {
			err = ls.ep.GetEp.AfterDB(&librarys[i])
			if err != nil {
				lw.Warning("LibraryModel GetLibrarys AfterDB() error: %s", err.Error())
			}
		}
		return librarys, 0

	case int64:
		return nil, uint64(result.(int64))

	case uint64:
		return nil, result.(uint64)

	default:
		return nil, 0

	}
}

//-------------------------------------------------------------------------------------------------------
// ORM db simple selector access methods
//-------------------------------------------------------------------------------------------------------
//
// Get all existing LibrarysByName from the db via the ORM
func (ls *librarySqac) GetLibrarysByName(op string, Name string) []Library {

	var librarys []Library
	var c string

	switch op {
	case "EQ":
		c = "name = ?"
	case "LIKE":
		c = "name like ?"
	default:
		return nil
	}
	qs := fmt.Sprintf("SELECT * FROM library WHERE %s;", c)
	err := ls.handle.Select(&librarys, qs, Name)
	if err != nil {
		lw.Warning("GetLibrarysByName got: %s", err.Error())
		return nil
	}

	if ls.handle.IsLog() {
		lw.Info("GetLibrarysByName found: %v based on (%s %v)", librarys, op, Name)
	}

	// call the extension-point
	for i := range librarys {
		err = ls.ep.GetEp.AfterDB(&librarys[i])
		if err != nil {
			lw.Warning("LibraryModel Getlibrarys AfterDB() error: %s", err.Error())
		}
	}
	return librarys
}

// Get all existing LibrarysByCity from the db via the ORM
func (ls *librarySqac) GetLibrarysByCity(op string, City string) []Library {

	var librarys []Library
	var c string

	switch op {
	case "EQ":
		c = "city = ?"
	case "LIKE":
		c = "city like ?"
	default:
		return nil
	}
	qs := fmt.Sprintf("SELECT * FROM library WHERE %s;", c)
	err := ls.handle.Select(&librarys, qs, City)
	if err != nil {
		lw.Warning("GetLibrarysByCity got: %s", err.Error())
		return nil
	}

	if ls.handle.IsLog() {
		lw.Info("GetLibrarysByCity found: %v based on (%s %v)", librarys, op, City)
	}

	// call the extension-point
	for i := range librarys {
		err = ls.ep.GetEp.AfterDB(&librarys[i])
		if err != nil {
			lw.Warning("LibraryModel Getlibrarys AfterDB() error: %s", err.Error())
		}
	}
	return librarys
}
