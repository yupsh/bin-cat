package main

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func TestRun(t *testing.T) {
	cases := []struct {
		files      map[string]string
		name       string
		version    string
		stdin      string
		wantOut    string
		wantErrSub string
		args       []string
		wantCode   int
	}{
		{
			name:    "stdin passthrough",
			args:    []string{"cat"},
			stdin:   "alpha\nbeta\n",
			wantOut: "alpha\nbeta\n",
		},
		{
			name:    "number all lines",
			args:    []string{"cat", "-n"},
			stdin:   "alpha\nbeta\n",
			wantOut: "     1\talpha\n     2\tbeta\n",
		},
		{
			name:    "number nonblank lines",
			args:    []string{"cat", "-b"},
			stdin:   "alpha\n\nbeta\n",
			wantOut: "     1\talpha\n\n     2\tbeta\n",
		},
		{
			name:    "file source",
			args:    []string{"cat", "/in.txt"},
			files:   map[string]string{"/in.txt": "one\ntwo\n"},
			wantOut: "one\ntwo\n",
		},
		{
			name:    "version flag reports injected version",
			version: "1.2.3",
			args:    []string{"cat", "--version"},
			wantOut: "cat version 1.2.3\n",
		},
		{
			name:       "unknown flag errors",
			args:       []string{"cat", "--nope"},
			wantCode:   1,
			wantErrSub: "cat:",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			for path, content := range tc.files {
				if err := afero.WriteFile(fs, path, []byte(content), 0o644); err != nil {
					t.Fatalf("write fixture %s: %v", path, err)
				}
			}

			var out, errOut bytes.Buffer
			code := run(tc.version, tc.args, strings.NewReader(tc.stdin), &out, &errOut, fs)

			if code != tc.wantCode {
				t.Fatalf("exit code = %d, want %d (stderr=%q)", code, tc.wantCode, errOut.String())
			}
			if tc.wantErrSub == "" && out.String() != tc.wantOut {
				t.Fatalf("stdout = %q, want %q", out.String(), tc.wantOut)
			}
			if tc.wantErrSub != "" && !strings.Contains(errOut.String(), tc.wantErrSub) {
				t.Fatalf("stderr = %q, want substring %q", errOut.String(), tc.wantErrSub)
			}
		})
	}
}

func Test_main(t *testing.T) {
	origExit, origRun := osExit, runCLI
	t.Cleanup(func() { osExit, runCLI = origExit, origRun })

	gotCode := -1
	osExit = func(code int) { gotCode = code }
	runCLI = func(string, []string, io.Reader, io.Writer, io.Writer, afero.Fs) int { return 7 }

	main()

	if gotCode != 7 {
		t.Fatalf("main propagated exit code %d, want 7", gotCode)
	}
}
