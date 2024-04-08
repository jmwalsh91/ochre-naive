# Ochre

Ochre is an application for extracting text from PDFs. First, Ochre will attempt OCR with Gosseract. If this fails, Ochre will attempt to extract the text with Poppler utils' pdftotext. Ochre can be run in the CLI with 
```ochre-naive --input *input directory* --output *output directory*```

## Prerequisites

Before you begin, ensure you have met the following requirements:

- Go (version 1.15 or higher recommended)
- Poppler utilities for pdftoppm and pdftotext
- Tesseract-OCR (for gosseract)

## Installation

Follow these steps to get your development environment running:

### Install Go

Download and install Go by following the instructions on the official Go website.

### Install Poppler Utilities

#### Linux (Debian/Ubuntu)

sudo apt-get update
sudo apt-get install poppler-utils

#### macOS

Using Homebrew:

brew install poppler

#### Windows

Windows users can download pre-compiled binaries of Poppler from here. Extract the files and add the directory containing pdftoppm and pdftotext to your system's PATH.

### Install Tesseract-OCR

Tesseract-OCR is required by gosseract. Installation instructions for various platforms can be found in the Tesseract GitHub repository.

### Clone the Repository

Clone this repository to your local machine

### Install Go Dependencies

Run the following command inside the project directory:

go get -d ./...

## Running the Program

To run the program, use the following command:

if not build, run go build ./ at the root of this project
then:
ochre-naive -input "/path/to/input/directory" -output "/path/to/output/directory"

alternatively:
go run main.go -input "/path/to/input/directory" -output "/path/to/output/directory"

Replace "/path/to/input/directory" and "/path/to/output/directory" with the actual paths to your input and output directories, respectively.
