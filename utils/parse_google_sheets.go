package utils

import (
	"log"
	"strconv"
	"strings"
	"time"

	"main/models"

	"gorm.io/gorm"
)

func UpdateUserFromSheet(db *gorm.DB, userID string, values []interface{}) {
	var user models.User
	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		log.Printf("Invalid user ID: %v", err)
		return
	}

	db.First(&user, "user_id = ?", id)

	if user.ID == 0 {
		// Create a new user if they do not exist in the database
		user = models.User{UserID: id}
	}

	// Function to clean and convert a string to float64
	convertToFloat := func(value string) (float64, error) {
		// Remove currency symbols and commas
		cleaned := strings.ReplaceAll(value, "฿", "")
		cleaned = strings.ReplaceAll(cleaned, ",", "")
		return strconv.ParseFloat(cleaned, 64)
	}

	// Check and convert each value from the sheet.
	if len(values) > 2 {
		if v, ok := values[2].(string); ok {
			if v != "" {
				if val, err := strconv.ParseInt(v, 10, 64); err == nil {
					user.ReferralCount = int(val)
				} else {
					log.Printf("Invalid ReferralCount for user %d: %v", id, err)
				}
			}
		}
	}

	if len(values) > 3 {
		if v, ok := values[3].(string); ok {
			if v != "" {
				if val, err := convertToFloat(v); err == nil {
					user.IncomeRate = val
				} else {
					log.Printf("Invalid IncomeRate for user %d: %v", id, err)
				}
			}
		}
	}

	if len(values) > 4 {
		if v, ok := values[4].(string); ok {
			if v != "" {
				if val, err := convertToFloat(v); err == nil {
					user.ReferralTotal = val
				} else {
					log.Printf("Invalid ReferralTotal for user %d: %v", id, err)
				}
			}
		}
	}

	if len(values) > 5 {
		if v, ok := values[5].(string); ok {
			if v != "" {
				if val, err := convertToFloat(v); err == nil {
					user.TotalBonus = val
				} else {
					log.Printf("Invalid TotalBonus for user %d: %v", id, err)
				}
			}
		}
	}

	if len(values) > 6 {
		if v, ok := values[6].(string); ok {
			if v != "" {
				if val, err := convertToFloat(v); err == nil {
					user.BonusToWithdraw = val
				} else {
					log.Printf("Invalid BonusToWithdraw for user %d: %v", id, err)
				}
			}
		}
	}

	// Update or create the user record in the database.
	db.Save(&user)
}

func UpdatePartnersSheet(db *gorm.DB, sheetID, writeRange string) {
	var users []models.User
	db.Select("user_id").Find(&users)

	var values [][]interface{}
	for _, user := range users {
		values = append(values, []interface{}{user.UserID})
	}

	if err := UpdateGoogleSheet(sheetID, writeRange, values); err != nil {
		log.Println("Error updating partner sheet:", err)
	}
}

func UpdateCRMReferralsSheet(db *gorm.DB, sheetID, writeRange string) {
	var referrals []models.Referral
	db.Select("user_id, referred_by").Find(&referrals)

	var values [][]interface{}
	for _, ref := range referrals {
		values = append(values, []interface{}{ref.UserID, ref.ReferredBy})
	}

	if err := UpdateGoogleSheet(sheetID, writeRange, values); err != nil {
		log.Println("Error updating CRM referrals sheet:", err)
	}
}

func UpdateReferralTradeAmount(db *gorm.DB, sheetID, readRange string) {
	values, err := GetGoogleSheetData(sheetID, readRange)
	if err != nil {
		log.Println("Error reading Google Sheets:", err)
		return
	}

	for i, row := range values {
		if len(row) >= 3 {
			userIDStr, ok := row[0].(string)
			if !ok {
				log.Printf("Invalid UserID format in row %d", i+1)
				continue
			}

			userID, err := strconv.ParseInt(userIDStr, 10, 64)
			if err != nil {
				log.Printf("Invalid UserID in row %d: %v", i+1, err)
				continue
			}

			tradeAmountStr, ok := row[2].(string)
			if !ok {
				log.Printf("Invalid TradeAmount format in row %d", i+1)
				continue
			}

			tradeAmountStr = strings.ReplaceAll(tradeAmountStr, "฿", "")
			tradeAmount, err := strconv.ParseFloat(tradeAmountStr, 64)
			if err != nil {
				log.Printf("Invalid TradeAmount in row %d: %v", i+1, err)
				continue
			}

			var referral models.Referral
			if err := db.First(&referral, "user_id = ?", userID).Error; err != nil {
				log.Printf("Referral not found for UserID %d in row %d", userID, i+1)
				continue
			}

			referral.TradeAmount = tradeAmount

			if err := db.Save(&referral).Error; err != nil {
				log.Printf("Failed to update TradeAmount for UserID %d in row %d: %v", userID, i+1, err)
			}
		}
	}
}

func StartUpdateRoutine(db *gorm.DB, sheetID1, range1, sheetID2, range2, sheetID3, range3 string) {
	ticker := time.NewTicker(5 * time.Second)
	log.Printf("Starting scheduled update every 5 seconds...")
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Starting scheduled update...")

		values, err := GetGoogleSheetData(sheetID1, range1)
		if err != nil {
			log.Println("Error reading Google Sheets:", err)
			continue
		}

		for i, row := range values {
			if i == 1 {
				continue
			}

			if len(row) > 1 {
				userID := row[1].(string)
				UpdateUserFromSheet(db, userID, row)
			}
		}

		// Update Partner IDs in "УПРАВЛЕНИЕ"
		UpdatePartnersSheet(db, sheetID2, range2)

		// Update Referrals and Partners in "БАЗА СДЕЛОК, CRM"
		UpdateCRMReferralsSheet(db, sheetID3, range3)

		UpdateReferralTradeAmount(db, sheetID3, range3)
	}
}
