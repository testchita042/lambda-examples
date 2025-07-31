package factorial

// Esta función hace la formula matematica Factorial
func factorial(n int) int {
	if n == 0 {
		return 1
	}
	return n * factorial(n-1)
}
