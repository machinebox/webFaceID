package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/machinebox/sdk-go/facebox"
	"github.com/matryer/way"
)

// Server is the app server.
type Server struct {
	assets  string
	facebox *facebox.Client
	router  *way.Router
}

// NewServer makes a new Server.
func NewServer(assets string, facebox *facebox.Client) *Server {
	srv := &Server{
		assets:  assets,
		facebox: facebox,
		router:  way.NewRouter(),
	}
	srv.router.Handle(http.MethodGet, "/assets/", Static("/assets/", assets))

	srv.router.HandleFunc(http.MethodPost, "/webFaceID", srv.handlewebFaceID)
	srv.router.HandleFunc(http.MethodGet, "/", srv.handleIndex)
	return srv
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(s.assets, "index.html"))
}

func (s *Server) handlewebFaceID(w http.ResponseWriter, r *http.Request) {
	img := r.FormValue("imgBase64")
	b64data := img[strings.IndexByte(img, ',')+1:]
	imgDec, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		log.Printf("[ERROR] Error decoding the image %v\n", err)
		http.Error(w, "can not decode the image", http.StatusInternalServerError)
		return
	}
	faces, err := s.facebox.Check(bytes.NewReader(imgDec))
	if err != nil {
		log.Printf("[ERROR] Error on facebox %v\n", err)
		http.Error(w, "something went wrong verifying the faces", http.StatusInternalServerError)
		return
	}
	var response struct {
		FaceLen int    `json:"faces_len"`
		Matched bool   `json:"matched"`
		Name    string `json:"name"`
	}
	response.FaceLen = len(faces)
	if len(faces) == 1 {
		response.Matched = faces[0].Matched
		response.Name = faces[0].Name
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Static gets a static file server for the specified path.
func Static(stripPrefix, dir string) http.Handler {
	h := http.StripPrefix(stripPrefix, http.FileServer(http.Dir(dir)))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}
