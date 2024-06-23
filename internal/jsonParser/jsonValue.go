package jsonParser

import (
	"errors"
	"fmt"
)

/*
	JsonValue is returned by the JSON parser & is used to access items within the parsed result.

	Call `GetType("key")` to get values for keys, & `GetType("")` (with blank str arg) to get values
	when there is no key (for example in an array access).

	Example JSON:
	```
	{
		"theString": "a",
		"theInt": 1,
		"theFloat": 2.222,
		"theObj": {
			"objA": "b",
			"objB": 1,
			"objC": 2.222
		},
		"theArray": [
			[1,2]
		]
	}
	```

	To get data out after parsing:
	```
	jsonResult.GetString("theString") 	// Gets "a"
	jsonResult.GetInt("theInt") 		// Gets 1
	jsonResult.GetFloat("theFloat") 	// Gets 2.222

	theObj := jsonResult.GetObject("theObj") // Gets theObj as a JsonValue
	theObj.GetString("objA")			 	 // Gets "b"
	theObj.GetInt("objB")			 	 	 // Gets 1
	theObj.GetFloat("objC")			 	 	 // Gets 2.222

	theArray := jsonResult.GetArray("theArray") // Gets theArray as a slice of JsonValue
	for _, item := range theArray {
		v, _ := item.GetInt("")					// Gets 1, then 2
	}
	```
*/

type JsonValue struct {
	data any
}

func NewJsonValue(data any) *JsonValue {
	return &JsonValue{data: data}
}

func (j *JsonValue) getKeyValue(key string) (any, error) {
	val, ok := j.data.(map[string]any)[key]
	if !ok {
		msg := fmt.Sprintf(`Key "%s" not found"`, key)
		return "", errors.New(msg)
	}

	return val, nil
}


// Returns a string for the given key, or if key is blank, returns own data as string
func (j *JsonValue) GetString(key string) (string, error) {
	var val any
	var err error

	if key != "" {
		val, err = j.getKeyValue(key)
		if err != nil {
			return "", errors.New(err.Error())
		}
	} else {
		val = j.data
	}

	strVal, strOk := val.(string)
	if !strOk {
		strMsg := fmt.Sprintf(`Error casting "%s" to string`, val)
		return "", errors.New(strMsg)
	}
	return strVal, nil
}

// Returns an int for the given key, or if key is blank, returns own data as int
func (j *JsonValue) GetInt(key string) (int, error) {
	var val any
	var err error

	if key != "" {
		val, err = j.getKeyValue(key)
		if err != nil {
			return 0, errors.New(err.Error())
		}
	} else {
		val = j.data
	}

	intVal, intOk := val.(int)
	if !intOk {
		intMsg := fmt.Sprintf(`Error casting "%s" to int`, val)
		return 0, errors.New(intMsg)
	}

	return intVal, nil
}

// Returns a float64 for the given key, or if key is blank, returns own data as float64
func (j *JsonValue) GetFloat(key string) (float64, error) {
	var val any
	var err error

	if key != "" {
		val, err = j.getKeyValue(key)
		if err != nil {
			return 0, errors.New(err.Error())
		}
	} else {
		val = j.data
	}

	floatVal, floatOk := val.(float64)
	if !floatOk {
		floatMsg := fmt.Sprintf(`Error casting "%s" to float`, val)
		return 0, errors.New(floatMsg)
	}
	return floatVal, nil
}

// Returns a *JsonValue for the given key, or if key is blank, returns own data as *JsonValue
func (j *JsonValue) GetObject(key string) (*JsonValue, error) {
	var val any
	var err error

	if key != "" {
		val, err = j.getKeyValue(key)
		if err != nil {
			return nil, errors.New(err.Error())
		}
	} else {
		val = j.data
	}

	objectVal, objectOk := val.(map[string]any)
	if !objectOk {
		objectMsg := fmt.Sprintf(`Error casting "%s" to object`, val)
		return nil, errors.New(objectMsg)
	}

	return &JsonValue{objectVal}, nil
}

// Returns a []*JsonValue for the given key, or if key is blank, returns own data as []*JsonValue
func (j *JsonValue) GetArray(key string) ([]*JsonValue, error) {
	var val any
	var err error

	if key != "" {
		val, err = j.getKeyValue(key)
		if err != nil {
			return nil, errors.New(err.Error())
		}
	} else {
		val = j.data
	}

	arrayVal, arrayOk := val.([]any)
	if !arrayOk {
		arrayMsg := fmt.Sprintf(`Error casting "%s" to array`, val)
		return nil, errors.New(arrayMsg)
	}

	resultArray := make([]*JsonValue, len(arrayVal))
	for i, v := range arrayVal {
		resultArray[i] = &JsonValue{v}
	}

	return resultArray, nil
}
