package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"strconv"
	"strings"
)

type Screenshot struct {
	Size        int64
	Type        string
	Filename    string
	filePointer *os.File
}

func getHandler(writer http.ResponseWriter, request *http.Request) {
	filename := parseFilenameFromRequest(request)

	screenshot := findScreenshot(filename)
	fileStats, _ := screenshot.Stat()

	writer.Header().Set("Content-Type", "image/png")
	writer.Header().Set("Content-Length", strconv.FormatInt(fileStats.Size(), 10))

	io.Copy(writer, screenshot)
}

func parseFilenameFromRequest(request *http.Request) string {
	uri := request.URL.RequestURI()
	uriParts := strings.Split(uri, "/")

	return uriParts[len(uriParts)-1]
}

func findScreenshot(filename string) *os.File {
	file, _ := os.Open("/tmp/" + filename)

	return file
}

func postHandler(responseWriter http.ResponseWriter, request *http.Request) {
	request.ParseMultipartForm(32 << 20)
	file, requestHandler, err := request.FormFile("screenshot")
	defer file.Close()

	if err != nil {
		return
	}

	if isImage(file) == false {
		return
	}

	screenshot := writeFileToDisk(file, requestHandler.Header)
	json, err := json.Marshal(screenshot)

	if err != nil {
		return
	}

	fmt.Fprint(responseWriter, string(json))
}

func isImage(file multipart.File) bool {
	filetypeBuffer := make([]byte, 512)
	file.Read(filetypeBuffer)

	contentType := http.DetectContentType(filetypeBuffer)

	return contentType == "image/png"
}

func writeFileToDisk(screenshot multipart.File, requestHandler textproto.MIMEHeader) *Screenshot {
	filename := createHashForFile(screenshot)

	osFile, _ := os.OpenFile("/tmp/"+filename, os.O_WRONLY|os.O_CREATE, 0666)
	defer osFile.Close()
	io.Copy(osFile, screenshot)

	return &Screenshot{Filename: filename}
}

func createHashForFile(file multipart.File) string {
	hash := sha1.New()
	io.Copy(hash, file)
	defer file.Seek(0, io.SeekStart)

	bytes := hash.Sum(nil)[:4]

	return hex.EncodeToString(bytes)
}

func main() {
	http.HandleFunc("/screenshot", postHandler)
	http.HandleFunc("/screenshot/", getHandler)
	http.ListenAndServe(":9000", nil)
}
