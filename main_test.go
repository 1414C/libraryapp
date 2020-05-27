package main_test

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/1414C/libraryapp/appobj"
	"github.com/1414C/libraryapp/models"
)

// SessionData contains session management vars
type SessionData struct {
	jwtToken     string
	client       *http.Client
	log          bool
	ID           uint64
	baseURL      string
	testURL      string
	testEndPoint string
	usrName      string
	usrID        uint64
}

var (
	sessionData SessionData
	certFile    = flag.String("cert", "mycert1.cer", "A PEM encoded certificate file.")
	keyFile     = flag.String("key", "mycert1.key", "A PEM encoded private key file.")
	caFile      = flag.String("CA", "myCA.cer", "A PEM encoded CA's certificate file.")
)

var a appobj.AppObj

func TestMain(m *testing.M) {

	// parse flags
	logFlag := flag.Bool("log", false, "extended log")
	useHttpsFlag := flag.Bool("https", false, "true == use https")
	addressFlag := flag.String("address", "localhost:3000", "address:port to connect to")
	u := flag.String("u", "admin", "user name")
	passwd := flag.String("passwd", "", "passwd")
	flag.Parse()

	sessionData.log = *logFlag

	// initialize client / transport
	err := sessionData.initializeClient(*useHttpsFlag)
	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}

	// build base url
	err = sessionData.buildURL(*useHttpsFlag, *addressFlag)
	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}

	// this method was implemented prior to the end-point authorization build
	// and does not presently work.  I think the test should run with a user
	// and password specified from the command line. :(
	// // create test usr
	// err = sessionData.createUsr()
	// if err != nil {
	// 	log.Fatalf("%s\n", err.Error())
	// }

	// login / get jwt
	err = sessionData.getJWT(*u, *passwd)
	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}

	code := m.Run()

	// // delete test usr
	// err = sessionData.deleteUsr()
	// if err != nil {
	//	log.Fatalf("%s\n", err.Error())
	//}

	os.Exit(code)

}

// initialize client / transport
func (sd *SessionData) initializeClient(useHttps bool) error {

	// https
	if useHttps {
		// Load client cert
		cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
		if err != nil {
			return err
		}

		// Load CA cert
		caCert, err := ioutil.ReadFile(*caFile)
		if err != nil {
			return err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		// Setup HTTPS client
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
		}
		tlsConfig.BuildNameToCertificate()
		transport := &http.Transport{TLSClientConfig: tlsConfig}
		sd.client = &http.Client{Transport: transport,
			Timeout: time.Second * 10,
		}
	}
	// http
	sd.client = &http.Client{
		Timeout: time.Second * 10,
	}
	return nil
}

// buildURL builds a url based on flag parameters
//
// internal
func (sd *SessionData) buildURL(useHttps bool, address string) error {

	sd.baseURL = "http"
	if useHttps {
		sd.baseURL = sd.baseURL + "s"
	}
	sd.baseURL = sd.baseURL + "://" + address
	return nil
}

// createUsr creates a test usr for the application
//
// POST - /usr
func (sd *SessionData) createUsr() error {

	url := sd.baseURL + "/usr"

	// create unique usr name
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return err
	}
	sessionData.usrName = fmt.Sprintf("%X%X%X%X%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	jsonStr := fmt.Sprintf("{\"email\":\"%s@1414c.io\",\"password\":\"woofwoof\"}", sessionData.usrName)

	// var jsonBody = []byte(`{"email":"testusr123@1414c.io", "password":"woofwoof"}`)
	var jsonBody = []byte(jsonStr)
	fmt.Println("creating usr:", string(jsonBody))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	if sd.log {
		fmt.Println("POST request Headers:", req.Header)
	}

	resp, err := sd.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var usr models.Usr
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&usr); err != nil {
		return err
	}

	sessionData.usrID = usr.ID

	if sd.log {
		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
	}
	return nil
}

