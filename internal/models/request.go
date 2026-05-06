package models

type HTTPMethod string

const (
	GET     HTTPMethod = "GET"
	POST    HTTPMethod = "POST"
	PUT     HTTPMethod = "PUT"
	PATCH   HTTPMethod = "PATCH"
	DELETE  HTTPMethod = "DELETE"
	HEAD    HTTPMethod = "HEAD"
	OPTIONS HTTPMethod = "OPTIONS"
)

var AllMethods = []HTTPMethod{GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS}

type AuthType string

const (
	AuthNone   AuthType = "none"
	AuthBearer AuthType = "bearer"
	AuthBasic  AuthType = "basic"
)

type Header struct {
	Key     string
	Value   string
	Enabled bool
}

type Param struct {
	Key     string
	Value   string
	Enabled bool
}

type Auth struct {
	Type     AuthType
	Token    string
	Username string
	Password string
}

type Request struct {
	Method  HTTPMethod
	URL     string
	Headers []Header
	Params  []Param
	Body    string
	Auth    Auth
}

type SavedRequest struct {
	Name      string
	CreatedAt int64
	UpdatedAt int64
	Request   Request
}
