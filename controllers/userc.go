package controllers

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/1414C/libraryapp/group/gmcl"
	"github.com/1414C/libraryapp/group/gmcom"
	"github.com/1414C/libraryapp/models"
	"github.com/1414C/lw"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

// UsrController - the usr controller type
type UsrController struct {
	us              models.UsrService
	jwtKeyMap       map[string]interface{}
	jwtSignMethod   string
	jwtLifetime     uint
	internalAddress string
	ActUsrsH        *gmcom.ActUsrsH //cache
}

// Token is the jwt return type
type Token struct {
	Token string `json:"token"`
}

// NewUsrController creates a new UsrController
func NewUsrController(us models.UsrService, jwtKeyMap map[string]interface{}, jwtSignMethod string, jwtLifetime uint, internalAddress string) *UsrController {
	lw.Console("Login() signing jwt's with %s", jwtSignMethod)
	return &UsrController{
		us:              us,
		jwtKeyMap:       jwtKeyMap,
		jwtSignMethod:   jwtSignMethod,
		jwtLifetime:     jwtLifetime,
		internalAddress: internalAddress,
	}
}

// Login - used to verify the provided email address and password
//
// POST /login
func (uc *UsrController) Login(w http.ResponseWriter, r *http.Request) {

	// parse the Usr data in JSON format from the incoming request
	var u models.Usr
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "userc: Invalid request payload")
		return
	}
	defer r.Body.Close()

	// fill the Usr model with the data needed for authentication only
	usr := models.Usr{
		Email:    u.Email,
		Password: u.Password,
	}

	// attempt to authenticate the user.  note that writing to
	// the response-writer in advance of setting the cookie will
	// cause the call to http.SetCookie(&cookie) to silently fail.
	authenticatedUsr, err := uc.us.Authenticate(usr.Email, usr.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			lw.ErrorWithPrefixString("UsrController.Login():", err)
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		default:
			lw.ErrorWithPrefixString("UsrController.Login():", err)
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	// groups
	lw.Info("Login()->user groups: %v", *authenticatedUsr.Groups)

	var tokenString string
	var httpStatus int

	lw.Info("AUTHENTICATED USR: %v", authenticatedUsr)

	// prepare claims for the token
	claims := make(jwt.MapClaims)
	claims["email"] = authenticatedUsr.Email
	claims["id"] = authenticatedUsr.ID
	claims["iat"] = time.Now().Unix()
	if uc.jwtLifetime == 0 {
		claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	} else {
		claims["exp"] = time.Now().Add(time.Minute * time.Duration(uc.jwtLifetime)).Unix()
	}

	// set custom claims; auth groups as ; separated string
	claims["Groups"] = *authenticatedUsr.Groups
	claims["uid"] = authenticatedUsr.ID

	switch uc.jwtSignMethod {
	case "ES256", "ES384", "ES512":
		tokenString, httpStatus, err = uc.signECDSA(claims)

	case "RS256", "RS384", "RS512":
		tokenString, httpStatus, err = uc.signRSA(claims)

	// case "HS256": // not supported

	default:
		lw.Warning("Authentication failure for: %v", authenticatedUsr.Email)
		respondWithError(w, http.StatusBadRequest, "authentication failure")
		return
	}

	if err != nil {
		lw.Warning("Authentication failure for: %v", authenticatedUsr.Email)
		respondWithError(w, httpStatus, err.Error())
		return
	}

	response := Token{tokenString}
	respondWithJSON(w, http.StatusOK, response)

	// update the local process's Usr cache and then forward the information
	// to all other group-members presently in a non-failed state.
	uc.disseminateUsrChange(authenticatedUsr.ID, true, authenticatedUsr.Active)
}

