package transport

import (
	"context"
	"encoding/json"
	"magic-mouse-engine/model"
	"net/http"
)

func decodeVersionRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return model.VersionRequest{}, nil
}

func encodeVersionResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	r, ok := response.(model.VersionResponse)
	if ok && r.Err != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, r.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(r.Body)
}
