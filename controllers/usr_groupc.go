package controllers

//=============================================================================================
// base Book entity controller code generated on 27 May 20 17:57 CDT
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

// UsrGroupController is the usrGroup controller type for route binding
type UsrGroupController struct {
	us              models.UsrGroupService
	internalAddress string
	UsrGroupsH      *gmcom.UsrGroupsH // cache
}

// NewUsrGroupController creates a new UsrGroupController
func NewUsrGroupController(us models.UsrGroupService, internalAddress string) *UsrGroupController {
	return &UsrGroupController{
		us:              us,
		internalAddress: internalAddress,
	}
}

// Create facilitates the creation of a new UsrGroup.  This method is bound
// to the gorilla.mux router in main.go.
//
// POST /usrgroup
func (uc *UsrGroupController) Create(w http.ResponseWriter, r *http.Request) {

	var u models.UsrGroup
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		lw.ErrorWithPrefixString("User Group Create:", err)
		respondWithError(w, http.StatusBadRequest, "usrGroupc: Invalid request payload")
		return
	}
	defer r.Body.Close()

	// fill the model
	usrgroup := models.UsrGroup{
		GroupName:   u.GroupName,
		Description: u.Description,
	}

	// build a base urlString for the JSON Body self-referencing Href tag
	urlString := buildHrefStringFromCRUDReq(r, true)

	// call the Create method on the usrgroup model
	err := uc.us.Create(&usrgroup)
	if err != nil {
		lw.ErrorWithPrefixString("User Group Create:", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	usrgroup.Href = urlString + strconv.FormatUint(uint64(usrgroup.ID), 10)
	respondWithJSON(w, http.StatusCreated, usrgroup)

	// disseminate the new usrgroup info to self and group-members if any
	uc.disseminateUsrGroupChange(usrgroup.ID, true, usrgroup.GroupName, gmcom.COpCreate)
}

// Update facilitates the update of an existing UsrGroup.  This method is bound
// to the gorilla.mux router in main.go.
//
// PUT /usrgoup:id
func (uc *UsrGroupController) Update(w http.ResponseWriter, r *http.Request) {

	// get the parameter(s)
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		lw.Warning("User Group Update:", err)
		respondWithError(w, http.StatusBadRequest, "Invalid usrgroup id")
		return
	}

	var u models.UsrGroup
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		lw.Warning("User Group Update:", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// fill the model
	usrgroup := models.UsrGroup{
		GroupName:   u.GroupName,
		Description: u.Description,
	}

	// build a base urlString for the JSON Body self-referencing Href tag
	urlString := buildHrefStringFromCRUDReq(r, false)

	usrgroup.ID = id

	// call the update method on the model
	err = uc.us.Update(&usrgroup)
	if err != nil {
		lw.ErrorWithPrefixString("User Group Update:", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	usrgroup.Href = urlString
	respondWithJSON(w, http.StatusCreated, usrgroup)

	// disseminate the updated usrgroup info to self and group-members if any
	uc.disseminateUsrGroupChange(usrgroup.ID, true, usrgroup.GroupName, gmcom.COpUpdate)
}

// Get facilitates the retrieval of an existing UsrGroup.  This method is bound
// to the gorilla.mux router in main.go.
//
// GET /usrgroup/:id
func (uc *UsrGroupController) Get(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		lw.Warning("User Group Get:", err)
		respondWithError(w, http.StatusBadRequest, "Invalid usrgroup ID")
		return
	}

	// build a base urlString for the JSON Body self-referencing Href tag
	urlString := buildHrefStringFromCRUDReq(r, false)

	usrgroup := models.UsrGroup{
		ID: id,
	}

	err = uc.us.Get(&usrgroup)
	if err != nil {
		lw.Warning(err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	usrgroup.Href = urlString
	respondWithJSON(w, http.StatusCreated, usrgroup)
}

// Delete facilitates the deletion of an existing UsrGroup.  This method is bound
// to the gorilla.mux router in main.go.
//
// DELETE /usrgroup/:id
func (uc *UsrGroupController) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		lw.Warning("User Group Delete:", err)
		respondWithError(w, http.StatusBadRequest, "Invalid UsrGroup ID")
		return
	}

	usrgroup := models.UsrGroup{
		ID: id,
	}

	err = uc.us.Delete(&usrgroup)
	if err != nil {
		lw.ErrorWithPrefixString("User Group Delete:", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithHeader(w, http.StatusAccepted)

	// disseminate the deleted usrgroup info to self and group-members if any
	uc.disseminateUsrGroupChange(usrgroup.ID, true, usrgroup.GroupName, gmcom.COpDelete)
}

// GetUsrGroups facilitates the retrieval of all existing UsrGroups.  This method is bound
// to the gorilla.mux router in main.go.
//
// GET /usrgroups
func (uc *UsrGroupController) GetUsrGroups(w http.ResponseWriter, r *http.Request) {

	// build base Href; common for each selected row
	urlString := buildHrefStringFromCRUDReq(r, true)
	urlString = strings.TrimSuffix(urlString, "s/")
	urlString = urlString + "/"

	usrgroups := uc.us.GetUsrGroups()
	if usrgroups != nil {
		for i, u := range usrgroups {
			usrgroups[i].Href = urlString + strconv.FormatUint(uint64(u.ID), 10)
		}
		respondWithJSON(w, http.StatusOK, usrgroups)
		return
	}
	respondWithJSON(w, http.StatusOK, "[]")
}

// GetUsrGroupsByGroupName facilitates the retrieval of existing
// UsrGroups based on GroupName.
// GET /usrgroups/group_name(OP 'searchString')
func (uc *UsrGroupController) GetUsrGroupsByGroupName(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	searchValue := vars["group_name"]
	if searchValue == "" {
		respondWithError(w, http.StatusBadRequest, "missing search criteria")
		return
	}

	op, predicate, err := buildStringQueryComponents(searchValue)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf(`{"GetUsrGroupsByGroupName": "%s"}`, err))
		return
	}

	// build base Href; common for each selected row
	urlString, err := buildHrefStringFromSimpleQueryReq(r, true)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf(`{"GetUsrGroupsByGroupName": "%s"}`, err))
		return
	}

	usrgroups := uc.us.GetUsrGroupsByGroupName(op, predicate)
	if usrgroups != nil {

		// add the base Href/{id}
		for i, b := range usrgroups {
			usrgroups[i].Href = urlString + strconv.FormatUint(uint64(b.ID), 10)
		}
		respondWithJSON(w, http.StatusOK, usrgroups)
		return
	}
	respondWithJSON(w, http.StatusOK, "[]")
}

