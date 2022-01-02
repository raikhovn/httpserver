package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// Define struct
type Server struct {
	Port int
	Host string
}

type DirResponse struct {
	Folder     string
	Subfolders []string
	Error      error
}

type FilesResponse struct {
	Folder string
	Files  []string
	Error  error
}

type FileAttrs struct {
	Name    string    // base name of the file
	Size    int64     // length in bytes for regular files; system-dependent for others
	Mode    uint64    // file mode bits
	ModTime time.Time // modification time
	Error   error
}

func ListDirs(root string) *DirResponse {

	var resp *DirResponse
	var files []string

	fileInfos, err := ioutil.ReadDir(root)

	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			files = append(files, fileInfo.Name())
		}
	}

	if err != nil {
		log.Fatal(err)
		resp = &DirResponse{
			Folder:     root,
			Subfolders: nil,
			Error:      err}
	} else {
		resp = &DirResponse{
			Folder:     root,
			Subfolders: files,
			Error:      nil}
	}

	return resp
}

func GetFileInfo(file string) *FileAttrs {

	var resp *FileAttrs

	fileInfo, err := os.Stat(file)

	if err != nil {
		resp = new(FileAttrs)
		resp.Error = err

	} else {
		resp = &FileAttrs{
			Name:    fileInfo.Name(),
			Size:    fileInfo.Size(),
			Mode:    uint64(fileInfo.Mode()),
			ModTime: fileInfo.ModTime(),
			Error:   nil}
	}

	return resp
}

func ListFiles(dir string) *FilesResponse {

	var resp *FilesResponse
	var files []string

	fileInfos, err := ioutil.ReadDir(dir)

	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() {
			files = append(files, fileInfo.Name())
		}
	}

	if err != nil {
		log.Fatal(err)
		resp = &FilesResponse{
			Folder: dir,
			Files:  nil,
			Error:  err}
	} else {
		resp = &FilesResponse{
			Folder: dir,
			Files:  files,
			Error:  nil}
	}

	return resp
}

// Define constructor
func NewServer(Port int, Host string) *Server {
	s := new(Server)
	s.Host = Host
	s.Port = Port

	return s
}

// Define handlers
func ListDirsHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("RegisterListDirs() called")
	params, err := r.URL.Query()["dir"]

	if !err || len(params[0]) < 1 {
		log.Println("Url Param 'dir' is missing")
		err := fmt.Errorf("Url Param 'dir' is missing")
		panic(err)
	}
	log.Println("Calling ListDirs")
	resp := ListDirs(params[0])

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}
}

func ListFilesHandler(w http.ResponseWriter, r *http.Request) {

	params, err := r.URL.Query()["dir"]

	if !err || len(params[0]) < 1 {
		log.Println("Url Param 'dir' is missing")
		err := fmt.Errorf("Url Param 'dir' is missing")
		panic(err)
	}
	log.Println("Calling ListFiles")
	resp := ListFiles(params[0])

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}
}

func GetFileInfoHandler(w http.ResponseWriter, r *http.Request) {

	params, err := r.URL.Query()["file"]

	if !err || len(params[0]) < 1 {
		log.Println("Url Param 'file' is missing")
		err := fmt.Errorf("Url Param 'file' is missing")
		panic(err)
	}
	log.Println("Calling GetFileInfo")
	resp := GetFileInfo(params[0])

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}
}

// Define struct methods
func (s Server) RegisterListDirs() {
	log.Println("func (s Server) RegisterListDirs()")
	http.HandleFunc("/listdirs", ListDirsHandler)
}

func (s Server) RegisterListFiles() {
	log.Println("func (s Server) RegisterListFiles()")
	http.HandleFunc("/listfiles", ListFilesHandler)
}

func (s Server) RegisterGetFileInfo() {
	log.Println("func (s Server) RegisterGetFileInfo()")
	http.HandleFunc("/getfileinfo", GetFileInfoHandler)
}

func (s Server) RegisterIndex() {
	log.Println("func (s Server) RegisterIndex()")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("RegisterIndex() called")
		fmt.Fprintf(w, "RESTER server is up!")
	})
}

func (s Server) StartListener() {

	log.Printf("Starting server on: %s:%d", s.Host, s.Port)
	http.ListenAndServe(fmt.Sprintf("%s:%d", s.Host, s.Port), nil)
}
