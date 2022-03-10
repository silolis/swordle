package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"unicode"

	"github.com/gorilla/mux"
)

var cfreq = map[byte]float32{
	'E': 12.0,
	'T': 9.10,
	'A': 8.12,
	'O': 7.68,
	'I': 7.31,
	'N': 6.95,
	'S': 6.28,
	'R': 6.02,
	'H': 5.92,
	'D': 4.32,
	'L': 3.98,
	'U': 2.88,
	'C': 2.71,
	'M': 2.61,
	'F': 2.30,
	'Y': 2.11,
	'W': 2.09,
	'G': 2.03,
	'P': 1.82,
	'B': 1.49,
	'V': 1.11,
	'K': 0.69,
	'X': 0.17,
	'Q': 0.11,
	'J': 0.10,
	'Z': 0.07}

func Swordle(word [5]byte) (float32, error) {
	var score float32 = 0

	for _, c := range word {
		cost, ok := cfreq[c]
		score += cost
		if !ok {
			return 0, fmt.Errorf("no such character (%c) in cost map", c)
		}
	}
	return score, nil
}

func Word(str string) ([5]byte, error) {
	word := [5]byte{'w', 'o', 'r', 'd', 'l'}
	if len(str) != 5 || len([]byte(str)) != 5 {
		return word, fmt.Errorf("word length (%d) is not 5", len(str))
	}

	for i, c := range str {
		if (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') {
			return word, fmt.Errorf("word contains illegal character (%c)", c)
		}

		word[i] = byte(unicode.ToUpper(rune(c)))
	}
	return word, nil
}

func WordCountHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	str, ok := vars["word"]
	log.Printf("swordl: requesting score for %s", str)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		//fmt.Fprintf(w, `{}`, )
		return
	}

	word, err := Word(str)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{'error': '%v'}`, err)
		return
	}

	score, err := Swordle(word)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{'error': '%v'}`, err)
		return
	}
	log.Printf("swordl: score for %s: %f", str, score)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{
  '5word': '%v',
  'score': '%f'
}`, string(word[:]), score)
}

func main() {
	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
		log.Printf("defaulting to port %s", port)
	}

	r := mux.NewRouter()
	r.HandleFunc("/score/{word}", WordCountHandler).Methods("GET")

	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    ":" + port,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("✨ swordl started successfully on %s ✨", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
