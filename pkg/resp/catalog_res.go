package resp

import "fmt"

// Catalog represents the JSON response from the /catalog endpoint
type Catalog struct {
	ModsAndVersions []ModuleAndVersion `json:"modules"`
	NextPageToken   string             `json:"next,omitempty"`
}

// ModuleAndVersion holds the module and version in the path of a
// ?go-get=1 request, and a single entry in the response of the /catalog
// endpoint
type ModuleAndVersion struct {
	Module  string `json:"module"`
	Version string `json:"version"`
}

func (m ModuleAndVersion) Error() string {
	return fmt.Sprintf("%s@%s", m.Module, m.Version)
}
