package server

import (
	"encoding/json"
	"fmt"
	"github.com/michaeljoelphillips/scatter-shot/internal/storage"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type HttpServer struct {
	port    string
	storage *storage.Storage
}

func NewHttpServer(port string, storage *storage.Storage) *HttpServer {
	return &HttpServer{
		port:    port,
		storage: storage,
	}
}

func (server HttpServer) Start() {
	fmt.Printf("Starting the HTTP Server on port %s...\n", server.port)

	http.HandleFunc("/upload", server.postHandler)
	http.HandleFunc("/view/", server.getHandler)
	http.ListenAndServe(fmt.Sprintf(":%s", server.port), nil)
}

func (server HttpServer) getHandler(writer http.ResponseWriter, request *http.Request) {
	filename := parseFilenameFromRequest(request)

	file, err := server.storage.Find(filename)

	if err != nil {
		fmt.Println(err)

		return
	}

	defer file.Close()

	writer.Header().Set("Content-Type", file.ContentType)
	writer.Header().Set("Content-Length", strconv.FormatInt(file.Size, 10))

	io.Copy(writer, file)
}

func parseFilenameFromRequest(request *http.Request) string {
	uri := request.URL.RequestURI()
	uriParts := strings.Split(uri, "/")

	return uriParts[len(uriParts)-1]
}

func (server HttpServer) postHandler(responseWriter http.ResponseWriter, request *http.Request) {
	request.ParseMultipartForm(32 << 20)
	content, _, err := request.FormFile("file")
	defer content.Close()

	if err != nil {
		return
	}

	file, err := server.storage.Add(content)

	if err != nil {
		fmt.Println(err)

		return
	}

	json, err := json.Marshal(file)

	if err != nil {
		return
	}

	fmt.Fprint(responseWriter, string(json))
}