// deleteUsr deletes the test usr
//
// DELETE - /usr/:id
func (sd *SessionData) deleteUsr() error {

	idStr := fmt.Sprint(sessionData.usrID)
	// url := "https://localhost:8080/usr/" + idStr
	fmt.Println("deleting usr:", sessionData.usrName, sessionData.usrID)
	url := sessionData.baseURL + "/usr/" + idStr
	var jsonBody = []byte(`{}`)
	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(jsonBody))
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+sessionData.jwtToken)

	if sessionData.log {
		fmt.Println("sessionData.ID:", string(sessionData.ID))
		fmt.Println("DELETE URL:", url)
		fmt.Println("DELETE request Headers:", req.Header)
	}

	resp, err := sessionData.client.Do(req)
	if err != nil {
		fmt.Printf("Test was unable to DELETE /usr/%d. Got %s.\n", sessionData.usrID, err.Error())
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		fmt.Printf("DELETE /usr{:id} expected http status code of 201 - got %d", resp.StatusCode)
		return err
	}
	return nil
}

// getJWT authenticates and get JWT
//
// POST - /usr/login
func (sd *SessionData) getJWT(u, p string) error {

	type jwtResponse struct {
		Token string `json:"token"`
	}

	// url := "https://localhost:8080/usr/login"
	url := sessionData.baseURL + "/usr/login"

	jsonStr := ""
	if u != "" {
		jsonStr = fmt.Sprintf("{\"email\":\"%s\",\"password\":\"%s\"}", u, p)
	} else {
		jsonStr = fmt.Sprintf("{\"email\":\"%s@1414c.io\",\"password\":\"woofwoof\"}", sessionData.usrName)
	}

	// var jsonStr = []byte(`{"email":"bunnybear10@1414c.io", "password":"woofwoof"}`)
	// jsonStr := fmt.Sprintf("{\"email\":\"%s@1414c.io\",\"password\":\"woofwoof\"}", sessionData.usrName)
	fmt.Println("using usr:", jsonStr)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonStr)))
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	if sd.log {
		fmt.Println("POST request Headers:", req.Header)
	}

	resp, err := sd.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var j jwtResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&j); err != nil {
		return err
	}

	sd.jwtToken = j.Token

	if sd.log {
		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
	}
	return nil
}

// testSelectableField is used to test the endpoint access to an entity field
// that has been marked as Selectable in the model file.  access will be tested
// for each of the supported operations via multiple calls to this method.
// The selection data provided in the end-point string is representitive of
// the field data-type only, and it is not expected that the string or
// number types will return a data payload in the response body.  Consequently,
// only the http status code in the response is examined.
//
// GET - sd.testURL
func (sd *SessionData) testSelectableField(t *testing.T) {

	var jsonStr = []byte(`{}`)
	req, _ := http.NewRequest("GET", sd.testURL, bytes.NewBuffer(jsonStr))
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+sessionData.jwtToken)

	if sessionData.log {
		fmt.Println("GET URL:", sd.testURL)
		fmt.Println("GET request Headers:", req.Header)
	}

	resp, err := sessionData.client.Do(req)
	if err != nil {
		t.Errorf("Test was unable to GET %s. Got %s.\n", sd.testEndPoint, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("GET %s expected http status code of 200 - got %d", sd.testEndPoint, resp.StatusCode)
	}
}

