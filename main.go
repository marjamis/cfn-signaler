package main

import (
	"errors"
	"html/template"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"

	log "github.com/sirupsen/logrus"
)

type response struct {
	Signal string
	Error  error
}

func signal(send string) (err error) {
	session := session.New()
	metadata := ec2metadata.New(session)

	if !metadata.Available() {
		log.Error("EC2 Metadata not available")
		return errors.New("EC2 Metadata not available")
	}

	region, err := metadata.Region()
	if err != nil {
		log.Error(err)
		return err
	}

	instanceId, err := metadata.GetMetadata("instance-id")
	if err != nil {
		log.Error(err)
		return err
	}

	logicalId := os.Getenv("LOGICALID")
	stackName := os.Getenv("STACKNAME")

	cfnConfig := aws.NewConfig().WithRegion(region)
	svc := cloudformation.New(session, cfnConfig)

	params := &cloudformation.SignalResourceInput{
		LogicalResourceId: aws.String(logicalId),
		StackName:         aws.String(stackName),
		Status:            aws.String(send),
		UniqueId:          aws.String(instanceId),
	}

	_, err = svc.SignalResource(params)
	if err != nil {
		log.Error(err)
		return err
	}

	log.WithFields(log.Fields{
		"Function":  "signal",
		"LogicalId": logicalId,
		"StackName": stackName,
	}).Info()

	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/"):]
	var filename string

	if len(title) == 0 {
		filename = "templates/index.html"
	} else {
		if _, err := os.Stat(title); os.IsNotExist(err) {
			filename = "templates/http_404.html"
		} else {
			filename = title
		}
	}

	t, _ := template.ParseFiles(filename)
	t.Execute(w, nil)

	log.WithFields(log.Fields{
		"Function": "signal",
		"IP":       r.RemoteAddr,
		"File":     filename,
	}).Info()
}

func signalHandler(w http.ResponseWriter, r *http.Request) {
	value := r.FormValue("send")
	var text string
	var err error

	if value == "SUCCESS" {
		err = signal("SUCCESS")
		text = "Success signal"
	} else if value == "FAILURE" {
		err = signal("FAILURE")
		text = "Failed signal"
	} else {
		value = "INVALID"
		err = errors.New("Invalid signal type specified in POST")
		text = "Invalid signal type"
	}

	t, _ := template.ParseFiles("templates/signal.html")
	data := &response{Signal: text, Error: err}
	t.Execute(w, data)

	log.WithFields(log.Fields{
		"Function": "signalHandler",
		"IP":       r.RemoteAddr,
		"Signal":   value,
		"Error":    err,
	}).Info()
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/signal", signalHandler)
	http.HandleFunc("/signal/", signalHandler)

	log.Info("Listening on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Error(err)
		return
	}
}
