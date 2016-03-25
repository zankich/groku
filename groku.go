package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"
)

var CONFIG string

const (
	VERSION = "0.4"
	USAGE   = `usage: groku [--version] [--help] <command> [<args>]

CLI remote for your Roku

Commands:
	home            Return to the home screen
	rev             Reverse
	fwd             Fast Forward
	select          Select
	left            Left
	right           Right
	up              Up
	down            Down
	back            Back
	info            Info
	backspace       Backspace
	enter           Enter
	search          Search
	replay          Replay
	play            Play
	pause           Pause
	discover        Discover Roku devices on your local network
	list            List known Roku devices
	use             Set Roku name to use
	device-info     Display device info
	text            Send text to the Roku
	apps            List installed apps on your Roku
	app             Launch specified app
	`
)

type dictionary struct {
	XMLName xml.Name `xml:"apps"`
	Apps    []app    `xml:"app"`
}

type deviceinfo struct {
	XMLName    xml.Name `xml:"device-info"`
	UDN        string   `xml:"udn"`
	Serial     string   `xml:"serial-number"`
	DeviceID   string   `xml:"device-id"`
	ModelNum   string   `xml:"model-number"`
	ModelName  string   `xml:"model-name"`
	DeviceName string   `xml:"user-device-name"`
}

type app struct {
	Name string `xml:",chardata"`
	ID   string `xml:"id,attr"`
}

type roku struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

type grokuConfig struct {
	LastName  string `json:"lastname"`
	Current   roku   `json:"current"`
	Rokus     []roku `json:"rokus"`
	Timestamp int64  `json:"timestamp"`
}

func main() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println("Cannot find home directory")
		os.Exit(1)
	}
	CONFIG = fmt.Sprintf("%s/groku.json", home)

	if len(os.Args) == 1 || os.Args[1] == "--help" || os.Args[1] == "-help" ||
		os.Args[1] == "--h" || os.Args[1] == "-h" || os.Args[1] == "help" {
		fmt.Println(USAGE)
		os.Exit(0)
	}

	if len(os.Args) == 2 && (os.Args[1] == "-v" || os.Args[1] == "--version" ||
		os.Args[1] == "--version") {
		fmt.Printf("groku version %s\n", VERSION)
		os.Exit(0)
	}

	switch os.Args[1] {
	case "home", "rev", "fwd", "select", "left", "right", "down", "up",
		"back", "info", "backspace", "enter", "search":
		http.PostForm(fmt.Sprintf("%vkeypress/%v", getCurrentRokuAddress(), os.Args[1]), nil)
		os.Exit(0)
	case "replay":
		http.PostForm(fmt.Sprintf("%vkeypress/%v", getCurrentRokuAddress(), "InstantReplay"), nil)
		os.Exit(0)
	case "play", "pause":
		http.PostForm(fmt.Sprintf("%vkeypress/%v", getCurrentRokuAddress(), "Play"), nil)
		os.Exit(0)
	case "discover":
		config := getRokuConfig()
		if len(config.Rokus) > 0 {
			for _, r := range config.Rokus {
				fmt.Print("Found roku at ", r.Address)
				if r.Name != "" {
					fmt.Print(" named ", r.Name)
				}
				fmt.Println()
			}
		}
		os.Exit(0)
	case "list":
		config := getRokuConfig()
		for _, r := range config.Rokus {
			if r.Name != "" {
				fmt.Print(r.Name, ": ")
			}
			fmt.Println(r.Address)
		}
	case "use":
		config := getRokuConfig()
		for _, r := range config.Rokus {
			if strings.ToUpper(os.Args[2]) == strings.ToUpper(r.Name) {
				config.Current = r
				config.LastName = os.Args[2]
				writeConfig(config)
				fmt.Printf("Using Roku named %v at %v", r.Name, r.Address)
				os.Exit(0)
			}
		}
		fmt.Printf("Cannot find Roku named %v\n", os.Args[2])
	case "device-info":
		var info = queryInfo()
		if getCurrentRokuName() != "" {
			fmt.Printf("Name:\t\t%v\n", info.DeviceName)
		}
		fmt.Printf("Model:\t\t%v %v\n", info.ModelName, info.ModelNum)
		fmt.Printf("Serial:\t\t%v\n", info.Serial)
	case "text":
		if len(os.Args) < 3 {
			fmt.Println(USAGE)
			os.Exit(1)
		}

		roku := getCurrentRokuAddress()
		for _, c := range os.Args[2] {
			http.PostForm(fmt.Sprintf("%skeypress/Lit_%s", roku, url.QueryEscape(string(c))), nil)
		}
		os.Exit(0)
	case "apps":
		dict := queryApps()
		for _, a := range dict.Apps {
			fmt.Println(a.Name)
		}
		os.Exit(0)
	case "app":
		if len(os.Args) < 3 {
			fmt.Println(USAGE)
			os.Exit(1)
		}

		dict := queryApps()

		for _, a := range dict.Apps {
			if a.Name == os.Args[2] {
				http.PostForm(fmt.Sprintf("%vlaunch/%v", getCurrentRokuAddress(), a.ID), nil)
				os.Exit(0)
			}
		}
		fmt.Printf("App %q not found\n", os.Args[2])
		os.Exit(1)
	default:
		fmt.Println(USAGE)
		os.Exit(1)
	}
}

