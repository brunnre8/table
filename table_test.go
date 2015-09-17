package table

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatter(t *testing.T) {
	var formatter Formatter

	fn := func(a string, b ...interface{}) string { return "" }
	f := Formatter(fn)

	assert.IsType(t, formatter, f)
}

func TestTable_New(t *testing.T) {
	buf := bytes.Buffer{}
	New("foo", "bar").WithWriter(&buf).Print()
	out := buf.String()

	assert.Contains(t, out, "foo")
	assert.Contains(t, out, "bar")

	buf.Reset()
	New().WithWriter(&buf).Print()
	out = buf.String()

	assert.Empty(t, strings.TrimSpace(out))
}

func TestTable_WithHeaderColumnFormatter(t *testing.T) {
	uppercase := func(f string, v ...interface{}) string {
		return strings.ToUpper(fmt.Sprintf(f, v...))
	}
	buf := bytes.Buffer{}

	tbl := New("foo", "bar").WithWriter(&buf).WithHeaderColumnFormatter(uppercase)
	tbl.Print()
	out := buf.String()

	assert.Contains(t, out, "FOO")
	assert.Contains(t, out, "BAR")

	buf.Reset()
	tbl.WithHeaderColumnFormatter(nil).Print()
	out = buf.String()

	assert.Contains(t, out, "foo")
	assert.Contains(t, out, "bar")
}

func TestTable_WithFirstColumnFormatter(t *testing.T) {
	uppercase := func(f string, v ...interface{}) string {
		return strings.ToUpper(fmt.Sprintf(f, v...))
	}

	buf := bytes.Buffer{}

	tbl := New("foo", "bar").WithWriter(&buf).WithFirstColumnFormatter(uppercase).AddRow("fizz", "buzz")
	tbl.Print()
	out := buf.String()

	assert.Contains(t, out, "foo")
	assert.Contains(t, out, "bar")
	assert.Contains(t, out, "FIZZ")
	assert.Contains(t, out, "buzz")

	buf.Reset()
	tbl.WithFirstColumnFormatter(nil).Print()
	out = buf.String()

	assert.Contains(t, out, "fizz")
	assert.Contains(t, out, "buzz")
}

func TestTable_WithPadding(t *testing.T) {
	// zero value
	buf := bytes.Buffer{}
	tbl := New("foo", "bar").WithWriter(&buf).WithPadding(0)
	tbl.Print()
	out := buf.String()
	assert.Contains(t, out, "foobar")

	// positive value
	buf.Reset()
	tbl.WithPadding(4).Print()
	out = buf.String()
	assert.Contains(t, out, "foo    bar    ")

	// negative value
	buf.Reset()
	tbl.WithPadding(-1).Print()
	out = buf.String()
	assert.Contains(t, out, "foobar")
}

func TestTable_WithWriter(t *testing.T) {
	// not that we haven't been using it in all these tests but:
	buf := bytes.Buffer{}
	New("foo", "bar").WithWriter(&buf).Print()
	assert.NotEmpty(t, buf.String())
}

func TestTable_AddRow(t *testing.T) {
	buf := bytes.Buffer{}
	tbl := New("foo", "bar").WithWriter(&buf).AddRow("fizz", "buzz")
	tbl.Print()
	out := buf.String()
	assert.Contains(t, out, "fizz")
	assert.Contains(t, out, "buzz")
	lines := strings.Count(out, "\n")

	// empty should add empty line
	buf.Reset()
	tbl.AddRow().Print()
	assert.Equal(t, lines+1, strings.Count(buf.String(), "\n"))

	// less than one will fill left-to-right
	buf.Reset()
	tbl.AddRow("cat").Print()
	assert.Contains(t, buf.String(), "\ncat")

	// more than initial length are truncated
	buf.Reset()
	tbl.AddRow("bippity", "boppity", "boo").Print()
	assert.NotContains(t, buf.String(), "boo")
}