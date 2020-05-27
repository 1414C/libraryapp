package appobj

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/1414C/libraryapp/controllers"
	"github.com/1414C/libraryapp/group/gmcom"
	"github.com/1414C/libraryapp/group/gmsrv"
	"github.com/1414C/libraryapp/middleware"
	"github.com/1414C/libraryapp/models"
	"github.com/1414C/lw"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

// AppObj is the one and only application object
type AppObj struct {
	cfg        Config
	dbConfig   DBConfig
	services   *models.Services
	libraryC   *controllers.LibraryController
	bookC      *controllers.BookController
	usrC       *controllers.UsrController
	usrgroupC  *controllers.UsrGroupController
	authC      *controllers.AuthController
	groupauthC *controllers.GroupAuthController
	router     *mux.Router
	// jwt support
	jwtKeyMap map[string]interface{}
}

// RunMode defines the two basic operation modes
type RunMode int

const (
	cDev RunMode = iota // development settings (.dev.config.json)
	cPrd                // production settings (.prd.config.json)
	cDef                // default settings - uses compiled config in app.DefaultConfig()
)

// Initialize the application
// dev indicates run with the dev profile
// prd indicates run with the prod profile
// dr indicates run a destructive reset and exit
// rs indicates a requested rebuild of the Auth allocations to the Super UsrGroup
func (a *AppObj) Initialize(dev, prd, dr, rs bool) {

	// retrieve the app config based on production/test setting
	if prd && dev {
		fmt.Println("please specify only -dev or -prod, but not both.  Exiting.")
		os.Exit(-1)
	}

	// fallback to compiled config if neither -dev or -prod were specified
	if !prd && !dev {
		a.cfg = LoadConfig(cDef)
	}

	// try to load production configuration file (.prd.config.json)
	if prd {
		a.cfg = LoadConfig(cPrd)
	}

	// try to load development configuration file (.dev.config.json)
	if dev {
		a.cfg = LoadConfig(cDev)
	}

	// get the DB config
	a.dbConfig = a.cfg.Database

	// initialize application logging
	a.initializeLogging(a.cfg.Logging)

	// create Services
	a.createServices(a.dbConfig.ORMDebugTraceActive, a.dbConfig.ORMLogActive)

	// perform destructive reset of db table content?
	if dr {
		a.destructiveReset()
	}

	// perform automigration of positive db table changes
	a.automigrate()

	// initialize JWT keys for user-authentiction
	a.initializeJWTKeys()

	// create Controllers
	a.createControllers()

	// initialize active usr cache/buffer
	a.initializeCachedActiveUsrs()
	// a.actUsrs.ActiveUsrs["admin"] = false

	// initialize the auths cache/buffer
	a.initializeCachedAuths()

	a.initializeCachedUsrGroups()

	// intitialize the group auths cache/buffer
	a.initializeCachedGroupAuths()

	// initialize Routes
	a.initializeRoutes()

	// validate Auths against registered end-points
	auths := a.initializeAuthsByRoute()
	if auths == nil || len(auths) == 0 {
		lw.Console("Auth initialization did not return the current list of route authorizations.  The Super UsrGroup will not be updated...")
	} else {
		// assign all current Auths to the Super UsrGroup / rebuild Super UsrGroup Auths?
		a.initializeSuperGroup(auths, rs)
		if rs {
			a.initializeCachedGroupAuths()
			a.initializeRoutes()
		}
	}

	// if Usr admin does not exist, create the user and assign the super UsrGroup
	a.initializeAdminUsr()
}

// destructiveReset executes a destructive reset to refresh db
func (a *AppObj) destructiveReset() {
	a.services.DestructiveReset()
	lw.Console("Destructive reset has been carried out.  Exiting...")
	os.Exit(0)
}

// automigrate the db tables to ensure positive changes in
// the model structure have been applied.
func (a *AppObj) automigrate() {
	if err := a.services.AlterAllTables(); err != nil {
		panic(err)
	}
}

// initializeLogging sets up the logger to stdout.  replace nil with your
// own io.Writer if you wish to direct the log output to another location.
func (a *AppObj) initializeLogging(l LogConfig) {
	ls := lw.LogWriterState{
		Enabled:        l.Active,
		LocEnabled:     l.CallLocation,
		TraceEnabled:   l.TraceMsgs,
		InfoEnabled:    l.InfoMsgs,
		WarningEnabled: l.WarningMsgs,
		DebugEnabled:   l.DebugMsgs,
		ErrorEnabled:   l.ErrorMsgs,
		ColorEnabled:   l.ColorMsgTypes,
	}
	lw.InitWithSettings(ls, nil)
}

// CreateLeadSetGet reads the  keys related to the group-leader
// KVS configuration and returns an appropriate implementation of
// the gmcom.GMLeaderSetterGetter{} interface.
func (a *AppObj) CreateLeadSetGet() gmcom.GMLeaderSetterGetter {

	// check the number of active LeadSetGet configs - only one is permitted
	c := 0
	if a.cfg.LeadSetGet.StandAlone.Active {
		c++
	}
	if a.cfg.LeadSetGet.Redis.Active {
		c++
	}
	if a.cfg.LeadSetGet.Memcached.Active {
		c++
	}
	if a.cfg.LeadSetGet.Sluggo.Active {
		c++
	}

	if c == 0 || c > 1 {
		lw.Fatal(errors.New("only one group_leader_kvs subsystem may be set as active in the application server configuration file"))
	}

	// create the specified gmcom.GMLeaderSetterGetter{} interface implemenation
	if a.cfg.LeadSetGet.StandAlone.Active == true {
		return &StandAloneLeadSetGet{
			LocalLeaderIPAddress: a.cfg.LeadSetGet.StandAlone.InternalAddress,
		}
	}
	if a.cfg.LeadSetGet.Redis.Active {
		k := &RedisLeadSetGet{}
		err := k.InitializeRedisLeadSetGet(a.cfg.LeadSetGet.Redis)
		if err != nil {
			panic(err)
		}
		return k
	}
	if a.cfg.LeadSetGet.Memcached.Active {
		k := &MemcachedLeadSetGet{}
		err := k.InitializeMemcachedLeadSetGet(a.cfg.LeadSetGet.Memcached)
		if err != nil {
			panic(err)
		}
		return k
	}
	if a.cfg.LeadSetGet.Sluggo.Active {
		return &SluggoLeadSetGet{
			internalAddress: a.cfg.LeadSetGet.Sluggo.SluggoAddress,
		}
	}
	return nil
}