// GetUsrGroupsByDescription facilitates the retrieval of existing
// UsrGroups based on Description.
// GET /usrgroups/description(OP 'searchString')
func (uc *UsrGroupController) GetUsrGroupsByDescription(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	searchValue := vars["description"]
	if searchValue == "" {
		respondWithError(w, http.StatusBadRequest, "missing search criteria")
		return
	}

	op, predicate, err := buildStringQueryComponents(searchValue)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf(`{"GetUsrGroupsByDescription": "%s"}`, err))
		return
	}

	// build base Href; common for each selected row
	urlString, err := buildHrefStringFromSimpleQueryReq(r, true)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf(`{"GetUsrGroupsByDescription": "%s"}`, err))
		return
	}

	usrgroups := uc.us.GetUsrGroupsByDescription(op, predicate)
	if usrgroups != nil {

		// add the base Href/{id}
		for i, b := range usrgroups {
			usrgroups[i].Href = urlString + strconv.FormatUint(uint64(b.ID), 10)
		}
		respondWithJSON(w, http.StatusOK, usrgroups)
		return
	}
	respondWithJSON(w, http.StatusOK, "[]")
}

// disseminateUsrGroupChange updates the local UsrGroupController's Groups cache directly,
// and then forwards the updated UsrGroup information to all other group-members that
// are in a non-failed status.
func (uc *UsrGroupController) disseminateUsrGroupChange(id uint64, fwd bool, groupName string, op gmcom.OpType) {

	// make sure the usrgroup info is in the local and group caches
	aa := gmcom.UsrGroupD{
		Forward:   fwd,
		ID:        id,
		GroupName: groupName,
		Op:        op,
	}

	err := gmcl.AddUpdUsrGroupCache(aa, uc.internalAddress)
	if err != nil {
		lw.ErrorWithPrefixString("UsrGroupController cache update error message:", err)
	}
}
