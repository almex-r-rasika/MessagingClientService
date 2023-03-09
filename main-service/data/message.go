package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type MessageRequestBody struct {
    PatientNumList []string `json:"patient_num_list"`
	Subject string `json:"subject"`
	Body string `json:"body"`
	Url string `json:"url"`
}

type Message struct {
	Error string      `json:"error"`
	ErrorCode string      `json:"errorCode"`
	ErrorString string      `json:"errorString"`
}

/* function for send message to the given destination 
   @param --> messageId, sendTime, sendUserId, address, subject, line1, line2, line3, line4, line5, line6, line7, line8, line9, line10, patientList[], token, count
   @param value --> messageId, sendTime, sendUserId, address, subject, line1, line2, line3, line4, line5, line6, line7, line8, line9, line10, patientList[], token, count
   description --> for each message when the send time comes message will send to the destination
   @return --> Message Object
*/
func DoMessage(bulkMessageId int,subject string,messageJson string,url string,token string) Message{

	var patients = GetPatientList(token)

	requestBody := &MessageRequestBody{
        PatientNumList: patients,
		Subject: subject,
		Body: messageJson,
		Url: url,
    }

    jsonString, err := json.Marshal(requestBody)
    if err != nil {
        Log.Fatal(err.Error())
    }

	Log.Debug("message request body "+string(jsonString))

	endpoint := "https://msgc.smapa-checkout.jp/v1/hospital/messages/send"
    req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonString))
    if err != nil {
        Log.Fatal(err.Error())
    }

    req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	// add CA Certificate
	client := addCaCertificate()

    response, err := client.Do(req)
    if err != nil {
        Log.Fatal(err.Error())
    }

	defer response.Body.Close()
	
   // Read the response body on the line below.
   body, err := ioutil.ReadAll(response.Body)
   if err != nil {
      Log.Fatal(err.Error())
   }

   Log.Debug("message response body "+string(body))
   
   // Convert the body to type message object
   var message Message
   json.Unmarshal([]byte(body), &message)

   // validate the response
   if(message.Error == "0"){
	Log.Info("user sent message")
   }else{
    Log.Error("error "+message.Error+" errorcode "+message.ErrorCode+" errorstring "+message.ErrorString)
   }

   // save send messages to the database
   saveBulkMessagesBox(bulkMessageId, fmt.Sprint(patients), subject, messageJson, url, message.Error )

   return message
}

