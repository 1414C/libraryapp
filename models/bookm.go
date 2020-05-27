package models

//=============================================================================================
// base Book entity model code generated on 27 May 20 17:32 CDT
//=============================================================================================

import (
	"fmt"

	"github.com/1414C/lw"
	"github.com/1414C/sqac"
)

// Book structure
type Book struct {
	ID        uint64  `json:"id" db:"id" sqac:"primary_key:inc"`
	Href      string  `json:"href" db:"href" sqac:"-"`
	Title     string  `json:"title" db:"title" sqac:"nullable:false;default:unknown title;index:non-unique"`
	Author    *string `json:"author,omitempty" db:"author" sqac:"nullable:true;index:non-unique"`
	Hardcover bool    `json:"hardcover" db:"hardcover" sqac:"nullable:false"`
	Copies    uint64  `json:"copies" db:"copies" sqac:"nullable:false"`
	LibraryID uint64  `json:"library_id" db:"library_id" sqac:"nullable:false;index:non-unique"`
}

// BookDB is a CRUD-type interface specifically for dealing with Books.
type BookDB interface {
	Create(book *Book) error
	Update(book *Book) error
	Delete(book *Book) error
	Get(book *Book) error
	GetBooks(params []sqac.GetParam, cmdMap map[string]interface{}) ([]Book, uint64) // uint64 holds $count result
	GetBooksByTitle(op string, Title string) []Book
	GetBooksByAuthor(op string, Author string) []Book
	GetBooksByHardcover(op string, Hardcover bool) []Book
	GetBooksByCopies(op string, Copies uint64) []Book
	GetBooksByLibraryID(op string, LibraryID uint64) []Book
}

// bookValidator checks and normalizes data prior to
// db access.
type bookValidator struct {
	BookDB
}

// bookValFunc type is the prototype for discrete Book normalization
// and validation functions that will be executed by func runBookValidationFuncs(...)
type bookValFunc func(*Book) error

// BookService is the public interface to the Book entity
type BookService interface {
	BookDB
}

// private service for book
type bookService struct {
	BookDB
}

// bookSqac is a sqac-based implementation of the BookDB interface.
type bookSqac struct {
	handle sqac.PublicDB
	ep     BookMdlExt
}

var _ BookDB = &bookSqac{}

// newBookValidator returns a new bookValidator
func newBookValidator(bdb BookDB) *bookValidator {
	return &bookValidator{
		BookDB: bdb,
	}
}

