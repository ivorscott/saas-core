package web

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

// Respond send a response back to the client.
func Respond(ctx context.Context, w http.ResponseWriter, val interface{}, statusCode int) error {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	if statusCode >= 400 {
		w.Header().Set("Content-Type", "application/problem+json")
	}

	if v, ok := ctx.Value(KeyValues).(*Values); ok {
		v.StatusCode = statusCode
	}

	if val != nil {
		res, err := json.Marshal(val)
		if err != nil {
			return err
		}

		_, err = w.Write(res)
		if err != nil {
			return err
		}
	}

	return nil
}

// RespondError sends an error response back to the client.
func RespondError(ctx context.Context, w http.ResponseWriter, err error) error {
	var webErr *Error

	if errors.As(err, &webErr) {
		if len(webErr.Fields) > 0 {
			var event = make(map[string]string)
			for _, v := range webErr.Fields {
				event[v.Field] = v.Error
			}
		}

		er := ErrorResponse{
			Error:  webErr.Err.Error(),
			Fields: webErr.Fields,
		}

		return Respond(ctx, w, er, webErr.Status)
	}

	er := ErrorResponse{
		Error: http.StatusText(http.StatusInternalServerError),
	}

	return Respond(ctx, w, er, http.StatusInternalServerError)
}
