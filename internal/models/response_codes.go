package models

// depending on approach this could be a shared file in a schema repo which is nice when multiple services
// can have the same definition of a response an example its useful is if we have a validation engine
// it can be a seperate service fail the request for a special reason and give a response

type ResponseCode int

type Response struct {
	Code ResponseCode
}

const (
	Approved ResponseCode = 1000
	Declined ResponseCode = 2000
)

func (r *Response) AsInt() int {
	return int(r.Code)
}