// createServices creates new services for the application object
func (a *AppObj) createServices(dbDebugLog, dbLog bool) {
	var err error
	a.services, err = models.NewServices(
		models.WithSqac(a.dbConfig.Dialect(), a.dbConfig.ConnectionInfo(), dbLog),
		models.WithLogMode(dbDebugLog),
		models.WithUsr(a.cfg.Pepper),
		models.WithUsrGroup(),
		models.WithAuth(),
		models.WithGroupAuth(),
		models.WithLibrary(),
		models.WithBook(),
		// models.With<Entity>,
	)

	if err != nil {
		panic(err)
	}
}

// initialize the jwt keys
func (a *AppObj) initializeJWTKeys() {

	// add the jwtKeyMap
	a.jwtKeyMap = make(map[string]interface{})

	lw.Console("ECDSA256PrivKeyFile: %v", a.cfg.ECDSA256PrivKeyFile)
	lw.Console("ECDSA256PubKeyFile: %v", a.cfg.ECDSA256PubKeyFile)
	if a.cfg.ECDSA256PrivKeyFile != "" {
		signBytes, err := ioutil.ReadFile(a.cfg.ECDSA256PrivKeyFile)
		fatal(err)

		signKey, err := jwt.ParseECPrivateKeyFromPEM(signBytes)
		fatal(err)

		a.jwtKeyMap["ES256SignKey"] = signKey
	}

	if a.cfg.ECDSA256PubKeyFile != "" {
		verifyBytes, err := ioutil.ReadFile(a.cfg.ECDSA256PubKeyFile)
		fatal(err)

		verifyKey, err := jwt.ParseECPublicKeyFromPEM(verifyBytes)
		fatal(err)

		a.jwtKeyMap["ES256VerifyKey"] = verifyKey
	}

	lw.Console("ECDSA384PrivKeyFile: %v", a.cfg.ECDSA384PrivKeyFile)
	lw.Console("ECDSA384PubKeyFile: %v", a.cfg.ECDSA384PubKeyFile)
	if a.cfg.ECDSA384PrivKeyFile != "" {
		signBytes, err := ioutil.ReadFile(a.cfg.ECDSA384PrivKeyFile)
		fatal(err)

		signKey, err := jwt.ParseECPrivateKeyFromPEM(signBytes)
		fatal(err)

		a.jwtKeyMap["ES384SignKey"] = signKey
	}

	if a.cfg.ECDSA384PubKeyFile != "" {
		verifyBytes, err := ioutil.ReadFile(a.cfg.ECDSA384PubKeyFile)
		fatal(err)

		verifyKey, err := jwt.ParseECPublicKeyFromPEM(verifyBytes)
		fatal(err)

		a.jwtKeyMap["ES384VerifyKey"] = verifyKey
	}

	lw.Console("ECDSA521PrivKeyFile: %v", a.cfg.ECDSA521PrivKeyFile)
	lw.Console("ECDSA521PubKeyFile: %v", a.cfg.ECDSA521PubKeyFile)
	if a.cfg.ECDSA521PrivKeyFile != "" {
		signBytes, err := ioutil.ReadFile(a.cfg.ECDSA521PrivKeyFile)
		fatal(err)

		signKey, err := jwt.ParseECPrivateKeyFromPEM(signBytes)
		fatal(err)

		a.jwtKeyMap["ES521SignKey"] = signKey
	}

	if a.cfg.ECDSA521PubKeyFile != "" {
		verifyBytes, err := ioutil.ReadFile(a.cfg.ECDSA521PubKeyFile)
		fatal(err)

		verifyKey, err := jwt.ParseECPublicKeyFromPEM(verifyBytes)
		fatal(err)

		a.jwtKeyMap["ES521VerifyKey"] = verifyKey
	}

	lw.Console("RSA256PrivKeyFile: %v", a.cfg.RSA256PrivKeyFile)
	lw.Console("RSA256PubKeyFile: %v", a.cfg.RSA256PubKeyFile)
	if a.cfg.RSA256PrivKeyFile != "" {
		signBytes, err := ioutil.ReadFile(a.cfg.RSA256PrivKeyFile)
		fatal(err)

		signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
		fatal(err)

		a.jwtKeyMap["RS256SignKey"] = signKey
	}

	if a.cfg.RSA256PubKeyFile != "" {
		verifyBytes, err := ioutil.ReadFile(a.cfg.RSA256PubKeyFile)
		fatal(err)

		verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
		fatal(err)

		a.jwtKeyMap["RS256VerifyKey"] = verifyKey
	}

	lw.Console("RSA384PrivKeyFile: %v", a.cfg.RSA384PrivKeyFile)
	lw.Console("RSA384PubKeyFile: %v", a.cfg.RSA384PubKeyFile)
	if a.cfg.RSA384PrivKeyFile != "" {
		signBytes, err := ioutil.ReadFile(a.cfg.RSA384PrivKeyFile)
		fatal(err)

		signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
		fatal(err)

		a.jwtKeyMap["RS384SignKey"] = signKey
	}

	if a.cfg.RSA384PubKeyFile != "" {
		verifyBytes, err := ioutil.ReadFile(a.cfg.RSA384PubKeyFile)
		fatal(err)

		verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
		fatal(err)

		a.jwtKeyMap["RS384VerifyKey"] = verifyKey
	}

	lw.Console("RSA512PrivKeyFile: %v", a.cfg.RSA512PrivKeyFile)
	lw.Console("RSA512PubKeyFile: %v", a.cfg.RSA512PubKeyFile)
	if a.cfg.RSA512PrivKeyFile != "" {
		signBytes, err := ioutil.ReadFile(a.cfg.RSA512PrivKeyFile)
		fatal(err)

		signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
		fatal(err)

		a.jwtKeyMap["RS512SignKey"] = signKey
	}

	if a.cfg.RSA512PubKeyFile != "" {
		verifyBytes, err := ioutil.ReadFile(a.cfg.RSA512PubKeyFile)
		fatal(err)

		verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
		fatal(err)

		a.jwtKeyMap["RS512VerifyKey"] = verifyKey
	}

}