// runBookValFuncs executes a list of discrete validation
// functions against a book.
func runBookValFuncs(book *Book, fns ...bookValFunc) error {

	// iterate over the slice of function names and execute
	// each in-turn.  the order in which the lists are made
	// can matter...
	for _, fn := range fns {
		err := fn(book)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewBookService needs some work:
func NewBookService(handle sqac.PublicDB) BookService {

	bs := &bookSqac{
		handle: handle,
		ep:     *InitBookMdlExt(),
	}

	bv := newBookValidator(bs) // *db
	return &bookService{
		BookDB: bv,
	}
}

// ensure consistency (build error if delta exists)
var _ BookDB = &bookValidator{}

//-------------------------------------------------------------------------------------------------------
// CRUD-type model methods for Book
//-------------------------------------------------------------------------------------------------------
//
// Create validates and normalizes data used in the book creation.
// Create then calls the creation code contained in BookService.
func (bv *bookValidator) Create(book *Book) error {

	// perform normalization and validation -- comment out checks that are not required
	// note that the check calls are generated as a straight enumeration of the entity
	// structure.  It may be neccessary to adjust the calling order depending on the
	// relationships between the fields in the entity structure.
	err := runBookValFuncs(book,
		bv.normvalTitle,
		bv.normvalAuthor,
		bv.normvalHardcover,
		bv.normvalCopies,
		bv.normvalLibraryID,
	)

	if err != nil {
		return err
	}
	return bv.BookDB.Create(book)
}

// Update validates and normalizes the content of the Book
// being updated by way of executing a list of predefined discrete
// checks.  if the checks are successful, the entity is updated
// on the db via the ORM.
func (bv *bookValidator) Update(book *Book) error {

	// perform normalization and validation -- comment out checks that are not required
	// note that the check calls are generated as a straight enumeration of the entity
	// structure.  It may be neccessary to adjust the calling order depending on the
	// relationships between the fields in the entity structure.
	err := runBookValFuncs(book,
		bv.normvalTitle,
		bv.normvalAuthor,
		bv.normvalHardcover,
		bv.normvalCopies,
		bv.normvalLibraryID,
	)

	if err != nil {
		return err
	}
	return bv.BookDB.Update(book)
}

// Delete is passed through to the ORM with no real
// validations.  id is checked in the controller.
func (bv *bookValidator) Delete(book *Book) error {

	return bv.BookDB.Delete(book)
}

// Get is passed through to the ORM with no real
// validations.  id is checked in the controller.
func (bv *bookValidator) Get(book *Book) error {

	return bv.BookDB.Get(book)
}

// GetBooks is passed through to the ORM with no validation
func (bv *bookValidator) GetBooks(params []sqac.GetParam, cmdMap map[string]interface{}) ([]Book, uint64) {

	return bv.BookDB.GetBooks(params, cmdMap)
}

//-------------------------------------------------------------------------------------------------------
// internal bookValidator funcs
//-------------------------------------------------------------------------------------------------------
// These discrete functions are used to normalize and validate the Entity fields
// from with in the Create and Update methods.  See the comments in the model's
// Create and Update methods for details regarding use.

// normvalTitle normalizes and validates field Title
func (bv *bookValidator) normvalTitle(book *Book) error {

	// TODO: implement normalization and validation for Title
	return nil
}

// normvalAuthor normalizes and validates field Author
func (bv *bookValidator) normvalAuthor(book *Book) error {

	// TODO: implement normalization and validation for Author
	return nil
}

// normvalHardcover normalizes and validates field Hardcover
func (bv *bookValidator) normvalHardcover(book *Book) error {

	// TODO: implement normalization and validation for Hardcover
	return nil
}

// normvalCopies normalizes and validates field Copies
func (bv *bookValidator) normvalCopies(book *Book) error {

	// TODO: implement normalization and validation for Copies
	return nil
}

// normvalLibraryID normalizes and validates field LibraryID
func (bv *bookValidator) normvalLibraryID(book *Book) error {

	// TODO: implement normalization and validation for LibraryID
	return nil
}

//-------------------------------------------------------------------------------------------------------
// internal book Simple Query Validator funcs
//-------------------------------------------------------------------------------------------------------
// Simple query normalization and validation occurs in the controller to an
// extent, as the URL has to be examined closely in order to determine what to call
// in the model.  This section may be blank if no model fields were marked as
// selectable in the <models>.json file.
// GetBooksByTitle is passed through to the ORM with no validation.
func (bv *bookValidator) GetBooksByTitle(op string, title string) []Book {

	// TODO: implement normalization and validation for the GetBooksByTitle call.
	// TODO: typically no modifications are required here.
	return bv.BookDB.GetBooksByTitle(op, title)
}

// GetBooksByAuthor is passed through to the ORM with no validation.
func (bv *bookValidator) GetBooksByAuthor(op string, author string) []Book {

	// TODO: implement normalization and validation for the GetBooksByAuthor call.
	// TODO: typically no modifications are required here.
	return bv.BookDB.GetBooksByAuthor(op, author)
}

// GetBooksByHardcover is passed through to the ORM with no validation.
func (bv *bookValidator) GetBooksByHardcover(op string, hardcover bool) []Book {

	// TODO: implement normalization and validation for the GetBooksByHardcover call.
	// TODO: typically no modifications are required here.
	return bv.BookDB.GetBooksByHardcover(op, hardcover)
}

// GetBooksByCopies is passed through to the ORM with no validation.
func (bv *bookValidator) GetBooksByCopies(op string, copies uint64) []Book {

	// TODO: implement normalization and validation for the GetBooksByCopies call.
	// TODO: typically no modifications are required here.
	return bv.BookDB.GetBooksByCopies(op, copies)
}

// GetBooksByLibraryID is passed through to the ORM with no validation.
func (bv *bookValidator) GetBooksByLibraryID(op string, library_id uint64) []Book {

	// TODO: implement normalization and validation for the GetBooksByLibraryID call.
	// TODO: typically no modifications are required here.
	return bv.BookDB.GetBooksByLibraryID(op, library_id)
}

//-------------------------------------------------------------------------------------------------------
// ORM db CRUD access methods
//-------------------------------------------------------------------------------------------------------
//
// Create a new Book in the database via the ORM
func (bs *bookSqac) Create(book *Book) error {

	err := bs.ep.CrtEp.BeforeDB(book)
	if err != nil {
		return err
	}
	err = bs.handle.Create(book)
	if err != nil {
		return err
	}
	err = bs.ep.CrtEp.AfterDB(book)
	if err != nil {
		return err
	}
	return err
}

// Update an existng Book in the database via the ORM
func (bs *bookSqac) Update(book *Book) error {

	err := bs.ep.UpdEp.BeforeDB(book)
	if err != nil {
		return err
	}
	err = bs.handle.Update(book)
	if err != nil {
		return err
	}
	err = bs.ep.UpdEp.AfterDB(book)
	if err != nil {
		return err
	}
	return err
}

// Delete an existing Book in the database via the ORM
func (bs *bookSqac) Delete(book *Book) error {
	return bs.handle.Delete(book)
}

// Get an existing Book from the database via the ORM
func (bs *bookSqac) Get(book *Book) error {

	err := bs.ep.GetEp.BeforeDB(book)
	if err != nil {
		return err
	}
	err = bs.handle.GetEntity(book)
	if err != nil {
		return err
	}
	err = bs.ep.GetEp.AfterDB(book)
	if err != nil {
		return err
	}
	return err
}

// Get all existing Books from the db via the ORM
func (bs *bookSqac) GetBooks(params []sqac.GetParam, cmdMap map[string]interface{}) ([]Book, uint64) {

	var err error

	// create a slice to read into
	books := []Book{}

	// call the ORM
	result, err := bs.handle.GetEntitiesWithCommands(books, params, cmdMap)
	if err != nil {
		lw.Warning("BookModel GetBooks() error: %s", err.Error())
		return nil, 0
	}

	// check to see what was returned
	switch result.(type) {
	case []Book:
		books = result.([]Book)

		// call the extension-point
		for i := range books {
			err = bs.ep.GetEp.AfterDB(&books[i])
			if err != nil {
				lw.Warning("BookModel GetBooks AfterDB() error: %s", err.Error())
			}
		}
		return books, 0

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
// Get all existing BooksByTitle from the db via the ORM
func (bs *bookSqac) GetBooksByTitle(op string, Title string) []Book {

	var books []Book
	var c string

	switch op {
	case "EQ":
		c = "title = ?"
	case "LIKE":
		c = "title like ?"
	default:
		return nil
	}
	qs := fmt.Sprintf("SELECT * FROM book WHERE %s;", c)
	err := bs.handle.Select(&books, qs, Title)
	if err != nil {
		lw.Warning("GetBooksByTitle got: %s", err.Error())
		return nil
	}

	if bs.handle.IsLog() {
		lw.Info("GetBooksByTitle found: %v based on (%s %v)", books, op, Title)
	}

	// call the extension-point
	for i := range books {
		err = bs.ep.GetEp.AfterDB(&books[i])
		if err != nil {
			lw.Warning("BookModel Getbooks AfterDB() error: %s", err.Error())
		}
	}
	return books
}

// Get all existing BooksByAuthor from the db via the ORM
func (bs *bookSqac) GetBooksByAuthor(op string, Author string) []Book {

	var books []Book
	var c string

	switch op {
	case "EQ":
		c = "author = ?"
	case "LIKE":
		c = "author like ?"
	default:
		return nil
	}
	qs := fmt.Sprintf("SELECT * FROM book WHERE %s;", c)
	err := bs.handle.Select(&books, qs, Author)
	if err != nil {
		lw.Warning("GetBooksByAuthor got: %s", err.Error())
		return nil
	}

	if bs.handle.IsLog() {
		lw.Info("GetBooksByAuthor found: %v based on (%s %v)", books, op, Author)
	}

	// call the extension-point
	for i := range books {
		err = bs.ep.GetEp.AfterDB(&books[i])
		if err != nil {
			lw.Warning("BookModel Getbooks AfterDB() error: %s", err.Error())
		}
	}
	return books
}

// Get all existing BooksByHardcover from the db via the ORM
func (bs *bookSqac) GetBooksByHardcover(op string, Hardcover bool) []Book {

	var books []Book
	var c string

	switch op {
	case "EQ":
		c = "hardcover = ?"
	case "NE":
		c = "hardcover != ?"
	default:
		return nil
	}
	qs := fmt.Sprintf("SELECT * FROM book WHERE %s;", c)
	err := bs.handle.Select(&books, qs, Hardcover)
	if err != nil {
		lw.Warning("GetBooksByHardcover got: %s", err.Error())
		return nil
	}

	if bs.handle.IsLog() {
		lw.Info("GetBooksByHardcover found: %v based on (%s %v)", books, op, Hardcover)
	}

	// call the extension-point
	for i := range books {
		err = bs.ep.GetEp.AfterDB(&books[i])
		if err != nil {
			lw.Warning("BookModel Getbooks AfterDB() error: %s", err.Error())
		}
	}
	return books
}

// Get all existing BooksByCopies from the db via the ORM
func (bs *bookSqac) GetBooksByCopies(op string, Copies uint64) []Book {

	var books []Book
	var c string

	switch op {
	case "EQ":
		c = "copies = ?"
	case "LT":
		c = "copies < ?"
	case "GT":
		c = "copies > ?"
	default:
		return nil
	}
	qs := fmt.Sprintf("SELECT * FROM book WHERE %s;", c)
	err := bs.handle.Select(&books, qs, Copies)
	if err != nil {
		lw.Warning("GetBooksByCopies got: %s", err.Error())
		return nil
	}

	if bs.handle.IsLog() {
		lw.Info("GetBooksByCopies found: %v based on (%s %v)", books, op, Copies)
	}

	// call the extension-point
	for i := range books {
		err = bs.ep.GetEp.AfterDB(&books[i])
		if err != nil {
			lw.Warning("BookModel Getbooks AfterDB() error: %s", err.Error())
		}
	}
	return books
}

// Get all existing BooksByLibraryID from the db via the ORM
func (bs *bookSqac) GetBooksByLibraryID(op string, LibraryID uint64) []Book {

	var books []Book
	var c string

	switch op {
	case "EQ":
		c = "library_id = ?"
	default:
		return nil
	}
	qs := fmt.Sprintf("SELECT * FROM book WHERE %s;", c)
	err := bs.handle.Select(&books, qs, LibraryID)
	if err != nil {
		lw.Warning("GetBooksByLibraryID got: %s", err.Error())
		return nil
	}

	if bs.handle.IsLog() {
		lw.Info("GetBooksByLibraryID found: %v based on (%s %v)", books, op, LibraryID)
	}

	// call the extension-point
	for i := range books {
		err = bs.ep.GetEp.AfterDB(&books[i])
		if err != nil {
			lw.Warning("BookModel Getbooks AfterDB() error: %s", err.Error())
		}
	}
	return books
}
