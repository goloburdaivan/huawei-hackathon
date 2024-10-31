package services

import (
	"fmt"
	"github.com/go-toast/toast"
)

func SendToastNotification(message string) {
	notification := toast.Notification{
		AppID:   "Network stats",
		Title:   "Сетевое уведомление",
		Message: message,
	}
	err := notification.Push()
	if err != nil {
		fmt.Println("Ошибка отправки уведомления:", err)
	}
}
