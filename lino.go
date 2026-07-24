package lino

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

var tmpl *template.Template

var funcMap template.FuncMap = template.FuncMap{
	"timeNow": time.Now,
}

func AddTemplateGlob(templateGlob string) {
	if tmpl == nil {
		tmpl = template.New("name").Funcs(funcMap)
	}

	var err error
	tmpl, err = tmpl.ParseGlob(templateGlob)
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}

func RenderTemplate(name string, data interface{}, w io.Writer) {
	err := tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}

func RenderMarkdown(data []byte) string {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(data)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return string(markdown.Render(doc, renderer))
}

func RenderMarkdownFile(path string) string {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}

	return RenderMarkdown(data)
}

func ListDirectory(path string) []os.DirEntry {
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
	return files
}

func CreateDirectory(p string) {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		err = os.Mkdir(p, 0755)
		if err != nil {
			fmt.Println("ERROR:", err)
			os.Exit(1)
		}
	}
}

func RunCommand(cmd string, args ...string) {
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}

func RunCommandAsync(cmd string, args ...string) *exec.Cmd {
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Start()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
	return c
}