func queryApps() dictionary {
	resp, err := http.Get(fmt.Sprintf("%squery/apps", getCurrentRokuAddress()))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	var dict dictionary
	if err := xml.NewDecoder(resp.Body).Decode(&dict); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return dict
}

func queryInfoForAddress(address string) deviceinfo {
	resp, err := http.Get(fmt.Sprintf("%squery/device-info", address))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	var info deviceinfo
	if err := xml.NewDecoder(resp.Body).Decode(&info); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return info
}

func queryInfo() deviceinfo {
	return queryInfoForAddress(getCurrentRokuAddress())
}

func findRokus() []roku {
	ssdp, err := net.ResolveUDPAddr("udp", "239.255.255.250:1900")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	addr, err := net.ResolveUDPAddr("udp", ":0")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	socket, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, err = socket.WriteToUDP([]byte("M-SEARCH * HTTP/1.1\r\n"+
		"HOST: 239.255.255.250:1900\r\n"+
		"MAN: \"ssdp:discover\"\r\n"+
		"ST: roku:ecp\r\n"+
		"MX: 3 \r\n\r\n"), ssdp)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var rokus []roku
	listentimer := time.Now().Add(5 * time.Second)
	for time.Now().Before(listentimer) {
		answerBytes := make([]byte, 1024)
		err = socket.SetReadDeadline(listentimer)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		_, _, err = socket.ReadFromUDP(answerBytes[:])

		if err == nil {
			ret := strings.Split(string(answerBytes), "\r\n")
			location := strings.TrimPrefix(ret[len(ret)-3], "LOCATION: ")
			name := queryInfoForAddress(location).DeviceName

			r := roku{Name: name, Address: location}
			rokus = append(rokus, r)
		}
	}

	return rokus
}

func getCurrentRokuAddress() string {
	return getRokuConfig().Current.Address
}

func getCurrentRokuName() string {
	return getRokuConfig().Current.Name
}

func getRokuConfigFor(name string) (*roku, error) {
	config := getRokuConfig()
	for _, e := range config.Rokus {
		if strings.ToUpper(e.Name) == strings.ToUpper(name) {
			return &e, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("%v not found", name))
}

func getRokuConfig() grokuConfig {
	var configFile *os.File
	var config grokuConfig

	configFile, err := os.Open(CONFIG)

	// the config file doesn't exist, but that's okay
	if err != nil {
		config.Rokus = findRokus()
		config.Timestamp = time.Now().Unix()
	} else {
		// the config file exists
		if err := json.NewDecoder(configFile).Decode(&config); err != nil {
			config.Rokus = findRokus()
		}

		//if the config file is over 60 seconds old, then replace it
		if config.Timestamp == 0 || time.Now().Unix()-config.Timestamp > 60 {
			config.Rokus = findRokus()
			config.Timestamp = time.Now().Unix()
		}
	}

	if config.LastName != "" {
		found := false
		for _, e := range config.Rokus {
			if strings.ToUpper(e.Name) == strings.ToUpper(config.LastName) {
				config.Current = e
				found = true
			}
		}
		if !found {
			config.Current = config.Rokus[0]
			fmt.Printf("Previously used Roku %v not found anymore, using %v as new default", config.LastName, config.Current.Name)
		}
	} else {
		config.Current = config.Rokus[0]
	}
	writeConfig(config)
	return config
}

func writeConfig(config grokuConfig) error {
	if b, err := json.Marshal(config); err == nil {
		ioutil.WriteFile(CONFIG, b, os.ModePerm)
	}

	return nil
}
