package worker

type MockWorker struct {
	ReturnChannel chan bool
}

func (mw *MockWorker) Start() {
	mw.ReturnChannel <- true
}
