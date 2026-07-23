package app

import "fmt"

// Stable exit categories. These are part of the CLI contract; automation may
// depend on them.
const (
	ExitOK      = 0
	ExitUsage   = 2 // usage or configuration error
	ExitCatalog = 3 // catalog or network error
	ExitSource  = 4 // source retrieval error
	ExitSkill   = 5 // invalid skill
	ExitFS      = 6 // filesystem collision or permission error
	ExitAdapter = 7 // partial adapter failure with no activation
	ExitState   = 8 // state or lock error
)

// ExitError carries an exit category alongside the underlying error.
type ExitError struct {
	Code int
	Err  error
}

func (e *ExitError) Error() string { return e.Err.Error() }
func (e *ExitError) Unwrap() error { return e.Err }

func exitErr(code int, err error) *ExitError { return &ExitError{Code: code, Err: err} }

func errf(code int, format string, args ...any) *ExitError {
	return &ExitError{Code: code, Err: fmt.Errorf(format, args...)}
}
