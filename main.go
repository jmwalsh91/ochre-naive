package main

import (
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
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

	log.Info("Starting PDF text extraction process...")

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".pdf") {
			pdfPath := filepath.Join(*inputDir, file.Name())
			log.Infof("Processing %s", pdfPath)

			if !extractTextWithPdftotext(pdfPath, *outputDir) {
				log.Warn("Direct text extraction failed, trying OCR...")
				tryOCR(pdfPath, *outputDir)
			}
		}
	}

	log.Info("PDF text extraction process completed.")
}

func extractTextWithPdftotext(pdfPath, outputDir string) bool {
	outputFilePath := filepath.Join(outputDir, strings.TrimSuffix(filepath.Base(pdfPath), filepath.Ext(pdfPath))+".txt")
	cmd := exec.Command("pdftotext", pdfPath, outputFilePath)
	err := cmd.Run()
	if err != nil {
		log.Errorf("Failed to extract text from %s using pdftotext: %s", pdfPath, err)
		return false
	}

	log.Infof("Direct text extraction successful for %s", pdfPath)
	return true
}

func tryOCR(pdfPath, outputDir string) bool {
	tempDir, err := os.MkdirTemp("", "pdf_images")
	if err != nil {
		log.Errorf("Failed to create temp directory: %s", err)
		return false
	}
	defer os.RemoveAll(tempDir)

	baseName := strings.TrimSuffix(filepath.Base(pdfPath), filepath.Ext(pdfPath))
	imagePattern := filepath.Join(tempDir, baseName+"-%03d.png")
	cmd := exec.Command("pdftoppm", "-png", pdfPath, imagePattern)
	if err := cmd.Run(); err != nil {
		log.Errorf("Failed to convert PDF to images: %s", err)
		return false
	}

	images, err := filepath.Glob(filepath.Join(tempDir, baseName+"-*.png"))
	if err != nil {
		log.Errorf("Failed to find converted images: %s", err)
		return false
	}

	var aggregatedText strings.Builder
	for _, imagePath := range images {
		client := gosseract.NewClient()
		defer client.Close()
		if err := client.SetImage(imagePath); err != nil {
			log.Errorf("Error setting image for OCR: %s", err)
			continue
		}
		text, err := client.Text()
		if err != nil || text == "" {
			log.Errorf("Error performing OCR or no text found: %s", err)
			continue
		}
		aggregatedText.WriteString(text + "\n")
	}

	outputFilePath := filepath.Join(outputDir, baseName+".txt")
	if err := os.WriteFile(outputFilePath, []byte(aggregatedText.String()), 0644); err != nil {
		log.Errorf("Failed to write OCR text to file: %s", err)
		return false
	}

	log.Infof("OCR text extraction successful for %s", pdfPath)
	return true
}
