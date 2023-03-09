package main

import (
	"service/data"
)

/* function for initialize and save login user/message objects to corresponding tables
   @param --> null
   @param value --> null
   description --> save login users and messages to db
   @return --> null
*/
func initializeService() {
	data.MakeLogger()
   data.ConnectToDatabase()
   data.DoLogin("install_11","install_11")
}

