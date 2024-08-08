package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/skip2/go-qrcode"
)

func GenerateReferralLink(userID int64) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	botLink := os.Getenv("BOT_LINK")
	return fmt.Sprintf(botLink+"?start=%d", userID)
}

func GenerateQRCode(link string) ([]byte, error) {
	png, err := qrcode.Encode(link, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}
	return png, nil
}
