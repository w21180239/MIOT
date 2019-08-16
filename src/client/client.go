package client

import (
	"bytes"
	"constant"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Client struct {
	UserID    string
	ProjectID string
}

func NewClient(userID, projectID string) *Client {

	return &Client{
		UserID:    userID,
		ProjectID: projectID,
	}

}

func (c *Client) Host() string {
	if constant.Local_debug {
		return "https://www.homepluscloud.com:8443/api/"
	} else {
		return "http://smart-home-service:8080/"
	}
}

func (c *Client) GetDevices() ([]Device, error) {
	var resp *http.Response
	var err error
	url := c.Host() + "projects/" + c.ProjectID + "/device-discovery"
	if constant.Local_debug {
		client := &http.Client{}
		req, _ := http.NewRequest("GET", url, nil)
		cookie := &http.Cookie{Name: "AccessToken", Value: constant.Cookie, HttpOnly: true}
		req.AddCookie(cookie)
		resp, err = client.Do(req)
	} else {
		resp, err = http.Get(url)
	}
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	var raw map[string][]Device
	//var raw interface{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&raw)
	if err != nil {
		return nil, err
	}
	devices := make([]Device, 0)
	for _, val := range raw {
		devices = append(devices, val...)
	}
	return devices, nil

}

func (c *Client) GetSpaces() (Space, error) {

	url := c.Host() + "projects/" + c.ProjectID + "/spaces?type=ENABLE"
	resp, err := http.Get(url)
	if err != nil {
		errSpace := Space{0, "err", 0, nil}
		fmt.Println(err)
		return errSpace, err
	}
	defer resp.Body.Close()
	var raw map[string]Space
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&raw)
	if err != nil {
		errSpace := Space{0, "err", 0, nil}
		return errSpace, err
	}
	spaces := raw["rootSpace"]
	return spaces, nil

}

func (c *Client) GetScenes() ([]Scene, error) {
	url := c.Host() + "projects/" + c.ProjectID + "/scenes"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	var raw []Scene
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&raw)
	if err != nil {
		return nil, err
	}
	spaces := raw
	return spaces, nil
}

func (c *Client) GetSpaceDevices(SpaceID int) (SpaceDetail, error) {
	url := c.Host() + "projects/" + c.ProjectID + "/spaces/" + strconv.Itoa(SpaceID)

	resp, err := http.Get(url)
	if err != nil {
		errSpaceDetail := SpaceDetail{0, "", 0, "", "", false, nil, nil, nil, nil, nil, nil}
		fmt.Println(err)
		return errSpaceDetail, err
	}
	defer resp.Body.Close()
	var raw SpaceDetail
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&raw)
	if err != nil {
		errSpaceDetail := SpaceDetail{0, "", 0, "", "", false, nil, nil, nil, nil, nil, nil}
		fmt.Println(err)
		return errSpaceDetail, err
	}
	spacedetail := raw
	return spacedetail, nil
}

func (c *Client) PostAction(deviceID int,
	action string, value interface{}) error {

	url := c.Host() + "/basic-control/action"
	postBody := make(map[string]interface{})
	postBody["projectId"] = c.ProjectID
	postBody["deviceId"] = deviceID
	postBody["action"] = action
	postBody["value"] = value
	postBodyBytes, err := json.Marshal(postBody)
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json",
		bytes.NewBuffer(postBodyBytes))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(resp.Status + " " + string(body))
	}
	return nil

}

func (c *Client) PostLightingGroupControl(lightingGroupID int,
	characteristic string, value interface{}) error {

	url := c.Host() + "/basic-control/lighting-group"
	postBody := make(map[string]interface{})
	postBody["projectId"] = c.ProjectID
	postBody["lightingGroupId"] = lightingGroupID
	postBody["characteristic"] = characteristic
	postBody["value"] = value
	postBodyBytes, err := json.Marshal(postBody)
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json",
		bytes.NewBuffer(postBodyBytes))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(resp.Status + " " + string(body))
	}
	return nil

}

func (c *Client) PostScene(sceneID int) error {

	url := c.Host() + "/basic-control/scene"
	postBody := make(map[string]interface{})
	postBody["projectId"] = c.ProjectID
	postBody["smartPanelId"] = 0
	postBody["spaceId"] = 0
	postBody["sceneId"] = sceneID
	postBodyBytes, err := json.Marshal(postBody)
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json",
		bytes.NewBuffer(postBodyBytes))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(resp.Status + " " + string(body))
	}
	return nil

}

func (c *Client) GetDeviceByID(deviceID int) (*Device, error) {

	deviceIDStr := strconv.Itoa(deviceID)
	url := c.Host() + "/projects/" + c.ProjectID + "/device-states"
	query := "?deviceIds=" + deviceIDStr + "&spaceIds=&groupIds="
	resp, err := http.Get(url + query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var raw map[string][]Device
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&raw)
	if err != nil {
		return nil, err
	}
	var device *Device
	for _, val := range raw {
		device = &val[0]
	}
	return device, nil

}

