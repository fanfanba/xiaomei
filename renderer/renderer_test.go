package renderer

import (
	"bytes"
	"path"
	"runtime"
	"testing"
)

type testData struct {
	Title, Content string
}

var r = New(
	path.Join(path.Dir(sourcePath()), `renderer_test`), `layout`, true, nil,
)

func TestRenderer1(t *testing.T) {
	var buf bytes.Buffer
	r.Render(&buf, `t1`, testData{`title`, `content`}, O{})
	got := buf.String()

	expect := `<html>
<head>
<title>t1: title</title>
</head>
<body>
<div>t1: content</div>
</body>
</html>
`
	if got != expect {
		t.Errorf(`expect:
"%s"
got:
"%s"
`, expect, got)
	}
}

func TestRender2(t *testing.T) {
	var buf bytes.Buffer
	r.Render(&buf, `t2`, testData{`title`, `content`}, O{})
	got := buf.String()

	expect := `<html>
<head>
<title>t2: title</title>
</head>
<body>
<div>t2: content</div>
t3
</body>
</html>
`
	if got != expect {
		t.Errorf(`expect:
"%s"
got:
"%s"
`, expect, got)
	}
}

func sourcePath() string {
	_, filename, _, _ := runtime.Caller(0)
	return filename
}
