package domain

import "github.com/dungxbuif/RuntimeRoasters/pkg/errs"

type Farm struct {
	ID      string
	Message string
}

func (d *Farm) Validate() error {
	if d.Message == "" {
		return errs.ErrValidation
	}
	return nil
}
