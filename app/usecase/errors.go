package usecase

import (
	"fmt"
)

// NoTargetError an error that a target does not exists for path
type NoTargetError struct {
	path string
}

// NewNoTargetError create a NoTargetError
func NewNoTargetError(path string) *NoTargetError {
	return &NoTargetError{path}
}

func (e *NoTargetError) Error() string {
	return fmt.Sprintf("no target for path %s", e.path)
}
