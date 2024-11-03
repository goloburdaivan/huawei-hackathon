package services

import (
	"fmt"
	"github.com/go-toast/toast"
)

func SendToastNotification(message string) {
	notification := toast.Notification{
		AppID:   "Network stats",
		Title:   "Network Notification",
		Message: message,
	}
	err := notification.Push()
	if err != nil {
		fmt.Println("Error sending notification:", err)
	}
}
