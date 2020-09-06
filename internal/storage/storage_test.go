package storage

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func baseDir() string {
	_, runtime, _, _ := runtime.Caller(0)

	return filepath.Dir(runtime) + "/../../test/data"
}

func TestFind(t *testing.T) {
	storage := NewStorage(baseDir())
	file, err := storage.Find("json_fixture")

	if err != nil {
		t.Error(err)
	}

	if file.Size != 37 {
		t.Errorf("Expected to receive a file with a filesize of %d, %d given", 37, file.Size)
	}

	if file.Name != "json_fixture" {
		t.Errorf("Expected to receive a file with the filename of %s, %s given", "abcd1234", file.Name)
	}

	if file.ContentType != "application/json" {
		t.Errorf("Expected to receive a file with a content type of %s, %s given", "application/json", file.ContentType)
	}
}

func TestFindReturnsCorrectMimeTypeForFileWithoutExtension(t *testing.T) {
	storage := NewStorage(baseDir())
	file, _ := storage.Find("8ce37613")

	if file.ContentType != "application/x-tar" {
		t.Errorf("Failed to correctly identify the file's content type, %s expect, %s given", "application/x-tar", file.ContentType)
	}
}

func TestFindReturnsErrorWhenFileDoesNotExist(t *testing.T) {
	storage := NewStorage(baseDir())
	_, err := storage.Find("nonexistent-file")

	if err == nil {
		t.Error("Expected to receive an error for a non-existent file")
	}
}

func TestAdd(t *testing.T) {
	file, _ := os.Open(baseDir() + "/json_fixture")
	storage := NewStorage(baseDir())
	newFile, err := storage.Add(file)

	if err != nil {
		t.Errorf("An error occurred while writing a file")
	}

	if newFile.Size != 37 {
		t.Errorf("Expected to receive a file with a filesize of %d, %d given", 37, newFile.Size)
	}

	if newFile.Name != "9bf48805" {
		t.Errorf("Expected to receive a file with the filename of %s, %s given", "9bf48805", newFile.Name)
	}

	if newFile.ContentType != "application/json" {
		t.Errorf("Expected to receive a file with a content type of %s, %s given", "application/json", newFile.ContentType)
	}
}