// TestCreateLibrary attempts to create a new Library on the db
//
// POST /library
func TestCreateLibrary(t *testing.T) {

	// url := "https://localhost:8080/library"
	url := sessionData.baseURL + "/library"

	var jsonStr = []byte(`{"name":"string_value",
"city":"string_value"}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+sessionData.jwtToken)

	if sessionData.log {
		fmt.Println("POST request Headers:", req.Header)
	}

	resp, err := sessionData.client.Do(req)
	if err != nil {
		t.Errorf("Test was unable to POST /library. Got %s.\n", err.Error())
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("GET /library response Body:", string(body))
		t.Errorf("Test was unable to POST /library. Got %s.\n", err.Error())
	}
	defer resp.Body.Close()

	var e models.Library
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&e); err != nil {
		t.Errorf("Test was unable to decode the result of POST /library. Got %s.\n", err.Error())
	}

	//============================================================================================
	// TODO: implement validation of the created entity here
	//============================================================================================
	if e.Name != "string_value" {
		t.Errorf("inconsistency detected in POST /library field Name.")
	}

	if e.City != "string_value" {
		t.Errorf("inconsistency detected in POST /library field City.")
	}

	if e.ID != 0 {
		sessionData.ID = e.ID
	} else {
		log.Printf("ID value of 0 detected - subsequent test cases will run with ID == 0!")
	}
}

// TestGetLibrarys attempts to read all librarys from the db
//
// GET /librarys
func TestGetLibrarys(t *testing.T) {

	// url := "https://localhost:8080/librarys"
	url := sessionData.baseURL + "/librarys"
	jsonStr := []byte(`{}`)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+sessionData.jwtToken)

	if sessionData.log {
		fmt.Println("GET /librarys request Headers:", req.Header)
	}

	// client := &http.Client{}
	resp, err := sessionData.client.Do(req)
	if err != nil {
		t.Errorf("Test was unable to GET /librarys. Got %s.\n", err.Error())
	}
	defer resp.Body.Close()

	if sessionData.log {
		fmt.Println("GET /librarys response Status:", resp.Status)
		fmt.Println("GET /librarys response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("GET /librarys response Body:", string(body))
	}
}

// TestGetLibrary attempts to read library/{:id} from the db
// using the id created in this entity's TestCreate function.
//
// GET /library/{:id}
func TestGetLibrary(t *testing.T) {

	idStr := fmt.Sprint(sessionData.ID)
	// url := "https://localhost:8080/library/" + idStr
	url := sessionData.baseURL + "/library/" + idStr
	jsonStr := []byte(`{}`)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+sessionData.jwtToken)

	if sessionData.log {
		fmt.Println("GET /library request Headers:", req.Header)
	}

	// client := &http.Client{}
	resp, err := sessionData.client.Do(req)
	if err != nil {
		t.Errorf("Test was unable to GET /library/%d. Got %s.\n", sessionData.ID, err.Error())
	}
	defer resp.Body.Close()

	if sessionData.log {
		fmt.Println("GET /library response Status:", resp.Status)
		fmt.Println("GET /library response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("GET /library response Body:", string(body))
	}
}

// TestUpdateLibrary attempts to update an existing Library on the db
//
// PUT /library/{:id}
func TestUpdateLibrary(t *testing.T) {

	idStr := fmt.Sprint(sessionData.ID)
	// url := "https://localhost:8080/library/" + idStr
	url := sessionData.baseURL + "/library/" + idStr

	var jsonStr = []byte(`{"name":"string_update",
"city":"string_update"}`)

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+sessionData.jwtToken)

	if sessionData.log {
		fmt.Println("POST request Headers:", req.Header)
	}

	resp, err := sessionData.client.Do(req)
	if err != nil {
		t.Errorf("Test was unable to PUT /library/{:id}. Got %s.\n", err.Error())
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("PUT /library{:id} expected http status code of 201 - got %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	var e models.Library
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&e); err != nil {
		t.Errorf("Test was unable to decode the result of PUT /library. Got %s.\n", err.Error())
	}

	//============================================================================================
	// TODO: implement validation of the updated entity here
	//============================================================================================
	if e.Name != "string_update" {
		t.Errorf("inconsistency detected in POST /library field Name.")
	}

	if e.City != "string_update" {
		t.Errorf("inconsistency detected in POST /library field City.")
	}

	if e.ID != 0 {
		sessionData.ID = e.ID
	} else {
		log.Printf("ID value of 0 detected - subsequent test cases will run with ID == 0!")
	}
}

// TestDeleteLibrary attempts to delete the new Library on the db
//
// DELETE /library/{:id}
func TestDeleteLibrary(t *testing.T) {

	idStr := fmt.Sprint(sessionData.ID)
	// url := "https://localhost:8080/library/" + idStr
	url := sessionData.baseURL + "/library/" + idStr
	var jsonStr = []byte(`{}`)
	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(jsonStr))
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+sessionData.jwtToken)

	if sessionData.log {
		fmt.Println("sessionData.ID:", string(sessionData.ID))
		fmt.Println("DELETE URL:", url)
		fmt.Println("DELETE request Headers:", req.Header)
	}

	resp, err := sessionData.client.Do(req)
	if err != nil {
		t.Errorf("Test was unable to DELETE /library/%d. Got %s.\n", sessionData.ID, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		t.Errorf("DELETE /library{:id} expected http status code of 201 - got %d", resp.StatusCode)
	}
}

func TestGetLibrarysByName(t *testing.T) {

	// http://127.0.0.1:<port>/librarys/name(OP '<sel_string>')
	sessionData.testEndPoint = "/librarys/name(EQ 'test_string')"
	sessionData.testURL = sessionData.baseURL + sessionData.testEndPoint
	sessionData.testSelectableField(t)

	sessionData.testEndPoint = "/librarys/name(LIKE 'test_string')"
	sessionData.testURL = sessionData.baseURL + sessionData.testEndPoint
	sessionData.testSelectableField(t)

} // end func {Name string name  false 0  false true false nonUnique EQ,LIKE     sqac:"nullable:false;index:non-unique;index:idx_library_name_city" json:"name" false false false}

func TestGetLibrarysByCity(t *testing.T) {

	// http://127.0.0.1:<port>/librarys/city(OP '<sel_string>')
	sessionData.testEndPoint = "/librarys/city(EQ 'test_string')"
	sessionData.testURL = sessionData.baseURL + sessionData.testEndPoint
	sessionData.testSelectableField(t)

	sessionData.testEndPoint = "/librarys/city(LT 'test_string')"
	sessionData.testURL = sessionData.baseURL + sessionData.testEndPoint
	sessionData.testSelectableField(t)

	sessionData.testEndPoint = "/librarys/city(GT 'test_string')"
	sessionData.testURL = sessionData.baseURL + sessionData.testEndPoint
	sessionData.testSelectableField(t)

	sessionData.testEndPoint = "/librarys/city(LIKE 'test_string')"
	sessionData.testURL = sessionData.baseURL + sessionData.testEndPoint
	sessionData.testSelectableField(t)

} // end func {City string city  false 0  false true false  EQ,LT,GT,LIKE     sqac:"nullable:false;index:idx_library_name_city" json:"city" false false false}

// TestCreateBook attempts to create a new Book on the db
//
// POST /book
func TestCreateBook(t *testing.T) {

	// url := "https://localhost:8080/book"
	url := sessionData.baseURL + "/book"

	var jsonStr = []byte(`{"title":"string_value",
"author":"string_value",
"hardcover":true,
"copies":50000,
"library_id":50000}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+sessionData.jwtToken)

	if sessionData.log {
		fmt.Println("POST request Headers:", req.Header)
	}

	resp, err := sessionData.client.Do(req)
	if err != nil {
		t.Errorf("Test was unable to POST /book. Got %s.\n", err.Error())
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("GET /book response Body:", string(body))
		t.Errorf("Test was unable to POST /book. Got %s.\n", err.Error())
	}
	defer resp.Body.Close()

	var e models.Book
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&e); err != nil {
		t.Errorf("Test was unable to decode the result of POST /book. Got %s.\n", err.Error())
	}

	//============================================================================================
	// TODO: implement validation of the created entity here
	//============================================================================================
	if e.Title != "string_value" {
		t.Errorf("inconsistency detected in POST /book field Title.")
	}

	if e.Author == nil || *e.Author != "string_value" {
		t.Errorf("inconsistency detected in POST /book field Author.")
	}

	if e.Hardcover != true {
		t.Errorf("inconsistency detected in POST /book field Hardcover.")
	}

	if e.Copies != 50000 {
		t.Errorf("inconsistency detected in POST /book field Copies.")
	}

	if e.LibraryID != 50000 {
		t.Errorf("inconsistency detected in POST /book field LibraryID.")
	}

	if e.ID != 0 {
		sessionData.ID = e.ID
	} else {
		log.Printf("ID value of 0 detected - subsequent test cases will run with ID == 0!")
	}
}

