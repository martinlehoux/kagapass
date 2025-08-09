package database

import (
	"time"

	"github.com/martinlehoux/kagapass/internal/types"
)

// CreateTestEntries creates sample entries for testing purposes
func CreateTestEntries() []types.Entry {
	now := time.Now()
	
	return []types.Entry{
		{
			Title:    "GitHub Personal",
			Username: "kagamino",
			Password: "super_secret_password_123",
			URL:      "https://github.com",
			Notes:    "Personal GitHub account with 2FA enabled",
			Group:    "Personal/Development",
			Created:  now.Add(-30 * 24 * time.Hour),
			Modified: now.Add(-7 * 24 * time.Hour),
		},
		{
			Title:    "GitHub Work",
			Username: "kagamino.work",
			Password: "work_password_456",
			URL:      "https://github.com",
			Notes:    "Work GitHub account",
			Group:    "Work/Development",
			Created:  now.Add(-60 * 24 * time.Hour),
			Modified: now.Add(-14 * 24 * time.Hour),
		},
		{
			Title:    "Gmail",
			Username: "kagamino@gmail.com",
			Password: "email_password_789",
			URL:      "https://gmail.com",
			Notes:    "Main email account",
			Group:    "Personal/Email",
			Created:  now.Add(-365 * 24 * time.Hour),
			Modified: now.Add(-30 * 24 * time.Hour),
		},
		{
			Title:    "Bank of Example",
			Username: "customer123",
			Password: "banking_password_000",
			URL:      "https://bankofexample.com",
			Notes:    "Online banking - remember security questions",
			Group:    "Personal/Banking",
			Created:  now.Add(-180 * 24 * time.Hour),
			Modified: now.Add(-60 * 24 * time.Hour),
		},
		{
			Title:    "VPN Service",
			Username: "vpnuser",
			Password: "vpn_password_111",
			URL:      "https://vpnservice.com",
			Notes:    "Premium subscription expires Dec 2024",
			Group:    "Personal/Services",
			Created:  now.Add(-90 * 24 * time.Hour),
			Modified: now.Add(-15 * 24 * time.Hour),
		},
		{
			Title:    "Work Database",
			Username: "admin",
			Password: "db_password_222",
			URL:      "https://db.company.com/admin",
			Notes:    "Production database - handle with care!\nConnection: postgres://db.company.com:5432/main",
			Group:    "Work/Databases",
			Created:  now.Add(-120 * 24 * time.Hour),
			Modified: now.Add(-3 * 24 * time.Hour),
		},
	}
}

// GetTestEntries is an alias for CreateTestEntries for backward compatibility
func GetTestEntries() []types.Entry {
	return CreateTestEntries()
}