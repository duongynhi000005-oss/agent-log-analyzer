package remediation

import (
	"archive/zip"
	"bytes"
	"testing"
)

func TestWriteZipWritesPluginFiles(t *testing.T) {
	artifact := Artifact{
		Files: []File{
			{Path: ".claude-plugin/plugin.json", Mode: "0644", Content: "{}\n"},
			{Path: "scripts/hook.py", Mode: "0755", Content: "#!/usr/bin/env python3\n"},
		},
	}
	var buf bytes.Buffer
	if err := WriteZip(&buf, artifact); err != nil {
		t.Fatal(err)
	}
	reader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}
	found := map[string]bool{}
	for _, file := range reader.File {
		found[file.Name] = true
	}
	for _, want := range []string{".claude-plugin/plugin.json", "scripts/hook.py"} {
		if !found[want] {
			t.Fatalf("zip missing %s: %#v", want, found)
		}
	}
}

func TestWriteZipRejectsUnsafePaths(t *testing.T) {
	for _, path := range []string{
		"../escape",
		"/absolute",
		`windows\path`,
		"safe/../escape",
		"",
	} {
		artifact := Artifact{Files: []File{{Path: path, Content: "x"}}}
		var buf bytes.Buffer
		if err := WriteZip(&buf, artifact); err == nil {
			t.Fatalf("expected unsafe path %q to fail", path)
		}
	}
}
