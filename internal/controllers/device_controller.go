package controllers

import (
	"Hackathon/internal/core/ssh"
	"Hackathon/internal/views"
	"log"
)

type DeviceController struct {
	parseDeviceInfo *ssh.SshService
	stopChannel     chan bool
}

func NewDeviceController(parseDeviceInfo *ssh.SshService) *DeviceController {
	return &DeviceController{
		parseDeviceInfo: parseDeviceInfo,
		stopChannel:     make(chan bool),
	}
}

func (pc *DeviceController) ShowDeviceInfo() {
	deviceStatus, err := pc.parseDeviceInfo.ParseDeviceStatus()
	if err != nil {
		log.Println("Ошибка при получении статуса устройства: %v", err)
		return
	}
	views.DisplayDeviceInfo(&deviceStatus)
}