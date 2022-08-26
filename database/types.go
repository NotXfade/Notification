package database

import "time"

// Accounts is a strucutre to stores user information
type Accounts struct {
	// auto generated
	ID int `json:"id" gorm:"primary_key"`
	//If password is empty then the user registered from career forms.
	//User can register password in future using forgot password.
	Password string `json:"-"`
	//email of user
	Email string `json:"email" gorm:"not null;unique;"`
	//name of user
	Name string `json:"name"`
	// contact number of user
	ContactNo string `json:"contact_no"`
	//this role id is used for auth portal management
	//value can be 'admin', 'user'
	Role string `json:"role"`
	//Level wil be the level of employees like L1 intern, L2 intern
	Level string `json:"level"`
	//account status can be active, new, deleted, blocked etc.
	AccountStatus string `json:"account_status"`
	//account creation date
	Timestamp int64     `json:"timestamp"`
	UpdatedAt time.Time `json:"-"`
	CreatedAt time.Time `json:"-"`
}

// Activities is a structure to record user activties
type Activities struct {
	ID int `json:"-" gorm:"primary_key"`
	//If username entered incorrect then these activities will also be recorded.
	Email string `json:"email" gorm:"index"`
	//can be login, failedlogin, signup
	ActivityName string    `json:"activity_name"`
	ClientIP     string    `json:"client_ip"`
	ClientAgent  string    `json:"client_agent"`
	Timestamp    int64     `json:"timestamp"`
	UpdatedAt    time.Time `json:"-"`
	CreatedAt    time.Time `json:"-"`
}

// ActiveSessions is a structure to stores active sessions
type ActiveSessions struct {
	ID          int `json:"-" gorm:"primary_key"`
	SessionID   string
	Userid      int `gorm:"index"`
	ClientAgent string
	Start       int64
	// if value is '0' then session is remembered.
	End       int64
	UpdatedAt time.Time `json:"-"`
	CreateAt  time.Time `json:"-"`
}

//InviteData defining structure for binding send code again data
type InviteData struct {
	Email []string `json:"email" `
}

//MailData is used to send as payload for notification service
type MailData struct {
	Email string `json:"email"`
	Link  string `json:"link"`
	Task  string `json:"task"`
}

//UploadLinks : This structure is used to save uploaded links
type UploadLinks struct {
	ID       int    `json:"-" gorm:"primary_key"`
	Userid   int    `json:"id" gorm:"index"`
	LinkType string `json:"link_type"`
	Link     string `json:"link"`
	//account creation date
	UpdatedAt time.Time `json:"-"`
	CreatedAt time.Time `json:"-"`
}

//PolicyLinks : This structure is used to store policy links
type PolicyLinks struct {
	ID         int       `json:"-" gorm:"primary_key"`
	AddedBy    int       `json:"id"`
	PolicyName string    `json:"policyname"`
	Level      string    `json:"level"`
	PolicyLink string    `json:"policylink"`
	UpdatedAt  time.Time `json:"-"`
	CreatedAt  time.Time `json:"-"`
}
