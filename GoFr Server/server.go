package main

import (
	"Go_backend/sendmail"
	"bytes"
	"errors"  // Standard Go errors package
	"os/exec" // Import os/exec for running external commands
	"strings"

	"gofr.dev/pkg/gofr"
)

func main() {
    app := gofr.New()

    // Chain with Sentiment Analysis API
    app.GET("/get_twit_trend", func(ctx *gofr.Context) (interface{}, error) {
		// Define the Python script to execute
		cmd := exec.Command("python", "twittrend.py")

		// Capture the output
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr

		// Run the Python script
		err := cmd.Run()
		if err != nil {
			// Return any error during script execution
			return nil, errors.New("Failed to run Python script: " + stderr.String())
		}

		// Get the output text from the script
		result := strings.TrimSpace(out.String())

		// Check if the result is empty
		if result == "" {
			return nil, errors.New("No tweets found or script output is empty")
		}

		// Return the concatenated tweets
		return map[string]string{
			"tweets": result,
		}, nil
	})


    app.POST("/create/x", func(ctx *gofr.Context) (interface{}, error) {
        // Get the text from request body
        var body struct {
            Text string `json:"text"`
        }
        if err := ctx.Bind(&body); err != nil {
            return nil, errors.New("invalid request body")
        }

        if body.Text == "" {
            return nil, errors.New("text is required")
        }

        ctx.Logger.Info("Received text: ", body.Text)

        cmd := exec.Command("python", "twit.py", body.Text)
        err := cmd.Run()
        if err != nil {
            return nil, errors.New("failed to run the Python script")
        }

        return map[string]string{
            "message": "Text received successfully and script executed",
        }, nil
    })

    app.POST("/create/email", func(ctx *gofr.Context) (interface{}, error) {
        var body struct {
            Context string `json:"context"`
        }
        if err := ctx.Bind(&body); err != nil {
            return nil, errors.New("invalid request body")
        }

        if body.Context == "" {
            return nil, errors.New("context is required")
        }

        ctx.Logger.Info("Received context: ", body.Context)

        // Call the Send_Main function from sendmail package
        err := sendmail.Send_Mail(body.Context)
        if err != nil {
            return nil, errors.New("failed to send email: " + err.Error())
        }

        return map[string]string{
            "message": "Email sent successfully",
        }, nil
    })

    app.Run()
}
