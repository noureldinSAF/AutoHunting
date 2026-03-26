package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func checkPermissions(url string) (map[string]int, error) {

	client := http.Client{Timeout: 10 * time.Second}

	results := make(map[string]int)

	resp, err := client.Get(fmt.Sprintf("%s?list-types=2", url))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	results["LIST"] = resp.StatusCode

	resp, err = client.Get(fmt.Sprintf("%s?policy", url))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	results["POLICY"] = resp.StatusCode

	resp, err = client.Get(fmt.Sprintf("%s?publicAccessBlock", url))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	results["PUBLIC_ACCESS_BLOCK"] = resp.StatusCode

	resp, err = client.Get(fmt.Sprintf("%s?acl", url))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	results["ACL"] = resp.StatusCode

	return results, nil

}

func main() {

	s3Url := os.Args[1]

	results, err := checkPermissions(s3Url)
	if err != nil {
		fmt.Println(err)
	}

	for per, sc := range results {
		fmt.Printf("%s :=> status code :=> %d \n", per, sc)
	}

}
