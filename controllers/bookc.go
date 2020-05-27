package controllers

//=============================================================================================
// base Book entity controller code generated on 27 May 20 17:57 CDT
//=============================================================================================

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/1414C/libraryapp/controllers/ext"
	"github.com/1414C/libraryapp/models"
	"github.com/1414C/lw"
	"github.com/1414C/sqac"
	"github.com/gorilla/mux"
)

// BookController is the book controller type for route binding
type BookController struct {
	bs   models.BookService
	ep   ext.BookCtrlExt
	svcs models.Services
}

// NewBookController creates a new BookController
func NewBookController(bs models.BookService, svcs models.Services) *BookController {
	return &BookController{
		bs:   bs,
		ep:   *ext.InitBookCtrlExt(),
		svcs: svcs,
	}
}

// Create facilitates the creation of a new Book.  This method is bound
// to the gorilla.mux router in main.go.
//
// POST /book
func (bc *BookController) Create(w http.ResponseWriter, r *http.Request) {

	var err error
	var bm models.Book

	// TODO: implement extension-point if required
	// TODO: safe to comment this block out if the extension-point is not needed
	err = bc.ep.CrtEp.BeforeFirst(w, r)
	if err != nil {
		lw.ErrorWithPrefixString("BookController CreateBeforeFirst() error:", err)
		respondWithError(w, http.StatusBadRequest, "bookc: Invalid request")
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&bm); err != nil {
		lw.ErrorWithPrefixString("Book Create:", err)
		respondWithError(w, http.StatusBadRequest, "bookc: Invalid request payload")
		return
	}
	defer r.Body.Close()

	// TODO: implement extension-point if required
	// TODO: safe to comment this block out if the extension-point is not needed
	err = bc.ep.CrtEp.AfterBodyDecode(&bm)
	if err != nil {
		lw.ErrorWithPrefixString("BookController CreateAfterBodyDecode() error:", err)
		respondWithError(w, http.StatusBadRequest, "bookc: Invalid request payload")
		return
	}

	// fill the model
	book := models.Book{
		Title:     bm.Title,
		Author:    bm.Author,
		Hardcover: bm.Hardcover,
		Copies:    bm.Copies,
		LibraryID: bm.LibraryID,
	}

	// build a base urlString for the JSON Body self-referencing Href tag
	urlString := buildHrefStringFromCRUDReq(r, true)

	// call the Create method on the book model
	err = bc.bs.Create(&book)
	if err != nil {
		lw.ErrorWithPrefixString("Book Create:", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	book.Href = urlString + strconv.FormatUint(uint64(book.ID), 10)

	// TODO: implement extension-point if required
	// TODO: safe to comment this block out if the extension-point is not needed
	err = bc.ep.CrtEp.BeforeResponse(&book)
	if err != nil {
		lw.ErrorWithPrefixString("BookController CreateBeforeResponse() error:", err)
		respondWithError(w, http.StatusBadRequest, "bookc: Invalid request")
		return
	}
	respondWithJSON(w, http.StatusCreated, book)
}

// Update facilitates the update of an existing Book.  This method is bound
// to the gorilla.mux router in main.go.
//
// PUT /book:id
func (bc *BookController) Update(w http.ResponseWriter, r *http.Request) {

	var err error
	var bm models.Book

	// TODO: implement extension-point if required
	// TODO: safe to comment this block out if the extension-point is not needed
	err = bc.ep.UpdEp.BeforeFirst(w, r)
	if err != nil {
		lw.ErrorWithPrefixString("BookController UpdateBeforeFirst() error:", err)
		respondWithError(w, http.StatusBadRequest, "bookc: Invalid request")
		return
	}

	// get the parameter(s)
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		lw.ErrorWithPrefixString("Book Update:", err)
		respondWithError(w, http.StatusBadRequest, "Invalid book id")
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&bm); err != nil {
		lw.ErrorWithPrefixString("Book Update:", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// TODO: implement extension-point if required
	// TODO: safe to comment this block out if the extension-point is not needed
	err = bc.ep.UpdEp.AfterBodyDecode(&bm)
	if err != nil {
		lw.ErrorWithPrefixString("BookController UpdateAfterBodyDecode() error:", err)
		respondWithError(w, http.StatusBadRequest, "bookc: Invalid request payload")
		return
	}

	// fill the model
	book := models.Book{
		Title:     bm.Title,
		Author:    bm.Author,
		Hardcover: bm.Hardcover,
		Copies:    bm.Copies,
		LibraryID: bm.LibraryID,
	}

	// build a base urlString for the JSON Body self-referencing Href tag
	urlString := buildHrefStringFromCRUDReq(r, false)
	book.ID = id

	// call the update method on the model
	err = bc.bs.Update(&book)
	if err != nil {
		lw.ErrorWithPrefixString("Book Update:", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	book.Href = urlString

	// TODO: implement extension-point if required
	// TODO: safe to comment this block out if the extension-point is not needed
	err = bc.ep.UpdEp.BeforeResponse(&book)
	if err != nil {
		lw.ErrorWithPrefixString("BookController UpdateBeforeResponse() error:", err)
		respondWithError(w, http.StatusBadRequest, "bookc: Invalid request")
		return
	}
	respondWithJSON(w, http.StatusCreated, book)
}

// Get facilitates the retrieval of an existing Book.  This method is bound
// to the gorilla.mux router in main.go.
//
// GET /book/:id
func (bc *BookController) Get(w http.ResponseWriter, r *http.Request) {

	var err error

	// TODO: implement extension-point if required
	// TODO: safe to comment this block out if the extension-point is not needed
	err = bc.ep.GetEp.BeforeFirst(w, r)
	if err != nil {
		lw.Warning("BookController GetBeforeFirst() error: %s", err.Error())
		respondWithError(w, http.StatusBadRequest, "bookc: Invalid request")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		lw.Warning("Book Get: %s", err.Error())
		respondWithError(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	// build a base urlString for the JSON Body self-referencing Href tag
	urlString := buildHrefStringFromCRUDReq(r, false)

	book := models.Book{
		ID: id,
	}

	// TODO: implement extension-point if required
	// TODO: safe to comment this block out if the extension-point is not needed
	err = bc.ep.GetEp.BeforeModelCall(&book)
	if err != nil {
		lw.Warning("BookController GetBeforeModelCall() error: %s", err.Error())
		respondWithError(w, http.StatusBadRequest, "bookc: Invalid request")
		return
	}

	err = bc.bs.Get(&book)
	if err != nil {
		lw.Warning(err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	book.Href = urlString

	// TODO: implement extension-point if required
	// TODO: safe to comment this block out if the extension-point is not needed
	err = bc.ep.GetEp.BeforeResponse(&book)
	if err != nil {
		lw.Warning("BookController GetBeforeResponse() error: %s", err.Error())
		respondWithError(w, http.StatusBadRequest, "bookc: Invalid request")
		return
	}
	respondWithJSON(w, http.StatusCreated, book)
}

// Delete facilitates the deletion of an existing Book.  This method is bound
// to the gorilla.mux router in main.go.
//
// DELETE /book/:id
func (bc *BookController) Delete(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		lw.ErrorWithPrefixString("Book Delete:", err)
		respondWithError(w, http.StatusBadRequest, "Invalid Book ID")
		return
	}

	book := models.Book{
		ID: id,
	}

	err = bc.bs.Delete(&book)
	if err != nil {
		lw.ErrorWithPrefixString("Book Delete:", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithHeader(w, http.StatusAccepted)
}

// getBookSet is used by all BookSet queries as a means of injecting parameters
// returns ([]Book, $count, countRequested, error)
func (bc *BookController) getBookSet(w http.ResponseWriter, r *http.Request, params []sqac.GetParam) ([]models.Book, uint64, bool, error) {

	var mapCommands map[string]interface{}
	var err error
	var urlString string
	var books []models.Book
	var count uint64
	countReq := false

	// check for mux.vars
	vars := mux.Vars(r)

	// parse commands ($cmd) if any
	if len(vars) > 0 && vars != nil {
		mapCommands, err = parseRequestCommands(vars)
		if err != nil {
			return nil, 0, false, err
		}
	}

	// $count trumps all other commands
	if mapCommands != nil {
		_, ok := mapCommands["count"]
		if ok {
			countReq = true
		}
		books, count = bc.bs.GetBooks(params, mapCommands)
	} else {
		books, count = bc.bs.GetBooks(params, nil)
	}

	// retrieved []Book and not asked to $count
	if books != nil && countReq == false {
		for i, l := range books {
			books[i].Href = urlString + strconv.FormatUint(uint64(l.ID), 10)
		}
		return books, 0, countReq, nil
	}

	// $count was requested, which trumps all other commands
	if countReq == true {
		return nil, count, countReq, nil
	}

	// fallthrough and return nothing
	return nil, 0, countReq, nil
}

// GetBooks facilitates the retrieval of all existing Books.  This method is bound
// to the gorilla.mux router in main.go.
//
// GET /books
// GET /books/$count | $limit=n $offset=n $orderby=<field_name> ($asc|$desc)
func (bc *BookController) GetBooks(w http.ResponseWriter, r *http.Request) {

	var books []models.Book
	var count uint64
	countReq := false

	// build base Href; common for each selected row
	urlString := buildHrefBasic(r, true)

	// call the common getBookSet method
	books, count, countReq, err := bc.getBookSet(w, r, nil)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf(`{"GetBooks": "%s"}`, err))
		return
	}

	// retrieved []Book and not asked to $count
	if books != nil && countReq == false {
		for i, l := range books {
			books[i].Href = urlString + strconv.FormatUint(uint64(l.ID), 10)
		}
		respondWithJSON(w, http.StatusOK, books)
		return
	}

	// $count was requested, which trumps all other commands
	if countReq == true {
		respondWithCount(w, http.StatusOK, count)
		return
	}

	// fallthrough and return nothing
	respondWithJSON(w, http.StatusOK, "[]")
}

// GetBooksByTitle facilitates the retrieval of existing
// Books based on Title.
// GET /books/title(OP 'searchString')
// GET /books/title(OP 'searchString')/$count | $limit=n $offset=n $orderby=<field_name> ($asc|$desc)
func (bc *BookController) GetBooksByTitle(w http.ResponseWriter, r *http.Request) {

	// get the title parameter
	vars := mux.Vars(r)
	searchValue := vars["title"]
	if searchValue == "" {
		respondWithError(w, http.StatusBadRequest, "missing search criteria")
		return
	}

	// adjust operator and predicate if neccessary
	op, predicate, err := buildStringQueryComponents(searchValue)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf(`{"GetBooksByTitle": "%s"}`, err))
		return
	}

	// build GetParam
	p := sqac.GetParam{
		FieldName:    "title",
		Operand:      op,
		ParamValue:   predicate,
		NextOperator: "",
	}
	params := []sqac.GetParam{}
	params = append(params, p)

	// build base Href; common for each selected row
	urlString := buildHrefBasic(r, true)

	// call the common Book GetSet method
	books, count, countReq, err := bc.getBookSet(w, r, params)
	if books != nil && countReq == false {
		for i, l := range books {
			books[i].Href = urlString + "book/" + strconv.FormatUint(uint64(l.ID), 10)
		}
		respondWithJSON(w, http.StatusOK, books)
		return
	}

	if countReq == true {
		respondWithCount(w, http.StatusOK, count)
		return
	}
	respondWithJSON(w, http.StatusOK, "[]")
}

// GetBooksByAuthor facilitates the retrieval of existing
// Books based on Author.
// GET /books/author(OP 'searchString')
// GET /books/author(OP 'searchString')/$count | $limit=n $offset=n $orderby=<field_name> ($asc|$desc)
func (bc *BookController) GetBooksByAuthor(w http.ResponseWriter, r *http.Request) {

	// get the author parameter
	vars := mux.Vars(r)
	searchValue := vars["author"]
	if searchValue == "" {
		respondWithError(w, http.StatusBadRequest, "missing search criteria")
		return
	}

	// adjust operator and predicate if neccessary
	op, predicate, err := buildStringQueryComponents(searchValue)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf(`{"GetBooksByAuthor": "%s"}`, err))
		return
	}

	// build GetParam
	p := sqac.GetParam{
		FieldName:    "author",
		Operand:      op,
		ParamValue:   predicate,
		NextOperator: "",
	}
	params := []sqac.GetParam{}
	params = append(params, p)

	// build base Href; common for each selected row
	urlString := buildHrefBasic(r, true)

	// call the common Book GetSet method
	books, count, countReq, err := bc.getBookSet(w, r, params)
	if books != nil && countReq == false {
		for i, l := range books {
			books[i].Href = urlString + "book/" + strconv.FormatUint(uint64(l.ID), 10)
		}
		respondWithJSON(w, http.StatusOK, books)
		return
	}

	if countReq == true {
		respondWithCount(w, http.StatusOK, count)
		return
	}
	respondWithJSON(w, http.StatusOK, "[]")
}

// GetBooksByHardcover facilitates the retrieval of existing
// Books based on Hardcover.

// GET /books/hardcover(OP searchValue)
// GET /books/hardcover(OP searchValue)/$count | $limit=n $offset=n $orderby=<field_name> ($asc|$desc)
func (bc *BookController) GetBooksByHardcover(w http.ResponseWriter, r *http.Request) {

	// get the hardcover parameter
	vars := mux.Vars(r)
	searchValue := vars["hardcover"]
	if searchValue == "" {
		respondWithError(w, http.StatusBadRequest, "missing search criteria")
		return
	}

	// adjust operator and predicate if neccessary
	op, predicate, err := buildBoolQueryComponents(searchValue)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf(`{"GetBooksByHardcover": "%s"}`, err))
		return
	}

	// build GetParam
	p := sqac.GetParam{
		FieldName:    "hardcover",
		Operand:      op,
		ParamValue:   predicate,
		NextOperator: "",
	}
	params := []sqac.GetParam{}
	params = append(params, p)

	// build base Href; common for each selected row
	urlString := buildHrefBasic(r, true)

	// call the common Book GetSet method
	books, count, countReq, err := bc.getBookSet(w, r, params)
	if books != nil && countReq == false {
		for i, l := range books {
			books[i].Href = urlString + "book/" + strconv.FormatUint(uint64(l.ID), 10)
		}
		respondWithJSON(w, http.StatusOK, books)
		return
	}

	if countReq == true {
		respondWithCount(w, http.StatusOK, count)
		return
	}
	respondWithJSON(w, http.StatusOK, "[]")
}

// GetBooksByCopies facilitates the retrieval of existing
// Books based on Copies.

// GET /books/copies(OP searchValue)
// GET /books/copies(OP searchValue)/$count | $limit=n $offset=n $orderby=<field_name> ($asc|$desc)
func (bc *BookController) GetBooksByCopies(w http.ResponseWriter, r *http.Request) {

	// get the copies parameter
	vars := mux.Vars(r)
	searchValue := vars["copies"]
	if searchValue == "" {
		respondWithError(w, http.StatusBadRequest, "missing search criteria")
		return
	}

	// adjust operator and predicate if neccessary
	op, predicate, err := buildUInt64QueryComponent(searchValue)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf(`{"GetBooksByCopies": "%s"}`, err))
		return
	}

	// build GetParam
	p := sqac.GetParam{
		FieldName:    "copies",
		Operand:      op,
		ParamValue:   predicate,
		NextOperator: "",
	}
	params := []sqac.GetParam{}
	params = append(params, p)

	// build base Href; common for each selected row
	urlString := buildHrefBasic(r, true)

	// call the common Book GetSet method
	books, count, countReq, err := bc.getBookSet(w, r, params)
	if books != nil && countReq == false {
		for i, l := range books {
			books[i].Href = urlString + "book/" + strconv.FormatUint(uint64(l.ID), 10)
		}
		respondWithJSON(w, http.StatusOK, books)
		return
	}

	if countReq == true {
		respondWithCount(w, http.StatusOK, count)
		return
	}
	respondWithJSON(w, http.StatusOK, "[]")
}

