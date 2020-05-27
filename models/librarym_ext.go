package models

import (
	// "fmt"
	// "reflect"
	"github.com/1414C/libraryapp/models/ext"
)

//===================================================================================================
// base Library entity model extension-point code generated on 27 May 20 17:32 CDT
//===================================================================================================

// MdlLibraryCreateExt provides access to the ModelCreateExt extension-point interface
type MdlLibraryCreateExt struct {
	ext.ModelCreateExt
}

// MdlLibraryUpdateExt provides access to the ControllerUpdateExt extension-point interface
type MdlLibraryUpdateExt struct {
	ext.ModelUpdateExt
}

// MdlLibraryGetExt provides access to the ControllerGetExt extension-point interface
type MdlLibraryGetExt struct {
	ext.ModelGetExt
}

// LibraryMdlExt provides access to the Library implementations of the following interfaces:
//   MdlCreateExt
//   MdlUpdateExt
//   MdlGetExt
type LibraryMdlExt struct {
	CrtEp MdlLibraryCreateExt
	UpdEp MdlLibraryUpdateExt
	GetEp MdlLibraryGetExt
}

var libraryMdlExp LibraryMdlExt

// InitLibraryMdlExt initializes the library entity's model
// extension-point interface implementations.
func InitLibraryMdlExt() *LibraryMdlExt {
	libraryMdlExp = LibraryMdlExt{}
	return &libraryMdlExp
}

//----------------------------------------------------------------------------
// ModelCreateExt interface implementation for entity Library
//----------------------------------------------------------------------------

// BeforeDB model extension-point implementation for entity Library
// TODO: implement pre-ORM call logic and document it here
func (crtEP *MdlLibraryCreateExt) BeforeDB(ent interface{}) error {

	// fmt.Println("TypeOf ent:", reflect.TypeOf(ent))
	// fmt.Println("ValueOf ent:", reflect.ValueOf(ent))
	// l := ent.(*models.Library)

	// make changes / validate the content struct pointer (l) here
	// l.Name = "A new field value"
	return nil
}

// AfterDB model extension-point implementation for entity Library
// TODO: implement post-ORM call logic and document it here
func (crtEP *MdlLibraryCreateExt) AfterDB(ent interface{}) error {

	// fmt.Println("TypeOf ent:", reflect.TypeOf(ent))
	// fmt.Println("ValueOf ent:", reflect.ValueOf(ent))
	// l := ent.(*models.Library)

	// make changes / validate the content struct pointer (l) here
	// l.Name = "A new field value"
	return nil
}

//----------------------------------------------------------------------------
// ModelUpdateExt interface implementation for entity Library
//----------------------------------------------------------------------------

// BeforeDB extension-point implementation for entity Library
// TODO: implement pre-ORM call logic and document it here
func (updEP *MdlLibraryUpdateExt) BeforeDB(ent interface{}) error {

	// fmt.Println("TypeOf ent:", reflect.TypeOf(ent))
	// fmt.Println("ValueOf ent:", reflect.ValueOf(ent))
	// l := ent.(*models.Library)

	// make changes / validate the content struct pointer (l) here
	// l.Name = "A new field value"
	return nil
}

// AfterDB extension-point implementation for entity Library
// TODO: implement post-ORM call logic and document it here
func (updEP *MdlLibraryUpdateExt) AfterDB(ent interface{}) error {

	// fmt.Println("TypeOf ent:", reflect.TypeOf(ent))
	// fmt.Println("ValueOf ent:", reflect.ValueOf(ent))
	// l := ent.(*models.Library)

	// make changes / validate the content struct pointer (l) here
	// l.Name = "A new field value"
	return nil
}

//----------------------------------------------------------------------------
// ModelGetExt interface implementation for entity Library
//----------------------------------------------------------------------------

// BeforeDB extension-point implementation for entity Library
// TODO: implement pre-ORM call logic and document it here
func (getEP *MdlLibraryGetExt) BeforeDB(ent interface{}) error {

	// fmt.Println("TypeOf ent:", reflect.TypeOf(ent))
	// fmt.Println("ValueOf ent:", reflect.ValueOf(ent))
	// l := ent.(*models.Library)

	// make changes / validate the content struct pointer (l) here
	// l.Name = "A new field value"
	return nil
}

// AfterDB extension-point implementation for entity Library
// TODO: implement post-ORM call logic and document it here
func (getEP *MdlLibraryGetExt) AfterDB(ent interface{}) error {

	// fmt.Println("TypeOf ent:", reflect.TypeOf(ent))
	// fmt.Println("ValueOf ent:", reflect.ValueOf(ent))
	// l := ent.(*Library)

	// make changes / validate the content struct pointer (l) here
	// l.Name = "A new field value"
	return nil
}
