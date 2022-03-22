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
	argLength := len(os.Args[1:])
	if argLength < 1 {
		panic(fmt.Errorf("Failed to detect file"))
	}

	pdfFile := os.Args[1]
	outFile := strings.TrimSuffix(pdfFile, ".pdf") + ".json"

	content, err := readPdf(pdfFile)
	if err != nil {
		panic(err)
	}

  content = strings.ReplaceAll(content, "~", "")

	var formattedQuestions []Question

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

	questionJson, err := json.MarshalIndent(formattedQuestions, " ", " ")
	if err != nil {
		panic(err)
	}

	ioutil.WriteFile(outFile, questionJson, os.ModePerm)

	return
}

func getQuestionObj(question string) (*Question, error) {
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

	TU := strings.SplitN(question, "BONUS ", 2)[0]
	B := strings.SplitN(question, "BONUS ", 2)[1]

	category := ""
	for _, cat := range categories {
    catExp := regexp.MustCompile(`(?i)` + cat)
		if catExp.MatchString(TU) {
			category = catExp.FindString(TU)
			break
		}
	}

	if category == "" {
		return nil, fmt.Errorf("Category not found")
	}

	if category == "MATH" {
		return nil, fmt.Errorf("math")
	}

  if strings.Contains(TU, category + " –") {
    TU = strings.Replace(TU, " –", "", 1)
  }

  if strings.Contains(B, category + " –") {
    B = strings.Replace(B, " –", "", 1)
  }
    
	tuFormat := strings.Join(strings.SplitN(strings.SplitN(TU, category+" ", 2)[1], " ", 4)[0:2], " ")
	tuQuestion := strings.SplitN(strings.SplitN(TU, " ANSWER:", 2)[0], tuFormat+" ", 2)[1]
	tuAnswer := strings.SplitN(TU, "ANSWER: ", 2)[1]

	bFormat := strings.Join(strings.SplitN(strings.SplitN(B, category+" ", 2)[1], " ", 4)[0:2], " ")
	bQuestion := strings.SplitN(strings.SplitN(B, " ANSWER:", 2)[0], bFormat+" ", 2)[1]
	bAnswer := strings.SplitN(B, "ANSWER: ", 2)[1]

	footerExp := regexp.MustCompile(`(?i)\s*(High School |Middle School |\d+\s*Regional.*)?Round \d.*`)

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

func readPdf(path string) (string, error) {
	f, r, err := pdf.Open(path)
	defer f.Close()
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	buf.ReadFrom(b)
	return buf.String(), nil
}
