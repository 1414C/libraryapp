package controllers

//=============================================================================================
// Auth entity controller code generated on 27 May 20 17:32 CDT
//=============================================================================================

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/1414C/libraryapp/group/gmcl"
	"github.com/1414C/libraryapp/group/gmcom"
	"github.com/1414C/libraryapp/models"
	"github.com/1414C/lw"
	"github.com/gorilla/mux"
)

// AuthController is the Auth controller type for route binding
type AuthController struct {
	as              models.AuthService
	internalAddress string
	AuthsH          *gmcom.AuthsH // cache
}

// NewAuthController creates a new AuthController
func NewAuthController(as models.AuthService, internalAddress string) *AuthController {
	return &AuthController{
		as:              as,
		internalAddress: internalAddress,
	}
}

// Create facilitates the creation of a new Auth.  This method is bound
// to the gorilla.mux router in main.go.
//
// POST /auth
func (ac *AuthController) Create(w http.ResponseWriter, r *http.Request) {

	var a models.Auth
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&a); err != nil {
		lw.ErrorWithPrefixString("Auth Resource Create() got:", err)
		respondWithError(w, http.StatusBadRequest, "Authc: Invalid request payload")
		return
	}
	defer r.Body.Close()

	// fill the model
	auth := models.Auth{
		AuthName:    a.AuthName,
		AuthType:    a.AuthType,
		Description: a.Description,
	}

	// build a base urlString for the JSON Body self-referencing Href tag
	urlString := buildHrefStringFromCRUDReq(r, true)

	// call the Create method on the usrgroup model
	err := ac.as.Create(&auth)
	if err != nil {
		lw.ErrorWithPrefixString("Auth Resource Create() got:", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	auth.Href = urlString + strconv.FormatUint(uint64(auth.ID), 10)
	respondWithJSON(w, http.StatusCreated, auth)

	// disseminate the new auth info to self and group-members if any
	ac.disseminateAuthChange(auth.ID, true, auth.AuthName, gmcom.COpCreate)
}

// Update facilitates the update of an existing Auth.  This method is bound
// to the gorilla.mux router in main.go.
//
// PUT /auth:id
func (ac *AuthController) Update(w http.ResponseWriter, r *http.Request) {

	// get the parameter(s)
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		lw.ErrorWithPrefixString("Auth Resource Update() got:", err)
		respondWithError(w, http.StatusBadRequest, "Invalid auth id")
		return
	}

	var a models.Auth
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&a); err != nil {
		lw.ErrorWithPrefixString("Auth Resource Update() got:", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// fill the model
	auth := models.Auth{
		AuthName:    a.AuthName,
		AuthType:    a.AuthType,
		Description: a.Description,
	}

	// build a base urlString for the JSON Body self-referencing Href tag
	urlString := buildHrefStringFromCRUDReq(r, false)

	auth.ID = id

	// call the update method on the model
	err = ac.as.Update(&auth)
	if err != nil {
		lw.ErrorWithPrefixString("Auth Resource Update() got:", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	auth.Href = urlString
	respondWithJSON(w, http.StatusCreated, auth)

	// disseminate the updated auth info to self and group-members if any
	ac.disseminateAuthChange(auth.ID, true, auth.AuthName, gmcom.COpUpdate)
}

// Get facilitates the retrieval of an existing Auth.  This method is bound
// to the gorilla.mux router in main.go.
//
// GET /auth/:id
func (ac *AuthController) Get(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		lw.ErrorWithPrefixString("Auth Resource Get() got:", err)
		respondWithError(w, http.StatusBadRequest, "Invalid auth resource ID")
		return
	}

	// build a base urlString for the JSON Body self-referencing Href tag
	urlString := buildHrefStringFromCRUDReq(r, false)

	auth := models.Auth{
		ID: id,
	}

	err = ac.as.Get(&auth)
	if err != nil {
		lw.Error(err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	auth.Href = urlString
	respondWithJSON(w, http.StatusCreated, auth)
}

// Delete facilitates the deletion of an existing Auth.  This method is bound
// to the gorilla.mux router in main.go.
//
// DELETE /auth/:id
func (ac *AuthController) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		lw.ErrorWithPrefixString("Auth Resource Delete() got:", err)
		respondWithError(w, http.StatusBadRequest, "Invalid Auth ID")
		return
	}

	auth := models.Auth{
		ID: id,
	}

	err = ac.as.Delete(&auth)
	if err != nil {
		lw.ErrorWithPrefixString("Auth Resource Delete() got:", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithHeader(w, http.StatusAccepted)

	// disseminate the deleted auth info to self and group-members if any
	ac.disseminateAuthChange(auth.ID, true, auth.AuthName, gmcom.COpDelete)
}

// GetAuths facilitates the retrieval of all existing Auths.  This method is
// bound to the gorilla.mux router in main.go.
//
// GET /auths
func (ac *AuthController) GetAuths(w http.ResponseWriter, r *http.Request) {

	// build base Href; common for each selected row
	urlString := buildHrefStringFromCRUDReq(r, true)
	urlString = strings.TrimSuffix(urlString, "s/")
	urlString = urlString + "/"

	auths := ac.as.GetAuths()
	if auths != nil {
		for i, u := range auths {
			auths[i].Href = urlString + strconv.FormatUint(uint64(u.ID), 10)
		}
		respondWithJSON(w, http.StatusOK, auths)
		return
	}
	respondWithJSON(w, http.StatusOK, "[]")
}

// GetAuthsByAuthName facilitates the retrieval of existing
// Auths based on AuthName.
// GET /auths/auth_name(OP 'searchString')
func (ac *AuthController) GetAuthsByAuthName(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	searchValue := vars["auth_name"]
	if searchValue == "" {
		respondWithError(w, http.StatusBadRequest, "missing search criteria")
		return
	}

	op, predicate, err := buildStringQueryComponents(searchValue)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf(`{"GetAuthByAuthName": "%s"}`, err))
		return
	}

	// build base Href; common for each selected row
	urlString, err := buildHrefStringFromSimpleQueryReq(r, true)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf(`{"GetAuthsByAuthName": "%s"}`, err))
		return
	}

	auths := ac.as.GetAuthsByAuthName(op, predicate)
	if auths != nil {

		// add the base Href/{id}
		for i, b := range auths {
			auths[i].Href = urlString + strconv.FormatUint(uint64(b.ID), 10)
		}
		respondWithJSON(w, http.StatusOK, auths)
		return
	}
	respondWithJSON(w, http.StatusOK, "[]")
}

