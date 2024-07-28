package schemas

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
