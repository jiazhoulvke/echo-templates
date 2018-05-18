package templates

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/flosch/pongo2"
	//"github.com/labstack/echo"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTemplates(t *testing.T) {
	Convey("测试模板", t, func() {
		tempDir, err := ioutil.TempDir("", "tqtms_test_")
		So(err, ShouldBeNil)
		So(tempDir, ShouldNotEqual, "")
		err = ioutil.WriteFile(filepath.Join(tempDir, "test.html"), []byte("hello,{{ name }}!"), 0666)
		So(err, ShouldBeNil)
		Convey("编译模板", func() {
			err := Compile(tempDir)
			So(err, ShouldBeNil)
		})
		Convey("执行模板", func() {
			buf := bytes.NewBuffer([]byte(""))
			err := Templates.Render(buf, "test.html", pongo2.Context{"name": "world"}, nil)
			So(err, ShouldBeNil)
			So(buf.String(), ShouldEqual, "hello,world!")
		})
		err = exec.Command("rm", "-fr", tempDir).Run()
		So(err, ShouldBeNil)
	})
}
