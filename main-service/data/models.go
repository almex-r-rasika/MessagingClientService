package data

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

var db *gorm.DB
var isConnect = false

type ImportData struct {
  Id   int `gorm:"primaryKey"`
  ImportId   int
  LineNum int
  HeaderId int
  Value string
}

type ImportHistory struct {
  ImportId int `gorm:"primaryKey"`
  StaffNum   int
  Title   string
  Header string
  CreatedAt time.Time
  DeletedAt sql.NullTime
}

type MessageTemplate struct {
  Id int `gorm:"primaryKey"`
  StaffNum   int
  TemplateType string `gorm:"default:S"`
  TemplateName string
  Subject string
  MessagesJson string
  URL string
  CreatedAt time.Time
  DeletedAt sql.NullTime
}

type APIMessageTemplate struct {
  Id int `json:"templateId"`
  StaffNum   int `json:"staffNum"`
  TemplateType string `json:"templateType"`
  TemplateName string `json:"name"`
  Subject string `json:"subject"`
  CreatedAt time.Time `json:"createdAt"`
}

type APIMessageTemplate2 struct {
  Id int `json:"templateId"`
  StaffNum   int `json:"staffNum"`
  TemplateType string `json:"templateType"`
  TemplateName string `json:"name"`
  Subject string `json:"subject"`
  MessagesJson string `json:"messagesJson"`
  URL string `json:"url"`
  CreatedAt time.Time `json:"createdAt"`
}

type APITemplateId struct {
  Id int `json:"templateId"`
}

type BulkMessage struct {
  BulkMessageId int `gorm:"primaryKey"`
  StaffNum   int
  BulkType   string `gorm:"default:0"`
  ImportId string
  FilterJson   string
  Subject   string
  MessagesJson string
  URL string
  CreatedAt time.Time
  DeletedAt sql.NullTime
}

type APIBulkMessage struct {
  BulkMessageId int `json:"bulkMessageId"`
  StaffNum   int `json:"staffNum"`
  BulkType   string `json:"bulkType"`
  ImportId string `json:"importId"`
  FilterJson   string `json:"filterJson"`
  Subject   string `json:"subject"`
  MessagesJson string `json:"messagesJson"`
  CreatedAt time.Time `json:"createdAt"`
}

type APIBulkMessage2 struct {
  BulkMessageId int `json:"bulkMessageId"`
  StaffNum   int `json:"staffNum"`
  BulkType   string `json:"bulkType"`
  ImportId string `json:"importId"`
  FilterJson   string `json:"filterJson"`
  Subject   string `json:"subject"`
  MessagesJson string `json:"messagesJson"`
  URL string `json:"url"`
  CreatedAt time.Time `json:"createdAt"`
}

type BulkMessagesBox struct {
  BulkMessagesBoxId int `gorm:"primaryKey"`
  BulkMessageId int
  JID   string
  Subject   string 
  Body string
  URL  string
  CreatedAt time.Time
  SendAt sql.NullTime
  DeletedAt sql.NullTime
  Status   string `gorm:"default:0"`
}

type APIBulkMessageId struct {
  BulkMessageId int `json:"bulkMessageId"`
}

type APIImportHistory struct {
  ImportId int `json:"importId"`
  StaffNum   int `json:"staffNum"`
  Title   string `json:"title"`
  CreatedAt time.Time `json:"createdAt"`
}

var importData []ImportData
var importHistory []APIImportHistory
var messageTemplate []APIMessageTemplate
var messageTemplate2 APIMessageTemplate2
var templateId APITemplateId
var bulkMessageId APIBulkMessageId
var bulkMessages []APIBulkMessage
var bulkMessages2 []APIBulkMessage2
var importHistoryObj ImportHistory

/* function for make a connection to the database
    @param --> null
    @param value --> null
    description --> if connection has established before, ignore the making new connection
	@return --> null
*/
func ConnectToDatabase(){
	if !isConnect{
		time.Sleep(time.Duration(10000 * time.Millisecond))
		db = makeDbConnection()
	}
	isConnect = true
}

/* function for pagination
    @param --> null
    @param value --> null
    description --> set offset and limit for sql query 
	@return --> null
*/
func Paginate(page string) func(db *gorm.DB) *gorm.DB {
  return func (db *gorm.DB) *gorm.DB {
    pageId, _ := strconv.Atoi(page)
    pageSize := 25
    offset := (pageId - 1) * pageSize
    return db.Offset(offset).Limit(pageSize)
  }
}

