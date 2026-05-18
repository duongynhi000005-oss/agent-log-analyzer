package remediation

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path"
	"strconv"
	"strings"
)

func WriteZip(w io.Writer, artifact Artifact) error {
	if len(artifact.Files) == 0 {
		return errors.New("artifact has no files")
	}
	zipWriter := zip.NewWriter(w)
	for _, file := range artifact.Files {
		if err := validateArchivePath(file.Path); err != nil {
			_ = zipWriter.Close()
			return err
		}
		mode, err := parseMode(file.Mode)
		if err != nil {
			_ = zipWriter.Close()
			return fmt.Errorf("invalid mode for %s: %w", file.Path, err)
		}
		header := &zip.FileHeader{Name: file.Path, Method: zip.Deflate}
		header.SetMode(mode)
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			_ = zipWriter.Close()
			return err
		}
		if _, err := io.WriteString(writer, file.Content); err != nil {
			_ = zipWriter.Close()
			return err
		}
	}
	return zipWriter.Close()
}

func validateArchivePath(name string) error {
	clean := path.Clean(name)
	if name == "" || strings.HasPrefix(name, "/") || strings.HasPrefix(name, `\`) || strings.Contains(name, `\`) {
		return fmt.Errorf("unsafe archive path: %q", name)
	}
	if clean == "." || clean != name || strings.HasPrefix(clean, "../") || strings.Contains(clean, "/../") {
		return fmt.Errorf("unsafe archive path: %q", name)
	}
	return nil
}

func parseMode(raw string) (fs.FileMode, error) {
	if raw == "" {
		raw = "0644"
	}
	parsed, err := strconv.ParseUint(raw, 8, 32)
	return fs.FileMode(parsed), err
}
