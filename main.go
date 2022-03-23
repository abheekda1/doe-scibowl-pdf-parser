package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/ledongthuc/pdf"
)

type Question struct {
	Category       string `json:"category"`
	TossupFormat   string `json:"tossupFormat"`
	TossupQuestion string `json:"tossupQuestion"`
	TossupAnswer   string `json:"tossupAnswer"`
	BonusFormat    string `json:"bonusFormat"`
	BonusQuestion  string `json:"bonusQuestion"`
	BonusAnswer    string `json:"bonusAnswer"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/pdf", pdfHandler).Methods("POST")

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func pdfHandler(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(file)
	readerAt := bytes.NewReader(buf.Bytes())
	//fmt.Println(buf.String())

	content, err := readPdf(readerAt)
	if err != nil {
		log.Fatal(err)
	}

	var formattedQuestions []Question

	// Splitting by "TOSS-UP " will allow looping though each question
	questionList := strings.Split(strings.ReplaceAll(strings.ReplaceAll(content, "\n", ""), "  ", " "), "TOSS-UP ")
	for i := 1; i < len(questionList); i++ {
		q, err := getQuestionObj(questionList[i])
		if err != nil {
			if err.Error() == "math" {
				continue
			} else {
				log.Fatal(err)
			}
		}

		formattedQuestions = append(formattedQuestions, *q)
	}

	// Create a JSON array object from the array of question
	questionJson, err := json.MarshalIndent(formattedQuestions, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(questionJson)
}

func getQuestionObj(question string) (*Question, error) {
	// List of possible categories, could maybe be expanded
	categories := []string{
		"BIOLOGY",
		"CHEMISTRY",
		"PHYSICS",
		"EARTH AND SPACE",
		"EARTH SCIENCE",
		"ENERGY",
		"MATH",
		"GENERAL SCIENCE",
		"ASTRONOMY",
	}

	TU := strings.SplitN(question, "BONUS ", 2)[0] // Toss-Up question
	B := strings.SplitN(question, "BONUS ", 2)[1]  // Bonus question

	// Category is an empty string to begin with
	// and is defined by looping through valid categories
	// until a case-insensitive regex match is found
	category := ""
	for _, cat := range categories {
		catExp := regexp.MustCompile(`(?i)` + cat) // Case-insensitive version of the category
		if catExp.MatchString(TU) {
			category = catExp.FindString(TU)
			break
		}
	}

	if category == "" {
		return nil, fmt.Errorf("category not found")
	}

	// Return if it's math as math is currently broken
	if category == "MATH" {
		return nil, fmt.Errorf("math")
	}

	// Replace the first dash from Toss-Ups and bonuses
	// if it appears after the category, helps with
	// format detection later on
	if strings.Contains(TU, category+" –") {
		TU = strings.Replace(TU, " –", "", 1)
	}

	if strings.Contains(B, category+" –") {
		B = strings.Replace(B, " –", "", 1)
	}

	// Parse the Toss-Up and bonus strings for their data
	tuFormat := strings.Join(strings.SplitN(strings.SplitN(TU, category+" ", 2)[1], " ", 4)[0:2], " ")
	tuQuestion := strings.SplitN(strings.SplitN(TU, " ANSWER:", 2)[0], tuFormat+" ", 2)[1]
	tuAnswer := strings.SplitN(TU, "ANSWER: ", 2)[1]

	bFormat := strings.Join(strings.SplitN(strings.SplitN(B, category+" ", 2)[1], " ", 4)[0:2], " ")
	bQuestion := strings.SplitN(strings.SplitN(B, " ANSWER:", 2)[0], bFormat+" ", 2)[1]
	bAnswer := strings.SplitN(B, "ANSWER: ", 2)[1]

	// Remove the footer which often trails the bonus answer as it's on the bottom of the page
	// FIXME: some formats still don't get picked up by this regex
	footerExp := regexp.MustCompile(`(?i)\s*((20|19)\d+.*Regional.*)?(High School |Middle School )?(Round|Questions).*\s+Page\s*\d+.*`)

	// Create a question object to return
	q := Question{
		Category:       strings.ToUpper(strings.TrimSpace(category)),
		TossupFormat:   strings.TrimSpace(tuFormat),
		TossupQuestion: strings.TrimSpace(tuQuestion),
		TossupAnswer:   strings.TrimSpace(footerExp.ReplaceAllString(tuAnswer, "")),
		BonusFormat:    strings.TrimSpace(bFormat),
		BonusQuestion:  strings.TrimSpace(bQuestion),
		BonusAnswer:    strings.TrimSpace(footerExp.ReplaceAllString(bAnswer, "")),
	}

	return &q, nil
}

// Read PDF data from a bytes reader and return raw text
func readPdf(file *bytes.Reader) (string, error) {
	r, err := pdf.NewReader(file, file.Size())
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	pdfData, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	buf.ReadFrom(pdfData)
	return buf.String(), nil
}
