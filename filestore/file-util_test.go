package filestore

import "fmt"

// ExampleSplitFilename demonstrates how to split a file name and extension.
// Run with: go test ./... and see output of examples.
//
// Output:
// photo.profile .JPG
// README
func ExampleSplitFilename() {
	name, ext := SplitFilename("photo.profile.JPG")
	fmt.Println(name, ext)

	name, ext = SplitFilename("README")
	fmt.Println(name, ext)
}
