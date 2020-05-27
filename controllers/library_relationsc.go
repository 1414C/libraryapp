package controllers

//=============================================================================================
// base Library entity controller_relations code generated on 27 May 20 17:57 CDT
//=============================================================================================

import (
	"net/http"
	"strconv"

	"github.com/1414C/libraryapp/models"
	"github.com/1414C/lw"
	"github.com/1414C/sqac"
	"github.com/gorilla/mux"
)

// GetLibraryToBooks facilitates the retrieval of Books related to Library
// by way of modeled 'hasMany' relationship ToBooks.
// This method is bound to the gorilla.mux router in appobj.go.
// 1:N
//
// GET /Library/:id/ToBooks
// GET /Library/:id/ToBooks/$count | $limit=n $offset=n $orderby=<field_name> ($asc|$desc)
// GET /Library/:id/ToBooks/:id
func (lc *LibraryController) GetLibraryToBooks(w http.ResponseWriter, r *http.Request) {

	var mapCommands map[string]interface{}
	var bookID uint64
	bSingle := false
	books := []models.Book{}
	countReq := false

	vars := mux.Vars(r)

	// check that a library_id has been provided (root entity id)
	libraryID, err := strconv.ParseUint(vars["library_id"], 10, 64)
	if err != nil {
		lw.Warning("Library Get: %s", err.Error())
		respondWithError(w, http.StatusBadRequest, "Invalid library number")
		return
	}

	// check to see if a book_id was provided
	_, ok := vars["book_id"]
	if ok {
		bookID, err = strconv.ParseUint(vars["book_id"], 10, 64)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid bookID")
			return
		}
		bSingle = true
	}

	// in all cases the library must be retrieved, as the validity of the
	// the access-path must be verified.  Also consider that the book
	// :id may not have been provided.
	library := models.Library{
		ID: libraryID,
	}

	// retrieve the root entity
	err = lc.ls.Get(&library)
	if err != nil {
		lw.Warning(err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// parse commands ($cmd) for the toMany selection
	if vars != nil {
		_, ok := vars["cmd"]
		if ok {
			mapCommands, err = parseRequestCommands(vars)
			if err != nil {
				respondWithError(w, http.StatusBadRequest, err.Error())
				return
			}
		}
	}

	// was $count requested?
	_, ok = mapCommands["count"]
	if ok {
		countReq = true
	}

	// add the root entity-key to the selection parameter list
	bookParams := []sqac.GetParam{}
	bookParam := sqac.GetParam{
		FieldName:    "LibraryID",
		Operand:      "=",
		ParamValue:   libraryID,
		NextOperator: "",
	}

	// if the child-entity-key was not provided, append the selection
	// parameter to the parameter list and then call the GET.
	// if the child-entity-key was provide, set the NextOperator to
	// 'AND', then add the child-entity-key to the parameter-list.
	if bookID == 0 {
		bookParams = append(bookParams, bookParam)
	} else {
		bookParam.NextOperator = "AND"
		bookParams = append(bookParams, bookParam)

		bookParam = sqac.GetParam{
			FieldName:    "ID",
			Operand:      "=",
			ParamValue:   bookID,
			NextOperator: "",
		}
		bookParams = append(bookParams, bookParam)
	}

	// build the root href for each book
	urlString := buildHrefBasic(r, true)
	lw.Debug("urlString: %s", urlString)
	urlString = urlString + "book/"

	// call the ORM to retrieve the books or count
	books, count := lc.svcs.Book.GetBooks(bookParams, mapCommands)
	lw.Debug("mapCommands: %v", mapCommands)
	lw.Debug("countReq: %v", countReq)
	// retrieved []Book and not asked to $count
	if books != nil && countReq == false {
		for i, l := range books {
			books[i].Href = urlString + strconv.FormatUint(uint64(l.ID), 10)
		}

		// send the result(s)
		if bSingle {
			respondWithJSON(w, http.StatusOK, books[0])
			return
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
