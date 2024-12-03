package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// ValidationError representa un error de validación con detalles
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationErrors es una colección de errores de validación
type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	var errMsgs []string
	for _, err := range e {
		errMsgs = append(errMsgs, err.Error())
	}
	return strings.Join(errMsgs, "; ")
}

// Validator es una interfaz que deben implementar las estructuras que requieren validación
type HttpValidator interface {
	Validate() ValidationErrors
}

// DecodeAndValidate decodes and validates an HTTP request based on the Content-Type.
func DecodeAndValidate[T HttpValidator](r *http.Request) (T, error) {
	var payload T

	switch r.Header.Get("Content-Type") {
	case "application/json":
		if err := decodeJSON(r, &payload); err != nil {
			return payload, err
		}
	case "application/x-www-form-urlencoded", "multipart/form-data":
		if err := decodeForm(r, &payload); err != nil {
			return payload, err
		}
	default:
		return payload, ValidationError{
			Field:   "Content-Type",
			Message: "unsupported content type",
		}
	}

	// Validate the decoded struct.
	if errs := payload.Validate(); len(errs) > 0 {
		return payload, errs
	}

	return payload, nil
}

// decodeJSON decodes a JSON request body into the given struct.
func decodeJSON[T any](r *http.Request, target *T) error {
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		return ValidationError{
			Field:   "body",
			Message: "invalid JSON payload",
		}
	}
	return nil
}

// decodeForm decodes form data and populates the given struct.
func decodeForm[T any](r *http.Request, target *T) error {
	if err := r.ParseForm(); err != nil {
		return ValidationError{
			Field:   "body",
			Message: "invalid form data",
		}
	}
	return populateStructFromForm(r.Form, target)
}

// populateStructFromForm populates the struct fields from form data using reflection.
func populateStructFromForm[T any](form url.Values, outputStruct *T) error {
	v := reflect.ValueOf(outputStruct)
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
		return errors.New("outputStruct must be a pointer to a struct")
	}

	v = v.Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		formTag := field.Tag.Get("form")
		if formTag == "" {
			continue
		}

		formValue, exists := form[formTag]
		if !exists || len(formValue) == 0 {
			continue
		}

		err := setFieldValue(v.Field(i), formValue[0])
		if err != nil {
			return fmt.Errorf("failed to set field %q: %v", field.Name, err)
		}
	}
	return nil
}

// setFieldValue sets the value of a reflect.Value from a string value.
func setFieldValue(field reflect.Value, value string) error {
	if !field.CanSet() {
		return errors.New("cannot set field")
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(intValue)
	case reflect.Float32, reflect.Float64:
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(floatValue)
	case reflect.Bool:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(boolValue)
	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind().String())
	}

	return nil
}
