package ext

import (
	"net/http"
)

// This file contains the interface declarations for all entity controller extension-points.
// Extension-points are provided as an escape hatch in the generated code for the application
// developer to implement their own logic without changing the core (generated) application
// code.  The use of extension-points for implementation logic also safeguards an application
// developer's work from being over-written if the application is regenerated.
//
// The default implementations of the interface methods return nil, and provide some seed
// code that has been commented out.  Method signatures containing (ent interface{})
// imply that a pointer to an entity is expected, and this is how the calls to such
// interface methods have been generated in the controllers.
//
// While it may seem that reflection is a questionable choice here, the true overhead
// in this use-case is quite small, as the interface implementations are generated
// and the ent can be confidently 'cast' with a simple type-assertion, then operated
// on.  For example:
//
// 	l := ent.(*models.Library)
//
//	// make changes / validate the content struct pointer (l) here
//	l.Name = "Bunny Dumpling Public Library"
//  return
//

// ControllerCreateExt is an interface used to define the extension-points available
// to be implemented within the Create method of an entity's controller.
type ControllerCreateExt interface {

	// BeforeFirst is an extension-point that can be implemented in order
	// to examine and potentially reject a Create entity request. This extension-
	// point is the first code executed in the controller's Create method. Authentication
	// and Authorization checks should be performed upstream in the route middleware-
	// layer and detailed checks of a request.Body should be carried out by the validator
	// in the model-layer.
	BeforeFirst(w http.ResponseWriter, r *http.Request) error

	// AfterBodyDecode is an extension-point that can be implemented to perform
	// preliminary checks and changes to the unmarshalled content of the request.Body.
	// Detailed checks of the unmarshalled data from the request.Body should be carried
	// out by the validator in the model-layer. This extension-point should only be
	// used to carry out deal-breaker checks and perhaps to default data in the entity
	// struct prior to calling the validator/normalization methods in the model-layer.
	AfterBodyDecode(ent interface{}) error

	// BeforeResponse is an extension-point that can be implemented to perform
	// checks following the return of the call to the model-layer. At this point,
	// changes to the db will have been made, so failing the call should take this
	// into consideration.
	BeforeResponse(ent interface{}) error
}

// ControllerUpdateExt is an interface used to define the extension-points available
// to be implemented within the Update method of an entity's controller.
type ControllerUpdateExt interface {

	// BeforeFirst is an extension-point that can be implemented in order
	// to examine and potentially reject an Update entity request. This extension-
	// point is the first code executed in the controller's Update method.
	// Authentication and Authorization checks should be performed upstream in the
	// route middleware-layer and detailed checks of a request.Body should be carried
	// out by the validator in the model-layer.
	BeforeFirst(w http.ResponseWriter, r *http.Request) error

	// AfterBodyDecode is an extension-point that can be implemented to perform
	// preliminary checks and changes to the unmarshalled content of the request.Body.
	// Detailed checks of the unmarshalled data from the request.Body should be carried
	// out by the validator in the model-layer. This extension-point should only be
	// used to carry out deal-breaker checks and perhaps to default data in the entity
	// struct prior to calling the validator/normalization methods in the model-layer.
	AfterBodyDecode(ent interface{}) error

	// BeforeResponse is an extension-point that can be implemented to perform
	// checks following the return of the call to the model-layer. At this point,
	// changes to the db will have been made, so failing the call should take this
	// into consideration.
	BeforeResponse(ent interface{}) error
}

// ControllerGetExt is an interface used to define the extension-points available
// to be implemented within the Get method of an entity's controller.
type ControllerGetExt interface {

	// BeforeFirst is an extension-point that can be implemented in order
	// to examine and potentially reject a Get entity request. This extension-
	// point is the first code executed in the controller's Create method.
	// Authentication and Authorization checks should be performed upstream in
	// the route middleware-layer.
	BeforeFirst(w http.ResponseWriter, r *http.Request) error

	// BeforeModelCall is an extension-point that can be implemented in order
	// to make changes to the content of the entity structure prior to calling
	// the model-layer. By default the controller's Get method will populate the
	// ID field of the entity structure using the :id value provided in the request
	// URL. The use of this extension-point would be seemingly rare and any values
	// added to the struct would be over-written in the model-layer when the call
	// to the DBMS is made. The added values would however be available for use in
	// the validation/normalization and DBMS access methods prior to the call to
	// the ORM.
	BeforeModelCall(ent interface{}) error

	// BeforeResponse is an extension-point that can be implemented to perform
	// checks / changes following the return of the call to the model-layer. At this
	// point, the db has been read and the populated entity structure is about to be
	// marshalled into JSON and passed back to the router/mux.
	BeforeResponse(ent interface{}) error
}

// type ControllerDeleteExt interface {
// 	//
// 	DeleteBeforeFirst
// 	DeleteAfterBodyDecode
// 	DeleteBeforeModelCall
// 	DeleteBeforeResponse
// }
// type ControllerGetsExt interface {
// 	//
// 	GetsBeforeFirst
// 	GetsAfterBodyDecode
// 	GetsBeforeModelCall
// 	GetsBeforeResponse
// }