// TestGetBooks attempts to read all books from the db
//
// GET /books
func TestGetBooks(t *testing.T) {

	// url := "https://localhost:8080/books"
	url := sessionData.baseURL + "/books"
	jsonStr := []byte(`{}`)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+sessionData.jwtToken)

	if sessionData.log {
		fmt.Println("GET /books request Headers:", req.Header)
	}

	// client := &http.Client{}
	resp, err := sessionData.client.Do(req)
	if err != nil {
		t.Errorf("Test was unable to GET /books. Got %s.\n", err.Error())
	}
	defer resp.Body.Close()

	if sessionData.log {
		fmt.Println("GET /books response Status:", resp.Status)
		fmt.Println("GET /books response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("GET /books response Body:", string(body))
	}
}

// TestGetBook attempts to read book/{:id} from the db
// using the id created in this entity's TestCreate function.
//
// GET /book/{:id}
func TestGetBook(t *testing.T) {

	idStr := fmt.Sprint(sessionData.ID)
	// url := "https://localhost:8080/book/" + idStr
	url := sessionData.baseURL + "/book/" + idStr
	jsonStr := []byte(`{}`)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+sessionData.jwtToken)

	if sessionData.log {
		fmt.Println("GET /book request Headers:", req.Header)
	}

	// client := &http.Client{}
	resp, err := sessionData.client.Do(req)
	if err != nil {
		t.Errorf("Test was unable to GET /book/%d. Got %s.\n", sessionData.ID, err.Error())
	}
	defer resp.Body.Close()

	if sessionData.log {
		fmt.Println("GET /book response Status:", resp.Status)
		fmt.Println("GET /book response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("GET /book response Body:", string(body))
	}
}

