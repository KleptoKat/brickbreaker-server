package response

type Response struct {
	Code int `json:"code"`
	Data interface {} `data:"data"`
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