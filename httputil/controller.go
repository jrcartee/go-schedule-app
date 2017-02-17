package httputil

import (
	"net/http"
	"github.com/jrcartee/scheduling/database"
)

type Model interface {
	// querying
	Instance(http.ResponseWriter, *http.Request) interface{}
	SelectAll() (interface{}, error)
	Insert(interface{}) error
	Update(interface{}) error
	// entity info
	New() interface{}
}


type ValidationErrors map[string]string
type Validator interface {
	ValidateInsert(interface{}, database.Queryer) ValidationErrors
	ValidateUpdate(interface{}, database.Queryer) ValidationErrors
}

func GenericIndex(m Model) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		set, err := m.SelectAll()
		if err != nil {
			InternalError(w, err)
			return
		}
		EncodeJSONResponse(w, set)

	}
}

func GenericRetrieve(m Model) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		x := m.Instance(w, r)
		if x == nil {
			return
		}
		EncodeJSONResponse(w, x)

	}
}

func GenericInsert(m Model, db database.Queryer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmp := m.New()
		err := DecodeJSONRequest(r, tmp)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}


		if v, ok := m.(Validator); ok  {
			errs := v.ValidateInsert(tmp, db)
			if !isValid(errs, w) {				
				return
			}
		}


		err = m.Insert(tmp)
		if err != nil {
			InternalError(w, err)
			return
		}

		err = EncodeJSONResponse(w, tmp)
		if err != nil {
			InternalError(w, err)
			return
		}

	}
}

func GenericUpdate(m Model, db database.Queryer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		x := m.Instance(w, r)
		if x == nil {
			return
		}

		err := DecodeJSONRequest(r, x)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		
		if v, ok := m.(Validator); ok {
			errs := v.ValidateUpdate(x, db)
			if !isValid(errs, w) {				
				return
			}
		}

		err = m.Update(x)
		if err != nil {
			InternalError(w, err)
			return
		}

		err = EncodeJSONResponse(w, x)
		if err != nil {
			InternalError(w, err)
			return
		}
	}
}

func isValid(e map[string]string, w http.ResponseWriter) bool {
	if len(e) != 0 {
		w.WriteHeader(http.StatusBadRequest)
		err := EncodeJSONResponse(w, e)
		if err != nil {
			InternalError(w, err)
		}
		return false
	}
	return true
}


func MergeErrors(l ...ValidationErrors) ValidationErrors {
	tmp := make(ValidationErrors, 0)
	for _, errs := range l {
		for k, v := range errs {
			tmp[k] = v
		}
	}
	return tmp
}