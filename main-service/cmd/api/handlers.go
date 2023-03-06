package main

import (
	"log"
	"net/http"
	"service/data"
	"strings"
)

type SaveMessageTemplateRequest struct {
	Name string `json:"name"`
	Subject string `json:"subject"`
	MessagesJson string `json:"messagesJson"`
	URL string `json:"url"`
}

type UpdateMessageTemplateRequest struct {
	Name string `json:"name"`
	Subject string `json:"subject"`
	MessagesJson string `json:"messagesJson"`
	URL string `json:"url"`
}

type SaveBulkMessageRequest struct {
	BulkType string `json:"bulkType"`
	ImportId string `json:"importId"`
	FilterJson string `json:"filterJson"`
	Name string `json:"name"`
	Subject string `json:"subject"`
	MessagesJson string `json:"messagesJson"`
	URL string `json:"url"`
}

type GetMessageTemplatesResponse struct {
	Error string `json:"error"`
	ErrorCode string `json:"errorCode"`
	ErrorString string `json:"errorString"`
	Count string `json:"count"`
	List []data.APIMessageTemplate `json:"list"`
}

type GetMessageTemplateResponse struct {
	Error string `json:"error"`
	ErrorCode string `json:"errorCode"`
	ErrorString string `json:"errorString"`
	Data data.APIMessageTemplate2 `json:"data"`
}

type SaveMessageTemplateResponse struct {
	Error string `json:"error"`
	ErrorCode string `json:"errorCode"`
	ErrorString string `json:"errorString"`
	Data data.APITemplateId `json:"data"`
}

type SaveBulkMessageResponse struct {
	Error string `json:"error"`
	ErrorCode string `json:"errorCode"`
	ErrorString string `json:"errorString"`
	Data data.APIBulkMessageId `json:"data"`
}

type GetBulkMessagesResponse struct {
	Error string `json:"error"`
	ErrorCode string `json:"errorCode"`
	ErrorString string `json:"errorString"`
	Count string `json:"count"`
	List []data.APIBulkMessage `json:"list"`
}

type GetBulkMessageResponse struct {
	Error string `json:"error"`
	ErrorCode string `json:"errorCode"`
	ErrorString string `json:"errorString"`
	Count string `json:"count"`
	List []data.APIBulkMessage2 `json:"list"`
}

type GetImportDataHistoryResponse struct {
	Error string `json:"error"`
	ErrorCode string `json:"errorCode"`
	ErrorString string `json:"errorString"`
	Count string `json:"count"`
	List []data.APIImportHistory `json:"list"`
}

type DeleteImportDataHistoryResponse struct {
	Error string `json:"error"`
	ErrorCode string `json:"errorCode"`
	ErrorString string `json:"errorString"`
}

type UpdateMessageTemplateResponse struct {
	Error string `json:"error"`
	ErrorCode string `json:"errorCode"`
	ErrorString string `json:"errorString"`
}

type DeleteMessageTemplateResponse struct {
	Error string `json:"error"`
	ErrorCode string `json:"errorCode"`
	ErrorString string `json:"errorString"`
}

