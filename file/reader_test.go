package file_test

import (
	"strings"
	"testing"

	"github.com/patrickhuber/pavement/file"
	"github.com/stretchr/testify/require"
)

func TestCanRead(t *testing.T) {
	str := `FROM ubuntu:latest
RUN ls -la`
	expected := []*file.Directive{
		{
			Type: file.From,
			From: &file.FromDirective{
				Base:    "ubuntu",
				Version: "latest",
			},
		},
		{
			Type: file.Run,
			Run: &file.RunDirective{
				Command: "ls",
				Arguments: []string{
					"-la",
				},
			},
		},
	}
	testReader(t, str, expected)
}

func TestCanReadWithContinuationBetweenRun(t *testing.T) {
	str := `FROM ubuntu:latest
RUN ls \
-la
RUN echo "hello"`
	expected := []*file.Directive{
		{
			Type: file.From,
			From: &file.FromDirective{
				Base:    "ubuntu",
				Version: "latest",
			},
		},
		{
			Type: file.Run,
			Run: &file.RunDirective{
				Command: "ls",
				Arguments: []string{
					"-la",
				},
			},
		},
		{
			Type: file.Run,
			Run: &file.RunDirective{
				Command: "echo",
				Arguments: []string{
					`"hello"`,
				},
			},
		},
	}

	testReader(t, str, expected)
}

func TestFailsWithoutContinuation(t *testing.T) {
	str := `FROM ubuntu:latest
RUN ls
-la`
	testReaderFail(t, str, 2)
}

func testReaderFail(t *testing.T, str string, failOn int) {
	reader := strings.NewReader(str)
	fileReader := file.NewReader(file.NewScanner(reader))
	for i := 0; true; i++ {
		more, err := fileReader.Next()
		if err != nil {
			if i == failOn {
				return
			}
			require.Nil(t, err)
		}
		if !more {
			if i == failOn {
				break
			}
			require.Fail(t, "wrong", "expected to fail on %d actual %d", failOn, i)
		}
	}
}
func testReader(t *testing.T, str string, expected []*file.Directive) {

	reader := strings.NewReader(str)
	fileReader := file.NewReader(file.NewScanner(reader))

	for i := 0; i < len(expected); i++ {

		more, err := fileReader.Next()
		require.Nil(t, err)
		require.True(t, more && i < len(expected), "directive '%d' expected more found eof", i)

		current := fileReader.Current()
		e := expected[i]

		// compare structures not pointers
		directiveEqual(t, e, current)
	}
}

func directiveEqual(t *testing.T, expected, actual *file.Directive) {
	require.NotNil(t, expected)
	require.NotNil(t, actual)
	require.Equal(t, expected.Type, actual.Type)
	switch expected.Type {
	case file.From:
		fromEqual(t, expected.From, actual.From)
	case file.Run:
		runEqual(t, expected.Run, actual.Run)
	}
}

func fromEqual(t *testing.T, expected, actual *file.FromDirective) {
	require.NotNil(t, expected)
	require.NotNil(t, actual)
	require.Equal(t, expected, actual)
}

func runEqual(t *testing.T, expected, actual *file.RunDirective) {
	require.NotNil(t, expected)
	require.NotNil(t, actual)
	require.Equal(t, expected.Command, actual.Command)
	require.Equal(t, expected.Arguments, actual.Arguments)
}
