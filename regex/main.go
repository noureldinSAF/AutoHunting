package main

import (
	"fmt"
	"regexp"
)

var (
	// URLRegex matches various URL formats, including:
	endpointsRegex = regexp.MustCompile(`(?m)['"](?:\/|\.\.\/|\.\/)[^\/\>\< \)\(\{\}\,\'\"]([^\\>< \)\(\{\}\,\'\"])*?\\?['"]`)

	emailRegex = regexp.MustCompile(`(?m)^[\w-]+(\.[\w-]+)*@([a-z0-9-]+(\.[a-z0-9-]+)*?\.[a-z]{2,6}|(\d{1,3}\.){3}\d{1,3})(:\d{4})?$`)

	githubSecretRegex = regexp.MustCompile(`[gG][iI][tT][hH][uU][bB].*[''|"][0-9a-zA-Z]{35,40}[''|"]`)

	// paramterRegex matches URL query parameters like ?key= or &key=
	paramterRegex = regexp.MustCompile(`(?m)([?&])([a-zA-Z_][a-zA-Z0-9_]*)=`)


	nodeModulesPathRegex = regexp.MustCompile(`/node_modules/(@?[a-z-_.0-9]+)/`)

	ipRegex = regexp.MustCompile(`['"](?:([a-zA-Z0-9]+:)?\/\/)?\d{1,3}(?:\.\d{1,3}){3}(:\d{1,5})?(?:\/.*?)?['"]`)
	groupedEmailRegex = regexp.MustCompile(`(\w+)@(\w+)\.(\w+)`)
)


func main(){

	text := `
	
	zomasec@wearehackerone.com


	 apiEndpoint = '/api/v1/admin/getAllUsers'

	 endpoint2 = '/api/v3/delete/?id='

	 github_secret= ghp_IGvgCoPdOjhQANnEEXLxF7ulf3kGFa4gd7Jx

	ip = 192.168.1.1 
	
	`

	test1 := endpointsRegex.FindAllString(text, -1)

	test2 := endpointsRegex.MatchString(text)

	test3 := groupedEmailRegex.FindStringSubmatch(text)

	fmt.Println(test1)

	fmt.Println(test2)

	fmt.Printf("FullMatch: %s, Username => %s, CompanyName => %s, TLD => %s", test3[0],test3[1], test3[2], test3[3])
}