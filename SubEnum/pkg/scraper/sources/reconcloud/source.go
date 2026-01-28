package reconcloud

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Source struct{}

type reconCloudResponse struct {
	MsgType         string            `json:"msg_type"`
	RequestID       string            `json:"request_id"`
	OnCache         bool              `json:"on_cache"`
	Step            string            `json:"step"`
	CloudAssetsList []cloudAssetsList `json:"cloud_assets_list"`
}

type cloudAssetsList struct {
	Key           string `json:"key"`
	Domain        string `json:"domain"`
	CloudProvider string `json:"cloud_provider"`
}

func (s *Source) Name() string {
	return "reconcloud"
}

func (s *Source) RequiresAPIKey() bool {
	return false
}

func (s *Source) Search(query string, client *http.Client) ([]string, error) {
	var results []string

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://recon.cloud/api/search?domain=%s", query), nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response reconCloudResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	if len(response.CloudAssetsList) > 0 {
		for _, cloudAsset := range response.CloudAssetsList {
			results = append(results, cloudAsset.Domain)
		}
	}

	return results, nil
}

