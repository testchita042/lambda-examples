package factorial

func Factorial(input int) (map[string]interface{}, error) {

	return map[string]interface{}{
		"result": factorial(input),
	}, nil

}
