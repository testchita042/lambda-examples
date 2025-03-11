package helloworld

type Response struct {
	Message string `json:"message"`
}

func HelloWorld() Response {
	return Response{
		Message: "Hello, World!",
	}
}
