# exchange_rate_bitcoin


# Bitcoin Exchange Rate

This is a Go application that provides the current exchange rate of Bitcoin (BTC) to USD and allows users to subscribe to rate change notifications via email.

## Installation

1. Clone the repository:

   git clone https://github.com/anna-bilohan/exchange_rate_bitcoin
   
2. Change to the project directory:

   cd exchange_rate_bitcoin
   
3. Build the application:

  go build
  
4.Run the application:

  ./exchange_rate_bitcoin
  
## Usage
The application exposes several API endpoints:

GET /rate: Fetches the current BTC exchange rate in USD.
POST /subscribe: Subscribes an email address for rate change notifications.
GET /subscribe: Serves the HTML form for email subscription.
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
