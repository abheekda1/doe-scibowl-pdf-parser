package parse

import "testing"

func TestParse(t *testing.T) {
	got, err := GetQuestionObj("1) Chemistry – Multiple Choice In a polyelectronic atom, which of the following orbitals would have the lowest energy? W) 3s X) 4s Y) 3p Z) 3d ANSWER: W) 3s BONUS 1) Chemistry – Short Answer Rank the following three carboxylic [CAR-box-IHL-ik] acid derivatives in terms of increasing reactivity to nucleophilic substitution: 1) Amide; 2) Acid anhydride; 3) Ester. ANSWER: 1, 3, 2")
	want := Question{
		Category:       "CHEMISTRY",
		TossupFormat:   "Multiple Choice",
		TossupQuestion: "In a polyelectronic atom, which of the following orbitals would have the lowest energy? W) 3s X) 4s Y) 3p Z) 3d",
		TossupAnswer:   "W) 3s",
		BonusFormat:    "Short Answer",
		BonusQuestion:  "Rank the following three carboxylic [CAR-box-IHL-ik] acid derivatives in terms of increasing reactivity to nucleophilic substitution: 1) Amide; 2) Acid anhydride; 3) Ester.",
		BonusAnswer:    "1, 3, 2",
	}

	if err != nil {
		t.Errorf("encountered an error: %q", err)
	}

	if *got != want {
		t.Errorf("got: %q, wanted: %q", got, want)
	}
}

func TestParseInvalidCategory(t *testing.T) {
	got, err := GetQuestionObj("1) Banana – Multiple Choice In a polyelectronic atom, which of the following orbitals would have the lowest energy? W) 3s X) 4s Y) 3p Z) 3d ANSWER: W) 3s BONUS 1) Chemistry – Short Answer Rank the following three carboxylic [CAR-box-IHL-ik] acid derivatives in terms of increasing reactivity to nucleophilic substitution: 1) Amide; 2) Acid anhydride; 3) Ester. ANSWER: 1, 3, 2")

	if err != nil && got == nil {
		if err.Error() != "category not found" {
			t.Errorf("expected error \"category not found\", received %q", err)
		}
	} else {
		t.Errorf("expected a non-nil error and a nil question, received a(n) %q error and a(n) %q question", err, got)
	}
}

func TestParseMathCategory(t *testing.T) {
	got, err := GetQuestionObj("1) Math - Short Answer: What is the volume of a sphere of radius \"R\"? ANSWER: (4/3) p R3 BONUS 1) Short Answer - What is the surface area of a sphere of radius \"r\"? ANSWER: 4 p R2 or 12.566 x R2")

	if err != nil && got == nil {
		if err.Error() != "category is math" {
			t.Errorf("expected error \"category is math\", received %q", err)
		}
	} else {
		t.Errorf("expected a non-nil error and a nil question, received a(n) %q error and a(n) %q question", err, got)
	}
}
