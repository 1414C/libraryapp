package controllers

//=============================================================================================
// base Book entity controller_relations code generated on 27 May 20 17:57 CDT
//=============================================================================================

import (
	"net/http"
	"strconv"

	"github.com/1414C/libraryapp/models"
	"github.com/1414C/lw"
	"github.com/1414C/sqac"
	"github.com/gorilla/mux"
)

// GetBookToLibrary facilitates the retrieval of Librarys related to Book
// by way of modeled 'belongsTo' relationship ToLibrary.
// This method is bound to the gorilla.mux router in appobj.go.
// 1:1 by default...
//
// GET /Book/:id/ToLibrary
// GET /Book/:id/ToLibrary/:id
func (bc *BookController) GetBookToLibrary(w http.ResponseWriter, r *http.Request) {

	var mapCommands map[string]interface{}
	var libraryID uint64
	bHaveTargetKey := false
	librarys := []models.Library{}
	countReq := false

	// read the mux vars
	vars := mux.Vars(r)

	// check that a bookid has been provided
	bookID, err := strconv.ParseUint(vars["book_id"], 10, 64)
	if err != nil {
		lw.Warning("Book Get: %s", err.Error())
		respondWithError(w, http.StatusBadRequest, "Invalid book id")
		return
	}

	// check to see if a library_id was provided
	_, ok := vars["library_id"]
	if ok {
		libraryID, err = strconv.ParseUint(vars["library_id"], 10, 64)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid libraryID")
			return
		}
		bHaveTargetKey = true
	}

	// in all cases the book must be retrieved, as the validity of the
	// the access-path must be verified.  Also consider that the library_id
	// may not have been provided.
	book := models.Book{
		ID: bookID,
	}

	// retrieve the root entity
	err = bc.bs.Get(&book)
	if err != nil {
		lw.Warning(err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// if the parent-entity was provided, check it against the
	// value contained in book.  From this point on,
	// libraryID is used as the target-entity's key.
	if bHaveTargetKey {
		if libraryID != book.LibraryID {
			respondWithError(w, http.StatusBadRequest, "Invalid libraryID")
			return
		}
	} else {
		libraryID = book.LibraryID
	}

	//// parse commands ($cmd) for the belongsTo selection
	//if vars != nil {
	//	_, ok := vars["cmd"]
	//	if ok {
	//		mapCommands, err = parseRequestCommands(vars)
	//		if err != nil {
	//			respondWithError(w, http.StatusBadRequest, err.Error())
	//			return
	//		}
	//	}
	//}
	//
	//// was $count requested?
	//_, ok = mapCommands["count"]
	//if ok {
	//	countReq = true
	//}

	// if there is no relationship-key, return nothing
	if libraryID == 0 {
		// fallthrough and return nothing
		respondWithJSON(w, http.StatusOK, "[]")
		return
	}

	libraryParams := []sqac.GetParam{}
	libraryParam := sqac.GetParam{
		FieldName:    "ID",
		Operand:      "=",
		ParamValue:   libraryID,
		NextOperator: "",
	}
	libraryParams = append(libraryParams, libraryParam)

	// build the root href for each book
	urlString := buildHrefBasic(r, true)
	lw.Debug("urlString: %s", urlString)
	urlString = urlString + "library/"

	// call the ORM to retrieve the librarys or count
	librarys, count := bc.svcs.Library.GetLibrarys(libraryParams, mapCommands)
	lw.Debug("mapCommands: %v", mapCommands)
	lw.Debug("countReq: %v", countReq)

	// retrieved []Library and not asked to $count
	if librarys != nil && countReq == false {
		for i, l := range librarys {
			librarys[i].Href = urlString + strconv.FormatUint(uint64(l.ID), 10)
		}

		// send the result -
		if bHaveTargetKey || len(librarys) == 1 {
			respondWithJSON(w, http.StatusOK, librarys[0])
			return
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