/* Save import history to the database
    @param --> null
    @param value --> null
    description --> save import history to db
    @return --> null
*/
func SaveImportHistory() {
  ConnectToDatabase()
  absPath, _ := filepath.Abs("../main-service/data/sample.sqlite")
  cols, _ := executeSqliteFile(absPath)
  header := fmt.Sprint(cols)
	history := []ImportHistory{{StaffNum: 333, Title: "sample1", Header: header, CreatedAt: time.Now().Local()}}
  db.Create(&history)
  Log.Info("saved import history")
}

/* Save import data
    @param --> null
    @param value --> null
    description --> save import data to db
    @return --> null
*/
func SaveImportData() {
  ConnectToDatabase()
  absPath, _ := filepath.Abs("../main-service/data/sample.sqlite")
  cols, rows := executeSqliteFile(absPath)
  colLen := len(cols)
  vals := make([]interface{}, colLen)
  importId := GetLastImportId()
  var lineNum int;

  for rows.Next() {
        for i := 0; i < colLen; i++ {
            vals[i] = new(string)
        }
        err := rows.Scan(vals...)
        if err != nil {
            log.Fatal(err.Error())
        }

        for i := 0; i < colLen; i++ {
          data := []ImportData{{ImportId: importId, LineNum: lineNum, HeaderId: i, Value: *vals[i].(*string)}}
          db.Create(&data)
        }
        lineNum++
    }

  Log.Info("saved import data")
}

/* Save message template to the database
    @param --> null
    @param value --> null
    description --> save template to database
    @return --> null
*/
func SaveMessageTemplate(tmpType string, tmpName string, tmpSubject string, messagesJson string, url string) {
  ConnectToDatabase()
	template := []MessageTemplate{{StaffNum: 333, TemplateType: tmpType, TemplateName: tmpName, Subject: tmpSubject, MessagesJson: messagesJson, URL: url, CreatedAt: time.Now().Local()}}
  db.Create(&template)
  Log.Info("created message template")
}

/* Update message template to the database
    @param --> null
    @param value --> null
    description --> update template
    @return --> null
*/
func UpdateMessageTemplate(id string,tmpType string, tmpName string, tmpSubject string, messagesJson string, url string) {
  ConnectToDatabase()
  db.Model(&MessageTemplate{}).Where("id = ?", id).Updates(MessageTemplate{StaffNum: 333, TemplateType: tmpType, TemplateName: tmpName, Subject: tmpSubject, MessagesJson: messagesJson, URL: url, CreatedAt: time.Now().Local()})
  Log.Info("updated message template")
}

/* Delete message template from the database
    @param --> null
    @param value --> null
    description --> delete template
    @return --> null
*/
func DeleteMessageTemplate(id string,tmpType string) {
  ConnectToDatabase()
  db.Model(&MessageTemplate{}).Where("id = ? AND template_type = ?", id, tmpType).Update("deleted_at", time.Now().Local())
  Log.Info("deleted message template")
}

/* function for get import history from import history table
    @param --> null
    @param value --> null
    description --> get import history from import history table
    @return --> import history array
*/
func GetImportHistory(page string) ([]APIImportHistory, string){
  ConnectToDatabase()
  result := db.Model(&ImportHistory{}).Where("deleted_at IS NULL").Find(&importHistory)
  db.Model(&ImportHistory{}).Where("deleted_at IS NULL").Scopes(Paginate(page)).Find(&importHistory)
  count := strconv.FormatInt(result.RowsAffected,10)
	Log.Info("get import history from the database")
	return importHistory, count
}

/* function for delete import history from import history table
    @param --> null
    @param value --> null
    description --> delete import history from import history table
    @return --> null
*/
func DeleteImportHistory(importId string){
  ConnectToDatabase()
  db.Model(&ImportHistory{}).Where("import_id = ?",importId).Update("deleted_at", time.Now().Local())
	Log.Info("deleted import history from the database")
}

/* function for get import data from import data table
    @param --> null
    @param value --> null
    description --> get import data from import data table
    @return --> import data array
*/
func GetImportData() (([]ImportData)){
  ConnectToDatabase()
	db.Order("id").Find(&importData)
	Log.Info("get import data from the database")
	return importData
}

