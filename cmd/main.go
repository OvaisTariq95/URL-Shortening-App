package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

const (
	Host   = "127.0.0.1"
	Port   = "8081"
	Prefix = "http://"
)

type URLShortener struct {
	urls map[string]string
}

func NewURLShortener() *URLShortener {
	return &URLShortener{
		urls: make(map[string]string),
	}
}

func (s *URLShortener) sendJSONResponse(w http.ResponseWriter, status int, payload interface{}) {
	//creating json response to be used throught out the app
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "Error encoding JSON response: "+err.Error(), http.StatusInternalServerError)
	}
}

func (s *URLShortener) shortenURL(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		URL string `json:"url"`
	}

	//decoding url from the body
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		s.sendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Error parsing request body: " + err.Error()})
		return
	}

	// declaring and initializing the url which we are getting from body
	longURL := requestBody.URL
	if longURL == "" {
		s.sendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "URL is required in body"})
		return
	}

	//checking if it has http:// prefix as we need it to execute redirect
	if !strings.HasPrefix(longURL, Prefix) {
		longURL = Prefix + longURL
	}

	//generating shortcode from the urlprovided in the body
	shortCode := s.generateShortCode(longURL)

	//checking if url already exist in the map
	if _, found := s.urls[shortCode]; found {
		storedURL := Prefix + Host + ":" + Port + "/" + shortCode
		s.sendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("URL:%s already exists as: %s", longURL, storedURL)})
		return
	}

	//sending response in the api
	s.urls[shortCode] = longURL
	s.sendJSONResponse(w, http.StatusOK, map[string]string{"shortenedURL": Prefix + Host + ":" + Port + "/" + shortCode})
}

func (s *URLShortener) originalURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	//assigning value of the provided shortcode in the api param
	shortCode := vars["shortCode"]

	//checking if there is any url stored in the map according to the shortcode provided, and sending relevant response afterwards
	storedLongURL, found := s.urls[shortCode]
	if !found {
		s.sendJSONResponse(w, http.StatusNotFound, map[string]string{"error": "URL not found"})
	}

	s.sendJSONResponse(w, http.StatusOK, map[string]string{"originalURL": storedLongURL})
}

func (s *URLShortener) redirectURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	//assigning value of the provided shortcode in the api param
	shortCode := vars["shortCode"]

	//checking if there is any url stored in the map according to the shortcode provided, and redirecting afterwards
	storedLongURL, found := s.urls[shortCode]
	if !found {
		http.NotFound(w, r)
	}

	http.Redirect(w, r, storedLongURL, http.StatusSeeOther)
}

func (s *URLShortener) allURLs(w http.ResponseWriter, r *http.Request) {
	urlList := make([]string, 0)

	//storing all the urls in a list
	for _, longURL := range s.urls {
		urlList = append(urlList, longURL)
	}

	//creating a response to show those urls
	response := map[string]interface{}{"urls": urlList}

	//checking if url list is empty or not , and sending apropriate response
	if len(urlList) == 0 {
		s.sendJSONResponse(w, http.StatusNoContent, response)
		return
	}

	s.sendJSONResponse(w, http.StatusOK, response)
}

func (s *URLShortener) generateShortCode(longURL string) string {
	return fmt.Sprintf("%x", longURL)
}

func main() {
	//defining route
	r := mux.NewRouter()
	urlShortener := NewURLShortener()

	//defining all the relevant endoints
	r.HandleFunc("/all", urlShortener.allURLs).Methods("GET")
	r.HandleFunc("/shorten", urlShortener.shortenURL).Methods("POST")
	r.HandleFunc("/original/{shortCode}", urlShortener.originalURL).Methods("GET")
	r.HandleFunc("/{shortCode}", urlShortener.redirectURL).Methods("GET")

	//starting server on the provided host and port
	address := Host + ":" + Port
	fmt.Printf("Server is listening on %s\n", address)
	http.Handle("/", r)
	http.ListenAndServe(address, nil)
}
