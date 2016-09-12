package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

type Response struct {
	Signal string
}

func signal(send string) {
	session := session.New()
	metadata := ec2metadata.New(session)

	if !metadata.Available() {
		log.Info("Error: Metadata not available.")
		return
	}

	region, err := metadata.Region()
	if err != nil {
		log.Error(err)
		return
	}

	instance_id, err := metadata.GetMetadata("instance-id")
	if err != nil {
		log.Error(err)
		return
	}

	lid := os.Getenv("LOGICALID")
	sn := os.Getenv("STACKNAME")

	cfn_config := aws.NewConfig().WithRegion(region)
	svc := cloudformation.New(session, cfn_config)

	params := &cloudformation.SignalResourceInput{
		LogicalResourceId: aws.String(lid),
		StackName:         aws.String(sn),
		Status:            aws.String(send),
		UniqueId:          aws.String(instance_id),
	}

	resp, err := svc.SignalResource(params)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info(resp)
}

func handler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/"):]
	var filename string

	if len(title) == 0 {
		filename = "templates/index.html"
	} else {
		filename = title
	}

	body, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(w, "%s", "No such page")
		return
	}
	fmt.Fprintf(w, "%s", body)
	log.Info("handler IP: ", r.RemoteAddr)
}

func signalHandler(w http.ResponseWriter, r *http.Request) {
	value := r.FormValue("send")

	var text string
	if value == "SUCCESS" {
		signal("SUCCESS")
		text = "Success signal sent"
	} else if value == "FAILURE" {
		signal("FAILURE")
		text = "Failed signal sent"
	} else {
		text = "Invalid signal type"
	}

	t, _ := template.ParseFiles("templates/signal.html")
	data := &Response{Signal: text}
	t.Execute(w, data)
	log.Info("signalHandler IP: ", r.RemoteAddr, " signal: ", value)
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/signal/", signalHandler)
	log.Info("Listening on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Error(err)
	}
}