/* function for get last import id from import history table
    @param --> null
    @param value --> null
    description --> get last import id from import history table
    @return --> last import id
*/
func GetLastImportId() (int){
  ConnectToDatabase()
	_, c := GetImportHistory("")
  count, _ := strconv.Atoi(c)
  //lastImportId = len(history)
	return count
}

/* function for get message templates from message template table
    @param --> null
    @param value --> null
    description --> get message templates from message template table
    @return --> message template array
*/
func GetMessageTemplates() ([]APIMessageTemplate, string){
  ConnectToDatabase()
  result := db.Model(&MessageTemplate{}).Where("deleted_at IS NULL").Find(&messageTemplate)
  count := strconv.FormatInt(result.RowsAffected,10)
	Log.Info("get message templates from the database")
	return messageTemplate, count
}

/* function for get message template from message template table
    @param --> null
    @param value --> null
    description --> get message templates from message template table
    @return --> message template array
*/
func GetMessageTemplate(id string) (APIMessageTemplate2){
  ConnectToDatabase()
  db.Model(&MessageTemplate{}).Where("id = ?", id).Find(&messageTemplate2)
	Log.Info("get message template from the database")
	return messageTemplate2
}

/* function for get last template from message template table
    @param --> null
    @param value --> null
    description --> get last template from message template table
    @return --> message template array
*/
func GetLastTemplate() (APITemplateId){
  ConnectToDatabase()
  db.Model(&MessageTemplate{}).Last(&templateId)
	Log.Info("get last template from the database")
	return templateId
}

/* Save bulk messages to the database
    @param --> null
    @param value --> null
    description --> save bulk messages to db
    @return --> null
*/
func SaveBulkMessages(bulkType string, importId string, filterJson string, subject string, messageJson string, url string) {
  ConnectToDatabase()

  // get line no to find key value pair
  line_num := getLineNo(importId, filterJson)
  id := getKeyValue(line_num, importId, "ID")

  // create subject with its key --> value
  _, subject_key := readSubject(subject)
  subject_value := getKeyValue(line_num, importId, subject_key)
  subject = makeSubject(subject, subject_value)

  createMessageRequestBody(id, subject, url)
  makeMessageJson(line_num, importId, messageJson)

	bulk_message := []BulkMessage{{StaffNum: 333, BulkType: bulkType, ImportId: importId, FilterJson: filterJson, Subject: subject, MessagesJson: messageJson, URL: url, CreatedAt: time.Now().Local()}}
  db.Create(&bulk_message)
  Log.Info("saved bulk message to the database")
}

/* function for get last bulk message id
    @param --> null
    @param value --> null
    description --> get last bulk message id from bulk messages table
    @return --> message template array
*/
func GetLastBulkMessageId() (APIBulkMessageId){
  ConnectToDatabase()
  db.Model(&BulkMessage{}).Last(&bulkMessageId)
	Log.Info("get last bulk message Id from the database")
	return bulkMessageId
}

/* function for get bulk messages from bulk message table
    @param --> null
    @param value --> null
    description --> get bulk messages from bulk message table
    @return --> bulk message array
*/
func GetBulkMessages(page string) ([]APIBulkMessage, string){
  ConnectToDatabase()
  result := db.Model(&BulkMessage{}).Where("deleted_at IS NULL").Find(&bulkMessages)
  db.Model(&BulkMessage{}).Where("deleted_at IS NULL").Scopes(Paginate(page)).Find(&bulkMessages)
  count := strconv.FormatInt(result.RowsAffected,10)
	Log.Info("get bulk messages from the database")
	return bulkMessages, count
}

/* function for get message template from message template table
    @param --> null
    @param value --> null
    description --> get message templates from message template table
    @return --> message template array
*/
func GetBulkMessage(id string) ([]APIBulkMessage2, string){
  ConnectToDatabase()
  result := db.Model(&BulkMessage{}).Where("bulk_message_id = ?", id).Find(&bulkMessages2)
  count := strconv.FormatInt(result.RowsAffected,10)
	Log.Info("get bulk message from the database")
	return bulkMessages2, count
}

