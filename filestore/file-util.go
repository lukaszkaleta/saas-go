package filestore

import (
	"path/filepath"
	"strings"
)

// SplitFilename splits a file name into name and extension parts.
// Example:
//
//	name, ext := SplitFilename("photo.profile.JPG")
//	// name == "photo.profile", ext == ".JPG"
//
// Notes:
// - ext includes the leading dot (".jpg").
// - If there is no extension, ext will be an empty string and name will equal input.
func SplitFilename(filename string) (name string, ext string) {
	ext = filepath.Ext(filename)
	name = strings.TrimSuffix(filename, ext)
	return name, ext
}
