package liquid

import (
	"fmt"

	"github.com/osteele/liquid"
)

// TemplatingError is a error for loading or rendering a template.
type TemplatingError interface {
	error
	FilePath() string
	Cause() error
	Bindings() interface{}
}

type templatingError struct {
	filePath string
	cause    error
	binding  interface{}
}

func (e *templatingError) FilePath() string {
	return e.filePath
}

func (e *templatingError) Cause() error {
	return e.cause
}

func (e *templatingError) Bindings() interface{} {
	return e.binding
}

func (e *templatingError) Error() string {
	return ""
}

// NewLoadError creates a loadError.
func NewLoadError(filePath string, cause error, binding interface{}) TemplatingError {
	return &loadError{
		templatingError: templatingError{
			filePath: filePath,
			cause:    cause,
		},
	}
}

type loadError struct {
	templatingError
}

func (e *loadError) Error() string {
	return formatError(&e.templatingError, "render")
}

// NewRenderError creates a renderError.
func NewRenderError(filePath string, cause error, binding interface{}) TemplatingError {
	return &renderError{
		templatingError: templatingError{
			filePath: filePath,
			cause:    cause,
		},
	}
}

type renderError struct {
	templatingError
}

func (e *renderError) Error() string {
	return formatError(&e.templatingError, "render")
}

func formatErrorBak(te templatingError, stage string) string {
	cause := ""
	location := fmt.Sprintf("(%s)", te.filePath)
	if te.cause != nil {
		sourceError, ok := te.cause.(liquid.SourceError)
		if ok {
			if sourceError.LineNumber() > 0 {
				location = fmt.Sprintf("(%s, %d)", te.filePath, sourceError.LineNumber()+1)
			}
			c := sourceError.(error)
			for {
				if c != nil {
					fmt.Println("ERROR: ", c)
					if tc, ok := c.(liquid.SourceError); ok {
						c = tc.Cause()
					} else {
						break
					}
				}
			}
			if sourceError.Cause() != nil {
				cause = sourceError.Cause().Error()
			} else {
				cause = sourceError.Error()
			}
		} else {
			cause = te.cause.Error()
		}
	}

	return fmt.Sprintf("template %s error %s: %s", stage, location, cause)
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
		serr, ok := typedErr.Cause().(liquid.SourceError)
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
			cause = typedErr.Cause().Error()
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
