package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

//fyne has some requirements that can be found here: https://developer.fyne.io/started/

type ipInfo struct {
	IpAddress string
	Port      string
}

var Quit bool

func main() {
	myApp := app.New()
	w := myApp.NewWindow("Form Widget")
	w.Resize(fyne.NewSize(300, 300))

	ip, p := readFile()
	fmt.Print(ip, p)

	statustext := binding.NewString()
	statustext.Set("Set IP and Port")
	status := widget.NewLabelWithData(statustext)

	level := binding.NewString()
	level.Set("Then Press Start")
	data := widget.NewLabelWithData(level)

	ip_address := widget.NewEntry()
	ip_address.SetText(ip)
	port := widget.NewEntry()
	port.SetText(p)

	form := &widget.Form{}

	form.Append("IP address", ip_address)
	form.Append("Port", port)
	form.Append("Status", status)
	form.Append("Data Packet", data)

	button := widget.NewButton("Start", func() {})
	Quit = true
	button.OnTapped = func() {
		if Quit {
			fmt.Println("Start")
			Quit = false
			theIP := ip_address.Text
			thePort := port.Text
			fmt.Print(theIP, thePort)
			button.SetText("Stop")
			go connect(theIP, thePort, level, statustext, data, status)
			fmt.Println("Start End")
			w.SetContent(container.New(layout.NewVBoxLayout(), form, button))
		} else {
			fmt.Println("Stop")
			Quit = true
			time.Sleep(2 * time.Second)
			button.SetText("Start")
			w.SetContent(container.New(layout.NewVBoxLayout(), form, button))
		}
	}
	w.SetContent(container.New(layout.NewVBoxLayout(), form, button))
	w.ShowAndRun()
}

func SetButtonLabel(b *widget.Button) {

}

func Test() (one string) {
	fmt.Print("Testing")
	one = "One"
	return
}

func saveFile(ip, p string) {
	content := make(map[string]string)
	content["ipAddress"] = ip
	content["port"] = p

	f, err := os.Create(".ipInfo.dll")
	if err != nil {
		fmt.Println(err)
		return
	}

	contentString, _ := json.Marshal(content)
	fmt.Print(contentString)

	l, err := f.WriteString(string(contentString))
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func readFile() (ip, port string) {
	var ipInfo ipInfo
	dflt_ip := "127.0.0.1"
	dflt_port := "40000"
	//if file doesnt exist save new file with default IP and port
	file, err := ioutil.ReadFile(".ipInfo.dll")
	if err != nil {
		fmt.Println("Error: ", err)
		saveFile(dflt_ip, dflt_port)
	}

	//read file info and assign to variables
	data := string(file)
	fmt.Println("Contents of file:", data)

	json.Unmarshal(file, &ipInfo)
	ip = ipInfo.IpAddress
	port = ipInfo.Port
	return
}

func exit() {
	os.Exit(0)
}

func connect(ip, port string, level, status binding.String, levelLabel, statusLabel *widget.Label) {
	level.Set("Test Set")
	arguments := ip + ":" + port
	fmt.Print(arguments)
	if len(arguments) == 1 {
		fmt.Println("IP/Port Error")
		setStatus("IP/Port Error", status)
		return
	}

	CONNECT := arguments
	d := net.Dialer{Timeout: 1 * time.Second}
	c, err := d.Dial("tcp", CONNECT)
	if err != nil {
		fmt.Println(err)
		setStatus("Network Error\nCheck IP/Port", status)
		return
	}

	saveFile(ip, port)

	fmt.Fprintf(c, "start")
	setStatus("Starting", status)
	count := 0
	for {
		text := "continue"
		fmt.Fprintf(c, text+"\n")

		message, _ := bufio.NewReader(c).ReadString('\n')
		message2, _ := bufio.NewReader(c).ReadString('\n')
		fmt.Print("->: " + message)
		setStatus("Receiving Data", status)
		if len(message) > 0 {
			msg := saveJSON(message2)
			fmt.Println(msg)
		} else {
			setStatus("Q-Sys Script Err", status)
			fmt.Println("0 Bytes Received - Q-Sys Error")
			return
		}

		if Quit == true {
			fmt.Print("Quitting 2")
			setStatus("Stopped", status)
			return
		}
		count++

		if count%50 == 0 {
			level.Set(strconv.Itoa(count))
			levelLabel.Refresh()
		}
		statusLabel.Refresh()
	}
}

var GlobalStatus string

func setStatus(s string, status binding.String) {
	if s == GlobalStatus {
		return
	} else {
		GlobalStatus = s
		status.Set(s)
	}
}

func saveJSON(m string) (msg string) {
	f, err := os.Create("MeterData.json")
	if err != nil {
		fmt.Println(err)
		msg = "Error saving file"
		return
	}

	l, err := f.WriteString(m)
	if err != nil {
		fmt.Println(err)
		msg = "Error writing file"
		f.Close()
		return
	}
	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		msg = "Error closing file"
		return
	}
	msg = "Success"
	return
}
