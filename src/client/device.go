package client

import (
	"strconv"
	"time"
)

type Device struct {
	DeviceID           int              `json:"deviceId"`
	ServiceID          int              `json:"serviceId"`
	BrandName          string           `json:"brandName"`
	CategoryName       string           `json:"categoryName"`
	DeviceProductName  string           `json:"deviceProductName"`
	DeviceDescription  string           `json:"deviceDescription"`
	DeviceName         string           `json:"deviceName"`
	DeviceLocation     string           `json:"deviceLocation"`
	DeviceState        string           `json:"deviceState"`
	DeviceIcon         string           `json:"deviceIcon"`
	ServiceType        string           `json:"serviceType"`
	ServiceDescription string           `json:"serviceDescription"`
	ServiceName        string           `json:"serviceName"`
	ServiceState       string           `json:"serviceState"`
	Characteristics    []Characteristic `json:"characteristics"`
}

func (d *Device) GetType() []string {
	types := make([]string, 0)
	switch d.ServiceType {
	// case "FAN":
	// 	return "fan"
	case "LIGHTING":
		types = append(types, "LIGHT")
		return types
	case "CURTAIN":
		types = append(types, "CURTAIN")
		return types
	case "OUTLET":
		types = append(types, "SOCKET")
		return types
	case "LIGHTINGGROUP":
		switch d.CategoryName {
		case "LIGHT":
			types = append(types, "LIGHT")
			return types
		case "CURTAIN":
			types = append(types, "CURTAIN")
			return types
		case "OUTLET":
			types = append(types, "SOCKET")
			return types
		default:
			return types
		}
	default:
		return types
	}
}

func (d *Device) GetAttributesAndActions() (
	attributes []map[string]interface{}, actions []string) {

	attributes = make([]map[string]interface{}, 0)
	actions = make([]string, 0)

	name := make(map[string]interface{})
	name["name"] = "name"
	name["value"] = d.DeviceName
	name["scale"] = ""
	name["timestampOfSample"] = time.Now().Unix()
	name["uncertaintyInMilliseconds"] = 0
	attributes = append(attributes, name)

	connectivity := make(map[string]interface{})
	connectivity["name"] = "connectivity"
	connectivity["value"] = "REACHABLE"
	connectivity["scale"] = ""
	connectivity["timestampOfSample"] = time.Now().Unix()
	connectivity["uncertaintyInMilliseconds"] = 0
	attributes = append(attributes, connectivity)

	for _, chara := range d.Characteristics {
		attribute := make(map[string]interface{})
		switch chara.CharacteristicName {
		case "BRIGHTNESS":
			attribute["name"] = "brightness"
			attribute["value"], _ = strconv.Atoi(chara.CharacteristicValue)
			attribute["value"] = attribute["value"].(int) / 254 * 100
			attribute["scale"] = ""
			attribute["timestampOfSample"] = time.Now().Unix()
			attribute["uncertaintyInMilliseconds"] = 0
			if d.ServiceType == "LIGHTINGGROUP" {
				actions = append(actions,
					"setBrightnessPercentage")
			} else {
				actions = append(actions,
					"setBrightnessPercentage", "incrementBrightnessPercentage", "decrementBrightnessPercentage")
			}
		case "POWER":
			attribute["name"] = "turnOnState"
			if chara.CharacteristicValue == "true" {
				attribute["value"] = "ON"
			} else {
				attribute["value"] = "OFF"
			}
			attribute["scale"] = ""
			attribute["timestampOfSample"] = time.Now().Unix()
			attribute["uncertaintyInMilliseconds"] = 0
			//attribute["value"] = (chara.CharacteristicValue == "true")
			actions = append(actions,
				"turnOn", "turnOff")
		default:
			continue
		}
		attributes = append(attributes, attribute)
	}
	return attributes, actions

}

func (d *Device) CanUse() (ok bool) {

	if len(d.GetType()) == 0 {
		return false
	}

	ok = false
	for _, chara := range d.Characteristics {
		if chara.CharacteristicName == "WINDSPEED" {
			ok = true
			break
		}
		if chara.CharacteristicName == "BRIGHTNESS" {
			ok = true
			break
		}
		if chara.CharacteristicName == "POWER" {
			ok = true
			break
		}
		if chara.CharacteristicName == "CURTAINSTATE" {
			ok = true
			break
		}
		if chara.CharacteristicName == "COLORTEMPERATURE" {
			ok = true
			break
		}
	}
	if !ok {
		return ok
	}

	return true

}
func (d *Device) GetCharacteristicsByName(name string) string {
	for _, feature := range d.Characteristics {
		if feature.CharacteristicName == name {
			return feature.CharacteristicValue
		}
	}
	return "Can not find value"
}