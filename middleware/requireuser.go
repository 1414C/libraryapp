package middleware

//=============================================================================================
// start of generated code: please do not modify this section
// code generated on 27 May 20 17:57 CDT
//=============================================================================================

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"fmt"
	"net/http"
	"strings"

	"github.com/1414C/libraryapp/group/gmcom"
	"github.com/1414C/libraryapp/models"
	"github.com/1414C/lw"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gorilla/mux"
)

// RequireUsr offers a closure that can be called with the requested page's
// w,r in order to verify that the usr is logged in prior to rendering the
// requested page.
type RequireUsr struct {
	Usr               models.UsrService
	ECDSA256VerifyKey *ecdsa.PublicKey
	ECDSA384VerifyKey *ecdsa.PublicKey
	ECDSA521VerifyKey *ecdsa.PublicKey
	RSA256VerifyKey   *rsa.PublicKey
	RSA384VerifyKey   *rsa.PublicKey
	RSA512VerifyKey   *rsa.PublicKey
	// mapGA             map[string]map[string]bool // map[groupName]map[auth_name]bool
	// mapActiveUsrs     map[string]bool            // ref to a.activeUsrs!
	ActUsrsH    *gmcom.ActUsrsH
	GroupAuthsH *gmcom.GroupAuthsH
	AuthsH      *gmcom.AuthsH
	UsrGroupsH  *gmcom.UsrGroupsH
}

// InitMW is used to initialize the usr authorization middleware
func InitMW(Usr models.UsrService, jwtKeyMap map[string]interface{}, groupAuths *gmcom.GroupAuthsH, actUsrs *gmcom.ActUsrsH, auths *gmcom.AuthsH, usrGroups *gmcom.UsrGroupsH) (requireUser RequireUsr) {

	var esVerifyKey *ecdsa.PublicKey
	var rsVerifyKey *rsa.PublicKey
	var ok bool

	esVerifyKey, ok = jwtKeyMap["ES256VerifyKey"].(*ecdsa.PublicKey)
	if ok {
		requireUser.ECDSA256VerifyKey = esVerifyKey
		esVerifyKey = nil
	}

	esVerifyKey, ok = jwtKeyMap["ES384VerifyKey"].(*ecdsa.PublicKey)
	if ok {
		requireUser.ECDSA384VerifyKey = esVerifyKey
		esVerifyKey = nil
	}

	esVerifyKey, ok = jwtKeyMap["ES521VerifyKey"].(*ecdsa.PublicKey)
	if ok {
		requireUser.ECDSA521VerifyKey = esVerifyKey
		esVerifyKey = nil
	}

	rsVerifyKey, ok = jwtKeyMap["RS256VerifyKey"].(*rsa.PublicKey)
	if ok {
		requireUser.RSA256VerifyKey = rsVerifyKey
		rsVerifyKey = nil
	}

	rsVerifyKey, ok = jwtKeyMap["RS384VerifyKey"].(*rsa.PublicKey)
	if ok {
		requireUser.RSA384VerifyKey = rsVerifyKey
		rsVerifyKey = nil
	}

	rsVerifyKey, ok = jwtKeyMap["RS512VerifyKey"].(*rsa.PublicKey)
	if ok {
		requireUser.RSA512VerifyKey = rsVerifyKey
		rsVerifyKey = nil
	}

	requireUser.Usr = Usr
	requireUser.GroupAuthsH = groupAuths
	requireUser.ActUsrsH = actUsrs
	requireUser.AuthsH = auths
	requireUser.UsrGroupsH = usrGroups
	return requireUser
}

// CheckAuth verifies that the AuthName is contained in one of the supplied Groups
func (mw *RequireUsr) CheckAuth(authName string, groups []string, id uint64) bool {

	// no group auths setup - disallow all access
	if len(mw.GroupAuthsH.GroupAuths) == 0 {
		return false
	}

	// is the user active?
	active := mw.ActUsrsH.ActiveUsrs[id]
	lw.Debug("==============================================================")
	lw.Debug("mw.ActUsrsH.ActiveUsrs: %v", mw.ActUsrsH.ActiveUsrs)
	lw.Debug("==============================================================")
	lw.Debug("mw.GroupAuthsH.GroupAuths: %v", mw.GroupAuthsH.GroupAuths)
	lw.Debug("==============================================================")
	lw.Debug("mw.AuthsH.Auths: %v", mw.AuthsH.Auths)
	lw.Debug("==============================================================")
	lw.Debug("mw.UsrGroups.UsrGroupsH: %v", mw.UsrGroupsH.GroupNames)
	lw.Debug("==============================================================")
	if !active {
		return false
	}

	for _, g := range groups {

		// nil here is odd but plausible
		mapAuths := mw.GroupAuthsH.GroupAuths[g]
		// lw.Debug("mw.mapGA[%s] got: %v", g, mapAuths)
		if mapAuths == nil || len(mapAuths) == 0 {
			continue
		}

		// lw.Debug("mapAuths: %v", mapAuths)
		auth := mapAuths[authName]
		if auth == true {
			return true
		}
	}
	return false
}

