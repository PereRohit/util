package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
)

func LoadAsJson(filepath string, cfg interface{}) error {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	v := reflect.ValueOf(cfg)
	if v.Type().Kind() != reflect.Ptr || !v.Elem().CanSet() {
		return fmt.Errorf("unable to set into given type: must be a pointer")
	}
	err = json.Unmarshal(content, cfg)
	if err != nil {
		return err
	}
	return nil
}
