package controllers

//=============================================================================================
// base Library entity controller code generated on 27 May 20 17:57 CDT
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

// LibraryController is the library controller type for route binding
type LibraryController struct {
	ls   models.LibraryService
	ep   ext.LibraryCtrlExt
	svcs models.Services
}

// NewLibraryController creates a new LibraryController
func NewLibraryController(ls models.LibraryService, svcs models.Services) *LibraryController {
	return &LibraryController{
		ls:   ls,
		ep:   *ext.InitLibraryCtrlExt(),
		svcs: svcs,
	}
}

// Create facilitates the creation of a new Library.  This method is bound
// to the gorilla.mux router in main.go.
//
// POST /library
func (lc *LibraryController) Create(w http.ResponseWriter, r *http.Request) {

	var err error
	var lm models.Library

	// TODO: implement extension-point if required
	// TODO: safe to comment this block out if the extension-point is not needed
	err = lc.ep.CrtEp.BeforeFirst(w, r)
	if err != nil {
		lw.ErrorWithPrefixString("LibraryController CreateBeforeFirst() error:", err)
		respondWithError(w, http.StatusBadRequest, "libraryc: Invalid request")
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&lm); err != nil {
		lw.ErrorWithPrefixString("Library Create:", err)
		respondWithError(w, http.StatusBadRequest, "libraryc: Invalid request payload")
		return
	}
	defer r.Body.Close()

	// TODO: implement extension-point if required
	// TODO: safe to comment this block out if the extension-point is not needed
	err = lc.ep.CrtEp.AfterBodyDecode(&lm)
	if err != nil {
		lw.ErrorWithPrefixString("LibraryController CreateAfterBodyDecode() error:", err)
		respondWithError(w, http.StatusBadRequest, "libraryc: Invalid request payload")
		return
	}

	// fill the model
	library := models.Library{
		Name: lm.Name,
		City: lm.City,
	}

	// build a base urlString for the JSON Body self-referencing Href tag
	urlString := buildHrefStringFromCRUDReq(r, true)

	// call the Create method on the library model
	err = lc.ls.Create(&library)
	if err != nil {
		lw.ErrorWithPrefixString("Library Create:", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	library.Href = urlString + strconv.FormatUint(uint64(library.ID), 10)

	// TODO: implement extension-point if required
	// TODO: safe to comment this block out if the extension-point is not needed
	err = lc.ep.CrtEp.BeforeResponse(&library)
	if err != nil {
		lw.ErrorWithPrefixString("LibraryController CreateBeforeResponse() error:", err)
		respondWithError(w, http.StatusBadRequest, "libraryc: Invalid request")
		return
	}
	respondWithJSON(w, http.StatusCreated, library)
}

// Update facilitates the update of an existing Library.  This method is bound
// to the gorilla.mux router in main.go.
//
// PUT /library:id
func (lc *LibraryController) Update(w http.ResponseWriter, r *http.Request) {

	var err error
	var lm models.Library

	// TODO: implement extension-point if required
	// TODO: safe to comment this block out if the extension-point is not needed
	err = lc.ep.UpdEp.BeforeFirst(w, r)
	if err != nil {
		lw.ErrorWithPrefixString("LibraryController UpdateBeforeFirst() error:", err)
		respondWithError(w, http.StatusBadRequest, "libraryc: Invalid request")
		return
	}

	// get the parameter(s)
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		lw.ErrorWithPrefixString("Library Update:", err)
		respondWithError(w, http.StatusBadRequest, "Invalid library id")
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&lm); err != nil {
		lw.ErrorWithPrefixString("Library Update:", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// TODO: implement extension-point if required
	// TODO: safe to comment this block out if the extension-point is not needed
	err = lc.ep.UpdEp.AfterBodyDecode(&lm)
	if err != nil {
		lw.ErrorWithPrefixString("LibraryController UpdateAfterBodyDecode() error:", err)
		respondWithError(w, http.StatusBadRequest, "libraryc: Invalid request payload")
		return
	}

	// fill the model
	library := models.Library{
		Name: lm.Name,
		City: lm.City,
	}

	// build a base urlString for the JSON Body self-referencing Href tag
	urlString := buildHrefStringFromCRUDReq(r, false)
	library.ID = id

	// call the update method on the model
	err = lc.ls.Update(&library)
	if err != nil {
		lw.ErrorWithPrefixString("Library Update:", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	library.Href = urlString

	// TODO: implement extension-point if required
	// TODO: safe to comment this block out if the extension-point is not needed
	err = lc.ep.UpdEp.BeforeResponse(&library)
	if err != nil {
		lw.ErrorWithPrefixString("LibraryController UpdateBeforeResponse() error:", err)
		respondWithError(w, http.StatusBadRequest, "libraryc: Invalid request")
		return
	}
	respondWithJSON(w, http.StatusCreated, library)
}

// Get facilitates the retrieval of an existing Library.  This method is bound
// to the gorilla.mux router in main.go.
//
// GET /library/:id
func (lc *LibraryController) Get(w http.ResponseWriter, r *http.Request) {

	var err error

	// TODO: implement extension-point if required
	// TODO: safe to comment this block out if the extension-point is not needed
	err = lc.ep.GetEp.BeforeFirst(w, r)
	if err != nil {
		lw.Warning("LibraryController GetBeforeFirst() error: %s", err.Error())
		respondWithError(w, http.StatusBadRequest, "libraryc: Invalid request")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		lw.Warning("Library Get: %s", err.Error())
		respondWithError(w, http.StatusBadRequest, "Invalid library ID")
		return
	}

	// build a base urlString for the JSON Body self-referencing Href tag
	urlString := buildHrefStringFromCRUDReq(r, false)

	library := models.Library{
		ID: id,
	}

	// TODO: implement extension-point if required
	// TODO: safe to comment this block out if the extension-point is not needed
	err = lc.ep.GetEp.BeforeModelCall(&library)
	if err != nil {
		lw.Warning("LibraryController GetBeforeModelCall() error: %s", err.Error())
		respondWithError(w, http.StatusBadRequest, "libraryc: Invalid request")
		return
	}

	err = lc.ls.Get(&library)
	if err != nil {
		lw.Warning(err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	library.Href = urlString

	// TODO: implement extension-point if required
	// TODO: safe to comment this block out if the extension-point is not needed
	err = lc.ep.GetEp.BeforeResponse(&library)
	if err != nil {
		lw.Warning("LibraryController GetBeforeResponse() error: %s", err.Error())
		respondWithError(w, http.StatusBadRequest, "libraryc: Invalid request")
		return
	}
	respondWithJSON(w, http.StatusCreated, library)
}

// Delete facilitates the deletion of an existing Library.  This method is bound
// to the gorilla.mux router in main.go.
//
// DELETE /library/:id
func (lc *LibraryController) Delete(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		lw.ErrorWithPrefixString("Library Delete:", err)
		respondWithError(w, http.StatusBadRequest, "Invalid Library ID")
		return
	}

	library := models.Library{
		ID: id,
	}

	err = lc.ls.Delete(&library)
	if err != nil {
		lw.ErrorWithPrefixString("Library Delete:", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithHeader(w, http.StatusAccepted)
}

// getLibrarySet is used by all LibrarySet queries as a means of injecting parameters
// returns ([]Library, $count, countRequested, error)
func (lc *LibraryController) getLibrarySet(w http.ResponseWriter, r *http.Request, params []sqac.GetParam) ([]models.Library, uint64, bool, error) {

	var mapCommands map[string]interface{}
	var err error
	var urlString string
	var librarys []models.Library
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
		librarys, count = lc.ls.GetLibrarys(params, mapCommands)
	} else {
		librarys, count = lc.ls.GetLibrarys(params, nil)
	}

	// retrieved []Library and not asked to $count
	if librarys != nil && countReq == false {
		for i, l := range librarys {
			librarys[i].Href = urlString + strconv.FormatUint(uint64(l.ID), 10)
		}
		return librarys, 0, countReq, nil
	}

	// $count was requested, which trumps all other commands
	if countReq == true {
		return nil, count, countReq, nil
	}

	// fallthrough and return nothing
	return nil, 0, countReq, nil
}

// GetLibrarys facilitates the retrieval of all existing Librarys.  This method is bound
// to the gorilla.mux router in main.go.
//
// GET /librarys
// GET /librarys/$count | $limit=n $offset=n $orderby=<field_name> ($asc|$desc)
func (lc *LibraryController) GetLibrarys(w http.ResponseWriter, r *http.Request) {

	var librarys []models.Library
	var count uint64
	countReq := false

	// build base Href; common for each selected row
	urlString := buildHrefBasic(r, true)

	// call the common getLibrarySet method
	librarys, count, countReq, err := lc.getLibrarySet(w, r, nil)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf(`{"GetLibrarys": "%s"}`, err))
		return
	}

	// retrieved []Library and not asked to $count
	if librarys != nil && countReq == false {
		for i, l := range librarys {
			librarys[i].Href = urlString + strconv.FormatUint(uint64(l.ID), 10)
		}
		respondWithJSON(w, http.StatusOK, librarys)
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

// GetLibrarysByName facilitates the retrieval of existing
// Librarys based on Name.
// GET /librarys/name(OP 'searchString')
// GET /librarys/name(OP 'searchString')/$count | $limit=n $offset=n $orderby=<field_name> ($asc|$desc)
func (lc *LibraryController) GetLibrarysByName(w http.ResponseWriter, r *http.Request) {

	// get the name parameter
	vars := mux.Vars(r)
	searchValue := vars["name"]
	if searchValue == "" {
		respondWithError(w, http.StatusBadRequest, "missing search criteria")
		return
	}

	// adjust operator and predicate if neccessary
	op, predicate, err := buildStringQueryComponents(searchValue)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf(`{"GetLibrarysByName": "%s"}`, err))
		return
	}

	// build GetParam
	p := sqac.GetParam{
		FieldName:    "name",
		Operand:      op,
		ParamValue:   predicate,
		NextOperator: "",
	}
	params := []sqac.GetParam{}
	params = append(params, p)

	// build base Href; common for each selected row
	urlString := buildHrefBasic(r, true)

	// call the common Library GetSet method
	librarys, count, countReq, err := lc.getLibrarySet(w, r, params)
	if librarys != nil && countReq == false {
		for i, l := range librarys {
			librarys[i].Href = urlString + "library/" + strconv.FormatUint(uint64(l.ID), 10)
		}
		respondWithJSON(w, http.StatusOK, librarys)
		return
	}

	if countReq == true {
		respondWithCount(w, http.StatusOK, count)
		return
	}
	respondWithJSON(w, http.StatusOK, "[]")
}

// GetLibrarysByCity facilitates the retrieval of existing
// Librarys based on City.
// GET /librarys/city(OP 'searchString')
// GET /librarys/city(OP 'searchString')/$count | $limit=n $offset=n $orderby=<field_name> ($asc|$desc)
func (lc *LibraryController) GetLibrarysByCity(w http.ResponseWriter, r *http.Request) {

	// get the city parameter
	vars := mux.Vars(r)
	searchValue := vars["city"]
	if searchValue == "" {
		respondWithError(w, http.StatusBadRequest, "missing search criteria")
		return
	}

	// adjust operator and predicate if neccessary
	op, predicate, err := buildStringQueryComponents(searchValue)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf(`{"GetLibrarysByCity": "%s"}`, err))
		return
	}

	// build GetParam
	p := sqac.GetParam{
		FieldName:    "city",
		Operand:      op,
		ParamValue:   predicate,
		NextOperator: "",
	}
	params := []sqac.GetParam{}
	params = append(params, p)

	// build base Href; common for each selected row
	urlString := buildHrefBasic(r, true)

	// call the common Library GetSet method
	librarys, count, countReq, err := lc.getLibrarySet(w, r, params)
	if librarys != nil && countReq == false {
		for i, l := range librarys {
			librarys[i].Href = urlString + "library/" + strconv.FormatUint(uint64(l.ID), 10)
		}
		respondWithJSON(w, http.StatusOK, librarys)
		return
	}

	if countReq == true {
		respondWithCount(w, http.StatusOK, count)
		return
	}
	respondWithJSON(w, http.StatusOK, "[]")
}
