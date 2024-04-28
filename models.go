package main

type DataFile struct {
	Devices []DeviceDefinition `yaml:"devices"`
}

type DeviceType string

const (
	DeviceRouter DeviceType = "r"
	DeviceHost   DeviceType = "h"
)

type DeviceDefinition struct {
	DeviceName string            `yaml:"name"`
	DeviceType DeviceType        `yaml:"type"`
	Interfaces map[string]string `yaml:"interfaces"`
	Table      []string          `yaml:"forwarding_table"`
}
