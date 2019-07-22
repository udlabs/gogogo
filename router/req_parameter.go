package router

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"gogogo/util"
	"io/ioutil"
	"net/http"
	"strings"
)

// Apply filter on request object and pipe out only defined properties
func getVariables(req *http.Request, route Route) (error, map[string]interface{}) {
	var variables = make(map[string]interface{})
	var err error

	// extract url
	err = getParams(req, variables, route.Parameters.Url, ReqParamUrlLabel)
	if err != nil {
		return err, nil
	}

	// extract query strings
	err = getParams(req, variables, route.Parameters.Query, ReqParamQueryLabel)
	if err != nil {
		return err, nil
	}

	// extract http request header
	err = getParams(req, variables, route.Parameters.Headers, ReqParamHeaderLabel)
	if err != nil {
		return err, nil
	}


	contentType := strings.ToLower(req.Header.Get(HttpContentTypeLabel))

	if !strings.HasPrefix(contentType, HttpContentTypeJson) {
		return nil, variables
	}

	// extract req body
	var body []byte
	body, err = ioutil.ReadAll(req.Body)
	if err != nil {
		return err, nil
	}

	if len(body) == 0 {
		return nil, variables
	}

	var bJson interface{}
	err = json.Unmarshal(body, &bJson)
	if err != nil {
		return err, nil
	}
	variables[HttpReqBodyLabel] = bJson
	return nil, variables
}

func getParams(req *http.Request, variables map[string]interface{}, params []Parameter, paramType string) error {
	for _, parameter := range params {
		parameter.ParamType = paramType
		err, value := validate(getParamValue(req, paramType, parameter), parameter)
		if err != nil {
			return err
		}
		keyName := getKeyName(HttpReqParamsPrefix, paramType, parameter.Name)
		variables[keyName] = value
	}
	return nil
}

func getKeyName(keys ...string) string {
	return util.JoinBy(HttpReqParamsSeparator, keys...)
}

func getParamValue(req *http.Request, paramType string, param Parameter) interface{} {
	switch paramType {
	case ReqParamUrlLabel:
		return chi.URLParam(req, param.Name)
	case ReqParamHeaderLabel:
		return req.Header.Get(param.Name)
	case ReqParamQueryLabel:
		return req.URL.Query().Get(param.Name)
	}
	return nil
}

func validate(value interface{}, param Parameter) (error, interface{}) {
	isEmpty := false
	switch value.(type) {
	case string:
		isEmpty = value.(string) == ""
	}

	if param.Mandatory && (value == nil || isEmpty) {
		return errors.New("missing mandatory field " + param.Name + " in " + param.ParamType), nil
	}
	if value == nil || isEmpty {
		return nil, param.Default
	}
	return nil, value
}

