package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	


	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func Get_llm_response(string context) (string, string) {
	prompt := "Your task is to draft a mail to some users of gofr make it playfull and engaging give output as a json where the subject of mail is key subject and body of mail is body"
	
	generate_post(prompt,"Email about Citcuit breakers in http")
	match, _ := regexp.MatchString("mail", prompt)
	re := regexp.MustCompile((?s)\{.*?\})
	// Find the JSON block
	match := re.FindString(text)
	if match != "" {
		var result map[string]string
		err := json.Unmarshal([]byte(match), &result)
		if err != nil {
			log.Fatalf("failed to unmarshal JSON: %s", err)
		}
		text = result["body"]
		title = result["subject"]
	} else {
		text = "_ERROR_"
		title = "_ERROR_"
	}
	return text, title
}

func Send_mail(context string) {
	file, err := os.Open("emails.csv")
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("failed to read file: %s", err)
	}

	text, title := Get_llm_response(context)
	if text == "_ERROR" || title == "ERROR_" {
		log.Fatalf("failed to get response")
	}
	else {
	for _, record := range records {
		toEmail := record[0]
		from := mail.NewEmail("Suvan Banerje", "suvan@burdenoff.com")
		subject := title
		to := mail.NewEmail("User", toEmail)
		plainTextContent := text
		htmlContent := "<strong>" + text + "</strong>"
		message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
		client := sendgrid.NewSendClient("")
		response, err := client.Send(message)
		if err != nil {
			log.Println(err)
		} else {
			fmt.Println(response.StatusCode)
			fmt.Println(response.Body)
			fmt.Println(response.Headers)
		}
	}
}
}

func GetAnalytics() {
	apiKey := ""
	host := "https://api.sendgrid.com"
	request := sendgrid.GetRequest(apiKey, "/v3/stats", host)
	request.Method = "GET"
	queryParams := make(map[string]string)
	queryParams["start_date"] = "2024-11-22"
	request.QueryParams = queryParams
	response, err := sendgrid.API(request)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}

func main() {
	GetAnalytics()
}