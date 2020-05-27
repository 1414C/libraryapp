package controllers

//=============================================================================================
// GroupAuth entity controller code generated on 27 May 20 17:32 CDT
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

// GroupAuthController is the GroupAuth controller type for route binding
type GroupAuthController struct {
	gs              models.GroupAuthService
	internalAddress string
	GroupAuthsH     *gmcom.GroupAuthsH // cache
}

// NewGroupAuthController creates a new GroupAuthController
func NewGroupAuthController(gs models.GroupAuthService, internalAddress string) *GroupAuthController {
	return &GroupAuthController{
		gs:              gs,
		internalAddress: internalAddress,
	}
}

// Create facilitates the creation of a new GroupAuth.  This method is bound
// to the gorilla.mux router in main.go.
//
// POST /groupauth
func (gc *GroupAuthController) Create(w http.ResponseWriter, r *http.Request) {

	var g models.GroupAuth
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&g); err != nil {
		lw.ErrorWithPrefixString("Group Auth Create:", err)
		respondWithError(w, http.StatusBadRequest, "GroupAuthc: Invalid request payload")
		return
	}
	defer r.Body.Close()

	// fill the model
	groupauth := models.GroupAuth{
		GroupID: g.GroupID,
		AuthID:  g.AuthID,
	}

	// build a base urlString for the JSON Body self-referencing Href tag
	urlString := buildHrefStringFromCRUDReq(r, true)

	// call the Create method on the usrgroup model
	err := gc.gs.Create(&groupauth)
	if err != nil {
		lw.ErrorWithPrefixString("Group Auth Create:", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	groupauth.Href = urlString + strconv.FormatUint(uint64(groupauth.ID), 10)
	respondWithJSON(w, http.StatusCreated, groupauth)

	// make sure the groupauth info is in the local and group caches
	// get the groupName and authNames
	gu := gmcom.GroupAuthD{
		Forward:   true,
		ID:        g.ID,
		GroupName: g.GroupName, // empty
		AuthName:  g.AuthName,  // empty
		GroupID:   g.GroupID,
		AuthID:    g.AuthID,
		Op:        gmcom.COpCreate,
	}

	err = gmcl.AddUpdGroupAuthCache(gu, gc.internalAddress)
	if err != nil {
		lw.ErrorWithPrefixString("GroupAuthController cache update error message:", err)
	}
}

// Update facilitates the update of an existing GroupAuth.  This method is bound
// to the gorilla.mux router in main.go.
//
// PUT /groupauth:id
func (gc *GroupAuthController) Update(w http.ResponseWriter, r *http.Request) {

	// get the parameter(s)
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		lw.ErrorWithPrefixString("Group Auth Update:", err)
		respondWithError(w, http.StatusBadRequest, "Invalid groupauth id")
		return
	}

	var g models.GroupAuth
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&g); err != nil {
		lw.ErrorWithPrefixString("Group Auth Update:", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// fill the model
	groupauth := models.GroupAuth{
		GroupID: g.GroupID,
		AuthID:  g.AuthID,
	}

	// build a base urlString for the JSON Body self-referencing Href tag
	urlString := buildHrefStringFromCRUDReq(r, false)

	groupauth.ID = id

	// call the update method on the model
	err = gc.gs.Update(&groupauth)
	if err != nil {
		lw.ErrorWithPrefixString("Group Auth Update:", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	groupauth.Href = urlString
	respondWithJSON(w, http.StatusCreated, groupauth)

	// update the groupauth info in the local and group caches
	// get the groupName and authNames
	gu := gmcom.GroupAuthD{
		Forward:   true,
		ID:        g.ID,
		GroupName: g.GroupName, // empty
		AuthName:  g.AuthName,  // empty
		GroupID:   g.GroupID,
		AuthID:    g.AuthID,
		Op:        gmcom.COpUpdate,
	}

	err = gmcl.AddUpdGroupAuthCache(gu, gc.internalAddress)
	if err != nil {
		lw.ErrorWithPrefixString("GroupAuthController cache update error message:", err)
	}
}

// Get facilitates the retrieval of an existing GroupAuth.  This method is bound
// to the gorilla.mux router in main.go.
//
// GET /groupauth/:id
func (gc *GroupAuthController) Get(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		lw.ErrorWithPrefixString("Group Auth Get:", err)
		respondWithError(w, http.StatusBadRequest, "Invalid group auth ID")
		return
	}

	// build a base urlString for the JSON Body self-referencing Href tag
	urlString := buildHrefStringFromCRUDReq(r, false)

	groupauth := models.GroupAuth{
		ID: id,
	}

	err = gc.gs.Get(&groupauth)
	if err != nil {
		lw.Error(err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	groupauth.Href = urlString
	respondWithJSON(w, http.StatusCreated, groupauth)
}

// Delete facilitates the deletion of an existing GroupAuth.  This method is bound
// to the gorilla.mux router in main.go.
//
// DELETE /groupauth/:id
func (gc *GroupAuthController) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		lw.ErrorWithPrefixString("Group Auth Delete:", err)
		respondWithError(w, http.StatusBadRequest, "Invalid GroupAuth ID")
		return
	}

	groupauth := models.GroupAuth{
		ID: id,
	}

	err = gc.gs.Delete(&groupauth)
	if err != nil {
		lw.ErrorWithPrefixString("Group Auth Delete:", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithHeader(w, http.StatusAccepted)

	// make sure the groupauth info is in the local and group caches
	// get the groupName and authNames
	gu := gmcom.GroupAuthD{
		Forward: true,
		ID:      groupauth.ID,
		Op:      gmcom.COpDelete,
	}

	err = gmcl.AddUpdGroupAuthCache(gu, gc.internalAddress)
	if err != nil {
		lw.ErrorWithPrefixString("GroupAuthController.Create cache update error message:", err)
	}
}

// GetGroupAuths facilitates the retrieval of all existing GroupAuths.  This method is
// bound to the gorilla.mux router in main.go.
//
// GET /groupauths
func (gc *GroupAuthController) GetGroupAuths(w http.ResponseWriter, r *http.Request) {

	// build base Href; common for each selected row
	urlString := buildHrefStringFromCRUDReq(r, true)
	urlString = strings.TrimSuffix(urlString, "s/")
	urlString = urlString + "/"

	groupauths := gc.gs.GetGroupAuths()
	if groupauths != nil {
		for i, u := range groupauths {
			groupauths[i].Href = urlString + strconv.FormatUint(uint64(u.ID), 10)
		}
		respondWithJSON(w, http.StatusOK, groupauths)
		return
	}
	respondWithJSON(w, http.StatusOK, "[]")
}

// GetGroupAuthsByAuthName facilitates the retrieval of existing
// GroupAuths based on AuthName.
// GET /groupauths/auth_name(OP 'searchString')
func (gc *GroupAuthController) GetGroupAuthsByAuthName(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	searchValue := vars["auth_name"]
	if searchValue == "" {
		respondWithError(w, http.StatusBadRequest, "missing search criteria")
		return
	}

	op, predicate, err := buildStringQueryComponents(searchValue)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf(`{"GetGroupAuthByAuthName": "%s"}`, err))
		return
	}

	// build base Href; common for each selected row
	urlString, err := buildHrefStringFromSimpleQueryReq(r, true)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf(`{"GetGroupAuthsByAuthName": "%s"}`, err))
		return
	}

	groupauths := gc.gs.GetGroupAuthsByAuthName(op, predicate)
	if groupauths != nil {

		// add the base Href/{id}
		for i, b := range groupauths {
			groupauths[i].Href = urlString + strconv.FormatUint(uint64(b.ID), 10)
		}
		respondWithJSON(w, http.StatusOK, groupauths)
		return
	}
	respondWithJSON(w, http.StatusOK, "[]")
}

