package main

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/chdb-io/chdb-go/chdb"
)

// Embedding the HTML file
//
//go:embed play.html
var content embed.FS

var (
	path string
)

func init() {
	path = os.Getenv("DATA_PATH")
	if path == "" {
		path = ".chdb_data"
	}
}

type server struct {
	sess *chdb.Session
}

func newServer(path string) *server {
	s := &server{}
	sess, err := chdb.NewSession(path)
	if err != nil {
		log.Fatalf("Error creating session: %v", err)
	}
	s.sess = sess
	return s
}

func (s *server) handleRootPost(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	format := r.URL.Query().Get("default_format")
	if format == "" {
		format = "TSV"
	}
	database := r.URL.Query().Get("database")
	var body string
	if r.Body != nil {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		body = string(bodyBytes)
	}
	if body != "" {
		body = strings.Join(strings.Fields(body), " ")
		if query != "" {
			query += " "
		}
		query += body
	}
	if database != "" {
		query = "USE " + database + "; " + query
	}
	log.Printf("Query: %s\n", query)
	output, err := s.sess.Query(query, format)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", output)
}

func (s *server) handlePing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Ok")
}

func (s *server) handlePlay(w http.ResponseWriter, r *http.Request) {
	data, err := content.ReadFile("play.html")
	if err != nil {
		http.Error(w, "Unable to open play.html", http.StatusInternalServerError)
		return
	}
	http.ServeContent(w, r, "play.html", time.Now(), bytes.NewReader(data))
}

func (s *server) handleNotFound(w http.ResponseWriter, r *http.Request) {
	data, err := content.ReadFile("play.html")
	if err != nil {
		http.Error(w, "Unable to open play.html", http.StatusInternalServerError)
		return
	}
	http.ServeContent(w, r, "play.html", time.Now(), bytes.NewReader(data))
}

func main() {
	s := newServer(path)
	mux := http.NewServeMux()

	// Register specific handlers
	mux.HandleFunc("/play", s.handlePlay)
	mux.HandleFunc("/ping", s.handlePing)

	// Catch-all handler: place this last
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			if r.Method == http.MethodPost {
				s.handleRootPost(w, r)
				return
			}
			http.Redirect(w, r, "/play", http.StatusFound)
		default:
			s.handleNotFound(w, r) // Use this for 404 responses
		}
	})

	log.Println("Listening on :8123")
	log.Fatal(http.ListenAndServe(":8123", mux))
}
