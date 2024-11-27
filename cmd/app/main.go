package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/rjshrjndrn/templater/pkg/utils"
	"gopkg.in/yaml.v3"
)

// parseYAMLValues reads a YAML file and returns the parsed data.
func parseYAMLValues(valuesPath string) (map[string]interface{}, error) {
	var values map[string]interface{}

	if valuesPath == "" {
		return values, nil // No values file provided
	}

	data, err := os.ReadFile(valuesPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &values)
	if err != nil {
		return nil, err
	}

	// Support helm structure
	// {{ .key }} -> {{ .Values.key }}
	return map[string]interface{}{"Values": values}, nil
}

// processTemplate reads the input file, applies the template with the given values, and outputs to outputPath or stdout.
func processTemplate(inputPath, outputPath string, values map[string]interface{}) error {
	var content []byte
	var err error
	if inputPath == "-" {
		// Read from stdin
		scanner := bufio.NewScanner(os.Stdin)
		// Scan line by line or token by token
		for scanner.Scan() {
			content = append(content, scanner.Bytes()...) // Append current line's bytes
			content = append(content, '\n')               // Add newline to preserve input structure
		}
		content = append(content, scanner.Bytes()...)
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
			return err
		}
		if outputPath != "" {
			fmt.Println("Input is from stdin. Output path will be stdout. If you need to save to a file, >file.yaml .")
			outputPath = ""
		}
	} else {
		content, err = os.ReadFile(inputPath)
		if err != nil {
			return err
		}
	}

	fm := sprig.FuncMap()
	fm["toYaml"] = utils.ToYAMLFunc

	tpl, err := template.New(filepath.Base(inputPath)).Funcs(fm).Parse(string(content))
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Output to os.Stdout by default
	output := os.Stdout
	if outputPath != "" { // Create file if outputPath is provided
		if err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm); err != nil {
			return err
		}
		file, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		defer file.Close()
		output = file // Use the file as output
	}

	if err := tpl.Execute(output, values); err != nil {
		return err
	}

	return nil
}

// Default version for dev builds
var appVersion = "dev"

func main() {
	var inputPath, outputPath, valuesPath string
	var showVersion bool
	flag.StringVar(&inputPath, "i", "", "Path to input file or directory")
	flag.StringVar(&outputPath, "o", "", "Output directory or file path (optional)")
	flag.StringVar(&valuesPath, "f", "", "Path to values YAML file (optional)")
	flag.BoolVar(&showVersion, "v", false, "Prints the version of the app and exits")
	flag.Parse()

	if showVersion {
		fmt.Println("App Version:", appVersion)
		return
	}

	if inputPath == "" {
		fmt.Println("Input path is required.")
		os.Exit(1)
	}

	values, err := parseYAMLValues(valuesPath)
	if err != nil {
		fmt.Printf("Error parsing values file: %v\n", err)
		os.Exit(1)
	}

	// input is stdin
	if inputPath != "-" {

		info, err := os.Stat(inputPath)
		if err != nil {
			fmt.Printf("Error accessing input path: %v\n", err)
			os.Exit(1)
		}

		if info.IsDir() {
			err := filepath.Walk(inputPath, func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					relPath, err := filepath.Rel(inputPath, path)
					if err != nil {
						return err
					}
					var outputFilePath string
					if outputPath != "" {
						outputFilePath = filepath.Join(outputPath, relPath)
					}
					return processTemplate(path, outputFilePath, values)
				}
				return nil
			})
			if err != nil {
				fmt.Printf("Error processing directory: %v\n", err)
				os.Exit(1)
			}
		}
	}
	if err := processTemplate(inputPath, outputPath, values); err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}
}
