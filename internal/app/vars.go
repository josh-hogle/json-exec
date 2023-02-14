package app

// These variables, while global, are set at compile-time and do not change.
var (
	// Build is the first 8 characters of the git commit hash.
	Build string

	// CommandName is the name of the command to display in help and version info.
	CommandName string

	// IsDevBuild is the boolean equivalent of IsDevelopment since we cannot pass bools in via the -X option to LDFLAGS.
	IsDevBuild bool

	// IsDevelopment is used to flag whether or not this is a development build.
	IsDevelopment string

	// ReleaseDate is the date the binary was released in the form DD MMM YYYY.
	ReleaseDate string

	// Version is the current semver-compatible version of the product.
	Version string
)
