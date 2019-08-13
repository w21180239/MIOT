package client

import "time"

type Scene struct {
	SceneID          int    `json:"id"`
	SceneName        string `json:"sceneName"`
	SceneDescription string `json:"sceneDescription"`
	SceneIconID      int    `json:"sceneIconId"`
	SceneState       string `json:"sceneState"`
	Enabled          bool   `json:"enabled"`
}

func (s *Scene) GetAttributesAndActions() (
	attributes []map[string]interface{}, actions []string) {

	attributes = make([]map[string]interface{}, 0)
	actions = make([]string, 0)

	name := make(map[string]interface{})
	name["name"] = "name"
	name["value"] = s.SceneName
	name["scale"] = ""
	name["timestampOfSample"] = time.Now().Unix()
	name["uncertaintyInMilliseconds"] = 0
	attributes = append(attributes, name)

	state := make(map[string]interface{})
	state["name"] = "turnOnState"
	state["value"] = "OFF"
	state["scale"] = ""
	state["timestampOfSample"] = time.Now().Unix()
	state["uncertaintyInMilliseconds"] = 0
	attributes = append(attributes, state)

	actions = append(actions, "turnOn", "turnOff")
	return attributes, actions

}
