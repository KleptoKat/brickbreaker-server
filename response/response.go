package response

type Response struct {
	Code int `json:"code"`
	Description string `json:"description,omitempty"`
	Data interface {} `json:"data,omitempty"`
}

type ErrorData struct {
	Result string `json:"result"`
}

func OKWithData(data interface{}) Response {
	return Response {
		Code:200,
		Data:data,
	}
}

func OK() Response {
	return Response {
		Code:200,
	}
}

func InternalError() Response {
	return Response {
		Code:500,
	}
}

func BadRequest() Response {
	return Response {
		Code:400,
	}
}

func BadRequestWithDescription(desc string) Response {
	return Response {
		Code:400,
		Description:desc,
	}
}

func BadRequestWithError(err error) Response {
	return Response {
		Code:400,
		Data:ErrorData{
			Result:err.Error(),
		},
	}
}