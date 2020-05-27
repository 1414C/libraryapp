package models

import (
	// "fmt"
	// "reflect"
	"github.com/1414C/libraryapp/models/ext"
)

//===================================================================================================
// base Book entity model extension-point code generated on 27 May 20 17:32 CDT
//===================================================================================================

// MdlBookCreateExt provides access to the ModelCreateExt extension-point interface
type MdlBookCreateExt struct {
	ext.ModelCreateExt
}

// MdlBookUpdateExt provides access to the ControllerUpdateExt extension-point interface
type MdlBookUpdateExt struct {
	ext.ModelUpdateExt
}

// MdlBookGetExt provides access to the ControllerGetExt extension-point interface
type MdlBookGetExt struct {
	ext.ModelGetExt
}

// BookMdlExt provides access to the Book implementations of the following interfaces:
//   MdlCreateExt
//   MdlUpdateExt
//   MdlGetExt
type BookMdlExt struct {
	CrtEp MdlBookCreateExt
	UpdEp MdlBookUpdateExt
	GetEp MdlBookGetExt
}

var bookMdlExp BookMdlExt

// InitBookMdlExt initializes the book entity's model
// extension-point interface implementations.
func InitBookMdlExt() *BookMdlExt {
	bookMdlExp = BookMdlExt{}
	return &bookMdlExp
}

//----------------------------------------------------------------------------
// ModelCreateExt interface implementation for entity Book
//----------------------------------------------------------------------------

// BeforeDB model extension-point implementation for entity Book
// TODO: implement pre-ORM call logic and document it here
func (crtEP *MdlBookCreateExt) BeforeDB(ent interface{}) error {

	// fmt.Println("TypeOf ent:", reflect.TypeOf(ent))
	// fmt.Println("ValueOf ent:", reflect.ValueOf(ent))
	// b := ent.(*models.Book)

	// make changes / validate the content struct pointer (b) here
	// b.Name = "A new field value"
	return nil
}

// AfterDB model extension-point implementation for entity Book
// TODO: implement post-ORM call logic and document it here
func (crtEP *MdlBookCreateExt) AfterDB(ent interface{}) error {

	// fmt.Println("TypeOf ent:", reflect.TypeOf(ent))
	// fmt.Println("ValueOf ent:", reflect.ValueOf(ent))
	// b := ent.(*models.Book)

	// make changes / validate the content struct pointer (b) here
	// b.Name = "A new field value"
	return nil
}

//----------------------------------------------------------------------------
// ModelUpdateExt interface implementation for entity Book
//----------------------------------------------------------------------------

// BeforeDB extension-point implementation for entity Book
// TODO: implement pre-ORM call logic and document it here
func (updEP *MdlBookUpdateExt) BeforeDB(ent interface{}) error {

	// fmt.Println("TypeOf ent:", reflect.TypeOf(ent))
	// fmt.Println("ValueOf ent:", reflect.ValueOf(ent))
	// b := ent.(*models.Book)

	// make changes / validate the content struct pointer (b) here
	// b.Name = "A new field value"
	return nil
}

// AfterDB extension-point implementation for entity Book
// TODO: implement post-ORM call logic and document it here
func (updEP *MdlBookUpdateExt) AfterDB(ent interface{}) error {

	// fmt.Println("TypeOf ent:", reflect.TypeOf(ent))
	// fmt.Println("ValueOf ent:", reflect.ValueOf(ent))
	// b := ent.(*models.Book)

	// make changes / validate the content struct pointer (b) here
	// b.Name = "A new field value"
	return nil
}

//----------------------------------------------------------------------------
// ModelGetExt interface implementation for entity Book
//----------------------------------------------------------------------------

// BeforeDB extension-point implementation for entity Book
// TODO: implement pre-ORM call logic and document it here
func (getEP *MdlBookGetExt) BeforeDB(ent interface{}) error {

	// fmt.Println("TypeOf ent:", reflect.TypeOf(ent))
	// fmt.Println("ValueOf ent:", reflect.ValueOf(ent))
	// b := ent.(*models.Book)

	// make changes / validate the content struct pointer (b) here
	// b.Name = "A new field value"
	return nil
}

// AfterDB extension-point implementation for entity Book
// TODO: implement post-ORM call logic and document it here
func (getEP *MdlBookGetExt) AfterDB(ent interface{}) error {

	// fmt.Println("TypeOf ent:", reflect.TypeOf(ent))
	// fmt.Println("ValueOf ent:", reflect.ValueOf(ent))
	// b := ent.(*Book)

	// make changes / validate the content struct pointer (b) here
	// b.Name = "A new field value"
	return nil
}
