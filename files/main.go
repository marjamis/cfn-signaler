package main

import (
  "fmt"
  "io/ioutil"
  "net/http"
  "html/template"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/cloudformation"
  "github.com/aws/aws-sdk-go/aws/ec2metadata"
  "encoding/json"
)

type Response struct {
  Signal string
}

type JsonData struct {
  LogicalResourceId string `json:"LogicalResourceId"`
  StackName string `json:"StackName"`
}

func signal(send string) {
  metadata := ec2metadata.New(nil)  
  region, err := metadata.Region()  
  if err !=nil {   
    fmt.Println("Error: ", err)
  }

  instance_id, err := metadata.GetMetadata("instance-id")
  if err !=nil {   
    fmt.Println("Error: ", err)
  }

  file, err := ioutil.ReadFile("./config/cfn-signaler.json")
  if err != nil {
    fmt.Println("Error: ", err)
  }
  
  var conf JsonData
  err=json.Unmarshal(file, &conf)
  if err!=nil{
    fmt.Println("Error: ", err)
  }
  
  cfn_config := aws.NewConfig().WithRegion(region)  
  svc := cloudformation.New(nil, cfn_config)

  params := &cloudformation.SignalResourceInput { 
    LogicalResourceId: aws.String(conf.LogicalResourceId),
    StackName: aws.String(conf.StackName),
    Status: aws.String(send),
    UniqueId: aws.String(instance_id),
  } 

  resp, err := svc.SignalResource(params)
  if err != nil { 
    // Print the error, cast err to awserr.Error to get the Code and 
    // Message from an error.
    fmt.Println("Error: ", err)
    return  
  }
 
  fmt.Println(resp)
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
  fmt.Println("handler IP: ", r.RemoteAddr)
}

func signalHandler(w http.ResponseWriter, r *http.Request) {  
  value := r.FormValue("send")
  
  var text string
  if value == "SUCCESS" {
    signal("SUCCESS")
    text = "Success signal sent"
  } else if  value == "FAILURE" {
    signal("FAILURE")
    text = "Failed signal sent"
  } else {
    text = "Invalid signal type"
  }  
  
  t, _ := template.ParseFiles("templates/signal.html")
  data := &Response{Signal: text}
  t.Execute(w, data)
  fmt.Println("signalHandler IP: ", r.RemoteAddr, "signal: ", value)
}

func main() {
    http.HandleFunc("/", handler)
    http.HandleFunc("/signal/", signalHandler)
    fmt.Println("Listening on port 8080...")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
      fmt.Println("Error: ", err)
    }
}