// signECDSA creates a jwt.Token, set the claims and the signs via the specified curve
func (uc *UsrController) signECDSA(claims jwt.MapClaims) (tokenString string, httpStatus int, err error) {

	var token *jwt.Token
	var signKey *ecdsa.PrivateKey
	var ok bool

	switch uc.jwtSignMethod {
	case "ES256":
		token = jwt.New(jwt.SigningMethodES256)
		signKey, ok = uc.jwtKeyMap["ES256SignKey"].(*ecdsa.PrivateKey)
	case "ES384":
		token = jwt.New(jwt.SigningMethodES384)
		signKey, ok = uc.jwtKeyMap["ES384SignKey"].(*ecdsa.PrivateKey)
	case "ES521":
		token = jwt.New(jwt.SigningMethodES512)
		signKey, ok = uc.jwtKeyMap["ES521SignKey"].(*ecdsa.PrivateKey)
	default:
		e := fmt.Errorf("error: unknown ECDSA signing-method %v", uc.jwtSignMethod)
		lw.Error(e)
		return "", http.StatusForbidden, e
	}
	if !ok {
		e := fmt.Errorf("error: could not read signKey in Login()")
		lw.Error(e)
		return "", http.StatusForbidden, e
	}

	token.Claims = claims
	tokenString, err = token.SignedString(signKey)
	if err != nil {
		e := fmt.Errorf("error: failed to sign jwt with ECDSA signing-method %v", uc.jwtSignMethod)
		lw.Error(e)
		return "", http.StatusInternalServerError, e
	}
	return tokenString, http.StatusOK, nil
}

// signRSA creates a jwt.Token, set the claims and the signs via the specified hash
func (uc *UsrController) signRSA(claims jwt.MapClaims) (tokenString string, httpStatus int, err error) {

	var token *jwt.Token
	var signKey *rsa.PrivateKey
	var ok bool

	switch uc.jwtSignMethod {
	case "RS256":
		token = jwt.New(jwt.SigningMethodRS256)
		signKey, ok = uc.jwtKeyMap["RS256SignKey"].(*rsa.PrivateKey)
	case "RS384":
		token = jwt.New(jwt.SigningMethodRS384)
		signKey, ok = uc.jwtKeyMap["RS384SignKey"].(*rsa.PrivateKey)
	case "RS512":
		token = jwt.New(jwt.SigningMethodRS512)
		signKey, ok = uc.jwtKeyMap["RS512SignKey"].(*rsa.PrivateKey)
	default:
		e := fmt.Errorf("error: unknown RSA signing-method %v", uc.jwtSignMethod)
		lw.Error(e)
		return "", http.StatusForbidden, e
	}
	if !ok {
		e := fmt.Errorf("error: could not read signKey in Login()")
		lw.Error(e)
		return "", http.StatusForbidden, e
	}

	token.Claims = claims
	tokenString, err = token.SignedString(signKey)
	if err != nil {
		e := fmt.Errorf("error: failed to sign jwt with ECDSA signing-method %v", uc.jwtSignMethod)
		lw.Error(e)
		return "", http.StatusInternalServerError, e
	}
	return tokenString, http.StatusOK, nil
}

// signHMac will one day create a jwt.Token, set the claims and the sign via the shared-secret
func (uc *UsrController) signHmac() {

}

