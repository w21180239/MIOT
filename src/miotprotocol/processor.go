package miotprotocol

import (
	"client"
	"constant"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strconv"
	"time"
)

// Processor class can process a single Protocol
type Processor struct {
	Protocol map[string]interface{}
	Token    string
	Client   *client.Client
}

// NewProcessor Creates a new processor instance
func NewProcessor(r *http.Request) *Processor {
	var p map[string]interface{}
	json.NewDecoder(r.Body).Decode(&p)

	return &Processor{
		Protocol: p,
		Token:    r.Header["User-Token"][0],
		Client:   nil,
	}
}

// Process process the protocol and return the response
func (p *Processor) Process() interface{} {
	userID, projectID, valid := p.extractToken()
	fmt.Println(userID, projectID, valid)
	if !valid {
		return p.InternalError(errorCodeAccessTokenInvalidate,
			errorMsgAccessTokenInvalidate)
	}
	p.Client = client.NewClient(userID, projectID)
	var resp interface{}
	switch p.Protocol["intent"] {
	case "get-devices":
		resp = p.GetDevices()
	case "get-properties":
		resp = p.GetProperties()
	case "set-properties":
		resp = p.SetProperties()
	case "invoke-action":
		resp = p.InvokeAction()
	case "subscribe":
		resp = p.Subscribe()
	case "unsubscribe":
		resp = p.Unsubscribe()
	case "get-device-status":
		resp = p.GetDeviceStatus()
	}

	if constant.Output_debug_information {
		fmt.Printf("[DEBUG INFO]UserID=%v,ProjectID=%v\n", userID, projectID)
		fmt.Println("[DEBUG INFO]Input data begin")
		jsonStr, _ := json.MarshalIndent(p.Protocol, "", "  ")
		fmt.Println(string(jsonStr))
		fmt.Println("[DEBUG INFO]Input data end")
		jsonStr, _ = json.MarshalIndent(resp, "", "  ")
		fmt.Println("[DEBUG INFO]Output data start")
		fmt.Println(string(jsonStr))
		fmt.Println("[DEBUG INFO]Output data end")
	}
	return resp
}
func (p *Processor) GetDevices() interface{} {
	resp := make(map[string]interface{})
	resp["requestId"] = p.Protocol["requestId"]
	resp["intent"] = p.Protocol["intent"]
	devices, err := p.Client.GetDevices()
	if err != nil {
		return p.InternalError(errorCodeServiceError,
			"Fail to get devices:"+err.Error())
	}
	resp["devices"] = make([]map[string]interface{}, 0)
	for _, device := range devices {
		if len(device.GetType()) == 0 {
			continue
		} else {
			deviceData := make(map[string]interface{})
			if device.ServiceType == "LIGHTINGGROUP" {
				deviceData["did"] = "LIGHTINGGROUP" + strconv.Itoa(device.DeviceID)
				//continue
			} else {
				deviceData["did"] = strconv.Itoa(device.DeviceID)
			}
			if device.DeviceProductName == "可调光调色LED灯Color版" {
				deviceData["type"] = "urn:miot-spec-v2:device:light:0000A001:lico-test:1:0000C802"
			} else {
				deviceData["type"] = "urn:miot-spec-v2:device:light:0000A001:lico-0000:1:0000C801"
			}
			deviceData["name"] = device.DeviceName
			if device.ServiceName != "" {
				deviceData["name"] = deviceData["name"].(string) + "(" + device.ServiceName + ")"
			}
			if device.ServiceState == "正常" && device.BrandName != "小米" && device.ServiceType != "LIGHTINGGROUP" {
				resp["devices"] = append(
					resp["devices"].([]map[string]interface{}),
					deviceData,
				)
			}
		}
	}
	return resp
}
func (p *Processor) GetProperties() interface{} {
	resp := p.Protocol
	for i := 0; i < len(resp["properties"].([]interface{})); i++ {
		deviceID, _ := strconv.Atoi(resp["properties"].([]interface{})[i].(map[string]interface{})["did"].(string))
		device, _ := p.Client.GetDeviceByID(deviceID)
		tmp_pro := resp["properties"].([]interface{})[i].(map[string]interface{})
		siid := int(tmp_pro["siid"].(float64))
		piid := int(tmp_pro["piid"].(float64))
		switch device.ServiceType {
		case "LIGHTING":
			switch siid {
			case 1:
				switch piid {
				case 1:
					tmp_pro["value"] = device.BrandName
				default:
					tmp_pro["value"] = "undefined"
				}
			case 2:
				switch piid {
				case 1:
					tmp_pro["value"] = device.GetCharacteristicsByName("POWER")
				case 2:
					value, _ := strconv.Atoi(device.GetCharacteristicsByName("BRIGHTNESS"))
					float_value := float64(value)
					value = int((float_value-1)/(254-1)*(100-1)+1) + 1
					tmp_pro["value"] = uint8(value)
				case 3:
					value, _ := strconv.Atoi(device.GetCharacteristicsByName("COLORTEMPERATURE"))
					float_value := float64(value)
					value = int((float_value-3000)/(6000-3000)*(20000-800) + 800)
					tmp_pro["value"] = uint32(value)
				}
			}
		}

		resp["properties"].([]interface{})[i].(map[string]interface{})["status"] = 0
	}
	return resp
}
func (p *Processor) SetProperties() interface{} {
	resp := p.Protocol
	for i := 0; i < len(resp["properties"].([]interface{})); i++ {
		action := ""
		var value int
		deviceID, _ := strconv.Atoi(resp["properties"].([]interface{})[i].(map[string]interface{})["did"].(string))
		device, _ := p.Client.GetDeviceByID(deviceID)
		tmp_pro := resp["properties"].([]interface{})[i].(map[string]interface{})
		siid := int(tmp_pro["siid"].(float64))
		piid := int(tmp_pro["piid"].(float64))
		if device.ServiceType == "LIGHTING" {
			switch siid {
			case 2:
				switch piid {
				case 1:
					if tmp_pro["value"].(bool) {
						action = "TURNONLIGHT"
					} else {
						action = "TURNOFFLIGHT"
					}
				case 2:
					action = "BRIGHTNESSSET"
					tmp_value := tmp_pro["value"].(float64)
					value = int((tmp_value-1)/(100-1)*(254-1) + 1)
				case 3:
					action = "COLORTEMPERATURESET"
					tmp_value := tmp_pro["value"].(float64)
					value = int((tmp_value-800)/(20000-800)*(6000-3000) + 3000)
				}

			}
		}
		err := p.Client.PostAction(deviceID, action, value)
		if err != nil {
			return p.InternalError(errorCodeServiceError,
				err.Error())
		}
		resp["properties"].([]interface{})[i].(map[string]interface{})["status"] = 0
	}
	return resp
}
func (p *Processor) InvokeAction() interface{} {
	resp := p.Protocol
	action := ""
	var value int
	deviceID, _ := strconv.Atoi(resp["action"].(map[string]interface{})["did"].(string))
	device, _ := p.Client.GetDeviceByID(deviceID)
	tmp_pro := resp["action"].(map[string]interface{})
	aiid := int(tmp_pro["aiid"].(float64))
	if device.ServiceType == "LIGHTING" {
		switch aiid {
		case 1:
			if tmp_pro["in"].([]interface{})[0].(bool) {
				action = "TURNONLIGHT"
			} else {
				action = "TURNOFFLIGHT"
			}
			resp["action"].(map[string]interface{})["out"] = tmp_pro["in"]
		case 2:
			action = "BRIGHTNESSSET"
			tmp_value := tmp_pro["in"].([]interface{})[0].(float64) + 1
			value = int((tmp_value-1)/(100-1)*(254-1) + 1)
			var out = []float64{tmp_value}
			resp["action"].(map[string]interface{})["out"] = out
		case 3:
			action = "BRIGHTNESSSET"
			tmp_value := tmp_pro["in"].([]interface{})[0].(float64) - 1
			value = int((tmp_value-1)/(100-1)*(254-1) + 1)
			var out = []float64{tmp_value}
			resp["action"].(map[string]interface{})["out"] = out
		}
	}
	err := p.Client.PostAction(deviceID, action, value)
	if err != nil {
		return p.InternalError(errorCodeServiceError,
			err.Error())
	}
	return resp
}
func (p *Processor) Subscribe() interface{} {
	resp := p.Protocol
	for i := 0; i < len(resp["devices"].([]interface{})); i++ {
		resp["devices"].([]interface{})[i].(map[string]interface{})["status"] = 0
	}
	return resp
}
func (p *Processor) Unsubscribe() interface{} {
	resp := p.Protocol
	for i := 0; i < len(resp["devices"].([]interface{})); i++ {
		resp["devices"].([]interface{})[i].(map[string]interface{})["status"] = 0
	}
	return resp
}
func (p *Processor) GetDeviceStatus() interface{} { return nil }

