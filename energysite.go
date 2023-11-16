package tesla

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// this represents site_info endpoint
type EnergySite struct {
	ID                   string `json:"id"`
	SiteName             string `json:"site_name"`
	BackupReservePercent int64  `json:"backup_reserve_percent,omitempty"`
	DefaultRealMode      string `json:"default_real_mode,omitempty"`

	productId int64
	c         *Client
}

type EnergySiteStatus struct {
	ResourceType      string  `json:"resource_type"`
	SiteName          string  `json:"site_name"`
	GatewayId         string  `json:"gateway_id"`
	EnergyLeft        float64 `json:"energy_left"`
	TotalPackEnergy   uint64  `json:"total_pack_energy"`
	PercentageCharged float64 `json:"percentage_charged"`
	BatteryType       string  `json:"battery_type"`
	BackupCapable     bool    `json:"backup_capable"`
	BatteryPower      int64   `json:"battery_power"`

	c *Client
}

type SiteInfoResponse struct {
	Response *EnergySite `json:"response"`
}

type SiteStatutsResponse struct {
	Response *EnergySiteStatus `json:"response"`
}

// SiteCommandResponse is the response from the Tesla API after POSTing a command.
type SiteCommandResponse struct {
	Response struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"response"`
}

// return fetches the energy site for the given product ID
func (c *Client) EnergySite(productID int64) (*EnergySite, error) {
	siteInfoResponse := &SiteInfoResponse{}
	if err := c.getJSON(c.baseURL+"/energy_sites/"+strconv.FormatInt(productID, 10)+"/site_info", siteInfoResponse); err != nil {
		return nil, err
	}
	siteInfoResponse.Response.c = c
	siteInfoResponse.Response.productId = productID
	return siteInfoResponse.Response, nil
}

func (s *EnergySite) EnergySiteStatus() (*EnergySiteStatus, error) {
	siteStatusResponse := &SiteStatutsResponse{}
	if err := s.c.getJSON(s.statusPath(), siteStatusResponse); err != nil {
		return nil, err
	}
	siteStatusResponse.Response.c = s.c
	return siteStatusResponse.Response, nil
}

func (s *EnergySite) basePath() string {
	return strings.Join([]string{s.c.baseURL, "energy_sites", strconv.FormatInt(s.productId, 10)}, "/")
}

func (s *EnergySite) statusPath() string {
	return strings.Join([]string{s.basePath(), "site_status"}, "/")
}

func (s *EnergySite) SetBatteryReserve(percent uint64) error {
	url := s.basePath() + "/backup"
	payload := fmt.Sprintf(`{"backup_reserve_percent":%d}`, percent)
	body, err := s.sendCommand(url, []byte(payload))
	if err != nil {
		return err
	}

	response := SiteCommandResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	if response.Response.Code != 201 {
		return fmt.Errorf("batteryReserve failed: %s", response.Response.Message)
	}

	return nil
}

// Sends a command to the vehicle
func (s *EnergySite) sendCommand(url string, reqBody []byte) ([]byte, error) {
	body, err := s.c.post(url, reqBody)
	if err != nil {
		return nil, err
	}
	if len(body) > 0 {
		response := &CommandResponse{}
		if err := json.Unmarshal(body, response); err != nil {
			return nil, err
		}
		if !response.Response.Result && response.Response.Reason != "" {
			return nil, errors.New(response.Response.Reason)
		}
	}
	return body, nil
}
