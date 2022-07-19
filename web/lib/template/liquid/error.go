package liquid

import (
	"fmt"

	losherrors "losh/internal/errors"

	"github.com/osteele/liquid"
)

// TemplatingError is a error for loading or rendering a template.
type TemplatingError interface {
	error
	Unwrap() error
	FilePath() string
	Bindings() interface{}
}

type templatingError struct {
	losherrors.AppError
	filePath string
	binding  interface{}
}

func (e *templatingError) FilePath() string {
	return e.filePath
}

func (e *templatingError) Bindings() interface{} {
	return e.binding
}

func (e *templatingError) Error() string {
	return ""
}

// NewLoadError creates a loadError.
func NewLoadError(filePath string, cause error, binding interface{}) TemplatingError {
	err := &loadError{
		templatingError: templatingError{
			AppError: *losherrors.NewAppErrorWrap(cause, ""),
			filePath: filePath,
		},
	}
	err.Add("file_path", filePath)
	return err
}

type loadError struct {
	templatingError
}

func (e *loadError) Error() string {
	return formatError(&e.templatingError, "render")
}

// NewRenderError creates a renderError.
func NewRenderError(filePath string, cause error, binding interface{}) TemplatingError {
	err := &renderError{
		templatingError: templatingError{
			AppError: *losherrors.NewAppErrorWrap(cause, ""),
			filePath: filePath,
		},
	}
	err.Add("file_path", filePath)
	return err
}

type renderError struct {
	templatingError
}

func (e *renderError) Error() string {
	return formatError(&e.templatingError, "render")
}

func formatError(te *templatingError, stage string) string {
	location, cause := formatLocation(te)
	return fmt.Sprintf("template %s error %s: %s", stage, location, cause)
}

func formatLocation(err error) (location, cause string) {
	if err == nil {
		return "", ""
	}

	switch typedErr := err.(type) {
	case TemplatingError:
		location = fmt.Sprintf("(%s)", typedErr.FilePath())
		serr, ok := typedErr.Unwrap().(liquid.SourceError)
		if ok {
			if serr.LineNumber() > 0 {
				location = fmt.Sprintf("(%s:%d)", typedErr.FilePath(), serr.LineNumber()+1)
			}
			subLocation, subCause := formatLocation(serr.Cause())
			if subLocation != "" {
				location = fmt.Sprintf("%s ⇒ %s", location, subLocation)
			}
			if subCause != "" {
				cause = subCause
			} else {
				cause = serr.Error()
			}
		} else {
			cause = typedErr.Unwrap().Error()
		}

	case liquid.SourceError:
		if typedErr.Cause() != nil {
			return formatLocation(typedErr.Cause())
		}

	default:
		return "", typedErr.Error()
	}

	return
}
