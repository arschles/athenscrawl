package main

type catalogRes struct {
	ModsAndVersions []allPathParams `json:"modules"`
	NextPageToken   string          `json:"next,omitempty"`
}

// AllPathParams holds the module and version in the path of a ?go-get=1
// request
type allPathParams struct {
	Module  string `json:"module"`
	Version string `json:"version"`
}
