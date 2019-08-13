package client

type Space struct {
	SpaceID   int     `json:"id"`
	SpaceName string  `json:"spaceName"`
	ParentID  int     `json:"parentId"`
	SubSpaces []Space `json:"subSpaces"`
}

type SubSpace struct {
	SpaceID   int    `json:"id"`
	SpaceName string `json:"spaceName"`
}

type SpaceDevice struct {
	DeviceID   int    `json:"id"`
	DeviceName string `json:"deviceName"`
}

type SpaceDetail struct {
	ID               int                      `json:"id"`
	SpaceName        string                   `json:"spaceName"`
	ParentID         int                      `json:"parentId"`
	SpaceDescription string                   `json:"SpaceDescription"`
	SpaceState       string                   `json:"spaceState"`
	Enabled          bool                     `json:"enabled"`
	SubSpaces        []SubSpace               `json:"subSpaces"`
	Devices          []SpaceDevice            `json:"devices"`
	LightingGroups   []map[string]interface{} `json:"lightingGroups"`
	ActionCommands   []map[string]interface{} `json:"actionCommands"`
	Sliders          []map[string]interface{} `json:"sliders"`
	Scenes           []map[string]interface{} `json:"scenes"`
}

func (s *Space) GetAllSpace() []Space {

	allspace := make([]Space, 0)
	if s.SubSpaces != nil {
		allspace = append(allspace, s.SubSpaces...)
		for _, subspace := range s.SubSpaces {
			allspace = append(allspace, subspace.GetAllSpace()...)
		}
	}
	return allspace
}