//func (p *Processor) processDiscovery() *Protocol {
//	resp := &Protocol{Header: p.Header, Payload: make(map[string]interface{})}
//	resp.Header.Name = "DiscoverAppliancesResponse"
//	resp.Header.PayloadVersion = "1"
//
//	devices, err := p.Client.GetDevices()
//	if err != nil {
//		return p.newErrorResp(errorCodeServiceError,
//			"Fail to get devices:"+err.Error())
//	}
//
//	resp.Payload["discoveredAppliances"] = make([]map[string]interface{}, 0)
//	for _, device := range devices {
//		/*if !device.CanUse() {
//			continue
//		}*/
//		fmt.Println(device)
//		if len(device.GetType()) == 0 {
//			continue
//		} else {
//			deviceData := make(map[string]interface{})
//			if device.ServiceType == "LIGHTINGGROUP" {
//				deviceData["applianceId"] = "LIGHTINGGROUP" + strconv.Itoa(device.DeviceID)
//				//continue
//			} else {
//				deviceData["applianceId"] = strconv.Itoa(device.DeviceID)
//			}
//			deviceData["version"] = "1"
//			deviceData["applianceTypes"] = device.GetType()
//			fmt.Println(device.GetType(), "test")
//			deviceData["friendlyName"] = device.DeviceName
//			deviceData["manufacturerName"] = device.BrandName
//			deviceData["modelName"] = device.DeviceProductName
//			deviceData["friendlyDescription"] = "Licotek"
//			deviceData["isReachable"] = true
//			deviceData["attributes"], deviceData["actions"] = device.GetAttributesAndActions()
//			deviceData["additionalApplianceDetails"] = map[string]interface{}{"deviceType": device.GetType()[0]}
//
//			resp.Payload["discoveredAppliances"] = append(
//				resp.Payload["discoveredAppliances"].([]map[string]interface{}),
//				deviceData,
//			)
//		}
//	}
//
//	scenes, err := p.Client.GetScenes()
//	for _, scene := range scenes {
//		deviceData := make(map[string]interface{})
//		deviceData["applianceId"] = strconv.Itoa(scene.SceneID)
//		deviceData["version"] = "1"
//		deviceData["applianceTypes"] = []string{"SCENE_TRIGGER"}
//		deviceData["friendlyName"] = scene.SceneName
//		deviceData["manufacturerName"] = "Licotek"
//		deviceData["modelName"] = "Licotek"
//		deviceData["friendlyDescription"] = "Licotek"
//		deviceData["isReachable"] = true
//		deviceData["attributes"], deviceData["actions"] = scene.GetAttributesAndActions()
//		deviceData["additionalApplianceDetails"] = map[string]interface{}{"deviceType": "SCENE_TRIGGER"}
//
//		resp.Payload["discoveredAppliances"] = append(
//			resp.Payload["discoveredAppliances"].([]map[string]interface{}),
//			deviceData,
//		)
//	}
//
//	rootspace, err := p.Client.GetSpaces()
//	fmt.Println(rootspace)
//	if err != nil {
//		return p.newErrorResp(errorCodeServiceError,
//			"Fail to get spaces:"+err.Error())
//	}
//	spaces := rootspace.GetAllSpace()
//	tempspace := []client.Space{rootspace}
//	spaces = append(tempspace, spaces...)
//	if len(spaces) > 10 {
//		spaces = spaces[0:10]
//	}
//	fmt.Println(spaces)
//	resp.Payload["discoveredGroups"] = make([]map[string]interface{}, 0)
//
//	for _, space := range spaces {
//		spaceData := make(map[string]interface{})
//		spaceData["groupName"] = space.SpaceName
//		spaceData["groupNotes"] = space.SpaceName
//		spaceData["additionalApplianceDetails"] = map[string]interface{}{"spaceId": space.SpaceID}
//		spaceData["applianceIds"] = make([]string, 0)
//		spaceDetail, _ := p.Client.GetSpaceDevices(space.SpaceID)
//		for _, spaceDevice := range spaceDetail.Devices {
//			spaceData["applianceIds"] = append(spaceData["applianceIds"].([]string), strconv.Itoa(spaceDevice.DeviceID))
//		}
//		resp.Payload["discoveredGroups"] = append(
//			resp.Payload["discoveredGroups"].([]map[string]interface{}),
//			spaceData,
//		)
//	}
//
//	return resp
//}
//
//func (p *Processor) processLightingGroup(controlType string) *Protocol {
//	appliance := p.Payload["appliance"].(map[string]interface{})
//	lightingGroupID, err := strconv.Atoi(appliance["applianceId"].(string)[13:])
//	if err != nil {
//		return p.newErrorResp(errorCodeServiceError, err.Error())
//	}
//	characteristic := ""
//	var value int
//	var postValue string
//	if controlType == "LIGHT" {
//		switch p.Header.Name {
//		case "TurnOnRequest":
//			characteristic = "POWER"
//			postValue = "true"
//		case "TurnOffRequest":
//			characteristic = "POWER"
//			postValue = "false"
//		case "SetBrightnessPercentageRequest":
//			if p.Payload["brightness"].(map[string]interface{})["value"].(float64) == 100 {
//				characteristic = "BRIGHTNESS"
//				postValue = "254"
//			} else if p.Payload["brightness"].(map[string]interface{})["value"].(float64) == 0 {
//				characteristic = "BRIGHTNESS"
//				postValue = "1"
//			} else {
//				characteristic = "BRIGHTNESS"
//				value = int(p.Payload["brightness"].(map[string]interface{})["value"].(float64) / 100 * 254)
//				if value < 50 {
//					postValue = "50"
//				} else {
//					postValue = strconv.Itoa(value)
//				}
//			}
//		default:
//			return p.newErrorResp(errorCodeDeviceNotSupportFunction,
//				errorMsgDeviceNotSupportFunction)
//		}
//	} else if controlType == "CURTAIN" {
//		switch p.Header.Name {
//		case "TurnOnRequest":
//			characteristic = "POWER"
//			postValue = "true"
//		case "TurnOffRequest":
//			characteristic = "POWER"
//			postValue = "false"
//		default:
//			return p.newErrorResp(errorCodeDeviceNotSupportFunction,
//				errorMsgDeviceNotSupportFunction)
//		}
//
//	} else if controlType == "SOCKET" {
//		switch p.Header.Name {
//		case "TurnOnRequest":
//			characteristic = "POWER"
//			postValue = "true"
//		case "TurnOffRequest":
//			characteristic = "POWER"
//			postValue = "false"
//		default:
//			return p.newErrorResp(errorCodeDeviceNotSupportFunction,
//				errorMsgDeviceNotSupportFunction)
//		}
//	}
//	if err != nil {
//		return p.newErrorResp(errorCodeServiceError, err.Error())
//	}
//	fmt.Println(value)
//	err = p.Client.PostLightingGroupControl(lightingGroupID,
//		characteristic, postValue)
//	if err != nil {
//		return p.newErrorResp(errorCodeServiceError,
//			err.Error())
//	}
//	resp := &Protocol{Header: p.Header, Payload: make(map[string]interface{})}
//	resp.Header.Name = (resp.Header.Name)[0:len(resp.Header.Name)-7] + "Confirmation"
//	resp.Header.PayloadVersion = "1"
//	switch p.Header.Name {
//	case "SetBrightnessPercentageRequest":
//		fmt.Println(value)
//		resp.Payload["brightness"] = map[string]interface{}{"value": value / 254 * 100}
//		brightness := make(map[string]interface{})
//		brightness["value"] = 0
//		resp.Payload["previousState"] = map[string]interface{}{"brightness": brightness}
//	default:
//	}
//	return resp
//}
//
//func (p *Processor) processControl() *Protocol {
//
//	appliance := p.Payload["appliance"].(map[string]interface{})
//
//	controlType := appliance["additionalApplianceDetails"].(map[string]interface{})["deviceType"].(string)
//
//	if controlType == "SCENE_TRIGGER" {
//		sceneID, err := strconv.Atoi(appliance["applianceId"].(string))
//		if err != nil {
//			return p.newErrorResp(errorCodeServiceError, err.Error())
//		}
//		switch p.Header.Name {
//		case "TurnOnRequest":
//			err = p.Client.PostScene(sceneID)
//		default:
//		}
//		resp := &Protocol{Header: p.Header, Payload: make(map[string]interface{})}
//		resp.Header.Name = (resp.Header.Name)[0:len(resp.Header.Name)-7] + "Confirmation"
//		resp.Header.PayloadVersion = "1"
//		return resp
//	}
//
//	if len(appliance["applianceId"].(string)) > 13 {
//		return p.processLightingGroup(controlType)
//	}
//
//	deviceID, err := strconv.Atoi(appliance["applianceId"].(string))
//	if err != nil {
//		return p.newErrorResp(errorCodeServiceError, err.Error())
//	}
//
//	device, err := p.Client.GetDeviceByID(deviceID)
//	attributes, _ := device.GetAttributesAndActions()
//	if err != nil {
//		return p.newErrorResp(errorCodeDeviceIsNotExist, err.Error())
//	}
//
//	action := ""
//	var value int
//
//	if controlType == "LIGHT" {
//		switch p.Header.Name {
//		case "TurnOnRequest":
//			action = "TURNONLIGHT"
//		case "TurnOffRequest":
//			action = "TURNOFFLIGHT"
//		case "SetBrightnessPercentageRequest":
//			if p.Payload["brightness"].(map[string]interface{})["value"].(float64) == 100 {
//				action = "BRIGHTNESSMAX"
//			} else if p.Payload["brightness"].(map[string]interface{})["value"].(float64) == 0 {
//				action = "BRIGHTNESSMIN"
//			} else {
//				action = "BRIGHTNESSSET"
//				value = int(p.Payload["brightness"].(map[string]interface{})["value"].(float64) / 100 * 254)
//				if value < 50 {
//					value = 50
//				}
//			}
//		case "IncrementBrightnessPercentageRequest":
//			action = "BRIGHTNESSUP"
//			value = int(p.Payload["deltaPercentage"].(map[string]interface{})["value"].(float64) / 100 * 254)
//			fmt.Println(value)
//			if value == 50 {
//				action = "BRIGHTNESSUPABIT"
//			}
//		case "DecrementBrightnessPercentageRequest":
//			action = "BRIGHTNESSDOWN"
//			value = int(p.Payload["deltaPercentage"].(map[string]interface{})["value"].(float64) / 100 * 254)
//			fmt.Println(value)
//			if value == 50 {
//				action = "BRIGHTNESSDOWNABIT"
//			}
//		default:
//			return p.newErrorResp(errorCodeDeviceNotSupportFunction,
//				errorMsgDeviceNotSupportFunction)
//		}
//	} else if controlType == "CURTAIN" {
//		switch p.Header.Name {
//		case "TurnOnRequest":
//			action = "TURNONCURTAIN"
//		case "TurnOffRequest":
//			action = "TURNOFFCURTAIN"
//		default:
//			return p.newErrorResp(errorCodeDeviceNotSupportFunction,
//				errorMsgDeviceNotSupportFunction)
//		}
//
//	} else if controlType == "SOCKET" {
//		switch p.Header.Name {
//		case "TurnOnRequest":
//			action = "TURNONOUTLET"
//		case "TurnOffRequest":
//			action = "TURNOFFOUTLET"
//		default:
//			return p.newErrorResp(errorCodeDeviceNotSupportFunction,
//				errorMsgDeviceNotSupportFunction)
//		}
//	}
//	if err != nil {
//		return p.newErrorResp(errorCodeServiceError, err.Error())
//	}
//	fmt.Println(value)
//	err = p.Client.PostAction(deviceID, action, value)
//	if err != nil {
//		return p.newErrorResp(errorCodeServiceError,
//			err.Error())
//	}
//	resp := &Protocol{Header: p.Header, Payload: make(map[string]interface{})}
//	resp.Header.Name = (resp.Header.Name)[0:len(resp.Header.Name)-7] + "Confirmation"
//	resp.Header.PayloadVersion = "1"
//	device, _ = p.Client.GetDeviceByID(deviceID)
//	resp.Payload["attributes"], _ = device.GetAttributesAndActions()
//
//	switch p.Header.Name {
//	case "SetBrightnessPercentageRequest":
//		fmt.Println(value)
//		resp.Payload["brightness"] = map[string]interface{}{"value": value / 254 * 100}
//		brightness := make(map[string]interface{})
//		for _, attribute := range attributes {
//			if attribute["name"] == "brightness" {
//				brightness["value"] = attribute["value"]
//			}
//		}
//		resp.Payload["previousState"] = map[string]interface{}{"brightness": brightness}
//	case "IncrementBrightnessPercentageRequest":
//		var brightness interface{}
//		for _, attribute := range attributes {
//			if attribute["name"] == "brightness" {
//				brightness = attribute["value"]
//			}
//		}
//		resp.Payload["previousState"] = map[string]interface{}{"brightness": brightness}
//		for _, attribute := range resp.Payload["attributes"].([]map[string]interface{}) {
//			if attribute["name"] == "brightness" {
//				brightness = attribute["value"]
//			}
//		}
//		resp.Payload["brightness"] = brightness
//	case "DecrementBrightnessPercentageRequest":
//		var brightness interface{}
//		for _, attribute := range attributes {
//			if attribute["name"] == "brightness" {
//				brightness = attribute["value"]
//			}
//		}
//		resp.Payload["previousState"] = map[string]interface{}{"brightness": brightness}
//		for _, attribute := range resp.Payload["attributes"].([]map[string]interface{}) {
//			if attribute["name"] == "brightness" {
//				brightness = attribute["value"]
//			}
//		}
//		resp.Payload["brightness"] = brightness
//	default:
//	}
//	return resp
//}

func (p *Processor) extractToken() (string, string, bool) {
	secret := []byte("secret")
	token, _ := jwt.Parse(p.Token,
		func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})
	claims, _ := token.Claims.(jwt.MapClaims)
	userID := claims["aud"].(string)
	projectID := claims["project_id"].(string)
	expiresIn := int64(claims["exp"].(float64))
	if time.Unix(expiresIn, 0).After(time.Now()) {
		return userID, projectID, true
	}
	return userID, projectID, false

}
