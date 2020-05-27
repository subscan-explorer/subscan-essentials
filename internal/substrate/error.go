package substrate

type MetadataModuleError struct {
	Module string   `json:"module"`
	Name   string   `json:"name"`
	Doc    []string `json:"doc"`
}
