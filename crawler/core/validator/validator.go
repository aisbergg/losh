package validator

import (
	"strings"

	"losh/internal/core/product/models"

	"github.com/gookit/validate"
	"golang.org/x/text/language"
)

var knownOKHV = [2]string{
	"OKHv1.0",
	"OKH-LOSHv1.0",
}

type ValidationError struct {
	messages []string
}

// Error returns the error message.
func (e *ValidationError) Error() string {
	return strings.Join(e.messages, "; ")
}

// Add adds an error message.
func (e *ValidationError) Add(msg string) {
	e.messages = append(e.messages, msg)
}

type Validator struct {
	licenses *license.LicenseCache
}

// NewValidator creates a new validator.
func NewValidator(licenses *license.LicenseCache) *Validator {
	return &Validator{
		licenses: licenses,
	}
}

// ValidateMandatory checks if the mandatory fields are set correctly.
func (v *Validator) ValidateMandatory(mdtFlds MandatoryFields) *ValidationError {
	vldErr := &ValidationError{}

	if !isValidOKHV(mdtFlds.OKHV) {
		vldErr.Add("invalid OKH version")
	}
	if strings.TrimSpace(mdtFlds.Name) == "" {
		vldErr.Add("invalid name")
	}
	if strings.TrimSpace(mdtFlds.Description) == "" {
		vldErr.Add("invalid description")
	}
	if strings.TrimSpace(mdtFlds.Version) == "" {
		vldErr.Add("invalid version")
	}
	if !validate.IsFullURL(mdtFlds.Repository) {
		vldErr.Add("invalid repository")
	}
	if v.licenses.GetByIDOrName(mdtFlds.License) == nil {
		vldErr.Add("invalid license")
	}
	if strings.TrimSpace(mdtFlds.Licensor) == "" {
		vldErr.Add("invalid licensor")
	}
	if !isValidLanguageTag(mdtFlds.DocumentationLanguage) {
		vldErr.Add("invalid documentation language")
	}

	if vldErr.messages != nil {
		return vldErr
	}
	return nil
}

// ValidateProduct checks if the product information conforms the LOSH
// specification.
func (v *Validator) ValidateProduct(product *models.Product) *ValidationError {
	// TODO: implement validation of product
	vldErr := &ValidationError{}

	if vldErr.messages != nil {
		return vldErr
	}
	return nil
}

type MandatoryFields struct {
	// The name of the component.
	OKHV string `json:"okhv"`
	// The name of the component.
	Name string `json:"name"`
	// The short description of the component.
	Description string `json:"description"`
	// The version string of the release.
	Version string `json:"version"`
	// The repository that this component is developed in.
	Repository string `json:"repository"`
	// The license used for the component.
	License string `json:"license"`
	// The license holder of the component.
	Licensor string `json:"licensor"`
	// The language in which the documentation is written.
	DocumentationLanguage string `json:"documentationLanguage"`
}

func isValidOKHV(okhv string) bool {
	okhv = strings.TrimSpace(okhv)
	if okhv == "" {
		return false
	}
	for _, v := range knownOKHV {
		if v == okhv {
			return true
		}
	}
	return false
}

func isValidLanguageTag(tag string) bool {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return false
	}
	_, err := language.Parse(tag)
	return err == nil
}
