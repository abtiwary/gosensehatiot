package main

/*
 * gosensehatiot.go
 * A Websocket based IoT framework for the Raspberry Pi and Sense Hat in Golang!
 *
 * Principal author(s) : Abhishek Tiwary
 *                       abhishek.tiwary@dolby.com
 *
 */

import (
	"fmt"
	"flag"
	"os"
	"bytes"
	"os/signal"
	"syscall"
	"net/http"
	"time"
	"context"
	"html/template"
	"encoding/json"
	
	"GoSenseHatIoT/SenseHatIoT"
	"github.com/gorilla/websocket"
)

// types
type ServerInfo struct {
	ServerIP string
	ServerPort string
}

type Message struct {
	Type string `json:"type"`
	Timestamp string `json:"timestamp,omitempty"`
	Text string `json:"text,omitempty"`
}

type TemplateDict map[string]interface{}


// channels
var exit_bool_channel = make(chan bool)
var exit_signal_channel = make(chan os.Signal, 1)
var exit_worker_channel = make(chan bool)
var measurements_channel = make(chan SenseHatIoT.Measurement, 8)

// globals
var server_info = ServerInfo{}

// state information
var measurements []SenseHatIoT.Measurement
var clients = make(map[*websocket.Conn]bool)


// websocker upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}


/*
 * Utility function to throw a panic if an error occurs
 */
func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}
}

// handlers
func HandleIndex(w http.ResponseWriter, r *http.Request) {
	templates := SenseHatIoT.HtmlTemplates{}
	tmpl := template.New("index")
	tmpl, errtmpl := tmpl.Parse(string(templates.IndexPageTemplate()))
	CheckError(errtmpl)
	
	var rendered = bytes.Buffer{}
	tdict := TemplateDict{
		"measurements" : measurements,
		"serverinfo" : server_info,
	}
	err := tmpl.Execute(&rendered, &tdict) 
	CheckError(err)
	
	w.Header().Set("Content-Type", "text/html")
	w.Write(rendered.Bytes())
}

func HandleExitApplication(w http.ResponseWriter, r *http.Request) {
	close(exit_worker_channel)
	exit_bool_channel <- true
	signal.Notify(exit_signal_channel, os.Interrupt, syscall.SIGTERM)
}

func HandleWsConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return
	}
	
	defer ws.Close()
	
    clients[ws] = true;

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			fmt.Println(err)
			delete(clients, ws)
			break
		}
		
		fmt.Fprintf(os.Stdout, "Received: Message Type: %s, Message: %s \n", msg.Type, msg.Text)
	}
}


func StartHttpServer(info *ServerInfo) *http.Server {
	srv := &http.Server{
		Addr: ":" + info.ServerPort,
	}
	
	// Routes
	http.HandleFunc("/", HandleIndex)
	http.HandleFunc("/exitapplication", HandleExitApplication)
	
	http.HandleFunc("/ws", HandleWsConnections)
	
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Fprintf(os.Stderr, "HTTP server error: %s", err.Error())
		}
	}()
	
	return srv
}


func Worker(exitworkerchan chan bool) {
	for {
		select {
		case <-exitworkerchan:
			break
		
		default:
			// take a measurement
			timeNow := time.Now()
			
			// if the date has changed, clear the database
			if len(measurements) > 0 {
                firstMeasurement := measurements[0]
			    if firstMeasurement.TimeObj.Day() != timeNow.Day() {
			  	    measurements = nil
			    }
            }
			
			measurement_now := SenseHatIoT.Measurement{
				TimeObj: timeNow,
				Timestamp: fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
					timeNow.Year(), timeNow.Month(), timeNow.Day(),
					timeNow.Hour(), timeNow.Minute(), timeNow.Second()),
				Temperature: SenseHatIoT.GetSenseHatTemperature(),
			}
			
			measurements_channel<-measurement_now
			
			//time.Sleep(5 * time.Second)
			time.Sleep(15 * time.Minute)
		}
	}
}

func MeasurementHandler(exitworkerchan chan bool) {
	for {
		select {
		case <-exitworkerchan:
			break
		
		case measurement := <-measurements_channel:
			measurements = append(measurements, measurement)
			// also do WS stuff here
			msg, _ := json.Marshal(measurement)
			fmt.Println(string(msg))
            for client := range(clients) {
				err := client.WriteJSON(string(msg))
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
			
		default: {}
		
		}
	}
}


func main() {
	portNumber := flag.String("port", "8080", "Port to serve on")
	
	flag.Parse()
	
	fmt.Println("Parsed these command-line flags:")
	fmt.Println("Port number:", *portNumber)
	
	// get the local IP
	local_port := *portNumber
	local_ip := SenseHatIoT.GetLocalIP()
	if local_ip == "" {
		local_ip = "0.0.0.0"
	}
	
	server_info.ServerIP = local_ip
	server_info.ServerPort = local_port
	
	// start a HTTP and WS server
	srv := StartHttpServer(&server_info)
	
	
	// start the temperature workers
	go Worker(exit_worker_channel)
	go MeasurementHandler(exit_worker_channel)
	
	
	// exit gracefully
	select {
	case <-exit_bool_channel:
	case <-exit_signal_channel:
		fmt.Println("Shutting down the server...")
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		srv.Shutdown(ctx)
	}
}