/* function to get line number from import data table
    @param --> null
    @param value --> null
    description --> get last template from message template table
    @return --> message template array
*/
func getLineNo(importId string, filterJson string) int{
  ConnectToDatabase()
  var importDataObj ImportData
  var header_id, line_no int
  header_list := getHeaderList(importId)
  header_key, header_value := readFilterJson(filterJson)
  
  for i, s := range header_list {
    if s == header_key{
      header_id = i
    }
  }

  db.Find(&importDataObj, "import_id = ? AND header_id = ? AND value = ?", importId, header_id, header_value)
  line_no = importDataObj.LineNum
  Log.Info("get Line Num from import data")
  return line_no
}

/* function to read filter json
    @param --> filterJson string
    @param value --> null
    description --> split filterJson to key and value 
    @return --> key and value
*/
func readFilterJson(filterJson string)(string, string){

  filter_json := strings.Split(filterJson, ":")

  filter_json_key := strings.TrimLeft(filter_json[0], "{")
  filter_json_key = strings.TrimRight(filter_json_key, "\"")
  filter_json_key = strings.TrimLeft(filter_json_key, "\"")


  filter_json_value := strings.TrimRight(filter_json[1], "}")
  filter_json_value = strings.TrimRight(filter_json_value, "\"")
  filter_json_value = strings.TrimLeft(filter_json_value, "\"")

	Log.Info("read filter json")
  return filter_json_key, filter_json_value
}


/* function to get line number from import data table
    @param --> null
    @param value --> null
    description --> get last template from message template table
    @return --> message template array
*/
func getKeyValue(lineNum int, importId string, key string) string {
  ConnectToDatabase()
  var importDataObj ImportData
  var header_id int
  var value string
  header_list := getHeaderList(importId)

  for i, s := range header_list {
    if s == key{
      header_id = i
    }
  }

  db.Find(&importDataObj, "import_id = ? AND header_id = ? AND line_num = ?", importId, header_id, lineNum)
  value = importDataObj.Value
  Log.Info("get key value from import data")
  return value
}

/* function to get header list from import history table
    @param --> import ID
    @param value --> string
    description --> get header list from import history table
    @return --> header list
*/

func getHeaderList(importId string) []string {
  ConnectToDatabase()
  db.Find(&importHistoryObj, "import_id = ?", importId)
  header_list := strings.Split(importHistoryObj.Header, " ")
  Log.Info("get header list from import history")
  return header_list
}

/* function to read subject
    @param --> subject string
    @param value --> null
    description --> split subject to find key
    @return --> key of subject
*/
func readSubject(subject string)([]string, string){

  re := regexp.MustCompile(`\$([^0-9A-Za-z_]+)\$`)
  subject_arr := re.FindStringSubmatch(subject)

  subject_key := subject_arr[1]

	Log.Info("read subject")
  return subject_arr, subject_key
}

/* function to make subject
    @param --> subject string, value string
    @param value --> null
    description --> split subject to find key
    @return --> modified subject
*/
func makeSubject(subject string, value string)(string){

  subject_array, _ := readSubject(subject)
  key := subject_array[0]
  subject = strings.Replace(subject, key, value, 1)
	Log.Info("make subject")
  return subject
}

/* function to read messageJson 
    @param --> messageJson string
    @param value --> string
    description --> split messageJson to find keys
    @return --> return array of keys
*/
func readMessagesJson(msgJson string)([]string){
  re := regexp.MustCompile(`\$([^ab]+)\$`)
  msg_arr := re.FindAllString(msgJson,-1)
  Log.Info("read message json")
  fmt.Println(msg_arr)
  return msg_arr
}

/* function to make subject
    @param --> subject string, value string
    @param value --> null
    description --> split subject to find key
    @return --> modified subject
*/
func makeMessageJson(lineNum int, importId string, messageJson string)(string){
  var arr[] string
  var result[] string
  msg_json := readMessagesJson(messageJson)

  for _, m := range msg_json { 
    v := strings.TrimLeft(m, "$")
    v = strings.TrimRight(v, "$")
    arr = append(arr,v)
	}

  for _, n := range arr { 
    value := getKeyValue(lineNum, importId, n)
    result = append(result, value)
	}

  for i,j := range msg_json {
    messageJson = strings.Replace(messageJson, j, result[i], 1)
  }
  
	Log.Info("make message json")
  fmt.Println(messageJson)
  return messageJson
}






