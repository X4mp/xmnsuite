package router

type response struct {
	header Header
	body   []byte
}

func createResponse(header Header, body []byte) Response {
	out := response{
		header: header,
		body:   body,
	}

	return &out
}

// Header returns the header
func (obj *response) Header() Header {
	return obj.header
}

// Body returns the body
func (obj *response) Body() []byte {
	return obj.body
}
