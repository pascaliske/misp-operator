package utils

import (
	"strings"
)

// extract the last path segment as module name
func ModuleNameFromImage(imageRef string) string {
	// remove any leading and trailing white space first
	ref := strings.TrimSpace(imageRef)

	// drop potential digest - must happen first, because the digest itself contains a colon
	if i := strings.Index(ref, "@"); i >= 0 {
		ref = ref[:i]
	}

	// keep only the last path segment - registry ports live in an earlier segment, they can't be confused with a tag
	if i := strings.LastIndex(ref, "/"); i >= 0 {
		ref = ref[i+1:]
	}

	// finally drop the tag
	if i := strings.Index(ref, ":"); i >= 0 {
		ref = ref[:i]
	}

	return ref
}
