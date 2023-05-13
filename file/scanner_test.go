package file_test

import (
	"strings"
	"testing"

	"github.com/patrickhuber/pavement/file"
	"github.com/stretchr/testify/require"
)

func TestCanScanFrom(t *testing.T) {
	str := "FROM ubuntu"
	expected := []*file.Token{
		{Type: file.IDENT, Position: 0, Content: "FROM"},
		{Type: file.WS, Position: 4, Content: " "},
		{Type: file.IDENT, Position: 5, Content: "ubuntu"}}
	testScanner(t, expected, str)
}

func TestCanScanFromWithVersion(t *testing.T) {
	str := "FROM ubuntu:latest"
	expected := []*file.Token{
		{Type: file.IDENT, Position: 0, Content: "FROM"},
		{Type: file.WS, Position: 4, Content: " "},
		{Type: file.IDENT, Position: 5, Content: "ubuntu"},
		{Type: file.COLON, Position: 11, Content: ":"},
		{Type: file.IDENT, Position: 12, Content: "latest"},
	}
	testScanner(t, expected, str)
}

func TestCanScanRunCommand(t *testing.T) {
	str := "RUN command arg1 arg2"
	expected := []*file.Token{
		{Type: file.IDENT, Position: 0, Content: "RUN"},
		{Type: file.WS, Position: 3, Content: " "},
		{Type: file.IDENT, Position: 4, Content: "command"},
		{Type: file.WS, Position: 11, Content: " "},
		{Type: file.IDENT, Position: 12, Content: "arg1"},
		{Type: file.WS, Position: 16, Content: " "},
		{Type: file.IDENT, Position: 17, Content: "arg2"},
	}
	testScanner(t, expected, str)
}

func TestCanScanContinuation(t *testing.T) {
	str := `RUN command \
arg1 arg2`
	expected := []*file.Token{
		{Type: file.IDENT, Position: 0, Content: "RUN"},
		{Type: file.WS, Position: 3, Content: " "},
		{Type: file.IDENT, Position: 4, Content: "command"},
		{Type: file.WS, Position: 11, Content: " "},
		{Type: file.CONTINUATION, Position: 12, Content: `\`},
		{Type: file.WS, Position: 13, Content: "\n"},
		{Type: file.IDENT, Position: 14, Content: "arg1"},
		{Type: file.WS, Position: 18, Content: " "},
		{Type: file.IDENT, Position: 19, Content: "arg2"}}
	testScanner(t, expected, str)
}

func TestCanScanFlags(t *testing.T) {
	str := `RUN command -flag1 --flag2 /flag3`
	expected := []*file.Token{
		{Type: file.IDENT, Position: 0, Content: "RUN"},
		{Type: file.WS, Position: 3, Content: " "},
		{Type: file.IDENT, Position: 4, Content: "command"},
		{Type: file.WS, Position: 11, Content: " "},
		{Type: file.FLAG, Position: 12, Content: "-flag1"},
		{Type: file.WS, Position: 18, Content: " "},
		{Type: file.FLAG, Position: 19, Content: "--flag2"},
		{Type: file.WS, Position: 26, Content: " "},
		{Type: file.FLAG, Position: 27, Content: "/flag3"}}

	testScanner(t, expected, str)
}

func testScanner(t *testing.T, expected []*file.Token, str string) {
	reader := strings.NewReader(str)
	scanner := file.NewScanner(reader)
	for i := 0; i < len(expected); i++ {
		token, err := scanner.Scan()
		require.Nil(t, err, "err not nil at '%d'", i)
		require.Equal(t, expected[i], token, "error at token '%d'", i)
	}
	token, err := scanner.Scan()
	require.Nil(t, err)
	require.NotNil(t, token)
	require.Equal(t, token.Type, file.EOF)
}
