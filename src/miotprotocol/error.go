package miotprotocol

const (
	errorCodeInvalidateControlOrder   string = "INVALIDATE_CONTROL_ORDER"
	errorCodeServiceError             string = "SERVICE_ERROR"
	errorCodeDeviceNotSupportFunction string = "DEVICE_NOT_SUPPORT_FUNCTION"
	errorCodeInvalidateParams         string = "INVALIDATE_PARAMS"
	errorCodeDeviceIsNotExist         string = "DEVICE_IS_NOT_EXIST"
	errorCodeIOTDeviceOffline         string = "IOT_DEVICE_OFFLINE"
	errorCodeAccessTokenInvalidate    string = "ACCESS_TOKEN_INVALIDATE"
)

const (
	errorMsgInvalidateControlOrder   string = "invalidate control order"
	errorMsgDeviceNotSupportFunction string = "device not support"
	errorMsgInvalidateParams         string = "invalidate params"
	errorMsgDeviceIsNotExist         string = "device is not exist"
	errorMsgIOTDeviceOffline         string = "device is offline"
	errorMsgAccessTokenInvalidate    string = "access_token is invalidate"
)

func (p *Processor) InternalError(errorcode string, errormes string) interface{} {
	resp := make(map[string]string)
	resp["code"] = "-101"
	resp["description"] = errorcode + "\n" + errormes
	return resp
}

//func (p *Processor) newErrorResp(errorCode, errorMsg string) *Protocol {
//	resp := &Protocol{Header: p.Header}
//	resp.Header.Name = "ErrorResponse"
//
//	resp.Payload = make(map[string]interface{})
//	resp.Payload["deviceId"] = p.Payload["deviceId"]
//	resp.Payload["errorCode"] = errorCode
//	resp.Payload["message"] = errorMsg
//
//	return resp
//}
