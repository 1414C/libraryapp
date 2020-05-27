package ext

import (
	// "fmt"
	"net/http"
	// "reflect"
)

//===================================================================================================
// base Library entity controller extension-point code generated on 27 May 20 17:32 CDT
//===================================================================================================

// CtrlLibraryCreateExt provides access to the ControllerCreateExt extension-point interface
type CtrlLibraryCreateExt struct {
	ControllerCreateExt
}

// CtrlLibraryGetExt provides access to the ControllerGetExt extension-point interface
type CtrlLibraryGetExt struct {
	ControllerGetExt
}

// CtrlLibraryUpdateExt provides access to the ControllerUpdateExt extension-point interface
type CtrlLibraryUpdateExt struct {
	ControllerUpdateExt
}

// LibraryCtrlExt provides access to the Library implementations of the following interfaces:
//   CtrlCreateExt
//   CtrlUpdateExt
//   CtrlGetExt
type LibraryCtrlExt struct {
	CrtEp CtrlLibraryCreateExt
	UpdEp CtrlLibraryUpdateExt
	GetEp CtrlLibraryGetExt
}

var libraryCtrlExp LibraryCtrlExt

// InitLibraryCtrlExt initializes the library entity's controller
// extension-point interface implementations.
func InitLibraryCtrlExt() *LibraryCtrlExt {
	libraryCtrlExp = LibraryCtrlExt{}
	return &libraryCtrlExp
}

//------------------------------------------------------------------------------------------
// ControllerCreateExt extension-point interface implementation for entity Library
//------------------------------------------------------------------------------------------

// BeforeFirst extension-point implementation for entity Library
// TODO: implement checks and document them here
func (crtEP *CtrlLibraryCreateExt) BeforeFirst(w http.ResponseWriter, r *http.Request) error {

	return nil
}

// AfterBodyDecode extension-point implementation for entity Library
// TODO: implement application logic and document it here
func (crtEP *CtrlLibraryCreateExt) AfterBodyDecode(ent interface{}) error {

	// fmt.Println("TypeOf ent:", reflect.TypeOf(ent))
	// fmt.Println("ValueOf ent:", reflect.ValueOf(ent))
	// l := ent.(*models.Library)
	//
	// make changes / validate the content struct pointer (l) here
	// l.<field_name> = "A new value"

	return nil
}

// BeforeResponse extension-point implementation for entity Library
// TODO: implement application logic and document it here
func (crtEP *CtrlLibraryCreateExt) BeforeResponse(ent interface{}) error {

	// fmt.Println("TypeOf ent:", reflect.TypeOf(ent))
	// fmt.Println("ValueOf ent:", reflect.ValueOf(ent))
	// l := ent.(*models.Library)
	//
	// make changes / validate the content struct pointer (l) here
	// l.<field_name> = l.<field_name> + "."

	return nil
}

//------------------------------------------------------------------------------------------
// ControllerUpdateExt extension-point interface implementation for entity Library
//------------------------------------------------------------------------------------------

// BeforeFirst extension-point implementation for entity Library
// TODO: implement checks and document them here
func (updEP *CtrlLibraryUpdateExt) BeforeFirst(w http.ResponseWriter, r *http.Request) error {

	return nil
}

// AfterBodyDecode extension-point implementation for entity Library
// TODO: implement application logic and document it here
func (updEP *CtrlLibraryUpdateExt) AfterBodyDecode(ent interface{}) error {

	// fmt.Println("TypeOf ent:", reflect.TypeOf(ent))
	// fmt.Println("ValueOf ent:", reflect.ValueOf(ent))
	// l := ent.(*models.Library)
	//
	// make changes / validate the content struct pointer (l) here
	// l.<field_name> = "An updated value"
	return nil
}

// BeforeResponse extension-point implementation for entity Library
// TODO: implement application logic and document it here
func (updEP *CtrlLibraryUpdateExt) BeforeResponse(ent interface{}) error {

	// fmt.Println("TypeOf ent:", reflect.TypeOf(ent))
	// fmt.Println("ValueOf ent:", reflect.ValueOf(ent))
	// l := ent.(*models.Library)
	//
	// make changes / validate the content struct pointer (l) here
	// l.<field_name> = l.<field_name> + "."

	return nil
}

//------------------------------------------------------------------------------------------
// ControllerGetExt extension-point interface implementation for entity Library
//------------------------------------------------------------------------------------------

// BeforeFirst extension-point implementation for entity Library
// TODO: implement checks and document them here
func (getEP *CtrlLibraryGetExt) BeforeFirst(w http.ResponseWriter, r *http.Request) error {

	return nil
}

// BeforeModelCall extension-point implementation for entity Library
// TODO: implement application logic and document it here
func (getEP *CtrlLibraryGetExt) BeforeModelCall(ent interface{}) error {

	// fmt.Println("TypeOf ent:", reflect.TypeOf(ent))
	// fmt.Println("ValueOf ent:", reflect.ValueOf(ent))
	// l := ent.(*models.Library)
	//
	// make changes / validate the content struct pointer (l) here
	// l.<field_name> = "A new value"

	return nil
}

// BeforeResponse extension-point implementation for entity Library
// TODO: implement application logic and document it here
func (getEP *CtrlLibraryGetExt) BeforeResponse(ent interface{}) error {

	// l := ent.(*models.Library)
	//
	// make changes / validate the content struct pointer (l) here
	// l.<field_name> = l.<field_name> + "."

	return nil
}
