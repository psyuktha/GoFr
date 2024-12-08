package sendmail

import (
	"bytes"
	// "encoding/csv"
	"encoding/json"
	// "fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type LLMRequest struct {
	Question string `json:"question"`
}

type LLMResponse struct {
	Status   string `json:"status"`
	Response string `json:"response"`
}

// Get_llm_response makes an API request to generate email content.
func Get_llm_response(context string) (string, string) {
	// Create the request body
	requestBody := LLMRequest{
		Question: context,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("Failed to marshal request: %s", err)
		return "_ERROR", "_ERROR"
	}

	// Make the request to the FastAPI endpoint
	resp, err := http.Post("http://localhost:8000/generate_content/",
		"application/json",
		bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Failed to make request: %s", err)
		return "_ERROR", "_ERROR"
	}
	defer resp.Body.Close()

	// Read and parse the response
	var llmResponse LLMResponse
	if err := json.NewDecoder(resp.Body).Decode(&llmResponse); err != nil {
		log.Printf("Failed to decode response: %s", err)
		return "_ERROR", "_ERROR"
	}

	// Parse the response content into subject and body
	// Assuming the response is formatted as "Subject: <subject>\nBody: <body>"
	parts := strings.Split(llmResponse.Response, "\n")
	var subject, body string

	for _, part := range parts {
		if strings.HasPrefix(part, "Subject: ") {
			subject = strings.TrimPrefix(part, "Subject: ")
		} else if strings.HasPrefix(part, "Body: ") {
			body = strings.TrimPrefix(part, "Body: ")
		}
	}

	if subject == "" || body == "" {
		log.Printf("Invalid response format")
		return "_ERROR", "_ERROR"
	}

	return body, subject
}

// Send_mail sends the generated email to a list of recipients from a CSV file.
func Send_mail(content string) error {
	from := mail.NewEmail("Suvan Banerjee", "suvan@burdenoff.com")
	to := mail.NewEmail("Recipient", "2109yukips@gmail.com")
	subject := "AI Generated Content"
	plainTextContent := content
	htmlContent := "<strong>" + content + "</strong>"

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)

	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	log.Printf("Email Status Code: %d", response.StatusCode)
	log.Printf("Email Body: %s", response.Body)

	return nil
}

// GetAnalytics fetches the analytics data from SendGrid.
func GetAnalytics() {

	host := "https://api.sendgrid.com"
	request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/stats", host)
	request.Method = "GET"
	queryParams := make(map[string]string)
	queryParams["start_date"] = "2024-11-22"
	request.QueryParams = queryParams

	// Make the API request
	response, err := sendgrid.API(request)
	if err != nil {
		log.Println(err)
	} else {
		// Log the analytics response
		log.Printf("SendGrid Analytics: StatusCode: %d, Body: %s, Headers: %s",
			response.StatusCode, response.Body, response.Headers)
	}
}
