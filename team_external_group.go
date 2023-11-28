package gapi

import (
	"encoding/json"
	"fmt"
)

// TeamGroup represents a Grafana TeamGroup.
type TeamGroup struct {
	OrgID   int64  `json:"orgId,omitempty"`
	TeamID  int64  `json:"teamId,omitempty"`
	GroupID string `json:"groupId,omitempty"`
}

// TeamGroups fetches and returns the list of Grafana team group whose Team ID it's passed.
func (c *Client) TeamGroups(teamId int64) ([]TeamGroup, error) {
	var teamGroups []TeamGroup
	err := c.request("GET", fmt.Sprintf("/api/teams/%d/groups", teamId), nil, nil, &teamGroups)
	if err != nil {
		return teamGroups, err
	}

	return teamGroups, nil
}

// NewTeamGroup creates a new Grafana Team Group .
func (c *Client) NewTeamGroup(teamId int64, groupID string) error {
	dataMap := map[string]string{
		"groupId": groupID,
	}
	data, err := json.Marshal(dataMap)
	if err != nil {
		return err
	}

	return c.request("POST", fmt.Sprintf("/api/teams/%d/groups", teamId), nil, data, nil)
}

// DeleteTeam deletes the Grafana team whose ID it's passed.
func (c *Client) DeleteTeamGroup(id int64, groupID string) error {
	return c.request("DELETE", fmt.Sprintf("/api/teams/%d/groups/%s", id, groupID), nil, nil, nil)
}
