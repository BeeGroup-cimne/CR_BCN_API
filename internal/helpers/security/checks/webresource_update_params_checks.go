package checks

import (
	"errors"
	"go-api/internal/helpers/utils"
	"go-api/internal/models"
)

func WebResourceUpdateParamsChecks(wrUpdate models.WebResourceUpdate, wr models.WebResource) error {
	if wrUpdate.Value == nil && wrUpdate.Type == nil && wrUpdate.Labels == nil && wrUpdate.Disabled == nil {
		return errors.New("Body is empty.")
	}
	if wrUpdate.Value != nil && *wrUpdate.Value == wr.Value {
		return errors.New("Value is repeated.")
	}
	if wrUpdate.Labels != nil && utils.ArraysEqual(*wrUpdate.Labels, wr.Labels) {
		return errors.New("Labels are repeated.")
	}
	if wrUpdate.Type != nil && *wrUpdate.Type == wr.Type {
		return errors.New("Type is repeated.")
	}
	if wrUpdate.Disabled != nil && *wrUpdate.Disabled == wr.Disabled {
		return errors.New("Disabled is repeated.")
	}
	return nil
}
