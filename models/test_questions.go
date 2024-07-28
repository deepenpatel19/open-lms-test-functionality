package models

type TestQuestionsSchema struct {
	Id           int64                             `json:"id"`
	TestId       int64                             `json:"test_id"`
	QuestionData QuestionResponseSchemaForTestTake `json:"question_data"`
}
