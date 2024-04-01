package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/otiai10/gosseract/v2"
)

func main() {
	inputDir := flag.String("input", "", "Directory that contains the PDF files")
	outputDir := flag.String("output", "output", "Directory to output the text results")
	flag.Parse()

	if *inputDir == "" {
		log.Fatal("Input directory is required")
	}

	files, err := os.ReadDir(*inputDir)
	if err != nil {
		log.Fatalf("Failed to read input directory: %s", err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".pdf") {
			pdfPath := filepath.Join(*inputDir, file.Name())
			fmt.Printf("Processing %s\n", pdfPath)
			if !tryOCR(pdfPath, *outputDir) {
				fmt.Println("OCR not suitable, trying direct text extraction with pdftotext...")
				extractTextWithPdftotext(pdfPath, *outputDir)
			}
		}
	}
}

func tryOCR(pdfPath, outputDir string) bool {
	client := gosseract.NewClient()
	defer client.Close()
	err := client.SetImage(pdfPath)
	if err != nil {
		log.Printf("Error setting image for OCR: %s", err)
		return false
	}

	text, err := client.Text()
	if err != nil || text == "" {
		log.Printf("Error performing OCR or no text found: %s", err)
		return false
	}

	outputFilePath := filepath.Join(outputDir, strings.TrimSuffix(filepath.Base(pdfPath), filepath.Ext(pdfPath))+".txt")
	err = os.WriteFile(outputFilePath, []byte(text), 0644)
	if err != nil {
		log.Printf("Failed to write OCR text to file: %s", err)
		return false
	}

	fmt.Printf("OCR text extraction successful for %s\n", pdfPath)
	return true
}

func extractTextWithPdftotext(pdfPath, outputDir string) {
	outputFilePath := filepath.Join(outputDir, strings.TrimSuffix(filepath.Base(pdfPath), filepath.Ext(pdfPath))+".txt")
	cmd := exec.Command("pdftotext", pdfPath, outputFilePath)
	err := cmd.Run()
	if err != nil {
		log.Printf("Failed to extract text from %s using pdftotext: %s", pdfPath, err)
	} else {
		fmt.Printf("Direct text extraction successful for %s\n", pdfPath)
	}
}
