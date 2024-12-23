// This package do database CRUD operations only

package items

type Item struct {
	Name  string `json:"Name"`
	Count string `json:"count"`
}

func CreateItem(name string, count int) *Item {
	// business logic for create item
	return &Item{}
}

func FindItem(name string) *Item {
	// business logic for find item
	return &Item{}
}
