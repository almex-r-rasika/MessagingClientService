package data

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type LoginRequestBody struct {
    ID string `json:"id"`
	Password string `json:"pw"`
}

type Login struct {
	Error string      `json:"error"`
	ErrorCode string      `json:"errorCode"`
	ErrorString string      `json:"errorString"`
	Token string `json:"token"`
	StaffName string `json:"staff_name"`
}

var logged_users = make(map[string]string)

/* function for each user to login
   @param --> id, password, login_time
   @param value --> loginuser id, loginuser password, login time
   description --> function for each user to login through a session
   @return --> Login object
*/
func DoLogin(id string, password string) Login{

	requestBody := &LoginRequestBody{
        ID: id,
		Password: password,
    }

    jsonString, err := json.Marshal(requestBody)
    if err != nil {
        Log.Fatal(err.Error())
    }

	Log.Debug("login request body "+string(jsonString))

	endpoint := "https://msgc.smapa-checkout.jp/v1/ctrl/login"
    req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonString))
    if err != nil {
        Log.Fatal(err.Error())
    }

    req.Header.Set("Content-Type", "application/json")
    
	// add CA Certificate
	client := addCaCertificate()

    response, err := client.Do(req)
    if err != nil {
        Log.Fatal(err.Error())
    }

	defer response.Body.Close()
	
   // Read the response body on the line below
   body, err := ioutil.ReadAll(response.Body)
   if err != nil {
        Log.Fatal(err.Error())
   }

   Log.Debug("login response body "+string(body))

   // Convert the body to type login object
   var login Login
   json.Unmarshal([]byte(body), &login)
   // add logged user token to logged_users map
   logged_users[id] = login.Token

   // validate the response
   if(login.Error == "0"){
	// add logged user token to logged_users map
    logged_users[id] = login.Token
	Log.Info("user "+id+" logged")
   }else{
    Log.Error("error "+login.Error+" errorcode "+login.ErrorCode+" errorstring "+login.ErrorString)
   }
   return login
}

/* function for get already logged users map
   @param --> null
   @param value --> null
   description --> get already logged users map
   @return --> map[user_id]token
*/
func GetLoggedUsers() map[string]string {
	return logged_users
}



