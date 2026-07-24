package main

import (
	"flag"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"
)

var funcMap template.FuncMap = template.FuncMap{
	"timeNow": time.Now,
}

func outputFiles(inputPath string, outputPath string, tmpl *template.Template) error {
	pagesPath := filepath.Join(inputPath, "pages")
	f, err := os.Open(pagesPath)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		err = os.Mkdir(outputPath, 0755)
		if err != nil {
			return err
		}
	}

	fileInfos, err := f.Readdir(0)
	if err != nil {
		return err
	}

	for _, info := range fileInfos {
		p := filepath.Join(outputPath, info.Name())
		f, err := os.Create(p)
		if err != nil {
			return err
		}
		defer f.Close()

		data := struct {
			FileName string
		}{
			FileName: info.Name(),
		}

		err = tmpl.ExecuteTemplate(f, info.Name(), data)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	inputPath := flag.String("input", ".", "templates input path")
	outputPath := flag.String("output", "./_site", "output path")
	flag.Parse()

	templateGlob := filepath.Join(*inputPath, "*/*.html")
	tmpl, err := template.New("name").Funcs(funcMap).ParseGlob(templateGlob)
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}

	err = outputFiles(*inputPath, *outputPath, tmpl)
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
