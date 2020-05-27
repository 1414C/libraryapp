package ext

import (
	// "fmt"
	"net/http"
	// "reflect"
)

//===================================================================================================
// base Book entity controller extension-point code generated on 27 May 20 17:32 CDT
//===================================================================================================

// CtrlBookCreateExt provides access to the ControllerCreateExt extension-point interface
type CtrlBookCreateExt struct {
	ControllerCreateExt
}

// CtrlBookGetExt provides access to the ControllerGetExt extension-point interface
type CtrlBookGetExt struct {
	ControllerGetExt
}

// CtrlBookUpdateExt provides access to the ControllerUpdateExt extension-point interface
type CtrlBookUpdateExt struct {
	ControllerUpdateExt
}

// BookCtrlExt provides access to the Book implementations of the following interfaces:
//   CtrlCreateExt
//   CtrlUpdateExt
//   CtrlGetExt
type BookCtrlExt struct {
	CrtEp CtrlBookCreateExt
	UpdEp CtrlBookUpdateExt
	GetEp CtrlBookGetExt
}

var bookCtrlExp BookCtrlExt

// InitBookCtrlExt initializes the book entity's controller
// extension-point interface implementations.
func InitBookCtrlExt() *BookCtrlExt {
	bookCtrlExp = BookCtrlExt{}
	return &bookCtrlExp
}

//------------------------------------------------------------------------------------------
// ControllerCreateExt extension-point interface implementation for entity Book
//------------------------------------------------------------------------------------------

// BeforeFirst extension-point implementation for entity Book
// TODO: implement checks and document them here
func (crtEP *CtrlBookCreateExt) BeforeFirst(w http.ResponseWriter, r *http.Request) error {

	return nil
}

// AfterBodyDecode extension-point implementation for entity Book
// TODO: implement application logic and document it here
func (crtEP *CtrlBookCreateExt) AfterBodyDecode(ent interface{}) error {

	// fmt.Println("TypeOf ent:", reflect.TypeOf(ent))
	// fmt.Println("ValueOf ent:", reflect.ValueOf(ent))
	// b := ent.(*models.Book)
	//
	// make changes / validate the content struct pointer (b) here
	// b.<field_name> = "A new value"

	return nil
}

// BeforeResponse extension-point implementation for entity Book
// TODO: implement application logic and document it here
func (crtEP *CtrlBookCreateExt) BeforeResponse(ent interface{}) error {

	// fmt.Println("TypeOf ent:", reflect.TypeOf(ent))
	// fmt.Println("ValueOf ent:", reflect.ValueOf(ent))
	// b := ent.(*models.Book)
	//
	// make changes / validate the content struct pointer (b) here
	// b.<field_name> = b.<field_name> + "."

	return nil
}

//------------------------------------------------------------------------------------------
// ControllerUpdateExt extension-point interface implementation for entity Book
//------------------------------------------------------------------------------------------

// BeforeFirst extension-point implementation for entity Book
// TODO: implement checks and document them here
func (updEP *CtrlBookUpdateExt) BeforeFirst(w http.ResponseWriter, r *http.Request) error {

	return nil
}

// AfterBodyDecode extension-point implementation for entity Book
// TODO: implement application logic and document it here
func (updEP *CtrlBookUpdateExt) AfterBodyDecode(ent interface{}) error {

	// fmt.Println("TypeOf ent:", reflect.TypeOf(ent))
	// fmt.Println("ValueOf ent:", reflect.ValueOf(ent))
	// b := ent.(*models.Book)
	//
	// make changes / validate the content struct pointer (b) here
	// b.<field_name> = "An updated value"
	return nil
}

// BeforeResponse extension-point implementation for entity Book
// TODO: implement application logic and document it here
func (updEP *CtrlBookUpdateExt) BeforeResponse(ent interface{}) error {

	// fmt.Println("TypeOf ent:", reflect.TypeOf(ent))
	// fmt.Println("ValueOf ent:", reflect.ValueOf(ent))
	// b := ent.(*models.Book)
	//
	// make changes / validate the content struct pointer (b) here
	// b.<field_name> = b.<field_name> + "."

	return nil
}

//------------------------------------------------------------------------------------------
// ControllerGetExt extension-point interface implementation for entity Book
//------------------------------------------------------------------------------------------

// BeforeFirst extension-point implementation for entity Book
// TODO: implement checks and document them here
func (getEP *CtrlBookGetExt) BeforeFirst(w http.ResponseWriter, r *http.Request) error {

	return nil
}

// BeforeModelCall extension-point implementation for entity Book
// TODO: implement application logic and document it here
func (getEP *CtrlBookGetExt) BeforeModelCall(ent interface{}) error {

	// fmt.Println("TypeOf ent:", reflect.TypeOf(ent))
	// fmt.Println("ValueOf ent:", reflect.ValueOf(ent))
	// b := ent.(*models.Book)
	//
	// make changes / validate the content struct pointer (b) here
	// b.<field_name> = "A new value"

	return nil
}

// BeforeResponse extension-point implementation for entity Book
// TODO: implement application logic and document it here
func (getEP *CtrlBookGetExt) BeforeResponse(ent interface{}) error {

	// b := ent.(*models.Book)
	//
	// make changes / validate the content struct pointer (b) here
	// b.<field_name> = b.<field_name> + "."

	return nil
}
