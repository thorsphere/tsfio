//go:build !windows

package tsfio

// Linux blocked Directories and Filenames.
// If a directory or their parents match invalDir,
// tsfio functions will return an error. If a Filename
// matches invalFile, tsfio functions will return an error.
var (
	// blocked directories
	invalDir = []Directory{
		"/boot",
		"/dev",
		"/lost+found",
		"/proc",
	}
	// blocked filenames
	invalFile = []Filename{
		"/",
		"/bin",
		"/etc",
		"/home",
		"/lib",
		"/media",
		"/mnt",
		"/opt",
		"/root",
		"/sbin",
		"/srv",
		"/tmp",
		"/usr",
		"/var",
	}
)

// InvalDir returns the slice of blocked directories. If a directory or their parents match InvalDir, tsfio
// functions will return an error.
func InvalDir() []Directory {
	return append([]Directory(nil), invalDir...)
}

// InvalFile returns the slice of blocked filenames. If a Filename matches InvalFile, tsfio
// functions will return an error.
func InvalFile() []Filename {
	return append([]Filename(nil), invalFile...)
}
