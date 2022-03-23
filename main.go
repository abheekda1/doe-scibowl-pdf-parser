package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

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
	// Check arg length to see if at least one argument was supplied, if not then panic
	argLength := len(os.Args[1:])
	if argLength < 1 {
		panic(fmt.Errorf("failed to detect file"))
	}

	// Set PDF file to the first argument, output file is
	// either the pdf file
	pdfFile := os.Args[1]
	outFile := strings.TrimSuffix(pdfFile, ".pdf") + ".json"

	content, err := readPdf(pdfFile)
	if err != nil {
		panic(err)
	}

	content = strings.ReplaceAll(content, "~", "")

	var formattedQuestions []Question

	// Splitting by "TOSS-UP " will allow looping though each question
	questionList := strings.Split(strings.ReplaceAll(strings.ReplaceAll(content, "\n", ""), "  ", " "), "TOSS-UP ")
	for i := 1; i < len(questionList); i++ {
		q, err := getQuestionObj(questionList[i])
		if err != nil {
			if err.Error() == "math" {
				continue
			} else {
				panic(err)
			}
		}

		formattedQuestions = append(formattedQuestions, *q)
	}

	// Create a JSON array object from the array of question
	questionJson, err := json.MarshalIndent(formattedQuestions, "", "  ")
	if err != nil {
		panic(err)
	}

	// Save the JSON to a file to be later used elsewhere
	ioutil.WriteFile(outFile, questionJson, os.ModePerm)
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

	/* Category is an empty string to begin with
	 * and is defined by looping through valid categories
	 * until a case-insensitive regex match is found
	 */
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

// This example is taken from the source of the PDF package in use
func readPdf(path string) (string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}

	defer f.Close()
	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	buf.ReadFrom(b)
	return buf.String(), nil
}
