package models

// Pref is the user's preferences
type Pref struct {
	UserID       int64 `gorm:"primary_key"`
	EmailOnError bool
}
