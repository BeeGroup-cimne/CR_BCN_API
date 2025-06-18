package checks

import (
	"errors"
	"go-api/internal/models"
)

func WebResourceCreateParamsChecks(wr models.WebResource) error {
	if wr.Value == "" {
		return errors.New("Value can't be empty.")
	}
	if wr.Type == "" {
		return errors.New("Type can't be empty.")
	}
	if len(wr.Labels) == 0 {
		return errors.New("Labels can't be empty.")
	}
	return nil
}
