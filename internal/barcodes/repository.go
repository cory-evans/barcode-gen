package barcodes

import (
	"encoding/json"
	"log"
	"os"
	"slices"
	"sync"

	"github.com/cory-evans/barcode-gen/internal/models"
)

var FilePath = "items.json"
var lock sync.Mutex

func loadItemsFromDB() ([]models.GenerateInput, error) {
	lock.Lock()
	defer lock.Unlock()

	f, err := os.OpenFile(FilePath, os.O_RDONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.GenerateInput{}, nil
		}

		return nil, err
	}

	defer f.Close()

	var items []models.GenerateInput
	err = json.NewDecoder(f).Decode(&items)

	if err != nil {
		log.Println("Error reading items from file:", err)
		return nil, err
	}

	slices.SortStableFunc(items, func(a, b models.GenerateInput) int {
		if a.SaveName < b.SaveName {
			return -1
		} else if a.SaveName > b.SaveName {
			return 1
		}

		return 0
	})

	return items, nil
}

func saveItemsToDB(items []models.GenerateInput) error {
	lock.Lock()
	defer lock.Unlock()

	f, err := os.OpenFile(FilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Println("Error opening file for writing:", err)
		return err
	}

	defer f.Close()

	err = json.NewEncoder(f).Encode(items)
	if err != nil {
		log.Println("Error writing items to file:", err)
	}

	return err
}

func SaveItem(item models.GenerateInput) {
	// filter out the item if it already exists
	if item.SaveName == "" {
		return
	}

	items, err := loadItemsFromDB()
	if err != nil {
		log.Println("Error loading items from DB:", err)
		return
	}

	for i, v := range items {
		if v.SaveName == item.SaveName {
			items = append(items[:i], items[i+1:]...)
			break
		}
	}

	log.Println("Saving item", item.SaveName)

	items = append(items, item)

	err = saveItemsToDB(items)
	if err != nil {
		log.Println("Error saving items to DB:", err)
	}
}

func GetItems() []models.GenerateInput {
	items, err := loadItemsFromDB()
	if err != nil {
		log.Println("Error loading items from DB:", err)
		return []models.GenerateInput{}
	}

	return items
}

func GetItem(name string) *models.GenerateInput {
	items, err := loadItemsFromDB()
	if err != nil {
		log.Println("Error loading items from DB:", err)
		return nil
	}

	for _, v := range items {
		if v.SaveName == name {
			return &v
		}
	}

	return nil
}

func DeleteItem(name string) {
	items, err := loadItemsFromDB()
	if err != nil {
		log.Println("Error loading items from DB:", err)
		return
	}

	for i, v := range items {
		if v.SaveName == name {
			log.Println("Deleting item", name)
			items = append(items[:i], items[i+1:]...)
			break
		}
	}

	err = saveItemsToDB(items)
	if err != nil {
		log.Println("Error saving items to DB:", err)
	}

}
