package ext

// This file contains the interface declarations for all entity model extension-points.
// Extension-points are provided as an escape hatch in the generated code for the application
// developer to implement their own logic without changing the core (generated) application
// code.  The use of extension-points for implementation logic also safeguards an application
// developer's work from being over-written if the application is regenerated.
//
// The default implementations of the interface methods return nil, and provide some seed
// code that has been commented out.  Method signatures containing (ent interface{})
// imply that a pointer to an entity is expected, and this is how the calls to such
// interface methods have been generated in the models.
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
//
// The model extension-point interfaces are identical to one another.  While the interfaces
// are the same, a distinction has been drawn between them based on expected use, as documented
// in their respective declarations.

// ModelCreateExt is an interface used to define the extension-points available
// to be implemented within the Create method of an entity's model.
type ModelCreateExt interface {

	// BeforeDB is an extension-point that can be implemented in order to examine
	// and potentially make changes to the values in the entity structure immediately
	// before the insertion request is made to the ORM. This extension-point is the first
	// code executed in the model's Create method. Authentication and Authorization checks
	// should be performed upstream in the route middleware-layer and detailed checks of an
	// entity's data should be carried out in the validator-layer.
	BeforeDB(ent interface{}) error

	// CreateAfterDB is an extension-point that can be implemented in order to examine
	// and potentially make changes to the values in the entity structure immediately
	// following the return of the ORM insertion request. This extension-point is the last
	// code executed in the model's Create method. As the insertion will have already
	// occurred at this point, care should be taken when deciding whether to issue an
	// error in this extension-point.  Augmentation of the of the Create result may be
	// carried out in this method in order to calculate non-persistent entity values
	// for example.
	AfterDB(ent interface{}) error
}

// ModelUpdateExt is an interface used to define the extension-points available
// to be implemented within the Update method of an entity's model.
type ModelUpdateExt interface {

	// BeforeDB is an extension-point that can be implemented in order to examine
	// and potentially make changes to the values in the entity structure immediately
	// before the update request is made to the ORM. This extension-point is the first
	// code executed in the model's Update method. Authentication and Authorization checks
	// should be performed upstream in the route middleware-layer and detailed checks and
	// normalization of the entity's data should be carried out in the validator-layer.
	BeforeDB(ent interface{}) error

	// AfterDB is an extension-point that can be implemented in order to examine
	// and potentially make changes to the values in the entity structure immediately
	// following the return of the ORM update request. This extension-point is the last
	// code executed in the model's Update method. As the update will have already
	// occurred at this point, care should be taken when deciding whether to issue an
	// error in this extension-point.  Augmentation of the of the Update result may be
	// carried out in this method in order to calculate non-persistent entity values
	// for example.
	AfterDB(ent interface{}) error
}

// ModelGetExt is an interface used to define the extension-points available
// to be implemented within the Get method of an entity's model.
type ModelGetExt interface {

	// BeforeDB is an extension-point that can be implemented in order to examine
	// and potentially make changes to the values in the entity structure immediately
	// before the read-entity request is made to the ORM. This extension-point is the
	// first code executed in the model's Get method. Authentication and Authorization
	// checks should be performed upstream in the route middleware-layer and detailed
	// checks of an entity's data should be carried out in the validator-layer.
	BeforeDB(ent interface{}) error

	// AfterDB is an extension-point that can be implemented in order to examine
	// and potentially make changes to the values in the entity structure immediately
	// following the return of the ORM read-entity request. This extension-point is the
	// last code executed in the model's Get method. As the read will have already
	// occurred at this point, care should be taken when deciding whether to issue an
	// error in this extension-point.  Augmentation of the of the Get result may be
	// carried out in this method in order to calculate non-persistent entity values
	// for example.
	AfterDB(ent interface{}) error
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
