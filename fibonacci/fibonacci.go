package fibonacci

type Result struct {
	Input  int `json:"input"`
	Result int `json:"result"`
}

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func Fibonacci(input int) Result {

	return Result{
		Input:  input,
		Result: fibonacci(input),
	}
}