// createControllers for each entity
func (a *AppObj) createControllers() {
	a.usrC = controllers.NewUsrController(a.services.Usr, a.jwtKeyMap, a.cfg.JWTSignMethod, a.cfg.JWTLifetime, a.cfg.InternalAddress)
	a.usrgroupC = controllers.NewUsrGroupController(a.services.UsrGroup, a.cfg.InternalAddress)
	a.authC = controllers.NewAuthController(a.services.Auth, a.cfg.InternalAddress)
	a.groupauthC = controllers.NewGroupAuthController(a.services.GroupAuth, a.cfg.InternalAddress)
	a.libraryC = controllers.NewLibraryController(a.services.Library, *a.services)
	a.bookC = controllers.NewBookController(a.services.Book, *a.services)
}

// initialize the list of cached active usrs
func (a *AppObj) initializeCachedActiveUsrs() {

	if a.usrC.ActUsrsH != nil {
		a.usrC.ActUsrsH.Lock()
		defer a.usrC.ActUsrsH.Unlock()
	}
	a.usrC.ActUsrsH = &gmcom.ActUsrsH{}
	a.usrC.ActUsrsH.ActiveUsrs = make(map[uint64]bool)
	u := a.services.Usr.GetUsrs()
	for _, v := range u {
		if v.Active == true {
			a.usrC.ActUsrsH.ActiveUsrs[v.ID] = true
		}
	}
}

// initialize the list of cached auths
func (a *AppObj) initializeCachedAuths() {

	// use mutex as a precautionary measure in case the method is called in a running process
	if a.authC.AuthsH != nil {
		a.authC.AuthsH.Lock()
		defer a.authC.AuthsH.Unlock()
	}
	a.authC.AuthsH = &gmcom.AuthsH{}
	a.authC.AuthsH.Auths = make(map[uint64]string)
	at := a.services.Auth.GetAuths()
	for _, v := range at {
		a.authC.AuthsH.Auths[v.ID] = v.AuthName
	}
}

// initialize the list of cached usrgroups
func (a *AppObj) initializeCachedUsrGroups() {

	// use mutex as a precautionary measure in case the method is called in a running process
	if a.usrgroupC.UsrGroupsH != nil {
		a.usrgroupC.UsrGroupsH.Lock()
		defer a.usrgroupC.UsrGroupsH.Unlock()
	}
	a.usrgroupC.UsrGroupsH = &gmcom.UsrGroupsH{}
	a.usrgroupC.UsrGroupsH.GroupNames = make(map[uint64]string)
	ug := a.services.UsrGroup.GetUsrGroups()
	for _, v := range ug {
		a.usrgroupC.UsrGroupsH.GroupNames[v.ID] = v.GroupName
	}
}

// initialize the list of cached groupauths
func (a *AppObj) initializeCachedGroupAuths() {

	// use mutex as a precautionary measure in case the method is called in a running process
	if a.groupauthC.GroupAuthsH != nil {
		a.groupauthC.GroupAuthsH.Lock()
		defer a.groupauthC.GroupAuthsH.Unlock()
	}
	a.groupauthC.GroupAuthsH = &gmcom.GroupAuthsH{}
	a.groupauthC.GroupAuthsH.GroupAuths = make(map[string]map[string]bool)
	a.groupauthC.GroupAuthsH.GroupAuthsID = make(map[uint64]gmcom.GroupAuthNames) // deletion support
	g := a.services.GroupAuth.GetGroupAuths()
	for _, v := range g {
		// create a new authMap for the group, then add the group and the auth
		// to mapGroupAuths
		mapAuth := a.groupauthC.GroupAuthsH.GroupAuths[v.GroupName]
		if mapAuth == nil {
			mapAuth = make(map[string]bool)
			mapAuth[v.AuthName] = true
			a.groupauthC.GroupAuthsH.GroupAuths[v.GroupName] = mapAuth
			continue
		}

		// if the groupName does exist in the top-level map, add the auth to
		// the group's auth map
		mapAuth[v.AuthName] = true

		// add the groupauth to the ID-based cache map
		a.groupauthC.GroupAuthsH.GroupAuthsID[v.ID] = gmcom.GroupAuthNames{GroupName: v.GroupName, AuthName: v.AuthName}
	}
}

