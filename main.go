package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

// Struct to hold the subscribed users' email addresses
type Subscription struct {
	Email string `json:"email"`
}

var subscribedUsers []Subscription

func main() {
	err := readEmailsFromFile()
	if err != nil {
		log.Fatal("Failed to read emails from file:", err)
	}

	router := gin.Default()

	// API endpoint to get the current BTC exchange rate in UAH
	router.GET("/rate", getCurrentExchangeRate)

	// API endpoint to subscribe an email for rate change notifications
	router.POST("/subscribe", subscribeEmail)

	// Serve the HTML form for email subscription
	router.GET("/subscribe", serveForm)

	// API endpoint to send the current rate to all subscribed users
	router.GET("/sendEmails", sendRatesToSubscribers)

	// Run the server
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func getCurrentExchangeRate(c *gin.Context) {
	// Make a request to the external API to fetch the exchange rate
	response, err := http.Get("https://blockchain.info/ticker")
	if err != nil {
		log.Println("Failed to fetch exchange rate:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch exchange rate",
		})
		return
	}
	defer response.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Failed to read response body:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read response body",
		})
		return
	}

	// Parse the response body into a JSON object
	var exchangeRates map[string]struct {
		Last float64 `json:"last"`
	}
	err = json.Unmarshal(body, &exchangeRates)
	if err != nil {
		log.Println("Failed to parse response body:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to parse response body",
		})
		return
	}

	// Extract the exchange rate for the desired currency (e.g., USD)
	usdRate := exchangeRates["USD"].Last

	c.JSON(http.StatusOK, gin.H{
		"exchange_rate": usdRate,
	})
}

func serveForm(c *gin.Context) {
	html := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Subscribe Email</title>
		</head>
		<body>
			<h1>Subscribe Email</h1>
			<form onsubmit="subscribeEmail(event)">
				<input type="email" id="email" placeholder="Enter your email" required>
				<button type="submit">Subscribe</button>
			</form>
			
			<script>
				function subscribeEmail(event) {
					event.preventDefault();
					
					const emailInput = document.getElementById('email');
					const email = emailInput.value.trim();
					
					fetch('/subscribe', {
						method: 'POST',
						headers: {
							'Content-Type': 'application/json'
						},
						body: JSON.stringify({ email: email })
					})
					.then(response => response.json())
					.then(data => {
						if (data.error) {
							alert(data.error);
						} else {
							alert(data.message);
							emailInput.value = '';
						}
					})
					.catch(error => {
						console.error('Error:', error);
						alert('An error occurred. Please try again.');
					});
				}
			</script>
		</body>
		</html>
	`
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

func subscribeEmail(c *gin.Context) {
	var subscription Subscription

	if err := c.ShouldBindJSON(&subscription); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return
	}

	// Check if the email already exists in subscribedUsers
	if isEmailSubscribed(subscription.Email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email already subscribed",
		})
		return
	}

	subscribedUsers = append(subscribedUsers, subscription)

	// Write the email to a file
	writeEmailToFile(subscription.Email)

	c.JSON(http.StatusOK, gin.H{
		"message": "Email subscribed successfully",
	})
}

func sendRatesToSubscribers(c *gin.Context) {
	// Get the current exchange rate
	usdRate, err := getCurrentExchangeRateInUSD()
	if err != nil {
		log.Println("Failed to get exchange rate:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get exchange rate",
		})
		return
	}

	// Compose the email body
	body := "The current exchange rate for BTC to USD is: " + fmt.Sprintf("%.2f", usdRate)

	// Compose the email message
	message := gomail.NewMessage()
	message.SetHeader("From", "bitcoin_exchange_rate_2023@outlook.com") // Replace with your email address
	message.SetHeader("To", getEmailsOfSubscribers()...)
	message.SetHeader("Subject", "Bitcoin Exchange Rate")
	message.SetBody("text/plain", body)

	// Configure the dialer with TLS
	dialer := gomail.NewDialer("smtp.office365.com", 587, "bitcoin_exchange_rate_2023@outlook.com", "bitcoin_bitcoin") // Replace with your SMTP server and credentials
	dialer.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         "smtp.office365.com",
	}

	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		log.Println("Failed to send email:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to send email",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Rates sent to subscribers",
	})
}

func getCurrentExchangeRateInUSD() (float64, error) {
	response, err := http.Get("https://blockchain.info/ticker")
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}

	var exchangeRates map[string]struct {
		Last float64 `json:"last"`
	}
	err = json.Unmarshal(body, &exchangeRates)
	if err != nil {
		return 0, err
	}

	return exchangeRates["USD"].Last, nil
}

func getEmailsOfSubscribers() []string {
	emails := make([]string, len(subscribedUsers))
	for i, user := range subscribedUsers {
		emails[i] = user.Email
	}
	return emails
}

func isEmailSubscribed(email string) bool {
	for _, user := range subscribedUsers {
		if strings.EqualFold(user.Email, email) {
			return true
		}
	}
	return false
}

func writeEmailToFile(email string) {
	file, err := os.OpenFile("emails.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Failed to open file:", err)
		return
	}
	defer file.Close()

	// Check if the email already exists in the file
	if checkDuplicateEmail(email, file) {
		log.Println("Email already exists in the file")
		return
	}

	if _, err := file.WriteString(email + "\n"); err != nil {
		log.Println("Failed to write email to file:", err)
	}
}

func checkDuplicateEmail(email string, file *os.File) bool {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.EqualFold(scanner.Text(), email) {
			return true
		}
	}
	return false
}

func readEmailsFromFile() error {
	file, err := os.Open("emails.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		email := scanner.Text()
		subscribedUsers = append(subscribedUsers, Subscription{Email: email})
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
