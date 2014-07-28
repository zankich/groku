package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/codegangsta/cli"
)

const VERSION = "0.3"

var CONFIG string

type dictonary struct {
	XMLName xml.Name `xml:"apps"`
	Apps    []app    `xml:"app"`
}

type app struct {
	Name string `xml:",chardata"`
	ID   string `xml:"id,attr"`
}

func main() {
	CONFIG = fmt.Sprintf("%v/groku", os.TempDir())
	app := cli.NewApp()
	app.Name = "groku"
	app.Version = VERSION
	app.Usage = "roku CLI remote"
	app.Commands = commands()
	app.Run(os.Args)
}

func queryApps() dictonary {
	resp, _ := http.Get(fmt.Sprintf("%vquery/apps", findRoku()))
	body := make([]byte, 2048)
	n, _ := resp.Body.Read(body)

	var dict dictonary
	if err := xml.Unmarshal(body[:n], &dict); err != nil {
		log.Fatalln(err)
	}
	return dict
}

func findRoku() string {
	fi, err := os.Open(CONFIG)
	defer fi.Close()
	if err != nil {
		ssdp, _ := net.ResolveUDPAddr("udp", "239.255.255.250:1900")
		addr, _ := net.ResolveUDPAddr("udp", ":0")
		socket, _ := net.ListenUDP("udp", addr)

		socket.WriteToUDP([]byte("M-SEARCH * HTTP/1.1\r\n"+
			"HOST: 239.255.255.250:1900\r\n"+
			"MAN: \"ssdp:discover\"\r\n"+
			"ST: roku:ecp\r\n"+
			"MX: 3 \r\n\r\n"), ssdp)

		answerBytes := make([]byte, 1024)
		socket.SetReadDeadline(time.Now().Add(3 * time.Second))
		_, _, err := socket.ReadFromUDP(answerBytes[:])
		if err != nil {
			fmt.Println("Could not find your Roku!")
			os.Exit(1)
		}

		ret := strings.Split(string(answerBytes), "\r\n")
		location := strings.TrimPrefix(ret[len(ret)-3], "LOCATION: ")

		fi, err := os.Create(CONFIG)
		defer fi.Close()
		if err != nil {
			return location
		}
		fi.Write([]byte(location))
		return location
	}
	buf := make([]byte, 1024)
	n, err := fi.Read(buf[:])
	if err != nil {
		os.Remove(CONFIG)
		return findRoku()
	} else {
		return string(buf[:n])
	}
}

func commands() []cli.Command {
	cmds := []cli.Command{}
	for _, cmd := range []string{
		"home",
		"rev",
		"fwd",
		"select",
		"left",
		"right",
		"down",
		"up",
		"back",
		"info",
		"backspace",
		"enter",
		"search",
	} {
		cmds = append(cmds, cli.Command{
			Name:  cmd,
			Usage: cmd,
			Action: func(c *cli.Context) {
				http.PostForm(fmt.Sprintf("%vkeypress/%v", findRoku(), c.Command.Name), nil)
			},
		})
	}
	cmds = append(cmds, cli.Command{
		Name:  "replay",
		Usage: "replay",
		Action: func(c *cli.Context) {
			http.PostForm(fmt.Sprintf("%vkeypress/%v", findRoku(), "InstantReplay"), nil)
		},
	})
	cmds = append(cmds, cli.Command{
		Name:  "play",
		Usage: "play/pause",
		Action: func(c *cli.Context) {
			http.PostForm(fmt.Sprintf("%vkeypress/%v", findRoku(), "Play"), nil)
		},
	})
	cmds = append(cmds, cli.Command{
		Name:  "discover",
		Usage: "discover roku on your local network",
		Action: func(c *cli.Context) {
			os.Remove(CONFIG)
			fmt.Println("Found roku at", findRoku())
		},
	})
	cmds = append(cmds, cli.Command{
		Name:  "text",
		Usage: "send text to the roku",
		Action: func(c *cli.Context) {
			roku := findRoku()
			for _, c := range c.Args()[0] {
				http.PostForm(fmt.Sprintf("%vkeypress/Lit_%v", roku, string(c)), nil)
			}
		},
	})
	cmds = append(cmds, cli.Command{
		Name:  "apps",
		Usage: "list installed apps on roku",
		Action: func(c *cli.Context) {
			dict := queryApps()
			fmt.Println("Installed apps:")
			for _, a := range dict.Apps {
				fmt.Println(a.Name)
			}
		},
	})
	cmds = append(cmds, cli.Command{
		Name:  "app",
		Usage: "launch specified app",
		Action: func(c *cli.Context) {
			dict := queryApps()
			for _, a := range dict.Apps {
				if a.Name == c.Args()[0] {
					http.PostForm(fmt.Sprintf("%vlaunch/%v", findRoku(), a.ID), nil)
					return
				}
			}
			fmt.Println("App not found!")
			os.Exit(1)
		},
	})
	return cmds
}
