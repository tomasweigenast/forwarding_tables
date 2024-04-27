package main

type JsonFile []DeviceDefinition

type DeviceType string

const (
	DeviceRouter DeviceType = "r"
	DeviceHost   DeviceType = "h"
)

type DeviceDefinition struct {
	DeviceName string                 `json:"name"`
	DeviceType DeviceType             `json:"type"`
	Interfaces map[string]string      `json:"interfaces"`
	Table      []ForwardingTableEntry `json:"table"`
}

type ForwardingTableEntry struct {
	Destination string `json:"destination"`
	NextHop     string `json:"next_hop"`
	Interface   string `json:"interface"`
}
