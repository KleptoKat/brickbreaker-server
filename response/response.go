package response

type Response struct {
	Code int `json:"code"`
	Data interface {} `json:"data"`
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

func BadRequestWithError(err error) Response {
	return Response {
		Code:400,
		Data:ErrorData{
			Result:err.Error(),
		},
	}
}