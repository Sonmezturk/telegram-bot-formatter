package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Sonmezturk/telegram-bot-formatter/structs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func prepareCsv(orders []structs.Order, filename string) {
	// Create a CSV file
	file, err := os.Create(filepath.Join("../", filename))
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)

	// Write the header row
	header := []string{"OrderDate", "OrderID", "Page", "Quantity", "SkuName", "Color", "Size", "Personalization", "ShipTo"}
	err = writer.Write(header)
	if err != nil {
		fmt.Println("Error writing header row:", err)
		return
	}

	// Write data for each order and its items
	for _, order := range orders {
		for _, item := range order.Items {
			row := []string{
				order.OrderDate,
				order.OrderID,
				item.Page,
				item.Quantity,
				item.SkuName,
				item.Color,
				item.Customizations.Size,
				item.Customizations.Personalization,
				order.ShipTo,
			}
			err = writer.Write(row)
			if err != nil {
				fmt.Println("Error writing data row:", err)
				return
			}
		}
	}

	// Flush the writer to ensure all data is written to the file
	writer.Flush()

	fmt.Println("CSV file created successfully!")
}

func prepareCsvForAggregatedItems(aggregatedItems []bson.M, filename string) error {
	// Create a CSV file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating CSV file: %v", err)
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)

	// Write the header row
	header := []string{"skuname", "color", "size", "totalQuantity", "orders", "users", "fileNames"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing header row: %v", err)
	}

	// Write data for each aggregated item
	for _, item := range aggregatedItems {
		// Convert each field that's expected to be a slice of strings
		orders := joinInterfaceSlice(item["orders"])
		users := joinInterfaceSlice(item["users"])
		fileNames := joinInterfaceSlice(item["fileNames"])

		row := []string{
			fmt.Sprint(item["skuname"]),
			fmt.Sprint(item["color"]),
			fmt.Sprint(item["size"]),
			fmt.Sprintf("%v", item["totalQuantity"]),
			orders,
			users,
			fileNames,
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("error writing data row: %v", err)
		}
	}

	// Flush the writer to ensure all data is written to the file
	writer.Flush()

	return nil
}

func joinInterfaceSlice(value interface{}) string {
	var strSlice []string
	if a, ok := value.(primitive.A); ok {
		for _, v := range a {
			strSlice = append(strSlice, fmt.Sprint(v))
		}
	}
	return strings.Join(strSlice, ";")
}