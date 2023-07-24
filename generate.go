package forge

// Build the shim binary which is often copied into containers
// and used as their entrypoint, then pack the shim binary
// so that copy times are as fast as possible.
//go:generate make fmt lint shim action
