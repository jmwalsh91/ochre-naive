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
	// Must have Poppler installed for pdftoppm
	tempDir, err := os.MkdirTemp("", "pdf_images")
	if err != nil {
		log.Printf("Failed to create temp directory: %s", err)
		return false
	}
	defer os.RemoveAll(tempDir)

	baseName := strings.TrimSuffix(filepath.Base(pdfPath), filepath.Ext(pdfPath))
	imagePattern := filepath.Join(tempDir, baseName+"-%03d.png")
	cmd := exec.Command("pdftoppm", "-png", pdfPath, imagePattern)
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to convert PDF to images: %s", err)
		return false
	}

	images, err := filepath.Glob(filepath.Join(tempDir, baseName+"-*.png"))
	if err != nil {
		log.Printf("Failed to find converted images: %s", err)
		return false
	}

	var aggregatedText strings.Builder

	for _, imagePath := range images {
		client := gosseract.NewClient()
		defer client.Close()
		if err := client.SetImage(imagePath); err != nil {
			log.Printf("Error setting image for OCR: %s", err)
			continue
		}

		text, err := client.Text()
		if err != nil || text == "" {
			log.Printf("Error performing OCR or no text found: %s", err)
			continue
		}
		aggregatedText.WriteString(text + "\n")
	}

	outputFilePath := filepath.Join(outputDir, baseName+".txt")
	if err := os.WriteFile(outputFilePath, []byte(aggregatedText.String()), 0644); err != nil {
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
