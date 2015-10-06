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
    return     
  }

  instance_id, err := metadata.GetMetadata("instance-id")
  if err !=nil {   
    return     
  }

  file, e := ioutil.ReadFile("./config/cfn_signaler.json")
  if e != nil {
    return
  }
  
  var conf JsonData
  err=json.Unmarshal(file, &conf)
  if err!=nil{
    fmt.Print("Error:",err)
  }
  

  fmt.Println(conf.StackName)
  fmt.Println(conf.LogicalResourceId)
  

  cfn_config := aws.NewConfig().WithRegion(region)  
  svc := cloudformation.New(cfn_config)

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
    fmt.Println(err.Error()) 
    return  
  }
 
  fmt.Println(resp)
}

func logging() {
  fmt.Println("Log")
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
  logging()
}

func signalHandler(w http.ResponseWriter, r *http.Request) {  
 fmt.Println(r.FormValue("send"))
  value := r.FormValue("send")
  
  var text string
  fmt.Println(text)
  if value == "SUCCESS" {
    signal("SUCCESS")
    text = "Success signal sent"
  } else if  value == "FAILURE" {
    signal("FAILURE")
    text = "Failed signal sent"
  } else {
    text = "Invalid signal type"
  }  
  
  fmt.Println(text)
  t, _ := template.ParseFiles("templates/signal.html")
  data := &Response{Signal: text}
  t.Execute(w, data)
  logging()
}

func main() {
    http.HandleFunc("/", handler)
    http.HandleFunc("/signal/", signalHandler)
    http.ListenAndServe(":8080", nil)
}
