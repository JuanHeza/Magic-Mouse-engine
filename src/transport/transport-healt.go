package transport

import (
	"context"
	"encoding/json"
	"magic-mouse-engine/model"
	"net/http"
)

func decodeHealthzRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return model.HealthzRequest{}, nil
}

func encodeHealthzResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	r, ok := response.(model.HealthzResponse)
	if ok && r.Err != nil {
		encodeError(ctx, r.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(r.Body)
}
