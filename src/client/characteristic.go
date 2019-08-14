package client

import "errors"

type Characteristic struct {
	ID                           int    `json:"id"`
	CharacteristicName           string `json:"characteristicName"`
	CharacteristicType           string `json:"characteristicType"`
	CharacteristicUnit           string `json:"characteristicUnit"`
	CharacteristicControlDisplay string `json:"characteristicControlDisplay"`
	CharacteristicDescription    string `json:"characteristicDescription"`
	CharacteristicParameter1     string `json:"characteristicParameter1"`
	CharacteristicParameter2     string `json:"characteristicParameter2"`
	CharacteristicParameter3     string `json:"characteristicParameter3"`
	CharacteristicValue          string `json:"characteristicValue"`
	CharacteristicState          string `json:"characteristicState"`
}

func (c *Characteristic) AliPropertyName() (string, error) {

	var name string
	switch c.CharacteristicName {
	case "POWER":
		name = "powerstate"
	case "BRIGHTNESS":
		name = "brightness"
	case "COLORTEMPERATURE":
		name = "colorTemperature"
	default:
		name = ""
	}
	if name == "" {
		return name, errors.New("unsupported name")
	}
	return name, nil

}

func (c *Characteristic) AliPropertyValue() interface{} {

	var value interface{}
	switch c.CharacteristicName {
	case "POWER":
		if c.CharacteristicValue == "true" {
			value = "on"
		} else {
			value = "off"
		}
	default:
		value = c.CharacteristicValue
	}
	return value

}

func AliPropertyName2CharaName(name string) (string, error) {

	var result string
	switch name {
	case "powerstate":
		result = "POWER"
	case "brightness":
		result = "BRIGHTNESS"
	case "colorTemperature":
		name = "COLORTEMPERATURE"
	default:
		result = ""
	}
	if result == "" {
		return name, errors.New("unsupported name")
	}
	return name, nil

}
