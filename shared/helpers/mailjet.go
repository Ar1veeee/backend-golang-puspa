package helpers

import (
    "backend-golang/shared/config"
    "fmt"

    "github.com/mailjet/mailjet-apiv3-go/v4"
)

var mailjetClient *mailjet.Client

func InitMailjet() error {
    apiKey := config.GetEnv("MAILJET_API_KEY", "")
    secretKey := config.GetEnv("MAILJET_SECRET_KEY", "")

    if apiKey == "" || secretKey == "" {
        return fmt.Errorf("mailjet API credentials not configured")
    }

    mailjetClient = mailjet.NewMailjetClient(apiKey, secretKey)
    return nil
}

func GetMailjetClient() *mailjet.Client {
    if mailjetClient == nil {
        panic("mailjet client not initialized")
    }
    return mailjetClient
}

func GetEmailSender() string {
    return config.GetEnv("MAILJET_SENDER", "")
}