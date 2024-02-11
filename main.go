package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Masterminds/sprig/v3"
)

func processTemplate(inputPath, outputPath string) error {
	// Read the input file
	content, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return err
	}

	// Create a new template with Sprig functions
	tpl, err := template.New(filepath.Base(inputPath)).Funcs(sprig.HtmlFuncMap()).Parse(string(content))
	if err != nil {
		return err
	}

	// If outputPath is not specified, use stdout
	var output *os.File
	if outputPath == "" {
		output = os.Stdout
	} else {
		// Ensure output directory exists
		if err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm); err != nil {
			return err
		}
		// Create the output file
		output, err = os.Create(outputPath)
		if err != nil {
			return err
		}
		defer output.Close()
	}

	// Execute the template
	if err := tpl.Execute(output, nil); err != nil {
		return err
	}

	return nil
}

func main() {
	var inputPath, outputPath string
	flag.StringVar(&inputPath, "i", "", "Path to input file or directory")
	flag.StringVar(&outputPath, "o", "out", "Output directory (optional, default: out)")
	flag.Parse()

	if inputPath == "" {
		fmt.Println("Input path is required.")
		os.Exit(1)
	}

	info, err := os.Stat(inputPath)
	if err != nil {
		fmt.Printf("Error accessing input path: %v\n", err)
		os.Exit(1)
	}

	if info.IsDir() {
		// Process each file in the directory
		err := filepath.Walk(inputPath, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				relPath, err := filepath.Rel(inputPath, path)
				if err != nil {
					return err
				}
				outputFilePath := filepath.Join(outputPath, relPath)
				return processTemplate(path, outputFilePath)
			}
			return nil
		})
		if err != nil {
			fmt.Printf("Error processing directory: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Process a single file
		if outputPath != "" {
			outputPath = filepath.Join(outputPath, filepath.Base(inputPath))
		}
		if err := processTemplate(inputPath, outputPath); err != nil {
			fmt.Printf("Error processing file: %v\n", err)
			os.Exit(1)
		}
	}
}
