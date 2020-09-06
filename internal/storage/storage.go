package storage

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"io"
	"os"
	"path/filepath"
)

type File struct {
	Name        string
	Size        int64
	Contents    io.ReadCloser
	ContentType string
}

type Storage struct {
	basePath string
}

func NewStorage(basePath string) *Storage {
	return &Storage{
		basePath: basePath,
	}
}

func (file File) Close() {
	file.Contents.Close()
}

func (file File) Read(p []byte) (n int, err error) {
	return file.Contents.Read(p)
}

func (storage Storage) Find(filename string) (*File, error) {
	file, err := os.Open(storage.basePath + "/" + filename)
	defer file.Seek(0, io.SeekStart)

	if err != nil {
		return &File{}, fmt.Errorf("No such file %s/%s", storage.basePath, filename)
	}

	stats, err := file.Stat()

	if err != nil {
		return &File{}, fmt.Errorf("Failed to read stats for file %s/%s", storage.basePath, filename)
	}

	contentType, err := mimetype.DetectReader(file)

	if err != nil {
		return &File{}, fmt.Errorf("Failed to detect the content type for file %s/%s", storage.basePath, filename)
	}

	return &File{
			Contents:    file,
			Size:        stats.Size(),
			ContentType: contentType.String(),
			Name:        filepath.Base(file.Name()),
		},
		nil
}

func (storage Storage) Add(contents io.ReadSeeker) (*File, error) {
	fileHash := createFileHash(contents)
	file, err := storage.writeContents(fileHash, contents)
	defer file.Close()

	if err != nil {
		return nil, fmt.Errorf("Could not write file contents to disk")
	}

	stats, err := file.Stat()

	if err != nil {
		return nil, fmt.Errorf("Failed to read stats for file %s/%s", storage.basePath, fileHash)
	}

	contentType, _ := mimetype.DetectReader(file)

	return &File{
			Contents:    file,
			Name:        fileHash,
			Size:        stats.Size(),
			ContentType: contentType.String(),
		},
		nil
}

func createFileHash(contents io.ReadSeeker) string {
	hash := sha1.New()
	io.Copy(hash, contents)
	defer contents.Seek(0, io.SeekStart)

	bytes := hash.Sum(nil)[:4]

	return hex.EncodeToString(bytes)
}

func (storage Storage) writeContents(filename string, contents io.Reader) (*os.File, error) {
	file, err := os.OpenFile(
		storage.basePath+"/"+filename,
		os.O_WRONLY|os.O_CREATE,
		0666,
	)

	if err != nil {
		return nil, fmt.Errorf("Could not open file %s/%s for writing", storage.basePath, filename)
	}

	io.Copy(file, contents)

	return file, nil
}
