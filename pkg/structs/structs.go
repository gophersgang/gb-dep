package structs

// extracted from packagefile prevent the diamond import problem

// PackageFile represent a packagefile
type PackageFile struct {
	Packages []Pkg `json:"packages"`
}

// Pkg represents a Golang package
type Pkg struct {
	Name      string   `json:"name"`                // name of the package
	Group     []string `json:"group,omitempty"`     // in which groups should this package be installed?
	Goos      []string `json:"goos,omitempty"`      // what Go OS is supported
	Insecure  bool     `json:"insecure,omitempty"`  // use insecure protocol for downloading
	Recursive bool     `json:"recursive,omitempty"` // should we also fetch the vendored packages for a package?
	Command   string   `json:"command,omitempty"`   // special command to executed for downloading the packag
	Private   bool     `json:"private,omitempty"`   // is this package private?
	Skipdep   bool     `json:"skipdep,omitempty"`   // shall we ignore dependencies?
	Target    string   `json:"target,omitempty"`    // folder to install package to
	// current status, commit > tag > branch
	Branch       string `json:"branch,omitempty"`
	Tag          string `json:"tag,omitempty"`
	Commit       string `json:"commit,omitempty"`
	LockedCommit string `json:"locked_commit,omitempty"`
	VcsType      string `json:"vcs_type"`
}
