package resp

func SliceToList(names []string) List {
	list := List{}
	for _, name := range names {
		list.Names = append(list.Names, name)
	}
	return list
}

type List struct {
	Names []string `json:"names"`
}
