package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	nlp "github.com/fluhus/gostuff/nlp"

	// Usage of the version of unipdf, before it became proprietary and required a license
	"github.com/oliverpool/unipdf/v3/extractor"
	pdf "github.com/oliverpool/unipdf/v3/model"
)

// Struct to hold the response
type Response struct {
	Topics map[string][]string `json:"topics"`
}

// Main function to start the server
func main() {
	http.HandleFunc("/process-pdf", corsMiddleware(processPDFHandler))
	http.Handle("/", http.FileServer(http.Dir("./static")))
	fmt.Println("App is listening on port 8080...")

	http.ListenAndServe(":8080", nil)
}

// CORS Middleware to handle CORS headers
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight request
		if r.Method == http.MethodOptions {
			return
		}

		// Call the next handler
		next(w, r)
	}
}

// Handler for processing the PDF
func processPDFHandler(w http.ResponseWriter, r *http.Request) {
	// Check the HTTP method
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the uploaded file
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to parse file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read the PDF content
	content, err := outputPdfText(file)
	if err != nil {
		http.Error(w, "Failed to read PDF", http.StatusInternalServerError)
		return
	}

	// Tokenize and process content
	tokenized := nlp.Tokenize(content, false)
	matrix := make([][]string, 1)
	matrix[0] = tokenized
	topics, _ := nlp.Lda(matrix, 10)

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

	// Create the response struct
	response := Response{Topics: topicGroups}

	// Convert response to JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// outputPdfText prints out contents of PDF file to stdout.
func outputPdfText(reader io.Reader) (string, error) {
	content := ""

	// Read the content into a byte slice
	contentBytes, err := io.ReadAll(reader)
	if err != nil {
		return content, err
	}

	pdfBytesReader := bytes.NewReader(contentBytes)
	pdfReader, err := pdf.NewPdfReader(pdfBytesReader)
	if err != nil {
		return content, err
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return content, err
	}

	for i := 0; i < numPages; i++ {
		pageNum := i + 1

		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return content, err
		}

		ex, err := extractor.New(page)
		if err != nil {
			return content, err
		}

		text, err := ex.ExtractText()
		if err != nil {
			return content, err
		}

		content += text + "\n"
	}

	// Enable for debugging
	// fmt.Println(content)

	return content, nil
}