// GetAuthsByDescription facilitates the retrieval of existing
// Auths based on Description.
// GET /auths/description(OP 'searchString')
func (ac *AuthController) GetAuthsByDescription(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	searchValue := vars["description"]
	if searchValue == "" {
		respondWithError(w, http.StatusBadRequest, "missing search criteria")
		return
	}

	op, predicate, err := buildStringQueryComponents(searchValue)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf(`{"GetAuthsByDescription": "%s"}`, err))
		return
	}

	// build base Href; common for each selected row
	urlString, err := buildHrefStringFromSimpleQueryReq(r, true)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf(`{"GetAuthsByDescription": "%s"}`, err))
		return
	}

	auths := ac.as.GetAuthsByDescription(op, predicate)
	if auths != nil {

		// add the base Href/{id}
		for i, b := range auths {
			auths[i].Href = urlString + strconv.FormatUint(uint64(b.ID), 10)
		}
		respondWithJSON(w, http.StatusOK, auths)
		return
	}
	respondWithJSON(w, http.StatusOK, "[]")
}

// disseminateAuthChange updates the local AuthController's Auths cache directly,
// and then forwards the updated Auth information to all other group-members that
// are in a non-failed status.
func (ac *AuthController) disseminateAuthChange(id uint64, fwd bool, authName string, op gmcom.OpType) {

	// make sure the auth info is in the local and group caches
	aa := gmcom.AuthD{
		Forward:  fwd,
		ID:       id,
		AuthName: authName,
		Op:       op,
	}

	err := gmcl.AddUpdAuthCache(aa, ac.internalAddress)
	if err != nil {
		lw.ErrorWithPrefixString("AuthController cache update error message:", err)
	}
}