// initializeRoutes creates routes and applies middleware
func (a *AppObj) initializeRoutes() {

	// create the RequireUsr middleware to ensure page access is secure.
	requireUserMw := middleware.InitMW(a.services.Usr, a.jwtKeyMap, a.groupauthC.GroupAuthsH, a.usrC.ActUsrsH, a.authC.AuthsH, a.usrgroupC.UsrGroupsH)

	// construct a map of local service activations
	svcActv := make(map[string]bool)
	if a.cfg.ServiceActivations != nil && len(a.cfg.ServiceActivations) > 0 {
		for i := range a.cfg.ServiceActivations {
			svcActv[a.cfg.ServiceActivations[i].ServiceName] = a.cfg.ServiceActivations[i].ServiceActive
		}
	}

	// get a gorilla router and create controllers
	a.router = mux.NewRouter()
	if a.router == nil {
		panic("appobj: failed to initialize mux")
	}

	// add usr routes
	a.router.HandleFunc("/usrs", requireUserMw.ApplyFn(a.usrC.GetUsrs)).Methods("GET").Name("usr.GET_SET")
	a.router.HandleFunc("/usr", a.usrC.Create).Methods("POST").Name("usr.CREATE")
	a.router.HandleFunc("/usr/{id:[0-9]+}", requireUserMw.ApplyFn(a.usrC.Get)).Methods("GET").Name("usr.GET_ID")
	a.router.HandleFunc("/usr/login", a.usrC.Login).Methods("POST").Name("usr.LOGIN")
	a.router.HandleFunc("/usr/{id:[0-9]+}", a.usrC.Delete).Methods("DELETE").Name("usr.DELETE")
	a.router.HandleFunc("/usr/{id:[0-9]+}", a.usrC.Update).Methods("PUT").Name("usr.UPDATE")

	// usrgroup CRUD routes
	a.router.HandleFunc("/usrgroups", requireUserMw.ApplyFn(a.usrgroupC.GetUsrGroups)).Methods("GET").Name("usrgroup.GET_SET")
	a.router.HandleFunc("/usrgroup", requireUserMw.ApplyFn(a.usrgroupC.Create)).Methods("POST").Name("usrgroup.CREATE")
	a.router.HandleFunc("/usrgroup/{id:[0-9]+}", requireUserMw.ApplyFn(a.usrgroupC.Get)).Methods("GET").Name("usrgroup.GET_ID")
	a.router.HandleFunc("/usrgroup/{id:[0-9]+}", requireUserMw.ApplyFn(a.usrgroupC.Update)).Methods("PUT").Name("usrgroup.CREATE")
	a.router.HandleFunc("/usrgroup/{id:[0-9]+}", requireUserMw.ApplyFn(a.usrgroupC.Delete)).Methods("DELETE").Name("usrgroup.DELETE")

	// usrgroup static filter routes
	// http://127.0.0.1:<port>/usrgroups/group_name(EQ '<sel_string>')
	a.router.HandleFunc("/usrgroups/group_name{group_name:[(]+(?:EQ|eq|LIKE|like)+[ ']+[a-zA-Z0-9_]+[')]+}",
		requireUserMw.ApplyFn(a.usrgroupC.GetUsrGroupsByGroupName)).Methods("GET").Name("usrgroup.STATICFLTR_ByGroupName")

	// http://127.0.0.1:<port>/usrgroups/description(EQ '<sel_string>')
	a.router.HandleFunc("/usrgroups/description{description:[(]+(?:EQ|eq|LIKE|like)+[ ']+[a-zA-Z0-9_]+[')]+}",
		requireUserMw.ApplyFn(a.usrgroupC.GetUsrGroupsByDescription)).Methods("GET").Name("usrgroup.STATICFLTR_ByDescription")

	// auth CRUD routes
	// a.router.HandleFunc("/auths", requireUserMw.ApplyFn(a.authC.GetAuths)).Methods("GET").Queries("desc", "{desc}")
	a.router.HandleFunc("/auths", requireUserMw.ApplyFn(a.authC.GetAuths)).Methods("GET").Name("auth.GET_SET")
	a.router.HandleFunc("/auth", requireUserMw.ApplyFn(a.authC.Create)).Methods("POST").Name("auth.CREATE")
	a.router.HandleFunc("/auth/{id:[0-9]+}", requireUserMw.ApplyFn(a.authC.Get)).Methods("GET").Name("auth.GET_ID")
	a.router.HandleFunc("/auth/{id:[0-9]+}", requireUserMw.ApplyFn(a.authC.Update)).Methods("PUT").Name("auth.UPDATE")
	a.router.HandleFunc("/auth/{id:[0-9]+}", requireUserMw.ApplyFn(a.authC.Delete)).Methods("DELETE").Name("auth.DELETE")

	// auth static filter routes
	// http://127.0.0.1:<port>/auths/auth_name(EQ '<sel_string>')
	a.router.HandleFunc("/auths/auth_name{auth_name:[(]+(?:EQ|eq|LIKE|like)+[ ']+[a-zA-Z0-9_]+[')]+}",
		requireUserMw.ApplyFn(a.authC.GetAuthsByAuthName)).Methods("GET").Name("auth.STATICFLTR_ByAuthName")

	// http://127.0.0.1:<port>/auths/description(EQ '<sel_string>')
	a.router.HandleFunc("/auths/description{description:[(]+(?:EQ|eq|LIKE|like)+[ ']+[a-zA-Z0-9_]+[')]+}",
		requireUserMw.ApplyFn(a.authC.GetAuthsByDescription)).Methods("GET").Name("auth.STATICFLTR_ByDescription")

	// groupauth CRUD routes
	a.router.HandleFunc("/groupauths", requireUserMw.ApplyFn(a.groupauthC.GetGroupAuths)).Methods("GET").Name("groupauth.GET_SET")
	a.router.HandleFunc("/groupauth", requireUserMw.ApplyFn(a.groupauthC.Create)).Methods("POST").Name("groupauth.CREATE")
	a.router.HandleFunc("/groupauth/{id:[0-9]+}", requireUserMw.ApplyFn(a.groupauthC.Get)).Methods("GET").Name("groupauth.GET_ID")
	a.router.HandleFunc("/groupauth/{id:[0-9]+}", requireUserMw.ApplyFn(a.groupauthC.Update)).Methods("PUT").Name("groupauth.UPDATE")
	a.router.HandleFunc("/groupauth/{id:[0-9]+}", requireUserMw.ApplyFn(a.groupauthC.Delete)).Methods("DELETE").Name("groupauth.DELETE")

	// http://127.0.0.1:<port>/groupauths/auth_name(EQ '<sel_string>')
	a.router.HandleFunc("/groupauths/auth_name{auth_name:[(]+(?:EQ|eq|LIKE|like)+[ ']+[a-zA-Z0-9_]+[')]+}",
		requireUserMw.ApplyFn(a.groupauthC.GetGroupAuthsByAuthName)).Methods("GET").Name("groupauth.STATICFLTR_ByAuthName")

	// http://127.0.0.1:<port>/groupauths/description(EQ '<sel_string>')
	a.router.HandleFunc("/groupauths/description{description:[(]+(?:EQ|eq|LIKE|like)+[ ']+[a-zA-Z0-9_]+[')]+}",
		requireUserMw.ApplyFn(a.groupauthC.GetGroupAuthsByDescription)).Methods("GET").Name("groupauth.STATICFLTR_ByDescription")

	// http://127.0.0.1:<port>/groupauths/group_id(EQ '<sel_string>')
	a.router.HandleFunc("/groupauths/group_id{group_id:[(]+(?:EQ|eq|LT|lt|GT|gt)+[ ']+[0-9]+[')]+}",
		requireUserMw.ApplyFn(a.groupauthC.GetGroupAuthsByGroupID)).Methods("GET").Name("groupauth.STATICFLTR_ByGroupID")

	// http://127.0.0.1:<port>/groupauths/group_id:
	a.router.HandleFunc("/groupauths/{group_id:[0-9]+}", requireUserMw.ApplyFn(a.groupauthC.DeleteGroupAuthsByGroupID)).Methods("DELETE").Name("groupauth.DELETE_ByGroupID")

	var pActive, ok bool

	// ====================== Library protected routes for standard CRUD access ======================
	pActive, ok = svcActv["Library"]
	if ok && pActive {
		a.router.HandleFunc("/librarys", requireUserMw.ApplyFn(a.libraryC.GetLibrarys)).Methods("GET").Name("library.GET_SET")
		a.router.HandleFunc("/librarys/{cmd:[$]+[a-zA-Z0-9_$=]+}", requireUserMw.ApplyFn(a.libraryC.GetLibrarys)).Methods("GET").Name("library.GET_SET_CMD")
		a.router.HandleFunc("/library", requireUserMw.ApplyFn(a.libraryC.Create)).Methods("POST").Name("library.CREATE")
		a.router.HandleFunc("/library/{id:[0-9]+}", requireUserMw.ApplyFn(a.libraryC.Get)).Methods("GET").Name("library.GET_ID")
		a.router.HandleFunc("/library/{id:[0-9]+}", requireUserMw.ApplyFn(a.libraryC.Update)).Methods("PUT").Name("library.UPDATE")
		a.router.HandleFunc("/library/{id:[0-9]+}", requireUserMw.ApplyFn(a.libraryC.Delete)).Methods("DELETE").Name("library.DELETE")

		//====================================== Library Relations ======================================
		// hasMany relation ToBooks for Library
		a.router.HandleFunc("/library/{library_id:[0-9]+}/tobooks", requireUserMw.ApplyFn(a.libraryC.GetLibraryToBooks)).Methods("GET").Name("library.REL_tobooks")
		a.router.HandleFunc("/library/{library_id:[0-9]+}/tobooks/{cmd:[$]+[a-zA-Z0-9_$=]+}", requireUserMw.ApplyFn(a.libraryC.GetLibraryToBooks)).Methods("GET").Name("library.REL_CMD_tobooks")

		a.router.HandleFunc("/library/{library_id:[0-9]+}/tobooks/{book_id:[0-9]+}", requireUserMw.ApplyFn(a.libraryC.GetLibraryToBooks)).Methods("GET").Name("library.REL_tobooks_id")

		//=================================== Library Static Filters ===================================
		// http://127.0.0.1:<port>/librarys/name(EQ '<sel_string>')
		a.router.HandleFunc("/librarys/name{name:[(]+(?:EQ|eq|LIKE|like)+[ ']+[a-zA-Z0-9_]+[')]+}",
			requireUserMw.ApplyFn(a.libraryC.GetLibrarysByName)).Methods("GET").Name("library.STATICFLTR_ByName")

		// http://127.0.0.1:<port>/librarys/name(EQ '<sel_string>')/$count |$limit=n $offset=n ($desc|$asc) $orderby=<col_name>
		a.router.HandleFunc("/librarys/name{name:[(]+(?:EQ|eq|LIKE|like)+[ ']+[a-zA-Z0-9_]+[')]+}/{cmd:[$]+[a-zA-Z0-9_$=]+}",
			requireUserMw.ApplyFn(a.libraryC.GetLibrarysByName)).Methods("GET").Name("library.STATICFLTR_CMD_ByName")

		// http://127.0.0.1:<port>/librarys/city(EQ '<sel_string>')
		a.router.HandleFunc("/librarys/city{city:[(]+(?:EQ|eq|LT|lt|GT|gt|LIKE|like)+[ ']+[a-zA-Z0-9_]+[')]+}",
			requireUserMw.ApplyFn(a.libraryC.GetLibrarysByCity)).Methods("GET").Name("library.STATICFLTR_ByCity")

		// http://127.0.0.1:<port>/librarys/city(EQ '<sel_string>')/$count |$limit=n $offset=n ($desc|$asc) $orderby=<col_name>
		a.router.HandleFunc("/librarys/city{city:[(]+(?:EQ|eq|LT|lt|GT|gt|LIKE|like)+[ ']+[a-zA-Z0-9_]+[')]+}/{cmd:[$]+[a-zA-Z0-9_$=]+}",
			requireUserMw.ApplyFn(a.libraryC.GetLibrarysByCity)).Methods("GET").Name("library.STATICFLTR_CMD_ByCity")

	}
	// ====================== Book protected routes for standard CRUD access ======================
	pActive, ok = svcActv["Book"]
	if ok && pActive {
		a.router.HandleFunc("/books", requireUserMw.ApplyFn(a.bookC.GetBooks)).Methods("GET").Name("book.GET_SET")
		a.router.HandleFunc("/books/{cmd:[$]+[a-zA-Z0-9_$=]+}", requireUserMw.ApplyFn(a.bookC.GetBooks)).Methods("GET").Name("book.GET_SET_CMD")
		a.router.HandleFunc("/book", requireUserMw.ApplyFn(a.bookC.Create)).Methods("POST").Name("book.CREATE")
		a.router.HandleFunc("/book/{id:[0-9]+}", requireUserMw.ApplyFn(a.bookC.Get)).Methods("GET").Name("book.GET_ID")
		a.router.HandleFunc("/book/{id:[0-9]+}", requireUserMw.ApplyFn(a.bookC.Update)).Methods("PUT").Name("book.UPDATE")
		a.router.HandleFunc("/book/{id:[0-9]+}", requireUserMw.ApplyFn(a.bookC.Delete)).Methods("DELETE").Name("book.DELETE")

		//====================================== Book Relations ======================================
		// belongsTo relation ToLibrary for Book
		a.router.HandleFunc("/book/{book_id:[0-9]+}/tolibrary", requireUserMw.ApplyFn(a.bookC.GetBookToLibrary)).Methods("GET").Name("book.REL_tolibrary")
		a.router.HandleFunc("/book/{book_id:[0-9]+}/tolibrary/{library_id:[0-9]+}", requireUserMw.ApplyFn(a.bookC.GetBookToLibrary)).Methods("GET").Name("book.REL_tolibrary_id")

		//=================================== Book Static Filters ===================================
		// http://127.0.0.1:<port>/books/title(EQ '<sel_string>')
		a.router.HandleFunc("/books/title{title:[(]+(?:EQ|eq|LIKE|like)+[ ']+[a-zA-Z0-9_]+[')]+}",
			requireUserMw.ApplyFn(a.bookC.GetBooksByTitle)).Methods("GET").Name("book.STATICFLTR_ByTitle")

		// http://127.0.0.1:<port>/books/title(EQ '<sel_string>')/$count |$limit=n $offset=n ($desc|$asc) $orderby=<col_name>
		a.router.HandleFunc("/books/title{title:[(]+(?:EQ|eq|LIKE|like)+[ ']+[a-zA-Z0-9_]+[')]+}/{cmd:[$]+[a-zA-Z0-9_$=]+}",
			requireUserMw.ApplyFn(a.bookC.GetBooksByTitle)).Methods("GET").Name("book.STATICFLTR_CMD_ByTitle")

		// http://127.0.0.1:<port>/books/author(EQ '<sel_string>')
		a.router.HandleFunc("/books/author{author:[(]+(?:EQ|eq|LIKE|like)+[ ']+[a-zA-Z0-9_]+[')]+}",
			requireUserMw.ApplyFn(a.bookC.GetBooksByAuthor)).Methods("GET").Name("book.STATICFLTR_ByAuthor")

		// http://127.0.0.1:<port>/books/author(EQ '<sel_string>')/$count |$limit=n $offset=n ($desc|$asc) $orderby=<col_name>
		a.router.HandleFunc("/books/author{author:[(]+(?:EQ|eq|LIKE|like)+[ ']+[a-zA-Z0-9_]+[')]+}/{cmd:[$]+[a-zA-Z0-9_$=]+}",
			requireUserMw.ApplyFn(a.bookC.GetBooksByAuthor)).Methods("GET").Name("book.STATICFLTR_CMD_ByAuthor")

		// http://127.0.0.1:<port>/books/hardcover(EQ TRUE)
		a.router.HandleFunc("/books/hardcover{hardcover:[(]+(?:EQ|eq|NE|ne)+[ ']+(?:true|TRUE|false|FALSE)+[')]+}",
			requireUserMw.ApplyFn(a.bookC.GetBooksByHardcover)).Methods("GET").Name("book.STATICFLTR_ByHardcover")

		// http://127.0.0.1:<port>/books/hardcover(EQ TRUE)/$count |$limit=n $offset=n ($desc|$asc) $orderby=<col_name>
		a.router.HandleFunc("/books/hardcover{hardcover:[(]+(?:EQ|eq|NE|ne)+[ ']+(?:true|TRUE|false|FALSE)+[')]+}/{cmd:[$]+[a-zA-Z0-9_$=]+}",
			requireUserMw.ApplyFn(a.bookC.GetBooksByHardcover)).Methods("GET").Name("book.STATICFLTR_CMD_ByHardcover")

		// http://127.0.0.1:<port>/books/copies(EQ 72.43)
		// http://127.0.0.1:<port>/books/copies(LT 110)
		// http://127.0.0.1:<port>/books/copies(GE -43)
		a.router.HandleFunc("/books/copies{copies:[(]+(?:EQ|eq|LT|lt|GT|gt)+[ ]+[0-9_]+[)]+}",
			requireUserMw.ApplyFn(a.bookC.GetBooksByCopies)).Methods("GET").Name("book.STATICFLTR_ByCopies")

		// http://127.0.0.1:<port>/books/copies(EQ 72.43)/$count |$limit=n $offset=n ($desc|$asc) $orderby=<col_name>
		// http://127.0.0.1:<port>/books/copies(LT 110)/$count |$limit=n $offset=n ($desc|$asc) $orderby=<col_name>
		// http://127.0.0.1:<port>/books/copies(GE -43)/$count |$limit=n $offset=n ($desc|$asc) $orderby=<col_name>
		a.router.HandleFunc("/books/copies{copies:[(]+(?:EQ|eq|LT|lt|GT|gt)+[ ]+[0-9_]+[)]+}/{cmd:[$]+[a-zA-Z0-9_$=]+}",
			requireUserMw.ApplyFn(a.bookC.GetBooksByCopies)).Methods("GET").Name("book.STATICFLTR_CMD_ByCopies")

		// http://127.0.0.1:<port>/books/library_id(EQ 72.43)
		// http://127.0.0.1:<port>/books/library_id(LT 110)
		// http://127.0.0.1:<port>/books/library_id(GE -43)
		a.router.HandleFunc("/books/library_id{library_id:[(]+(?:EQ|eq)+[ ]+[0-9_]+[)]+}",
			requireUserMw.ApplyFn(a.bookC.GetBooksByLibraryID)).Methods("GET").Name("book.STATICFLTR_ByLibraryID")

		// http://127.0.0.1:<port>/books/library_id(EQ 72.43)/$count |$limit=n $offset=n ($desc|$asc) $orderby=<col_name>
		// http://127.0.0.1:<port>/books/library_id(LT 110)/$count |$limit=n $offset=n ($desc|$asc) $orderby=<col_name>
		// http://127.0.0.1:<port>/books/library_id(GE -43)/$count |$limit=n $offset=n ($desc|$asc) $orderby=<col_name>
		a.router.HandleFunc("/books/library_id{library_id:[(]+(?:EQ|eq)+[ ]+[0-9_]+[)]+}/{cmd:[$]+[a-zA-Z0-9_$=]+}",
			requireUserMw.ApplyFn(a.bookC.GetBooksByLibraryID)).Methods("GET").Name("book.STATICFLTR_CMD_ByLibraryID")

	}
}

