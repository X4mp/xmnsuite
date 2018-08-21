package router

type header struct {
	statusCode int
	lines      map[string]string
}

func createHeader(statusCode int, lines map[string]string) Header {
	out := header{
		statusCode: statusCode,
		lines:      lines,
	}

	return &out
}

// StatusCode returns the status code
func (obj *header) StatusCode() int {
	return obj.statusCode
}

// Lines returns the header lines
func (obj *header) Lines() map[string]string {
	return obj.lines
}
