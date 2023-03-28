package file_test

import (
	"io"
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
	test(t, expected, str)
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
	test(t, expected, str)
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
	test(t, expected, str)
}

func test(t *testing.T, expected []*file.Token, str string) {
	reader := strings.NewReader(str)
	scanner := file.NewScanner(reader)
	for i := 0; i < len(expected); i++ {
		token, err := scanner.Scan()
		require.Nil(t, err)
		require.Equal(t, expected[i], token)
	}
	_, err := scanner.Scan()
	require.Equal(t, io.EOF, err)
}