/* PUT API to save message template
   @param --> http.ResponseWriter, *http.Request
   @param value --> w http.ResponseWriter, r *http.Request
   description --> save message template
   @return --> null
*/
func (app *Config) saveMessageTemplate(w http.ResponseWriter, r *http.Request) {

	var requestPayload SaveMessageTemplateRequest

	templateType := strings.TrimPrefix(r.URL.Path, "/v2/msgb/message_template/")

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	data.SaveMessageTemplate(templateType, requestPayload.Name, requestPayload.Subject, requestPayload.MessagesJson, requestPayload.URL)
	data := data.GetLastTemplate()

	payload := SaveMessageTemplateResponse{
		Error:   "0",
		ErrorCode: "",
		ErrorString: "",
		Data: data,
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

/* POST API to update message template
   @param --> http.ResponseWriter, *http.Request
   @param value --> w http.ResponseWriter, r *http.Request
   description --> update message template
   @return --> null
*/
func (app *Config) updateMessageTemplate(w http.ResponseWriter, r *http.Request) {

	var requestPayload UpdateMessageTemplateRequest

	template := strings.TrimPrefix(r.URL.Path, "/v2/msgb/message_template/")
	arr := strings.Split(template, "/")
	templateType := arr[0]
	templateId := arr[1]

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	data.UpdateMessageTemplate(templateId, templateType, requestPayload.Name, requestPayload.Subject, requestPayload.MessagesJson, requestPayload.URL)

	payload := UpdateMessageTemplateResponse{
		Error:   "0",
		ErrorCode: "",
		ErrorString: "",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

/* DEL API to update message template
   @param --> http.ResponseWriter, *http.Request
   @param value --> w http.ResponseWriter, r *http.Request
   description --> update message template
   @return --> null
*/
func (app *Config) deleteMessageTemplate(w http.ResponseWriter, r *http.Request) {

	template := strings.TrimPrefix(r.URL.Path, "/v2/msgb/message_template/")
	arr := strings.Split(template, "/")
	templateType := arr[0]
	templateId := arr[1]

	data.DeleteMessageTemplate(templateId, templateType)

	payload := DeleteMessageTemplateResponse{
		Error:   "0",
		ErrorCode: "",
		ErrorString: "",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

/* GET API to save import data
   @param --> http.ResponseWriter, *http.Request
   @param value --> w http.ResponseWriter, r *http.Request
   description --> save import data
   @return --> null
*/
func (app *Config) saveImportData(w http.ResponseWriter, r *http.Request) {

   data.SaveImportHistory()
   data.SaveImportData()

	payload := jsonResponse{
		Error:   false,
		ErrorCode: "",
		ErrorString: "",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

/* GET API to get import data history
   @param --> http.ResponseWriter, *http.Request
   @param value --> w http.ResponseWriter, r *http.Request
   description --> get import data history
   @return --> null
*/
func (app *Config) getImportDataHistory(w http.ResponseWriter, r *http.Request) {

	page := strings.TrimPrefix(r.URL.Path, "/v2/msgb/import_data/")

	list, count := data.GetImportHistory(page)

	payload := GetImportDataHistoryResponse{
		Error:   "0",
		ErrorCode: "",
		ErrorString: "",
		Count: count,
		List: list,
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

/* DELETE API to get import data history
   @param --> http.ResponseWriter, *http.Request
   @param value --> w http.ResponseWriter, r *http.Request
   description --> get import data history
   @return --> null
*/
func (app *Config) deleteImportDataHistory(w http.ResponseWriter, r *http.Request) {

	importId := strings.TrimPrefix(r.URL.Path, "/v2/msgb/import_data/")

	data.DeleteImportHistory(importId)

	payload := DeleteImportDataHistoryResponse{
		Error:   "0",
		ErrorCode: "",
		ErrorString: "",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

/* GET API to get message template
   @param --> http.ResponseWriter, *http.Request
   @param value --> w http.ResponseWriter, r *http.Request
   description --> get message template
   @return --> null
*/
func (app *Config) getMessageTemplates(w http.ResponseWriter, r *http.Request) {

	list, count := data.GetMessageTemplates()

	payload := GetMessageTemplatesResponse{
		Error:   "0",
		ErrorCode: "",
		ErrorString: "",
		Count: count,
		List: list,
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

/* GET API to get message template
   @param --> http.ResponseWriter, *http.Request
   @param value --> w http.ResponseWriter, r *http.Request
   description --> get message template
   @return --> null
*/
func (app *Config) getMessageTemplate(w http.ResponseWriter, r *http.Request) {

	templateId := strings.TrimPrefix(r.URL.Path, "/v2/msgb/message_template/")
	Id := strings.Split(templateId, "/")
	templateId = Id[1]

	list:= data.GetMessageTemplate(templateId)

	payload := GetMessageTemplateResponse{
		Error:   "0",
		ErrorCode: "",
		ErrorString: "",
		Data: list,
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

/* GET API to save bulk messages
   @param --> http.ResponseWriter, *http.Request
   @param value --> w http.ResponseWriter, r *http.Request
   description --> save bulk messages
   @return --> null
*/
func (app *Config) saveBulkMessage(w http.ResponseWriter, r *http.Request) {

	var requestPayload SaveBulkMessageRequest

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

   data.SaveBulkMessages(requestPayload.BulkType, requestPayload.ImportId, requestPayload.FilterJson, requestPayload.Subject, requestPayload.MessagesJson, requestPayload.URL)
   data := data.GetLastBulkMessageId()

	payload := SaveBulkMessageResponse{
		Error:   "0",
		ErrorCode: "",
		ErrorString: "",
		Data: data,
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

/* GET API to get bulk messages
   @param --> http.ResponseWriter, *http.Request
   @param value --> w http.ResponseWriter, r *http.Request
   description --> get bulk messages
   @return --> null
*/
func (app *Config) getBulkMessages(w http.ResponseWriter, r *http.Request) {

	page := strings.TrimPrefix(r.URL.Path, "/v2/msgb/bulk_message/")

	list, count := data.GetBulkMessages(page)

	payload := GetBulkMessagesResponse{
		Error:   "0",
		ErrorCode: "",
		ErrorString: "",
		Count: count,
		List: list,
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

/* GET API to get bulk message
   @param --> http.ResponseWriter, *http.Request
   @param value --> w http.ResponseWriter, r *http.Request
   description --> get message template
   @return --> null
*/
func (app *Config) getBulkMessage(w http.ResponseWriter, r *http.Request) {

	bulkMessageId := strings.TrimPrefix(r.URL.Path, "/v2/msgb/bulk_message/")
	Id := strings.Split(bulkMessageId, "/")
	bulkMessageId = Id[0]

	log.Println(bulkMessageId)

	list, count:= data.GetBulkMessage(bulkMessageId)

	payload := GetBulkMessageResponse{
		Error:   "0",
		ErrorCode: "",
		ErrorString: "",
		Count: count,
		List: list,
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}
