package lambda_with_env

type Req struct {
	IP string `json:"__lmbd_ip"`
	Bearer string `json:"__lmbd_bearer"`
	Name string `json:"__lmbd_name"`
	Method string `json:"__lmbd_method"`
	Path string `json:"__lmbd_path"`
	BasePath string `json:"__lmbd_base_path"`
	Payload []byte `json:"__lmbd_payload"`
	PayloadContentType string `json:"__lmbd_payload_content_type"`
	Env map[string]any `json:"__lmbd_env"`
}

type Resp struct {
	IP string
	Bearer string
	Name string
	Method string
	Path string
	BasePath string
	Payload []byte
	PayloadContentType string
	Env map[string]any
}

func LambdaWithEnv(req Req) Resp {
	return Resp{
IP: req.IP,
Bearer: req.Bearer,
Name: req.Name,
Method: req.Method,
Path: req.Path,
BasePath: req.BasePath,
Payload: req.Payload,
PayloadContentType: req.PayloadContentType,
Env: req.Env,
}	
}
