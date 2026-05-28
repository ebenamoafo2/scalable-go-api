package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/ebenamoafo2/scalable-go-api/internal/validator"
	"github.com/julienschmidt/httprouter"
)

type envelope map[string]any

func (app *application) readIdParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	//Encode the data to json, returning the error if it fails
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	for key, value := range headers {
		w.Header()[key] = value
	}

	js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(js)
	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	//Use the http.MaxBytesReader to cap the request body at 1MB before any reading begins,
	//preventing memory exhaustion from oversized payloads.
	r.Body = http.MaxBytesReader(w, r.Body, 1_048_576)

	//Create a new decoder from the request body
	//and disallow unknown fields in the JSON
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	//Decode the json request body into the target destination
	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshallTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {

		// Malformed JSON with a known character offset.
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly formed-JSON (at character %d)", syntaxError.Offset)

		// Malformed JSON with no recoverable offset.
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly formed-JSON")

		// Wrong JSON type for a field; include field name if available, else offset.
		case errors.As(err, &unmarshallTypeError):
			if unmarshallTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for the field %q", unmarshallTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshallTypeError.Offset)

		//Request body was empty
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		// Unknown field in the JSON body; extract and return the field name.
		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		//Request body was too large
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes ", maxBytesError.Limit)

		//dst was not a pointer to a struct- programming error, panic immediately
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	//Ensure that the request body is fully consumed
	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must contain a single JSON value")
	}
	return nil
}

// Read a string from the query string and returns it
func (app *application) readStrings(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}

	return s
}

// Read a comma-separated string from the query string and returns it as a slice of strings
func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {

	csv := qs.Get(key)

	if csv == "" {
		return defaultValue
	}

	return strings.Split(csv, ",")
}

// Read an integer from the query string and returns it
func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer")
		return defaultValue
	}
	return i
}
