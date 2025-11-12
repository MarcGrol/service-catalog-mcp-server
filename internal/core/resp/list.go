package resp

// SliceToList is used to wrap a string-slice into a single object
func SliceToList(names []string) List {
	list := List{}
	for _, name := range names {
		list.Names = append(list.Names, name)
	}
	return list
}

// List wraps a string-slice into a single object
type List struct {
	Names []string `json:"names"`
}
