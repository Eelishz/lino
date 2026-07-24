package lino

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

var tmpl *template.Template

var funcMap template.FuncMap = template.FuncMap{
	"timeNow": time.Now,
	"html": func(s any) template.HTML {
		switch v := s.(type) {
		case string:
			return template.HTML(v)
		case []byte:
			return template.HTML(v)
		default:
			return template.HTML(fmt.Sprint(v))
		}
	},
}

func reportError(err error) {
	_, file, lineNo, _ := runtime.Caller(2)
	fmt.Printf("%s:%d: %s\n", file, lineNo, err)
	os.Exit(1)
}

func AddTemplateGlob(templateGlob string) {
	if tmpl == nil {
		tmpl = template.New("name").Funcs(funcMap)
	}

	var err error
	tmpl, err = tmpl.ParseGlob(templateGlob)
	if err != nil {
		reportError(err)
	}
}

func RenderTemplate(name string, data interface{}, w io.Writer) {
	err := tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		reportError(err)
	}
}

func RenderTemplateFile(name string, data interface{}, path string) {
	f, err := os.Create(path)
	if err != nil {
		reportError(err)
	}
	defer f.Close()
	RenderTemplate(name, data, f)
}

func RenderMarkdown(data []byte) []byte {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(data)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func RenderMarkdownFile(path string) []byte {
	f, err := os.Open(path)
	if err != nil {
		reportError(err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		reportError(err)
	}

	return RenderMarkdown(data)
}

func ListDirectory(path string) []os.DirEntry {
	files, err := os.ReadDir(path)
	if err != nil {
		reportError(err)
	}
	return files
}

func CreateDirectory(p string) {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		err = os.Mkdir(p, 0755)
		if err != nil {
			reportError(err)
		}
	}
}

func RemoveDirectoryIfExists(p string) {
	if _, err := os.Stat(p); err == nil {
		err = os.RemoveAll(p)
		if err != nil {
			reportError(err)
		}
	}
}

func RunCommand(cmd string, args ...string) {
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	if err != nil {
		reportError(err)
	}
}

func RunCommandAsync(cmd string, args ...string) *exec.Cmd {
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Start()
	if err != nil {
		reportError(err)
	}
	return c
}