// ApplyFn assumes that Usr middleware has already been run - i.e. the application
// has attempted to authenticate the user in terms of their login credentials being
// valid.
// Check to see whether the Authentication was successful.
// Check to see whether the Usr has a Group assignment which permits access to the
// requested resource.
// ApplyFn
func (mw *RequireUsr) ApplyFn(next http.HandlerFunc) http.HandlerFunc {

	// CustomClaims are used to facilitate access to application-specific
	// claims that are not part of the JWT standard set.
	type CustomClaims struct {
		*jwt.StandardClaims
		TokenType string
		Groups    string
		UID       uint64
		Email     string
	}

	// http.HandlerFunc is casting the type of the closure here
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		noKey := false

		// verify the JWT content
		token, err := request.ParseFromRequestWithClaims(r, request.AuthorizationHeaderExtractor, &CustomClaims{},
			func(token *jwt.Token) (interface{}, error) {
				switch token.Header["alg"] {
				case "ES256":
					if mw.ECDSA256VerifyKey != nil {
						return mw.ECDSA256VerifyKey, nil
					}
					noKey = true

				case "ES384":
					if mw.ECDSA384VerifyKey != nil {
						return mw.ECDSA384VerifyKey, nil
					}
					noKey = true

				case "ES521":
					if mw.ECDSA521VerifyKey != nil {
						return mw.ECDSA521VerifyKey, nil
					}
					noKey = true

				case "RS256":
					if mw.RSA256VerifyKey != nil {
						return mw.RSA256VerifyKey, nil
					}
					noKey = true

				case "RS384":
					if mw.RSA384VerifyKey != nil {
						return mw.RSA384VerifyKey, nil
					}
					noKey = true

				case "RS512":
					if mw.RSA512VerifyKey != nil {
						return mw.RSA512VerifyKey, nil
					}
					noKey = true

				case "HS256":
					return nil, fmt.Errorf("hmac signed jwt's are not accepted")

				default:
					if noKey {
						return nil, fmt.Errorf("unable to verify access token for %v signing algorithm", token.Header["alg"])
					}
					return nil, fmt.Errorf("unknown 'alg': %v in JWT header", token.Header["alg"])
				}
				return nil, fmt.Errorf("unknown error validating access token")
			})

		if err == nil {
			if token.Valid {
				claims := token.Claims.(*CustomClaims)

				// lw.Debug("checking auth: %s with groups: %v", mux.CurrentRoute(r).GetName(), claims.Groups)
				// lw.Debug("token.Header: %v", token.Header)
				// lw.Debug("claims.IssuedAt: %v", claims.IssuedAt)
				// lw.Debug("claims.ExpiresAt: %v", claims.ExpiresAt)
				// lw.Debug("claims.NotBefore: %v", claims.NotBefore)
				// lw.Debug("claims.Issuer: %v", claims.Issuer)
				// lw.Debug("claims.Subject: %v", claims.Subject)
				// lw.Debug("claims.TokenType: %v", claims.TokenType)
				// lw.Debug("claims.Groups: %v", claims.Groups)
				// lw.Debug("claims.Id: %v", claims.Id)
				// lw.Debug("claims.UID: %v", claims.UID)
				// lw.Debug("claims.Email: %v", claims.Email)

				gps := strings.Split(claims.Groups, ";")
				for i := range gps {
					gps[i] = strings.TrimSpace(gps[i])
				}

				// check the user's authorization for the route
				if mw.CheckAuth(mux.CurrentRoute(r).GetName(), gps, claims.UID) {
					next(w, r)
					return
				}
			}
			w.WriteHeader(http.StatusUnauthorized)
			lw.Warning("Unauthorized access to this resource: %v", w)
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		lw.Warning("Unauthorized access to this resource: %v %s", w, err.Error())
		return
	})
}

// Apply assumes that Usr middleware has already been run
// otherwise it will not work correctly.
func (mw *RequireUsr) Apply(next http.Handler) http.HandlerFunc {
	// lw.Info("middleware: *RequireUser: Apply: %v", next)
	return mw.ApplyFn(next.ServeHTTP)
}
