package lead

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	fmt.Println("This message will show up in the CLI console.")

	// Get netlify environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dbConnect := getDbConnect(dbHost, dbPort, dbUser, dbPassword, dbName)

	// Connect to database
	db, err := dbConnect.connect()
	if err != nil {
		return nil, err
	}

	// Create lead
	lead := &Lead{
		Name:        request.QueryStringParameters["name"],
		Email:       request.QueryStringParameters["email"],
		NickName:    request.QueryStringParameters["nickname"],
		AcceptTerms: request.QueryStringParameters["acceptterms"],
		AcceptOptin: request.QueryStringParameters["acceptoptin"],
	}

	_, err = createLead(db, lead)
	if err != nil {
		return nil, err
	}

	// Send email
	mail := createMail(lead)
	err = mail.send()

	return &events.APIGatewayProxyResponse{
		StatusCode:      200,
		Headers:         map[string]string{"Content-Type": "text/plain"},
		Body:            "Hello, world!",
		IsBase64Encoded: false,
	}, nil
}

type dbConnect struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

type Lead struct {
	Name        string `json:"name"`
	NickName    string `json:"nickname"`
	Email       string `json:"email"`
	AcceptTerms string `json:"acceptTerms"`
	AcceptOptin string `json:"acceptOptin"`
}

type Mail struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func getDbConnect(host string, port string, username string, password string, database string) *dbConnect {
	return &dbConnect{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Database: database,
	}
}

func createLead(db *sql.DB, lead *Lead) (sql.Result, error) {
	// Insert lead into database
	result, err := db.Exec("INSERT INTO leads (name, email) VALUES ($1, $2)", lead.Name, lead.Email)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (db *dbConnect) connect() (*sql.DB, error) {
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", db.Host, db.Port, db.Username, db.Password, db.Database)
	return sql.Open("mysql", dbInfo)
}

func createMail(lead *Lead) *Mail {
	return &Mail{
		To:      lead.Email,
		Subject: "Welcome to our website!",
		Body:    "Thank you for registering!",
	}
}

func (mail *Mail) send() error {
	// Send email logic
	return nil
}

func main() {
	lambda.Start(handler)
}
