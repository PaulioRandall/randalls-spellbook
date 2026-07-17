package projectdata

// PathType determines the type of path in a SystemFile.
//
// Yes, 2026 and Windows file pathing is still a giant
// pain in the humble programmer's ass. I want to make it
// as easy as possible to move project files between
// different systems, hence I'm keeping track of the path
// type. A little bit more cumbersome but I'm just not
// feeling in a regex mood.
//
// PathType is used to determine how to access the media on
// the local file system. It also means migrating the
// project between systems will probably be a headache that
// we should accept and deal with rather than pushing it on
// to the user.
type PathType string

const (
	// PathTypeWindows means the value of the associated file
	// path field should work with standard Windows file
	// systems.
	//
	// It may be absolute or relative. Regardless, it should
	// be safe to throw straight into file APIs, such as 'os'
	// and 'io', when used with Windows.
	PathTypeOS PathType = "OS"

	// PathTypePOSIX means the value of the associated file
	// path field should work with POSIX file systems.
	//
	// It may be absolute or relative. Regardless, it should
	// be safe to throw straight into file APIs, such as 'os'
	// and 'io', when used on Unix-like platforms such as
	// Linux and iOS.
	PathTypePOSIX PathType = "POSIX"
)

// SystemFile holds a path to file on the file system and
// type (format) of that path.
//
// The path type should be checked for compatibility with
// the local file system before use. It may need
// converting, and possible corrective intervention by the
// user, if the project and its assests have been moved
// to a new machine with a different OS, e.g. moved
// between Windows to a Linux systems.
//
// It's called SystemFile to avoid naming confusion with
// standard libraries.
type SystemFile struct {
	// Type is the type of the path, typically Windows or
	// POSIX.
	Type PathType

	// Path is the absolute or relative file path to the
	// file including extension.
	//
	// Relative paths must always be relative to the root
	// project directory.
	Path string
}
