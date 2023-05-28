# exchange_rate_bitcoin

This is a Go application that provides the current exchange rate of Bitcoin (BTC) to USD and allows users to subscribe to rate change notifications via email.

## Installation

1. Clone the repository:
```git clone https://github.com/anna-bilohan/exchange_rate_bitcoin```
2. Change to the project directory:
```cd exchange_rate_bitcoin```
3. Run the application:
```go run main.go```

Alternatively, you can run the application with the help of Dockerfile. Please refer to the documentation: https://docs.docker.com/get-started/02_our_app/

## Usage
The application exposes several API endpoints:

GET /rate: Fetches the current BTC exchange rate in USD.\
POST /subscribe: Subscribes an email address for rate change notifications.\
GET /subscribe: Serves the HTML form for email subscription.\
GET /sendEmails: Sends the current rate to all subscribed users.

Access the application in your web browser:

http://localhost:8080

## Dependencies
This application uses the following external dependencies:

Gin: A web framework for Go.
gomail: A Go package for sending emails.
crypto/tls: Package tls provides support for SSL/TLS protocols.
Make sure to have these dependencies installed before running the application.

## Licence
This project is licensed under the MIT License.

## Warning

For testing purposes the email-adress was created in order to send out notifications to the subscribers. 
It is hardcodedn under: gomail.NewDialer("smtp.office365.com", 587, "bitcoin_exchange_rate_2023@outlook.com", "bitcoin_bitcoin") // Replace with your SMTP server and credentials. After testing the program multiple times, the spam warning was received. Thus, you might need to change the environment variables/hardcode ane-mail adress and coressponding credentials of your choice in orde to get an e-mail with the exchange rate.

If emails.txt is populated with one or more of these email addresses that are using a reserved domain, such as "example.com", "test.com", or "invalid.com", sending the current rate to all subscribed users might fail. Please, make sure there are only valid e-mails, that are not using a reserved domain, otherwise you will get an error message "Recipient address reserved by RFC 2606" in the command-line.
