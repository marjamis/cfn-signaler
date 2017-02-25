package main

import (
	"os"
	"errors"
	"html/template"
	"net/http"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	log "github.com/Sirupsen/logrus"
)


type Response struct {
	Signal string
	Error error
}

func signal(send string) (err error){
	session := session.New()
	metadata := ec2metadata.New(session)

	if !metadata.Available() {
		log.Error("Error: Metadata not available.")
		return errors.New("Error: Metadata not available.")
	}

	region, err := metadata.Region()
	if err != nil {
		log.Error(err)
		return err
	}

	instance_id, err := metadata.GetMetadata("instance-id")
	if err != nil {
		log.Error(err)
		return err
	}

	logical_id := os.Getenv("LOGICALID")
	stack_name := os.Getenv("STACKNAME")

	cfn_config := aws.NewConfig().WithRegion(region)
	svc := cloudformation.New(session, cfn_config)

	params := &cloudformation.SignalResourceInput{
		LogicalResourceId: aws.String(logical_id),
		StackName:         aws.String(stack_name),
		Status:            aws.String(send),
		UniqueId:          aws.String(instance_id),
	}

	_, err = svc.SignalResource(params)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Info("Function: handler - LogicalID: ", logical_id, " - StackName: ", stack_name)
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/"):]
	var filename string

	if len(title) == 0 {
		filename = "templates/index.html"
	} else {
		if _, err := os.Stat(title); os.IsNotExist(err){
			filename = "templates/http_404.html"
		} else {
			filename = title
		}
	}

	t, _ := template.ParseFiles(filename)
        t.Execute(w, nil)

	log.Info("Function: handler - IP: ", r.RemoteAddr, " - File: ", filename)
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
	data := &Response{Signal: text, Error: err}
	t.Execute(w, data)

	log.Info("Function: signalHandler - IP: ", r.RemoteAddr, " - signal: ", value)
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
