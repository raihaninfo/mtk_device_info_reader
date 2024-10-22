package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/tarm/serial"
)

// GetAvailablePorts detects COM ports by scanning common port names
func GetAvailablePorts() []string {
	ports := []string{}
	// Scan through a range of typical COM port names
	for i := 1; i <= 256; i++ {
		port := fmt.Sprintf("COM%d", i)
		_, err := serial.OpenPort(&serial.Config{Name: port, Baud: 9600})
		if err == nil {
			ports = append(ports, port)
		}
	}
	return ports
}

// DetectBaudRate tries common baud rates and returns the first successful one
func DetectBaudRate(portName string) (int, error) {
	commonBaudRates := []int{
		9600, 14400, 19200, 38400, 57600, 115200, 128000, 256000,
	}

	for _, baud := range commonBaudRates {
		config := &serial.Config{
			Name:        portName,
			Baud:        baud,
			ReadTimeout: time.Second * 2,
		}

		port, err := serial.OpenPort(config)
		if err != nil {
			continue // Try the next baud rate
		}
		defer port.Close()

		// Send a test command (modify according to your device's protocol)
		_, err = port.Write([]byte("AT\r")) // Example AT command for testing
		if err != nil {
			continue
		}

		buf := make([]byte, 128)
		_, err = port.Read(buf)
		if err == nil {
			// Successfully communicated with the device
			return baud, nil
		}
	}

	return 0, fmt.Errorf("could not detect baud rate")
}

func readDeviceInfo(portName string, baudRate int) (string, error) {
	config := &serial.Config{
		Name:        portName,
		Baud:        baudRate,
		ReadTimeout: time.Second * 2,
	}

	port, err := serial.OpenPort(config)
	if err != nil {
		return "", err
	}
	defer port.Close()

	// Send an example command to the device (adjust based on your device's protocol)
	_, err = port.Write([]byte("AT+DEVICEINFO\r")) // Example AT command, modify if necessary
	if err != nil {
		return "", err
	}

	buf := make([]byte, 128)
	n, err := port.Read(buf)
	if err != nil {
		return "", err
	}

	return string(buf[:n]), nil
}

func main() {
	a := app.New()
	w := a.NewWindow("MKT Device Info Reader")

	// Automatically detect available COM ports
	availablePorts := GetAvailablePorts()

	// if len(availablePorts) == 0 {
	//     log.Println("No COM ports detected")
	//     return
	// }

	// GUI Elements
	portDropdown := widget.NewSelect(availablePorts, func(value string) {})
	portDropdown.PlaceHolder = "Select COM port"

	infoLabel := widget.NewLabel("Device Info will be displayed here")

	readButton := widget.NewButton("Read Info", func() {
		port := portDropdown.Selected

		if port == "" {
			infoLabel.SetText("Please select a COM port")
			return
		}

		// Automatically detect baud rate
		baudRate, err := DetectBaudRate(port)
		if err != nil {
			infoLabel.SetText("Error detecting baud rate: " + err.Error())
			return
		}

		// Read device information
		info, err := readDeviceInfo(port, baudRate)
		if err != nil {
			infoLabel.SetText("Error: " + err.Error())
		} else {
			infoLabel.SetText("Device Info: " + info)
		}
	})

	content := container.NewVBox(
		widget.NewLabel("MediaTek (MKT) Device Info Tool"),
		widget.NewLabel("Detected COM Ports:"),
		portDropdown,
		readButton,
		infoLabel,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(400, 300))
	w.ShowAndRun()
}