// TestUpdateBook attempts to update an existing Book on the db
//
// PUT /book/{:id}
func TestUpdateBook(t *testing.T) {

	idStr := fmt.Sprint(sessionData.ID)
	// url := "https://localhost:8080/book/" + idStr
	url := sessionData.baseURL + "/book/" + idStr

	var jsonStr = []byte(`{"title":"string_update",
"author":"string_update",
"hardcover":false,
"copies":99999,
"library_id":99999}`)

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+sessionData.jwtToken)

	if sessionData.log {
		fmt.Println("POST request Headers:", req.Header)
	}

	resp, err := sessionData.client.Do(req)
	if err != nil {
		t.Errorf("Test was unable to PUT /book/{:id}. Got %s.\n", err.Error())
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("PUT /book{:id} expected http status code of 201 - got %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	var e models.Book
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&e); err != nil {
		t.Errorf("Test was unable to decode the result of PUT /book. Got %s.\n", err.Error())
	}

	//============================================================================================
	// TODO: implement validation of the updated entity here
	//============================================================================================
	if e.Title != "string_update" {
		t.Errorf("inconsistency detected in POST /book field Title.")
	}

	if e.Author == nil || *e.Author != "string_update" {
		t.Errorf("inconsistency detected in POST /book field Author.")
	}

	if e.Hardcover != false {
		t.Errorf("inconsistency detected in POST /book field Hardcover.")
	}

	if e.Copies != 99999 {
		t.Errorf("inconsistency detected in POST /book field Copies.")
	}

	if e.LibraryID != 99999 {
		t.Errorf("inconsistency detected in POST /book field LibraryID.")
	}

	if e.ID != 0 {
		sessionData.ID = e.ID
	} else {
		log.Printf("ID value of 0 detected - subsequent test cases will run with ID == 0!")
	}
}

// TestDeleteBook attempts to delete the new Book on the db
//
// DELETE /book/{:id}
func TestDeleteBook(t *testing.T) {

	idStr := fmt.Sprint(sessionData.ID)
	// url := "https://localhost:8080/book/" + idStr
	url := sessionData.baseURL + "/book/" + idStr
	var jsonStr = []byte(`{}`)
	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(jsonStr))
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+sessionData.jwtToken)

	if sessionData.log {
		fmt.Println("sessionData.ID:", string(sessionData.ID))
		fmt.Println("DELETE URL:", url)
		fmt.Println("DELETE request Headers:", req.Header)
	}

	resp, err := sessionData.client.Do(req)
	if err != nil {
		t.Errorf("Test was unable to DELETE /book/%d. Got %s.\n", sessionData.ID, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		t.Errorf("DELETE /book{:id} expected http status code of 201 - got %d", resp.StatusCode)
	}
}

