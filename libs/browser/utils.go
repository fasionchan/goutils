package browser

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/form/v4"
)

var (
	queryDecoder = form.NewDecoder()
)

type RequestParams[
	Body any,
	Query any,
	Path any,
] struct {
	Body Body
	Query Query
	Path Path
}

func ParseRequest[
	Request any,
](r *http.Request) (result Request, err error) {
	if err = queryDecoder.Decode(&result, r.URL.Query()); err != nil {
		return
	}

	switch r.Header.Get("Content-Type") {
	case "application/json":
		if err = json.NewDecoder(r.Body).Decode(&result); err != nil {
			return
		}
	case "application/x-www-form-urlencoded":
		if err = form.NewDecoder().Decode(&result, r.URL.Query()); err != nil {
			return
		}
	}

	return
}

func init() {
	queryDecoder.SetTagName("query")
	queryDecoder.SetMode(form.ModeExplicit)
}