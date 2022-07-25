package util

import "regexp"

// var _manifest_name_pattern = r"^okh([_\-\t ].+)*$"
var manifestNamePattern = regexp.MustCompile(`^okh([_\-\t ].+)*$`)

// IsAcceptedManifestFileName returns true if the given file name matches an accepted manifest name.
func IsAcceptedManifestFileName(path string) bool {
	return manifestNamePattern.MatchString(path)
}
