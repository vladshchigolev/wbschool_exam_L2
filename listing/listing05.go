package main

type customError struct {
	msg string
}
// Теперь тип *customError реализует интерфейс error
func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
}

func main() {
	var err error
	err = test() // {type:*customError, val: nil}
	if err.Error() != nil {
		println("error")
		return
	}
	println("ok")
}

