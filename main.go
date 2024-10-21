package main

import (
    "fmt"
    "log"
    "time"

    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
    "fyne.io/fyne/v2"
    "github.com/tarm/serial" // Use only this serial package
)

// GetAvailablePorts lists all available COM ports on the system
func GetAvailablePorts() ([]string, error) {
    // Hardcoding common COM port names for simplicity (adjust for your platform)
    return []string{"COM1", "COM2", "COM3", "COM4"}, nil
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

    // Get all available COM ports
    availablePorts, err := GetAvailablePorts()
    if err != nil {
        log.Fatalf("Error fetching COM ports: %v", err)
    }

    // Define common baud rates
    baudRates := []string{
        "9600", "14400", "19200", "38400", "57600", "115200", "128000", "256000",
    }

    // Dropdowns for COM ports and baud rates
    portDropdown := widget.NewSelect(availablePorts, func(value string) {})
    portDropdown.PlaceHolder = "Select COM port"

    baudDropdown := widget.NewSelect(baudRates, func(value string) {})
    baudDropdown.PlaceHolder = "Select Baud Rate"

    infoLabel := widget.NewLabel("Device Info will be displayed here")

    readButton := widget.NewButton("Read Info", func() {
        port := portDropdown.Selected
        baudRate := baudDropdown.Selected

        if port == "" || baudRate == "" {
            infoLabel.SetText("Please select both COM port and Baud Rate")
            return
        }

        // Convert baud rate from string to int
        var baud int
        fmt.Sscanf(baudRate, "%d", &baud)

        // Read device information
        info, err := readDeviceInfo(port, baud)
        if err != nil {
            infoLabel.SetText("Error: " + err.Error())
        } else {
            infoLabel.SetText("Device Info: " + info)
        }
    })

    content := container.NewVBox(
        widget.NewLabel("MediaTek (MKT) Device Info Tool"),
        widget.NewLabel("Select COM Port:"),
        portDropdown,
        widget.NewLabel("Select Baud Rate:"),
        baudDropdown,
        readButton,
        infoLabel,
    )

    w.SetContent(content)
    w.Resize(fyne.NewSize(400, 300))
    w.ShowAndRun()
}
