package easyfri

import (
	"fmt"
	"github.com/eriklupander/tradfri-go/tradfri"
	"github.com/sirupsen/logrus"
	"regexp"
)

var (
	tc *tradfri.Client
)

type Device struct {
	Name       string
	Id         int
	Type       string
	Alive      int
	State      bool
	Dimmer     int
	RGBHex     string
	Temperatur int
}

type Group struct {
	Name    string
	Id      int
	Devices []int
}

func CreateClient(ip string, user_id string, psk string) {
	tc = tradfri.NewTradfriClient(ip, user_id, psk)
}

func SetPower(device_id int, state bool) Device {

	var tmp int

	if state {
		tmp = 1
	} else {
		tmp = 0
	}

	_, err := tc.PutDevicePower(device_id, tmp)

	if err != nil {
		logrus.Debug("%v", err)
	}

	return GetDevice(device_id)

}

func SetRGB(device_id int, value string) Device {

	_, err := tc.PutDeviceColorRGB(device_id, value)

	if err != nil {
		logrus.Debug("%v", err)
	}

	return GetDevice(device_id)
}

func SetDim(device_id int, value int) Device {

	_, err := tc.PutDeviceDimming(device_id, value)

	if err != nil {
		logrus.Debug("%v", err)
	}

	return GetDevice(device_id)
}

func SetDimForGroupByName(regex string, value int) []Device {
	groups := GetGroups()

	var devices []Device
	for _, g := range groups {
		tmp, _ := regexp.MatchString(regex, g.Name)
		if tmp {
			for _, d := range g.Devices {
				devices = append(devices, SetDim(d, value))
			}
		}
	}
	return devices
}

func SetDimForGroup(group_id int, value int) []Device {
	groups := GetGroups()

	var devices []Device
	for _, g := range groups {
		if g.Id == group_id {
			for _, d := range g.Devices {
				devices = append(devices, SetDim(d, value))
			}
		}
	}

	return devices
}

func SetPowerForGroupByName(regex string, state bool) []Device {
	groups := GetGroups()

	var devices []Device
	for _, g := range groups {
		tmp, _ := regexp.MatchString(regex, g.Name)
		if tmp {
			for _, d := range g.Devices {
				devices = append(devices, SetPower(d, state))
			}
		}
	}
	return devices
}

func SetPowerForGroup(group_id int, state bool) []Device {
	groups := GetGroups()

	var devices []Device
	for _, g := range groups {
		if g.Id == group_id {
			for _, d := range g.Devices {
				devices = append(devices, SetPower(d, state))
			}
		}
	}

	return devices

}

func SwitchPowerForGroupByName(regex string) []Device {
	groups := GetGroups()

	var devices []Device
	for _, g := range groups {
		tmp, _ := regexp.MatchString(regex, g.Name)
		if tmp {
			for _, d := range g.Devices {
				devices = append(devices, SwitchPower(d))
			}
		}
	}
	return devices
}

func SwitchPowerForGroup(group_id int) []Device {
	groups := GetGroups()

	var devices []Device
	for _, g := range groups {
		if g.Id == group_id {
			for _, d := range g.Devices {
				devices = append(devices, SwitchPower(d))
			}
		}
	}

	fmt.Println(devices)

	return devices
}

func SwitchPower(device_id int) Device {

	device := GetDevice(device_id)

	device = SetPower(device_id, !device.State)

	return device
}

func GetGroups() []Group {

	rooms, err := tc.ListGroups()

	if err != nil {
		logrus.Info(err)
	}

	groups := []Group{}

	for _, r := range rooms {
		groups = append(groups, Group{r.Name, r.DeviceId, r.Content.DeviceList.DeviceIds})
	}

	return groups
}

func GetDevice(device_id int) Device {

	device, _ := tc.GetDevice(device_id)

	d := Device{}

	d.Type = device.Metadata.TypeName
	d.Name = device.Name
	d.Id = device.DeviceId
	d.Alive = device.Alive

	if !(d.Type == "TRADFRI remote control") {

		d.State = !(device.LightControl[0].Power == 0)
		d.Dimmer = device.LightControl[0].Dimmer
		d.RGBHex = device.LightControl[0].RGBHex
		d.Temperatur = device.LightControl[0].ColorTemperature
	}

	return d
}
