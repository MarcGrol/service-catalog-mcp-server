package servicecatalog

// ModuleDescriptor is the short version of a Module
type ModuleDescriptor struct {
	ModuleID        string
	Name            string
	Description     string
	ComplexityScore float32 `json:",omitempty"`
}

// ModuleDescriptorList wraps a list into a single object (because the API does not allow lists)
type ModuleDescriptorList struct {
	Modules []ModuleDescriptor `json:"modules"`
}

// InterfaceDescriptor is the short version of an Interface
type InterfaceDescriptor struct {
	InterfaceID     string
	Description     string
	Kind            string
	ComplexityScore int `yaml:",omitempty"`
}

// InterfaceDescriptorList wraps a list into a single object (because the API does not allow lists)
type InterfaceDescriptorList struct {
	Interfaces []InterfaceDescriptor `json:"interfaces"`
}