// GetGroupAuthsByDescription facilitates the retrieval of existing
// GroupAuths based on Description.
// GET /groupauths/description(OP 'searchString')
func (gc *GroupAuthController) GetGroupAuthsByDescription(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	searchValue := vars["description"]
	if searchValue == "" {
		respondWithError(w, http.StatusBadRequest, "missing search criteria")
		return
	}

	op, predicate, err := buildStringQueryComponents(searchValue)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf(`{"GetGroupAuthByDescription": "%s"}`, err))
		return
	}

	// build base Href; common for each selected row
	urlString, err := buildHrefStringFromSimpleQueryReq(r, true)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf(`{"GetGroupAuthsByDescription": "%s"}`, err))
		return
	}

	groupauths := gc.gs.GetGroupAuthsByDescription(op, predicate)
	if groupauths != nil {

		// add the base Href/{id}
		for i, b := range groupauths {
			groupauths[i].Href = urlString + strconv.FormatUint(uint64(b.ID), 10)
		}
		respondWithJSON(w, http.StatusOK, groupauths)
		return
	}
	respondWithJSON(w, http.StatusOK, "[]")
}

// GetGroupAuthsByGroupID facilitates the retrieval of existing
// GroupAuths based on GroupID.
// GET /groupauths/group_id(OP :id)
func (gc *GroupAuthController) GetGroupAuthsByGroupID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	searchValue := vars["group_id"]
	if searchValue == "" {
		respondWithError(w, http.StatusBadRequest, "missing search criteria")
		return
	}

	op, predicate, err := buildStringQueryComponents(searchValue)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf(`{"GetGroupAuthByGroupID": "%s"}`, err))
		return
	}

	// build base Href; common for each selected row
	urlString, err := buildHrefStringFromSimpleQueryReq(r, true)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf(`{"GetGroupAuthsByGroupID": "%s"}`, err))
		return
	}

	groupauths := gc.gs.GetGroupAuthsByGroupID(op, predicate)
	if groupauths != nil {

		// add the base Href/{id}
		for i, b := range groupauths {
			groupauths[i].Href = urlString + strconv.FormatUint(uint64(b.ID), 10)
		}
		respondWithJSON(w, http.StatusOK, groupauths)
		return
	}
	respondWithJSON(w, http.StatusOK, "[]")
}

// DeleteGroupAuthsByGroupID facilitates the deletion of all existing Auth assignments
// to the specified Group.
// DELETE /groupauths/group_id(OP :id)
func (gc *GroupAuthController) DeleteGroupAuthsByGroupID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	groupID := vars["group_id"]
	if groupID == "" {
		respondWithError(w, http.StatusBadRequest, "missing group_id")
		return
	}

	err := gc.gs.DeleteGroupAuthsByGroupID(groupID)
	if err != nil {
		lw.ErrorWithPrefixString("Group Auth DeleteGroupAuthsByGroupID error:", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithHeader(w, http.StatusAccepted)
}
