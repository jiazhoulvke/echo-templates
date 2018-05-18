package templates

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	//"strings"

	"github.com/flosch/pongo2"
	"github.com/labstack/echo"
)

var (
	//Templates 预编译模板
	Templates *templates
	//Debug 调试模式
	Debug bool
	//tplPath 模板路径
	tplPath string
	//Exclude 排除文件或目录
	Exclude []string
)

func init() {
	Templates = &templates{
		tplMap: make(map[string]*pongo2.Template),
	}
}

type templates struct {
	tplMap map[string]*pongo2.Template
}

//Compile 编译模板
func Compile(templatePath string) error {
	os.Chdir(templatePath)
	tplPath = templatePath
	excludePaths := make([]string, 0, 8)
	for _, p := range Exclude {
		paths, err := filepath.Glob(p)
		if err != nil {
			return err
		}
		excludePaths = append(excludePaths, paths...)
	}
	return filepath.Walk(tplPath, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			panic(tplPath)
			// return fmt.Errorf("path error")
		}
		relPath, _ := filepath.Rel(tplPath, path)
		for _, p := range excludePaths {
			if p == relPath {
				if Debug {
					fmt.Println("[templates] skip path:", p)
				}
				return filepath.SkipDir
			}
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".html" {
			return nil
		}
		t, err := pongo2.FromFile(path)
		if err != nil {
			return err
		}
		key := getKeyName(tplPath, path)
		Templates.Set(key, t)
		return nil
	})
}

func (t *templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	var ctx pongo2.Context
	switch v := data.(type) {
	case map[string]interface{}:
		ctx = pongo2.Context(v)
	case pongo2.Context:
		ctx = v
	}
	if Debug {
		tpl, err := pongo2.FromFile(filepath.Join(tplPath, name))
		if err != nil {
			return err
		}
		return tpl.ExecuteWriter(ctx, w)
	}
	if tpl, exists := t.tplMap[name]; exists {
		return tpl.ExecuteWriter(ctx, w)
	}
	return fmt.Errorf("template not found:%s", name)
}

func (t *templates) Set(name string, tpl *pongo2.Template) {
	t.tplMap[name] = tpl
}

func (t *templates) Get(name string) *pongo2.Template {
	return t.tplMap[name]
}

func getKeyName(d string, fPath string) string {
	relPath, err := filepath.Rel(d, fPath)
	if err != nil {
		return ""
	}
	return relPath
}