// getRouteNames walks the routes to get the route names for usr/group/auth lookup
func (a *AppObj) getRouteNames() (authNames []string) {

	walkFn := func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		name := route.GetName()
		if name == "" {
			path, _ := route.GetPathTemplate()
			return fmt.Errorf("walkRoutes found route %v with no name", path)
		}
		authNames = append(authNames, name)
		return nil
	}

	err := a.router.Walk(walkFn)
	if err != nil {
		lw.Console("warning: %v", err)
	}
	return authNames
}

// initializeAuthsByRoute loads the existing Auths from the db, then
// examines the current set of routes (route.Names).  Missing
// Auths are then created.  Auths that are no longer required
// could be reported to stdout at this point, but for now they
// are ignored.
func (a *AppObj) initializeAuthsByRoute() map[string]uint64 {

	var pendingAuths []string

	// load auths from db (auth objects only)
	auths := a.services.Auth.GetAuths()
	mapAuths := make(map[string]uint64)
	for _, a := range auths {
		mapAuths[a.AuthName] = a.ID
	}

	// load the registered route names
	routeNames := a.getRouteNames()
	for _, rn := range routeNames {
		id := mapAuths[rn]
		if id == 0 {
			lw.Warning("auth %s not found in the db Auth master data", rn)
			pendingAuths = append(pendingAuths, rn)
		}
	}

	// create the missing auths in the Auth master data
	for _, pa := range pendingAuths {
		lw.Console("creating auth %s in the db Auth master data...", pa)

		d := ""
		parts := strings.Split(pa, ".")
		if len(parts) != 2 {
			d = pa + "inserted during initialization"
		} else {
			d = fmt.Sprintf("Allow %s access to the %s entity", parts[1], parts[0])
		}

		newAuth := models.Auth{
			AuthName:    pa,
			AuthType:    "endpoint",
			Description: d,
		}

		err := a.services.Auth.Create(&newAuth)
		if err != nil {
			lw.Warning("attempted creation of auth %s failed. got: %s", pa, err)
			return nil
		}
		lw.Warning("new auth %s must be added to at least one group", pa)
		mapAuths[newAuth.AuthName] = newAuth.ID
	}
	return mapAuths
}

