package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	

)


// Pattern represents a regex pattern with metadata
// type Pattern struct {
// 	Sensitive bool   `yaml:"sensitive"`
// 	Name      string `yaml:"name"`
// 	Value     string `yaml:"value"`
// }

// Config represents the structure of the patterns file


func fetchTags(image string) ([]string, error) {
    url := fmt.Sprintf("https://registry.hub.docker.com/v2/repositories/%s/tags?page_size=100", image)

    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    var result struct {
        Results []struct {
            Name string `json:"name"`
        } `json:"results"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    tags := make([]string, len(result.Results))
    for i, r := range result.Results {
        tags[i] = r.Name
    }

    return tags, nil
}

// func main() {
//     image := "library/nginx" // Replace with your image name
//     tags, err := fetchTags(image)
//     if err != nil {
//         fmt.Fprintf(os.Stderr, "Error fetching tags: %v\n", err)
//     }

//     fmt.Printf("Available tags for %s:\n", image)
//     for _, tag := range tags {
//         fmt.Println(tag)
//     }
// }