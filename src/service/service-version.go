package service

import (
	"context"
	"magic-mouse-engine/model"
	"magic-mouse-engine/util"
)

// @Summary Version Endpoint
// @Description Version Service.
// @Tags Version
// @ID Version
// @Accept json
// @Produce json
// @Sucess  200 {object} model.VersionResponseBody{}
// @Failure 400 {object} model.VersionResponseBody{}
// @Failure 500 {object} model.VersionResponseBody{}
// @Router /version [get]
func (svc *PosDigitalReceiptService) Version(ctx context.Context) (model.VersionResponseBody, error) {
	response := model.VersionResponseBody{Version: util.Version.Version,
		Commit: util.Version.Commit, BuildDate: util.Version.BuildDate}
	return response, nil
}
