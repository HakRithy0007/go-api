package share

import (
	"time"
)

type UserContext struct {
	UserID       float64   `json:"user_id"`
	UserName     string    `json:"user_name"`
	LoginSession string    `json:"login_session"`
	Exp          time.Time `json:"exp"`
	UserAgent    string    `json:"user_agent"`
	Ip           string    `json:"ip"`
	MembershipId float64   `json:"membership_id"`
	RoleID       int   `json:"role_id"`
}

type Paging struct {
	Page    int `json:"page" query:"page" validate:"required,min=1"`
	Perpage int `json:"per_page" query:"per_page" validate:"required,min=1"`
}

type Sort struct {
	Property  string `json:"property" validate:"required"`
	Direction string `json:"direction" validate:"required,oneof=asc desc"`
}
type Filter struct {
	Property string      `json:"property" validate:"required"`
	Value    interface{} `json:"value" validate:"required"`
}

type Platform struct {
	ID                   uint64 `json:"id"`
	PlatformName         string `json:"platform_name"`
	PlatformHost         string `json:"platform_host"`
	PlatformToken        string `json:"platform_token"`
	PlatformExtraPayload string `json:"platform_extra_payload"`
	InternalToken        string `json:"internal_token"`
	StatusID             uint64 `json:"status_id"`
	Order                uint64 `json:"order"`
}

type Status struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// StatusData contains predefined status values.
var StatusData = []Status{
	{ID: 1, Name: "Active"},
	{ID: 2, Name: "Inactive"},
	{ID: 3, Name: "Suspended"},
	{ID: 4, Name: "Deleted"},
}