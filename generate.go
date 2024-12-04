package forge

// Build the shim binaries which are often copied into
// containers and used as their entrypoint, then pack the
// shim binaries so that copy times are as fast as possible.
//go:generate make fmt shims action
