package go_tenable

type IOError interface {
	error
}

type SCError interface {
	error
}

//type APIError struct {
//	StatusCode	int
//	Description	string
//}
