package internal

type APIResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type DataControl struct {
	Next     string `json:"next"`
	Previous string `json:"prev"`
}

type Meta struct {
	Count int `json:"count"`
}

type ErrorResponse struct {
	APIResponse
}

func NewErrorResponse(message string) ErrorResponse {
	return ErrorResponse{
		APIResponse: APIResponse{
			Status:  false,
			Message: message,
		},
	}
}

type SuccessResponse struct {
	APIResponse
	Data interface{} `json:"data,omitempty"`
}

type SuccessResponseOption func(*SuccessResponse)

func WithData(data interface{}) SuccessResponseOption {
	return func(r *SuccessResponse) {
		r.Data = data
	}
}

func NewSuccessResponse(message string, opts ...SuccessResponseOption) *SuccessResponse {
	response := &SuccessResponse{
		APIResponse: APIResponse{
			Status:  true,
			Message: message,
		},
	}
	for _, opt := range opts {
		opt(response)
	}

	return response
}

type ValidationError struct {
	Placement  string `json:"placement"`
	Detail     string `json:"detail"`
	Field      string `json:"field"`
	Code       string `json:"code"`
	Expression string `json:"expression"`
	Parameter  string `json:"parameter"`
	TraceId    string `json:"trace_id"`
}

type ValidationErrorResponse struct {
	APIResponse
	Errors []ValidationError `json:"errors"`
}

// "placement": "field",
// "title": "value too short",
// "detail": "field username must be at least 4 symbols",
// "location": "username",
// "field": "username",
// "code": "validation.min",
// "expression": "min",
// "argument": "4",
// "traceid": "74681b27-b1ea-454d-9847-d27059e19119",
