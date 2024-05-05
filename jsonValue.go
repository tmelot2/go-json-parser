package main

import (
	"errors"
	"fmt"
)

// TODO: Doc all this stuff mang

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

func (j *JsonValue) GetString(key string) (string, error) {
	val, err := j.getKeyValue(key)
	if err != nil {
		return "", errors.New(err.Error())
	}

	strVal, strOk := val.(string)
	if !strOk {
		strMsg := fmt.Sprintf(`Error casing "%s" to string`, val)
		return "", errors.New(strMsg)
	}
	return strVal, nil
}

func (j *JsonValue) GetInt(key string) (int, error) {
	val, err := j.getKeyValue(key)
	if err != nil {
		return 0, errors.New(err.Error())
	}

	intVal, intOk := val.(int)
	if !intOk {
		intMsg := fmt.Sprintf(`Error casing "%s" to int`, val)
		return 0, errors.New(intMsg)
	}
	return intVal, nil
}