// initializeSuperGroup performs the following tasks:
// 1. Checks for the existence of the Super Group.
// 2. If the Super Group is not found, it is created.
// 3. deletes the current Auth allocations to the Super Group (if any).
// 4. adds all the current Auths read in initializeAuths to the Super Group.
func (a *AppObj) initializeSuperGroup(Auths map[string]uint64, rebuildSuperGroup bool) {

	var sg models.UsrGroup

	// get UsrGroup 'Super'
	superGroup := a.services.UsrGroup.GetUsrGroupsByGroupName("EQ", "Super")
	if superGroup == nil || len(superGroup) == 0 {
		sg = models.UsrGroup{
			GroupName:   "Super",
			Description: "Super Group - use for admin only",
		}
		a.services.UsrGroup.Create(&sg)
	} else {
		if rebuildSuperGroup {
			sg = superGroup[0]
		}
	}

	if sg.ID == 0 {
		return
	}

	// delete current auth allocations to Super
	err := a.services.GroupAuth.DeleteGroupAuthsByGroupID(strconv.Itoa(int(sg.ID)))
	if err != nil {
		lw.Warning(err.Error())
		lw.Warning("Super GroupAuths will not be maintained.  Some end-points may not be available...")
		return
	}

	// Allocate all Auths to the Super UsrGroup
	for k, v := range Auths {

		// call the Create method on the groupauth model
		err := a.services.GroupAuth.CreateGroupAuthDirect(sg.ID, v)
		if err != nil {
			panic(fmt.Sprintf("CreateGroupAuthDirect error occured while rebuilding the Super Group Auth assignment for UsrGroup %v and Auth %v.\n", sg.ID, k))
		}
	}
	lw.Console("The Super UsrGroup has been initialized with %v Auth objects.", len(Auths))
	lw.Console("re-initializing local middleware to accomodate Super group changes.")
}

