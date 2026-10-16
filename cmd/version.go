// Package cmd provides command-line interface entry points.
package cmd

var version = "dev"

// Version returns the compiled application version.
func Version() string {
	return version
}
