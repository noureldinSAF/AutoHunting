package main 
import (
	"fmt"
	"github.com/likexian/whois"
	"github.com/likexian/whois-parser"
)

func main() {
	query := "google.com"

	result, err := whois.Whois(query)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	parsedResult, err := whoisparser.Parse(result)
	if err != nil {
		fmt.Println("Error parsing whois data:", err)
		return
	}

	// Print parsed whois information

	fmt.Println(parsedResult.Registrant.Organization)
	

}
