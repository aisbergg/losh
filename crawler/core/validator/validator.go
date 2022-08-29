package validator

import (
	"strings"

	"losh/internal/core/product/models"
	"losh/internal/core/product/services"

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

func newValidationError() *ValidationError {
	return &ValidationError{
		messages: []string{},
	}
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
	productService *services.Service
}

// NewValidator creates a new validator.
func NewValidator(productService *services.Service) *Validator {
	return &Validator{
		productService: productService,
	}
}

// ValidateMandatory checks if the mandatory fields are set correctly.
func (v *Validator) ValidateMandatory(mdtFlds MandatoryFields) error {
	vldErr := newValidationError()

	if !isValidOKHV(mdtFlds.OKHV) {
		vldErr.Add("invalid OKH version")
	}
	if strings.TrimSpace(mdtFlds.Name) == "" {
		vldErr.Add("missing name")
	}
	if strings.TrimSpace(mdtFlds.Description) == "" {
		vldErr.Add("missing description")
	}
	if strings.TrimSpace(mdtFlds.Version) == "" {
		vldErr.Add("missing version")
	}
	if !validate.IsFullURL(mdtFlds.Repository) {
		vldErr.Add("invalid repository")
	}
	if mdtFlds.License == "" {
		// allow empty license
		// vldErr.Add("missing license")
	} else if v.productService.GetCachedLicenseByIDOrName(mdtFlds.License) == nil {
		vldErr.Add("invalid license")
	}
	if strings.TrimSpace(mdtFlds.Licensor) == "" {
		vldErr.Add("missing licensor")
	}
	if !isValidLanguageTag(mdtFlds.DocumentationLanguage) {
		vldErr.Add("invalid documentation language")
	}
	if !containsReadme(mdtFlds.FilePaths) {
		vldErr.Add("missing readme")
	}

	if len(vldErr.messages) != 0 {
		return vldErr
	}
	return nil
}

// ValidateProduct checks if the product information conforms the LOSH
// specification.
func (v *Validator) ValidateProduct(product *models.Product) error {
	// TODO: implement validation of product
	vldErr := newValidationError()

	if len(vldErr.messages) != 0 {
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
	// Referenced files of the component.
	FilePaths []string `json:"filePaths"`
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

func containsReadme(filePaths []string) bool {
	if filePaths == nil {
		return false
	}
	for _, fp := range filePaths {
		if pos := strings.LastIndexByte(fp, '.'); pos != -1 {
			fp = fp[:pos]
		}
		fp = strings.TrimSpace(fp)
		fp = strings.Replace(fp, " ", "", -1)
		fp = strings.Replace(fp, "-", "", -1)
		fp = strings.Replace(fp, "_", "", -1)
		fp = strings.ToUpper(fp)
		if fp == "README" {
			return true
		}
	}
	return false
}
