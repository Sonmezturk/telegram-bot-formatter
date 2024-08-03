package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/ledongthuc/pdf"
	"github.com/Sonmezturk/telegram-bot-formatter/structs"
)

func parser(link string) []structs.Order {

	resp, err := http.Get(link)
	if err != nil {
		fmt.Println("Error downloading PDF:", err)
		return nil
	}
	defer resp.Body.Close()

	pdfFile, err := ioutil.TempFile("", "pdf")
	if err != nil {
		fmt.Println("Error creating temporary file:", err)
		return nil
	}
	defer os.Remove(pdfFile.Name())

	_, err = io.Copy(pdfFile, resp.Body)
	if err != nil {
		fmt.Println("Error writing to temporary file:", err)
		return nil
	}
	extracktedText, err := readPdf(pdfFile.Name())
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("output.txt", []byte(extracktedText), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return nil
	}

	text := string(extracktedText)
	orders := strings.Split(text, "Order #")
	orders = orders[1:]
	var result []structs.Order

	for _, order := range orders {
		var items []structs.Item
		orderLines := strings.Split(order, "\n")
		orderID := strings.TrimSpace(strings.Split(orderLines[0], "&")[0])
		skus := strings.Split(order, "SKU: ")[1:]
		quantities := strings.Split(order, "Quantity: ")[1:]
		sizes := strings.Split(order, "Size: ")[1:]
		colors := strings.Split(order, "Color: ")[1:]
		personalizations := strings.Split(order, "Personalization: ")[1:]
		currentPage := strings.Split(strings.Split(order, "Page: ")[1], "\n")[0]
		shipToIndex := strings.Index(order, "Ship to")
		scheduledToShipIndex := strings.Index(order, "Scheduled to ship")
		var shipToText = ""
		if shipToIndex != -1 && scheduledToShipIndex != -1 {
			shipToText = order[shipToIndex+len("Ship to") : scheduledToShipIndex]
			shipToText = strings.ReplaceAll(shipToText, "&", "")
		}
		orderDate := strings.ReplaceAll(strings.Split(strings.Split(order, "Order date")[1], "Payment")[0], "&", "")
		for i, sku := range skus {
			customizations := structs.Customizations{}
			if len(sizes) > i {
				sizeLines := strings.Split(sizes[i], "&Color")
				if len(sizeLines) > 0 {
					customizations.Size = strings.TrimSpace(sizeLines[0])
				}
			}
			if len(personalizations) > i {
				index := strings.Index(personalizations[0], "Disney")
				if index != -1 {
					customizations.Personalization = strings.ReplaceAll(personalizations[0][:index], "&", "")
				} else {
					personalizationLines := strings.Split(personalizations[i], "&")
					if len(personalizationLines) > 0 {
						customizations.Personalization = strings.TrimSpace(personalizationLines[0])
					}
				}
			}
			modifiedSku := ""
			skuLines := strings.Split(sku, "&")
			if len(skuLines) > 0 {
				modifiedSku = strings.TrimSpace(skuLines[0])
			}
			quantity := ""
			if len(quantities) > i {
				quantityLines := strings.Split(quantities[i], "&")
				if len(quantityLines) > 0 {
					quantity = strings.TrimSpace(quantityLines[0])
				}
			}
			color := ""
			if len(colors) > i {
				colorLines := strings.Split(colors[i], "&")
				if len(colorLines) > 0 {
					color = strings.TrimSpace(colorLines[0])
				}
			}
			if color == "Black" {
				modifiedSku += "W"
			}
			if strings.Contains(strings.ToLower(customizations.Personalization), "minnie") ||
				strings.Contains(strings.ToLower(customizations.Personalization), "minie") ||
				strings.Contains(strings.ToLower(customizations.Personalization), "mini") ||
				strings.Contains(strings.ToLower(customizations.Personalization), "minni") {
				modifiedSku += "F"
			}
			item := structs.Item{
				Page:           currentPage,
				Quantity:       quantity,
				SkuName:        modifiedSku,
				Color:          color,
				Customizations: customizations,
			}
			items = append(items, item)
		}
		result = append(result, structs.Order{OrderDate: orderDate, OrderID: orderID, Items: items, ShipTo: shipToText})
	}

	return result

	// resultJSON, _ := json.MarshalIndent(result, "", "  ")
	// ioutil.WriteFile("1.json", resultJSON, 0644)

	// resultJSON, _ = json.Marshal(result)
	// ioutil.WriteFile("2.json", resultJSON, 0644)
}

func readPdf(path string) (string, error) {
	f, r, err := pdf.Open(path)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return "", err
	}
	totalPage := r.NumPage()
	var content strings.Builder
	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		rows, _ := p.GetTextByRow()
		for _, row := range rows {
			for _, word := range row.Content {
				if word.S != "" {
					content.WriteString(word.S)
					content.WriteString(" &")
				}

			}
			content.WriteString("Page: ")
			content.WriteString(strconv.Itoa(pageIndex))
			content.WriteString("\n")
		}
	}
	return content.String(), nil
}
