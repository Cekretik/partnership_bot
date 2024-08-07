package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GenerateReferralLink(userID int64) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	botLink := os.Getenv("BOT_LINK")
	return fmt.Sprintf(botLink+"?start=%d", userID)
}
