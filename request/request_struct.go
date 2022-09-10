package request

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/PereRohit/util/validator"
)

func FromJson(r *http.Request, data interface{}) (int, error) {
	v := reflect.ValueOf(data)
	if v.Type().Kind() != reflect.Ptr || !v.Elem().CanSet() {
		return http.StatusInternalServerError, fmt.Errorf("unable to set into data: must be a pointer")
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("request body read : %s", err.Error())
	}
	defer r.Body.Close()
	err = json.Unmarshal(body, &data)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("put data into data: %s", err.Error())
	}
	err = validator.Validate(data)
	if err != nil {
		return http.StatusBadRequest, err
	}
	return 0, nil
}
