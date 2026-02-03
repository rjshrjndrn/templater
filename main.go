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
	"github.com/rjshrjndrn/templater/v6/pkg/helper"
	"github.com/rjshrjndrn/templater/v6/pkg/utils"
)

// processTemplate reads the input file, applies the template with the given values, and outputs to outputPath or stdout.
func processTemplate(inputPath, outputPath string, values map[string]any) error {
	var content []byte
	var err error
	var baseDir string

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
		// For stdin, use current working directory as base
		baseDir, _ = os.Getwd()
	} else {
		content, err = os.ReadFile(inputPath)
		if err != nil {
			return err
		}
		// Get absolute path and extract directory for resolving relative includes
		absPath, err := filepath.Abs(inputPath)
		if err != nil {
			return err
		}
		baseDir = filepath.Dir(absPath)
	}

	fm := sprig.FuncMap()
	fm["toYaml"] = utils.ToYAMLFunc
	fm["tpl"] = utils.MakeTplFunc()
	fm["include"] = utils.MakeIncludeFunc(baseDir, fm["tpl"].(func(string, any) (string, error)))

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

// stringSliceFlag implements flag.Value interface to collect multiple flag values
type stringSliceFlag []string

func (s *stringSliceFlag) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *stringSliceFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func main() {
	var valuesPaths stringSliceFlag

	var inputPath, outputPath string
	var showVersion bool
	flag.StringVar(&inputPath, "i", "", "Path to input file or directory. - for stdin. eg: echo {{ .Values | toJson }} | templater -i - -f values.yaml")
	flag.StringVar(&outputPath, "o", "", "Output directory or file path (optional)")
	flag.Var(&valuesPaths, "f", "Path to values YAML file (optional)")
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

	values, err := helper.ParseYAMLValues(valuesPaths)
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
			return
		}
	}
	if err := processTemplate(inputPath, outputPath, values); err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}
}
