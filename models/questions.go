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

type QuestionCreateSchema struct {
	Type         string `json:"type"`
	QuestionData string `json:"question_data"`
	AnswerData   string `json:"answer_data"`
}

type QuestionResponseSchemaForTestTake struct {
	Type         string `json:"type"`
	QuestionData string `json:"question_data"`
}

type QuestionResponseSchema struct {
	Id           int64  `json:"id"`
	Type         string `json:"type"`
	QuestionData string `json:"question_data"`
	AnswerData   string `json:"answer_data"`
}
