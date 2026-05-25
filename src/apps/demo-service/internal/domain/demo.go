package domain

import "RuntimeRoasters/pkg/errs"

type Demo struct {
	ID      string
	Message string
}

func (d *Demo) Validate() error {
	if d.Message == "" {
		return errs.ErrValidation
	}
	return nil
}
