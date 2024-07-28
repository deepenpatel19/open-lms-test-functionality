package models

const (
	TRUEORFALSE    int = 0
	MULTIPLECHOICE int = 1
)

func ValidateQuestionType(questionType int) string {
	if questionType == MULTIPLECHOICE {
		return "multiple_choice"
	} else if questionType == TRUEORFALSE {
		return "true_or_false"
	} else {
		return ""
	}
}
