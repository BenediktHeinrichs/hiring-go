package main

import (
	_ "embed"
	"fmt"

	nlp "github.com/fluhus/gostuff/nlp"
	"github.com/ledongthuc/pdf"
)

func main() {
	content, err := readPdf("CV.pdf") // Read local pdf file
	if err != nil {
		panic(err)
	}
	fmt.Println(content)

	tokenized := nlp.Tokenize(content, false)
	matrix := make([][]string, 1)
	matrix[0] = tokenized
	topics, ids := nlp.Lda(matrix, 10)
	fmt.Println(topics)
	fmt.Println(ids)

	// Group topics by non-zero column indices
	topicGroups := make(map[string][]string)

	for topic, values := range topics {
		// Iterate over values and find non-zero columns
		for i, val := range values {
			if val > 0 {
				columnKey := fmt.Sprintf("Column %d", i+1)
				topicGroups[columnKey] = append(topicGroups[columnKey], topic)
			}
		}
	}

	// Print topics grouped by non-zero column indices
	for column, topicList := range topicGroups {
		fmt.Println("Topics with non-zero values in", column)
		for _, topic := range topicList {
			fmt.Println("- Topic:", topic)
		}
		fmt.Println()
	}
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

	var content string

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		rows, _ := p.GetTextByRow()
		for _, row := range rows {
			for _, word := range row.Content {
				content += word.S + " "
			}
			content += "\n"
		}
		content += "\n"
	}
	return content, nil
}
