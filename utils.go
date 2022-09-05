package hopgo

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/hopinc/hop-go/types"
)

type errorResponse struct {
	Success bool             `json:"success"`
	Error   types.BadRequest `json:"error"`
}

func handleErrors(res *http.Response) error {
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var r errorResponse
	err = json.Unmarshal(b, &r)
	if err != nil {
		return types.ServerError("status code " + strconv.Itoa(res.StatusCode) + " (cannot turn into json): " +
			string(b))
	}
	if r.Success {
		return errors.New("api response error: error request was marked as success - please report this to " +
			"the go-hop github repository")
	}

	switch r.Error.Code {
	case "not_found", "project_not_found":
		return types.NotFound(r.Error.Message)
	case "invalid_auth":
		return types.NotAuthorized(r.Error.Message)
	}
	if res.StatusCode == 400 {
		// As a special case, for 400s return the r.Error object.
		return r.Error
	}
	if res.StatusCode >= 500 {
		// Infer that this is a internal server error.
		return types.ServerError(r.Error.Message)
	}
	return types.UnknownServerError{
		StatusCode: res.StatusCode,
		Message:    r.Error.Message,
		Code:       r.Error.Code,
	}
}