// Create - process the signup when a usr attempts to create a new usr account.
// It is debatable whether Create should be supported as a general usr creation
// mechanism.
//
// POST /signup
func (uc *UsrController) Create(w http.ResponseWriter, r *http.Request) {

	// parse the Usr data in JSON format from the incoming request
	var u models.Usr
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		lw.Warning("usrc: Invalid request payload")
		respondWithError(w, http.StatusBadRequest, "usrc: Invalid request payload")
		return
	}
	defer r.Body.Close()

	// fill the Usr model
	usr := models.Usr{
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
		Active:   u.Active,
		Groups:   u.Groups,
	}

	// build a base urlString for the JSON Body self-referencing Href tag
	urlString := buildHrefStringFromCRUDReq(r, false)

	// call the create method on the usr model
	err := uc.us.Create(&usr)
	if err != nil {
		lw.Warning(err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	usr.Password = ""
	usr.PasswordHash = ""
	usr.Href = urlString + "/" + strconv.FormatUint(uint64(usr.ID), 10)
	respondWithJSON(w, http.StatusCreated, usr)

	// disseminate the new user info to self and group-members if any
	uc.disseminateUsrChange(usr.ID, true, usr.Active)
}

// Update facilitates the update of an existing Usr.  This method is bound
// to the gorilla.mux router in main.go.
//
// PUT /usr:id
func (uc *UsrController) Update(w http.ResponseWriter, r *http.Request) {

	// get the parameter(s)
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		lw.Warning("Usr Update:", err)
		respondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	var u models.Usr
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&u); err != nil {
		lw.Warning("Usr Update:", err)
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	// when a user-id is provided in the body, make sure that
	// it matches the user-id sent in the request params
	if u.ID != 0 && u.ID != id {
		lw.Warning("usr update target %d had a non-zero and non-matching user-id in the request body\n", id)
		respondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	// get the existing Usr record from the db - byID
	staleUsr := models.Usr{
		ID: id,
	}

	err = uc.us.Get(&staleUsr)
	if err != nil {
		lw.Warning("Usr Update:", err)
		respondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	// fill the model
	uTime := time.Now()
	// utcTime := uTime.UTC()
	usr := models.Usr{
		ID:           id,
		Name:         u.Name,
		Email:        u.Email,
		PasswordHash: staleUsr.PasswordHash,
		CreatedOn:    staleUsr.CreatedOn,
		UpdatedOn:    &uTime,
		Active:       u.Active,
		Groups:       u.Groups,
	}

	// build a base urlString for the JSON Body self-referencing Href tag
	urlString := buildHrefStringFromCRUDReq(r, false)

	// call the update method on the model
	err = uc.us.Update(&usr)
	if err != nil {
		lw.ErrorWithPrefixString("User Update:", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// remove the password info from the response data
	usr.PasswordHash = ""
	usr.PasswordHash = ""
	usr.Href = urlString
	respondWithJSON(w, http.StatusCreated, usr)

	// disseminate the updated user info to self and group-members if any
	uc.disseminateUsrChange(usr.ID, true, usr.Active)
}

// Get facilitates the retrieval of an existing Usr.  This method is bound
// to the gorilla.mux router in main.go.
//
// GET /usr/:id
func (uc *UsrController) Get(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	// build a base urlString for the JSON Body self-referencing Href tag

	urlString := buildHrefStringFromCRUDReq(r, false)

	usr := models.Usr{
		ID: id,
	}

	err = uc.us.Get(&usr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	usr.PasswordHash = ""
	usr.PasswordHash = ""
	usr.Href = urlString
	respondWithJSON(w, http.StatusCreated, usr)
}

// Delete facilitates the deletion of an existing Usr.  This method is bound
// to the gorilla.mux router in main.go.
//
// DELETE /usr/:id
func (uc *UsrController) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	usr := models.Usr{
		ID: id,
	}

	err = uc.us.Delete(&usr)
	if err != nil {
		if err != nil {
			lw.Error(err)
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
	}
	respondWithHeader(w, http.StatusAccepted)

	// disseminate the user deletion info to self and group-members if any.
	uc.disseminateUsrChange(usr.ID, true, false)
}

// GetUsrs facilitates the retrieval of all existing Usrs.  This method is
// bound to the gorilla.mux router in main.go.
//
// GET /usrs
func (uc *UsrController) GetUsrs(w http.ResponseWriter, r *http.Request) {

	// build base Href; common for each selected row
	urlString := buildHrefStringFromCRUDReq(r, true)
	urlString = strings.TrimSuffix(urlString, "s/")
	urlString = urlString + "/"

	usrs := uc.us.GetUsrs()
	if usrs != nil {
		for i, u := range usrs {
			usrs[i].Href = urlString + strconv.FormatUint(uint64(u.ID), 10)
		}
		respondWithJSON(w, http.StatusOK, usrs)
		return
	}
	respondWithJSON(w, http.StatusOK, "[]")
}

func fatal(err error) {
	if err != nil {
		lw.Fatal(err)
	}
}

// disseminateUsrChange updates the local UsrController's ActUsrs cache directly,
// and then forwards the updated Usr information to all other group-members that
// are in a non-failed status.
func (uc *UsrController) disseminateUsrChange(id uint64, fwd, active bool) {

	// make sure the usr info is in the local and group caches
	au := gmcom.ActUsrD{
		Forward: fwd,
		ID:      id,
		Active:  active,
	}

	err := gmcl.AddUpdUsrCache(au, uc.internalAddress)
	if err != nil {
		lw.ErrorWithPrefixString("UsrController cache update error message:", err)
	}
}
