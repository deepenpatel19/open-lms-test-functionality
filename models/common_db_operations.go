package models

type QueryStructToExecute struct {
	Query     string
	QueryList []string
}

type Operations interface {
}

func (query QueryStructToExecute) InsertOperations(uuidString string) {

}

func (query QueryStructToExecute) UpdateOperation(uuidString string) {

}

func (query QueryStructToExecute) DeleteOperation(uuidString string) {

}
