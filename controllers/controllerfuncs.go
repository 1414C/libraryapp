package controllers

//=============================================================================================
//
//  controllerfuncs.go contains common discrete functions that are leveraged
//  throughout the controllers package.
//
// generated code: please do not modify
// code generated on 27 May 20 17:32 CDT
//=============================================================================================

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// send a generic err response in JSON
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// send an arbitrary response to the caller
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// send a count
func respondWithCount(w http.ResponseWriter, code int, count uint64) {

	w.Header().Set("Content-Type", "application/text")
	w.WriteHeader(code)
	cs := strconv.FormatUint(count, 10)
	w.Write([]byte(cs))
}

// send a http-header-only response
func respondWithHeader(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(nil)
}

// parseRequestCommands accepts the value of mux.Vars() from the calling
// entity-controller looks for a "cmd" key in the map, and then attempts
// to validate the cmd content (map value) and parse it into a command map.
//
// valid commands are as follows:
// $count               :return a count of records
// $top=n               :return the top n records (like DB limit)
// $skip=n              :skip the first n records in the selected set
// $limit=n             :set limit n for record retrieval
// $orderby=field_name  :orderby the db field name (default is id)
// $asc                 :sort selected records in ascending order based on key (ID? need to think about this one)
// $desc                :sort selected records in descending order based on key (ID? need to think about this one)
func parseRequestCommands(vars map[string]string) (map[string]interface{}, error) {

	strCmds := vars["cmd"]
	if strCmds == "" {
		// this feels like an error condition, but it can't be
		return nil, nil
	}

	commands := strings.Split(strCmds, "$")
	mapCommands := make(map[string]interface{})

	for _, v := range commands {
		if v != "" {
			if strings.ContainsAny(v, "=") {
				p := strings.Split(v, "=")
				mapCommands[p[0]] = p[1]
			} else {
				mapCommands[v] = nil
			}
		}
	}
	return mapCommands, nil
}

//---------------------------------------------------------------------------------------------
// support functions for Href creation
//---------------------------------------------------------------------------------------------
// buildHrefBasic builds a base URL string with no entity or path information
func buildHrefBasic(r *http.Request, withSuffix bool) string {

	// fill the urlString for the response Href
	var urlString string
	if r.TLS != nil {
		urlString = "https://"
	} else {
		urlString = "http://"
	}

	if withSuffix {
		return urlString + r.Host + "/"
	}
	return urlString + r.Host
}

// buildHrefStringFromCRUDReq builds a rawURLString from data in a CRUD request
func buildHrefStringFromCRUDReq(r *http.Request, withSuffix bool) string {

	// fill the urlString for the response Href
	var urlString string
	if r.TLS != nil {
		urlString = "https://"
	} else {
		urlString = "http://"
	}
	if withSuffix {
		return urlString + r.Host + r.URL.String() + "/"
	}
	return urlString + r.Host + r.URL.String()
}

// buildHrefStringFromSimpleQueryReq builds a rawURLString from data in a simple query request
// for example, http://localhost:3000/entity/name(LIKE 'tha')
// is converted to http://localhost:3000/entity/name/
func buildHrefStringFromSimpleQueryReq(r *http.Request, withSuffix bool) (string, error) {

	// fill the urlString for the response Href
	var urlString string
	if r.TLS != nil {
		urlString = "https://"
	} else {
		urlString = "http://"
	}

	// %-encoded - unescape the URL.String
	endPoint := r.URL.String()
	var err error
	if strings.ContainsAny(endPoint, "%") {
		endPoint, err = url.PathUnescape(endPoint)
		if err != nil {
			return "", err
		}
	}

	sl := strings.Split(endPoint, "(")
	if len(sl) != 2 {
		return "", fmt.Errorf("unable to determine base Href for request.  got %s", endPoint)
	}

	if withSuffix {
		return urlString + r.Host + sl[0] + "/", nil
	}
	return urlString + r.Host + sl[0], nil
}

