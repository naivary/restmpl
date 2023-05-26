package testutil

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func AreEqualJSON(s1, s2 string) (bool, error) {
	var o1 any
	var o2 any

	var err error
	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false, fmt.Errorf("error mashalling string 1 :: %s", err.Error())
	}
	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false, fmt.Errorf("error mashalling string 2 :: %s", err.Error())
	}

	return reflect.DeepEqual(o1, o2), nil
}
