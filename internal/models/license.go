package models

import "strings"

// AsLicenseType returns a license type from a string.
func AsLicenseType(s string) LicenseType {
	s = strings.TrimSpace(strings.ToUpper(s))
	switch s {
	case "STRONG":
		return LicenseTypeStrong
	case "WEAK":
		return LicenseTypeWeak
	case "PERMISSIVE":
		return LicenseTypePermissive
	default:
		return LicenseTypeUnknown
	}
}