//---------------------------------------------------------------------------------------------
// support functions for simple query construction
//---------------------------------------------------------------------------------------------
// buildStringQueryComponents separates the incoming simple query-string into its consituent
// parts; namely operation and predicate.  This is used for all generated entity simple accessors
// where the predicate field is of the string-type.
func buildStringQueryComponents(searchString string) (op, predicate string, err error) {

	searchString = strings.TrimPrefix(searchString, "(")
	searchString = strings.TrimSuffix(searchString, ")")
	sc := strings.SplitN(searchString, "'", 2)
	for i, s := range sc {
		if i == 0 {
			sc[i] = strings.ToUpper(strings.TrimSuffix(s, " "))
		} else {
			sc[i] = strings.TrimSuffix(s, "'")
		}
		fmt.Println([]byte(sc[i]))
	}

	op = sc[0]
	predicate = sc[1]

	switch op {
	case "NE":
		// do not modify predicate
	case "EQ":
		// do not modify predicate
	case "LIKE":
		if predicate != "" {
			predicate = "%" + predicate + "%"
		}
	case "LT":
		// do not modify predicate
	case "LE":
		// do not modify predicate
	case "GT":
		// do not modify predicate
	case "GE":
		// do not modify predicate
	default:
		// do nothing
	}
	return op, predicate, nil
}

// convert the model ops into sql operands
func convertOp(op string) (string, error) {

	switch op {
	case "NE":
		op = "!="
	case "EQ":
		op = "="
	case "LIKE":
		op = "LIKE"
	case "LT":
		op = "<"
	case "LE":
		op = "<="
	case "GT":
		op = ">"
	case "GE":
		op = ">="
	default:
		// unknown operand - error
		return "", fmt.Errorf("error: %s is an unknown selection operand", op)
	}
	return op, nil
}

// buildUIntQueryComponentGeneral separates the incoming simple query-string into its consituent
// parts; namely operation and predicate.  This is used as the foundation for all generated entity
// simple accessors where the predicate field is of the uint, uint8, uint16, uint32 and uint64-type.
func buildUIntQueryComponentGeneral(searchString string) (op string, predicate uint, err error) {

	searchString = strings.TrimPrefix(searchString, "(")
	searchString = strings.TrimSuffix(searchString, ")")
	sc := strings.SplitN(searchString, " ", 2)
	if len(sc) == 2 {
		sc[0] = strings.ToUpper(strings.TrimSuffix(sc[0], " "))
	} else {
		// this can't really happen based on the regex in the route definition
		return "", 0, fmt.Errorf("bad query construction")
	}

	op, err = convertOp(sc[0])
	if err != nil {
		return "", 0, err
	}

	p, err := strconv.ParseUint(sc[1], 10, 0)
	if err != nil {
		return "", 0, fmt.Errorf("unable to convert string %s to uint-value", sc[1])
	}

	predicate = uint(p)
	return op, predicate, nil
}

// buildUIntQueryComponent returns a uint predicate for the search parameter
func buildUIntQueryComponent(searchString string) (op string, predicate uint, err error) {

	return buildUIntQueryComponentGeneral(searchString)
}

// buildUInt8QueryComponent returns a uint8 predicate for the search parameter
func buildUInt8QueryComponent(searchString string) (op string, predicate uint8, err error) {

	op, interimPredicate, err := buildUIntQueryComponentGeneral(searchString)
	if err != nil {
		return "", 0, err
	}
	return op, uint8(interimPredicate), nil
}

// buildUInt16QueryComponent returns a uint16 predicate for the search parameter
func buildUInt16QueryComponent(searchString string) (op string, predicate uint16, err error) {

	op, interimPredicate, err := buildUIntQueryComponentGeneral(searchString)
	if err != nil {
		return "", 0, err
	}
	return op, uint16(interimPredicate), nil
}

// buildUInt32QueryComponent returns a uint32 predicate for the search parameter
func buildUInt32QueryComponent(searchString string) (op string, predicate uint32, err error) {

	op, interimPredicate, err := buildUIntQueryComponentGeneral(searchString)
	if err != nil {
		return "", 0, err
	}
	return op, uint32(interimPredicate), nil
}

// buildUInt64QueryComponent returns a uint64 predicate for the search parameter
func buildUInt64QueryComponent(searchString string) (op string, predicate uint64, err error) {

	op, interimPredicate, err := buildUIntQueryComponentGeneral(searchString)
	if err != nil {
		return "", 0, err
	}
	return op, uint64(interimPredicate), nil
}

// buildIntQueryComponentGeneral separates the incoming simple query-string into its consituent
// parts; namely operation and predicate.  This is used This is used as the foundation for all
// generated entity simple accessors where the predicate field is of the int, int8, int16, int32
// and int64-type.
func buildIntQueryComponentGeneral(searchString string) (op string, predicate int, err error) {

	searchString = strings.TrimPrefix(searchString, "(")
	searchString = strings.TrimSuffix(searchString, ")")
	sc := strings.SplitN(searchString, " ", 2)
	if len(sc) == 2 {
		sc[0] = strings.ToUpper(strings.TrimSuffix(sc[0], " "))
	} else {
		// this can't really happen based on the regex in the route definition
		return "", 0, fmt.Errorf("bad query construction")
	}

	op, err = convertOp(sc[0])
	if err != nil {
		return "", 0, err
	}

	predicate, err = strconv.Atoi(sc[1])
	if err != nil {
		return "", 0, fmt.Errorf("unable to convert string %s to int-value", sc[1])
	}
	return op, predicate, nil
}