// initializeAdminUsr checks for the existance of Usr 'admin'.  If 'admin' does not exist, a new Usr
// record will be created ('admin') and the 'super' UsrGroup will be assigned.  The check does not
// examine whether the 'admin' Usr is Active or Inactive, so if the 'admin' Usr has been set to
// Inactive, no changes will be made.
// Note that the 'ByEmail' method makes a special exception for 'admin' in its validation rules.
func (a *AppObj) initializeAdminUsr() {

	// does Usr 'admin' exist?
	ua, _ := a.services.Usr.ByEmail("admin")
	if ua != nil {
		return
	}

	strGroups := "Super"
	usrAdmin := models.Usr{
		Name:     "Admin",
		Email:    "admin",
		Password: "initpass",
		Active:   true,
		Groups:   &strGroups,
	}

	err := a.services.Usr.Create(&usrAdmin)
	if err != nil {
		panic(fmt.Sprintf("failed to create admin user in intializeAdminUsr: %v\n", err))
	}
	lw.Console("admin user created with ID: %v and initial password of %v", usrAdmin.ID, "initpass")

	// add the admin usr to the local cache
	a.usrC.ActUsrsH.Lock()
	a.usrC.ActUsrsH.ActiveUsrs[usrAdmin.ID] = true
	a.usrC.ActUsrsH.Unlock()
}

