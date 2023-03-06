package data

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type PatientData struct {
	Error       string   `json:"error"`
	ErrorCode   string   `json:"errorCode"`
	ErrorString string   `json:"errorString"`
	PersonList  []PersonList `json:"personList"`
}

type PersonList struct {
	PatientNum   string `json:"patient_num"`
	LocalId      string `json:"localid"`
	Name         string `json:"name"`
	KanaName     string `json:"kananame"`
	Relationship string `json:"relationship"`
	Birthday     string `json:"birthday"`
	Sex          string `json:"sex"`
	CreatedAt    string `json:"created_at"`
	ReceiveOk    bool   `json:"receive_ok"`
}

type PatientDataRequestBody struct {
	FreeWord string `json:"free-word"`
}

/* function for get patient list for given free word
   @param --> token
   @param value --> token
   description --> get patient list for given free word
   @return --> patient_list array
*/
func GetPatientList(address string,token string) []string{

	var person_list string

	address_list := getFreewordAndRelationship(address)

	requestBody := &PatientDataRequestBody{
		FreeWord:       address_list[0],
	}

	jsonString, err := json.Marshal(requestBody)
	if err != nil {
		Log.Fatal(err.Error())
	}

	Log.Debug("patient request body "+string(jsonString))

	endpoint := "https://msgc.smapa-checkout.jp/v1/hospital/users/search/freeword"

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

	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		Log.Fatal(err.Error())
	}

	Log.Debug("patient response body "+string(body))

	//Convert the body to type patient object
	var patient PatientData
	json.Unmarshal([]byte(body), &patient)

	person := patient.PersonList

	for _, p := range person{
        patient_num := p.PatientNum
		if(p.Relationship == address_list[1]){
			person_list = patient_num+" "+person_list
		}
	}

	patient_list := strings.Fields(person_list)

	// validate the response
   if(patient.Error == "0"){
	Log.Info("get patient list from the api")
   }else{
    Log.Error("error "+patient.Error+" errorcode "+patient.ErrorCode+" errorstring "+patient.ErrorString)
   }
	return patient_list
}

/* function for get free-word and relationship from the message address
   @param --> address
   @param value --> message object's address field
   description --> split address string to string array using , operator
   @return --> string array 
*/
func getFreewordAndRelationship(address string) []string{

	result := strings.Split(address, ",")
	return result
}