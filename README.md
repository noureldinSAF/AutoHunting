# Here is The AutoHunting tool 

# How to use asn2cidr tool 
go to inside the folder and run the following bash cmd 
```
./asnmap -asn AS32934    //ASN for facebook 
```

# How to use cidr2ips tool 
Get the List of cidrs in list.txt and run the following cmd inside the folder
```
cat list.txt | go run .
```

# How to use both asn2cidr and cidr2ips together. From ASN to get All the ips 
go inside cidr2ips and run the following command 
```
../asn2cidr/asnmap -asn AS32934 | go run .
```
you can use ``` | wc ```  to know the number of lines

# Domain Enumeration ( Enumerate every domains related to a company by its name ) 
go to this folder:  DomEnum/cmd/DomEnum

```go go run . -h ```

passive enumeration ```go go run . -q Swisscom -o swisscomDomains.txt```

passive and active enumeration ```go go run . -q Swisscom -o swisscomDomains.txt -active```

add timeout ( important ) ```go go run . -q Swisscom -o swisscomDomains.txt -active -t 60```

Note : For passive enumeration, we have three sources 
1- crtsh : free without api

2- whoisfreaks: free with api, it devides the results into pages. every page contains 50 domain. to grep the next page you must sleep for 1 min either you will get rate limit error. 

in whoisfreaks.go line 92  
```go
if ro.CurrentPage >= 1 {   // replace it with  ro.CurrentPage >= ro.ToTotalPages to get all pages
```
there is a subscribtion wihtout rate limit 

3- whoisxmlapi: is not free.

You can add or edit api for all tools in this file https://github.com/noureldinSAF/AutoHunting/tree/main/DomEnum/cmd/DomEnum/internal/config/config.yaml . whoisfreaks and whoisxmlapi will work only if the config file contains api for them



# DnsEnum 
Check the existence of targest and you can get all dns records A,AAAA,CNAME, .... 
go to ./main.go file 
run 
```go
go run ./main.go -h
```
usage example
```go
go run . -l swisscomDomains.txt -o swisscomDns.txt -of lines -c 20
```

# SubEnum
Enumerate subdomains
go inside this folder SubEnum/cmd/subenum
```bash
go run . -h
```
usage example
```bash
go run . -active -c 10 -i domains.txt -o subs.txt
```
mutation is time consuming, so it is important to limit it 
```bash
go run . -active -c 20 -i domains.txt -o subs.txt -e -max-mutations-size 50
```
To Be Updated -> The tool doesn't brute force in active enumeration 







