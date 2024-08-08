package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username        string
	UserID          int64
	IncomeRate      float64
	ReferralCount   int
	ReferralTrades  int
	ReferralTotal   float64
	FirstTradeDate  time.Time
	LastPayoutDate  time.Time
	TotalBonus      float64
	BonusToWithdraw float64
}

type Referral struct {
	gorm.Model
	UserID      int64
	UserName    string
	TradeAmount float64
	ReferredBy  int64
}
