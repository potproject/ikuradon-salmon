package notification

// N WebPush Notification Body
type N struct {
	NotificationID   int    `json:"notification_id"`
	AccessToken      string `json:"access_token"`
	PreferredLocale  string `json:"preferred_locale"`
	NotificationType string `json:"notification_type"`
	Icon             string `json:"icon"`
	Title            string `json:"title"`
	Body             string `json:"body"`
}
