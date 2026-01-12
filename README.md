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
