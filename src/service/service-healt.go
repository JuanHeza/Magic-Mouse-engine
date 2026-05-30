package service

import (
	"context"
	"magic-mouse-engine/model"
)

// Healthz Api endpoint to check service health
// @Summary Healthz Endpoint
// @Description this function is the end point for Healthz it will response ok.
// @Tags healtz
// @ID Healtz
// @Accept json
// @Produce json
// @Sucess  200 {object} model.HealthzResponseBody{}
// @Failure 400 {object} model.HealthzResponseBody{}
// @Failure 500 {object} model.HealthzResponseBody{}
// @Router /healthz [get]
func (svc *PosDigitalReceiptService) Healthz(ctx context.Context) (model.HealthzResponseBody, error) {
	var response model.HealthzResponseBody
	response.Status = "ok"
	return response, nil
}