// Run the application
func (a *AppObj) Run(lsg gmcom.GMLeaderSetterGetter) {

	// start the group-membership server
	gv := &gmsrv.GMServ{}
	go gv.Serve(a.cfg.InternalAddress, lsg, a.usrC.ActUsrsH, a.groupauthC.GroupAuthsH, a.authC.AuthsH, a.usrgroupC.UsrGroupsH, true, a.cfg.PingCycle, a.cfg.FailureThreshold)

	// close db connection later
	defer a.services.Close()

	// set basic http server values
	var err error
	srv := http.Server{
		Addr:    a.cfg.ExternalAddress,
		Handler: a.router,
	}

	// graceful shutdown channels
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	// specify signals of interest for shutdown
	signal.Notify(sigs, syscall.SIGKILL, syscall.SIGTERM, os.Interrupt)

	// create a context for the srv.Shutdown() call
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	// start a blocking goroutine to listen for signals. when an interesting signal arrives
	// push a true into the done channel.
	go func() {
		signal := <-sigs
		lw.Console("signal received: %v", signal)

		// sending departing message to group-members
		gv.SendDeparting()
		done <- true
	}()

	if a.cfg.IsProd() {
		lw.Console("Production settings selected...")

		tlsCfg := &tls.Config{
			MinVersion: tls.VersionTLS10, // good enough?
		}

		tlsCfg.Certificates = make([]tls.Certificate, 1)
		tlsCfg.Certificates[0], err = tls.LoadX509KeyPair(a.cfg.CertFile, a.cfg.KeyFile)
		if err != nil {
			lw.Fatal(err)
		}
		srv.TLSConfig = tlsCfg

		// create a tlsListener
		lsnr, err := net.Listen("tcp", a.cfg.ExternalAddress)
		if err != nil {
			lw.Fatal(err)
		}
		tlsListener := tls.NewListener(lsnr, tlsCfg)

		lw.Console("Starting https server on: %v", a.cfg.ExternalAddress)
		go func() {
			err := srv.Serve(tlsListener)
			if err != nil {
				lw.Fatal(err)
			}
		}()

	} else {
		if a.cfg.IsDev() {
			lw.Console("Development settings selected...")
		} else {
			lw.Console("Default settings selected...")
		}
		lw.Console("Starting http server on: %v", a.cfg.ExternalAddress)
		go func() {
			err := srv.ListenAndServe()
			if err != nil {
				lw.Fatal(err)
			}
		}()
	}

	// handle shutdown
	<-done
	srv.Shutdown(ctx)
	lsg.Cleanup()
	lw.Console("Shutdown complete")
}

func fatal(err error) {
	if err != nil {
		lw.Fatal(err)
	}
}
