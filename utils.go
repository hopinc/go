package hop

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"go.hop.io/sdk/types"
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
		return types.ServerError("status code " + strconv.Itoa(res.StatusCode) + " (cannot unmarshal from json): " +
			string(b))
	}
	if r.Success {
		return errors.New("api response error: error request was marked as success - please report this to " +
			"the go-hop github repository")
	}

	if r.Error.Code == "invalid_auth" {
		return types.NotAuthorized(r.Error.Message)
	}
	if res.StatusCode == 400 {
		// As a special case, for 400s return the r.Error object.
		return r.Error
	}
	if res.StatusCode == 404 {
		// Infer that this is a not found error.
		return types.NotFound(r.Error.Message)
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
