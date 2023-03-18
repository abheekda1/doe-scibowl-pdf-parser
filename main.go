package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"doe-scibowl-pdf-parser/parse"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/pdf", pdfHandler).Methods("POST")

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func pdfHandler(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("There was an error"))
		return
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(file)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("There was an error"))
		return
	}

	readerAt := bytes.NewReader(buf.Bytes())

	content, err := parse.ReadPdfToString(readerAt)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("There was an error"))
		return
	}

	var formattedQuestions []parse.Question

	// Splitting by "TOSS-UP " will allow looping though each question
	questionList := strings.Split(strings.ReplaceAll(strings.ReplaceAll(content, "\n", ""), "  ", " "), "TOSS-UP")
	fmt.Println(len(questionList))
	for i := 1; i < len(questionList); i++ {
		q, err := parse.GetQuestionObj(questionList[i])
		if err != nil {
			if err.Error() == "category is math" {
				continue
			} else {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("There was an error"))
				return
			}
		}

		formattedQuestions = append(formattedQuestions, *q)
	}

	// Create a JSON array object from the array of question
	questionJson, err := json.MarshalIndent(formattedQuestions, "", "  ")
	fmt.Println(questionJson)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("There was an error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(questionJson)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("There was an error"))
		return
	}
}
