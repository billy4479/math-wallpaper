package main

import (
	"bufio"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

//go:embed preamble.tex
var preamble string

//go:embed ending.tex
var ending string

func extract() {
	result := make(map[string][]string)

	err := filepath.WalkDir("../latex", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".tex" {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		f := bufio.NewReader(file)
		filename := d.Name()[:len(d.Name())-4]

		inEquation := false
		currentEq := ""
		for {
			line, err := f.ReadString('\n')
			if errors.Is(err, io.EOF) {
				break
			}

			if err != nil {
				return err
			}

			trimmed := strings.TrimSpace(line)

			if trimmed == "\\begin{equation}" {
				inEquation = true
				continue
			}

			if trimmed == "\\end{equation}" {
				inEquation = false
				if result[filename] == nil {
					result[filename] = []string{currentEq}
				} else {
					result[filename] = append(result[filename], currentEq)
				}
				currentEq = ""
				continue
			}

			if inEquation {
				currentEq += trimmed
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = os.MkdirAll("result", 0755)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for filename, eqs := range result {
		for i, eq := range eqs {
			outFile, err := os.Create(fmt.Sprintf("result/%s-%03d.tex", filename, i))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			defer outFile.Close()

			_, err = outFile.WriteString(preamble)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			_, err = outFile.WriteString(eq)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			_, err = outFile.WriteString(ending)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}
}