func TestGetBooksByTitle(t *testing.T) {

	// http://127.0.0.1:<port>/books/title(OP '<sel_string>')
	sessionData.testEndPoint = "/books/title(EQ 'test_string')"
	sessionData.testURL = sessionData.baseURL + sessionData.testEndPoint
	sessionData.testSelectableField(t)

	sessionData.testEndPoint = "/books/title(LIKE 'test_string')"
	sessionData.testURL = sessionData.baseURL + sessionData.testEndPoint
	sessionData.testSelectableField(t)

} // end func {Title string title  false 0  false true false nonUnique EQ,LIKE unknown title    sqac:"nullable:false;default:unknown title;index:non-unique" json:"title" false false false}

func TestGetBooksByAuthor(t *testing.T) {

	// http://127.0.0.1:<port>/books/author(OP '<sel_string>')
	sessionData.testEndPoint = "/books/author(EQ 'test_string')"
	sessionData.testURL = sessionData.baseURL + sessionData.testEndPoint
	sessionData.testSelectableField(t)

	sessionData.testEndPoint = "/books/author(LIKE 'test_string')"
	sessionData.testURL = sessionData.baseURL + sessionData.testEndPoint
	sessionData.testSelectableField(t)

} // end func {Author string author  false 0  false false false nonUnique EQ,LIKE     sqac:"nullable:true;index:non-unique" json:"author,omitempty" false false false}

func TestGetBooksByHardcover(t *testing.T) {

	// http://127.0.0.1:<port>/books/hardcover(OP true|false)
	sessionData.testEndPoint = "/books/hardcover(EQ true)"
	sessionData.testURL = sessionData.baseURL + sessionData.testEndPoint
	sessionData.testSelectableField(t)

	sessionData.testEndPoint = "/books/hardcover(EQ false)"
	sessionData.testURL = sessionData.baseURL + sessionData.testEndPoint
	sessionData.testSelectableField(t)

	sessionData.testEndPoint = "/books/hardcover(NE true)"
	sessionData.testURL = sessionData.baseURL + sessionData.testEndPoint
	sessionData.testSelectableField(t)

	sessionData.testEndPoint = "/books/hardcover(NE false)"
	sessionData.testURL = sessionData.baseURL + sessionData.testEndPoint
	sessionData.testSelectableField(t)

} // end func {Hardcover bool hardcover  false 0  false true false  EQ,NE     sqac:"nullable:false" json:"hardcover" false false false}

func TestGetBooksByCopies(t *testing.T) {

	// http://127.0.0.1:<port>/books/copies(OP XXX)
	sessionData.testEndPoint = "/books/copies(EQ 77)"
	sessionData.testURL = sessionData.baseURL + sessionData.testEndPoint
	sessionData.testSelectableField(t)

	sessionData.testEndPoint = "/books/copies(LT 77)"
	sessionData.testURL = sessionData.baseURL + sessionData.testEndPoint
	sessionData.testSelectableField(t)

	sessionData.testEndPoint = "/books/copies(GT 77)"
	sessionData.testURL = sessionData.baseURL + sessionData.testEndPoint
	sessionData.testSelectableField(t)

} // end func {Copies uint64 copies  false 0  false true false  EQ,LT,GT <nil>    sqac:"nullable:false" json:"copies" false false false}

func TestGetBooksByLibraryID(t *testing.T) {

	// http://127.0.0.1:<port>/books/library_id(OP XXX)
	sessionData.testEndPoint = "/books/library_id(EQ 77)"
	sessionData.testURL = sessionData.baseURL + sessionData.testEndPoint
	sessionData.testSelectableField(t)

	sessionData.testEndPoint = "/books/library_id(LIKE 77)"
	sessionData.testURL = sessionData.baseURL + sessionData.testEndPoint
	sessionData.testSelectableField(t)

} // end func {LibraryID uint64 library_id  false 0  false true false nonUnique EQ,LIKE <nil>    sqac:"nullable:false;index:non-unique" json:"library_id" false false false}