// GetBooksByLibraryID facilitates the retrieval of existing
// Books based on LibraryID.

// GET /books/library_id(OP searchValue)
// GET /books/library_id(OP searchValue)/$count | $limit=n $offset=n $orderby=<field_name> ($asc|$desc)
func (bc *BookController) GetBooksByLibraryID(w http.ResponseWriter, r *http.Request) {

	// get the library_id parameter
	vars := mux.Vars(r)
	searchValue := vars["library_id"]
	if searchValue == "" {
		respondWithError(w, http.StatusBadRequest, "missing search criteria")
		return
	}

	// adjust operator and predicate if neccessary
	op, predicate, err := buildUInt64QueryComponent(searchValue)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf(`{"GetBooksByLibraryID": "%s"}`, err))
		return
	}

	// build GetParam
	p := sqac.GetParam{
		FieldName:    "library_id",
		Operand:      op,
		ParamValue:   predicate,
		NextOperator: "",
	}
	params := []sqac.GetParam{}
	params = append(params, p)

	// build base Href; common for each selected row
	urlString := buildHrefBasic(r, true)

	// call the common Book GetSet method
	books, count, countReq, err := bc.getBookSet(w, r, params)
	if books != nil && countReq == false {
		for i, l := range books {
			books[i].Href = urlString + "book/" + strconv.FormatUint(uint64(l.ID), 10)
		}
		respondWithJSON(w, http.StatusOK, books)
		return
	}

	if countReq == true {
		respondWithCount(w, http.StatusOK, count)
		return
	}
	respondWithJSON(w, http.StatusOK, "[]")
}