// buildIntQueryComponent returns an int predicate for the search parameter
func buildIntQueryComponent(searchString string) (op string, predicate int, err error) {

	return buildIntQueryComponentGeneral(searchString)
}

// buildInt8QueryComponent returns an int8 predicate for the search parameter
func buildInt8QueryComponent(searchString string) (op string, predicate int8, err error) {

	op, interimPredicate, err := buildIntQueryComponentGeneral(searchString)
	if err != nil {
		return "", 0, err
	}
	return op, int8(interimPredicate), nil
}

// buildInt16QueryComponent returns an int16 predicate for the search parameter
func buildInt16QueryComponent(searchString string) (op string, predicate int16, err error) {

	op, interimPredicate, err := buildIntQueryComponentGeneral(searchString)
	if err != nil {
		return "", 0, err
	}
	return op, int16(interimPredicate), nil
}

// buildInt32QueryComponent returns an int32 predicate for the search parameter
func buildInt32QueryComponent(searchString string) (op string, predicate int32, err error) {

	op, interimPredicate, err := buildIntQueryComponentGeneral(searchString)
	if err != nil {
		return "", 0, err
	}
	return op, int32(interimPredicate), nil
}

// buildInt64QueryComponent returns an int64 predicate for the search parameter
func buildInt64QueryComponent(searchString string) (op string, predicate int64, err error) {

	op, interimPredicate, err := buildIntQueryComponentGeneral(searchString)
	if err != nil {
		return "", 0, err
	}
	return op, int64(interimPredicate), nil
}

// build simple query components for float search predicates
func buildFloatQueryComponentGeneral(searchString string) (op string, predicate float64, err error) {

	searchString = strings.TrimPrefix(searchString, "(")
	searchString = strings.TrimSuffix(searchString, ")")
	sc := strings.SplitN(searchString, " ", 2)
	if len(sc) == 2 {
		sc[0] = strings.ToUpper(strings.TrimSuffix(sc[0], " "))
	} else {
		// this can't really happen based on the regex in the route definition
		return "", 0, fmt.Errorf("bad query construction")
	}

	op, err = convertOp(sc[0])
	if err != nil {
		return "", 0, err
	}

	predicate, err = strconv.ParseFloat(sc[1], 64)
	if err != nil {
		return "", 0, fmt.Errorf("unable to convert string %s to float64-value", sc[1])
	}
	return op, predicate, nil
}

// buildFloat32QueryComponent returns a float32 predicate for the search parameter
func buildFloat32QueryComponent(searchString string) (op string, predicate float32, err error) {

	op, interimPredicate, err := buildFloatQueryComponentGeneral(searchString)
	if err != nil {
		return "", 0, err
	}
	return op, float32(interimPredicate), nil
}

// buildFloat64QueryComponent returns a float64 predicate for the search parameter
func buildFloat64QueryComponent(searchString string) (op string, predicate float64, err error) {

	op, interimPredicate, err := buildFloatQueryComponentGeneral(searchString)
	if err != nil {
		return "", 0, err
	}
	return op, float64(interimPredicate), nil
}

// build simple query components for bool search predicates
func buildBoolQueryComponents(searchString string) (op string, predicate bool, err error) {

	searchString = strings.TrimPrefix(searchString, "(")
	searchString = strings.TrimSuffix(searchString, ")")
	sc := strings.SplitN(searchString, " ", 2)

	// upper-case for operator and predicate
	for i, s := range sc {
		if i == 0 {
			sc[i] = strings.ToUpper(strings.TrimSuffix(s, " "))
		} else {
			sc[i] = strings.ToUpper(s)
		}
	}

	// this can't really happen based on the regex in the route definition
	if len(sc) != 2 {
		return "", false, fmt.Errorf("unable to convert string %s to bool value", searchString)
	}

	op, err = convertOp(sc[0])
	if err != nil {
		return "", false, err
	}

	if sc[1] == "TRUE" {
		predicate = true
	} else {
		predicate = false
	}
	return op, predicate, nil
}

//=============================================================================================
// end of generated code
//=============================================================================================
