package model

import (
	"bytes"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"

	"github.com/gorilla/context"
	gonanoid "github.com/matoous/go-nanoid"
)

//QUESTION - Module Name
var QUESTION = "FEEDBACK_QUESTION"

// AddQuestion - Add Meal Type
func AddQuestion(r *http.Request, reqMap data.FeedBackQuestion) bool {

	util.LogIt(r, "Model - V_Feedback_Question - AddQuestion")

	var Qry bytes.Buffer
	nanoid, _ := gonanoid.Nanoid()

	Qry.WriteString("INSERT INTO cf_feedback_questions(id, question, created_at, created_by) VALUES (?,?,?,?)")
	err := ExecuteNonQuery(Qry.String(), nanoid, reqMap.Question, util.GetIsoLocalDateTime(), context.Get(r, "UserId"))
	if util.CheckErrorLog(r, err) {
		return false
	}

	AddLog(r, "", QUESTION, "Create", nanoid, GetLogsValueMap(r, util.ToMap(reqMap), false, ""))

	return true
}

// UpdateQuestion - Update Meal Type
func UpdateQuestion(r *http.Request, reqMap data.FeedBackQuestion) bool {

	util.LogIt(r, "Model - V_Feedback_Question - UpdateQuestion")

	var Qry bytes.Buffer

	BeforeUpdate, _ := GetModuleFieldByID(r, QUESTION, reqMap.ID, "question")

	Qry.WriteString("UPDATE cf_feedback_questions SET question = ? WHERE id = ?")
	err := ExecuteNonQuery(Qry.String(), reqMap.Question, reqMap.ID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	AddLog(r, BeforeUpdate.(string), QUESTION, "Update", reqMap.ID, GetLogsValueMap(r, util.ToMap(reqMap), true, "ID"))

	return true
}

// QuestionListing - Get Meal Type Listing
func QuestionListing(r *http.Request, reqMap data.JQueryTableUI) (map[string]interface{}, error) {
	util.LogIt(r, "Model - V_Feedback_Question - QuestionListing")

	var Qry, QryCnt, QryFilter, GroupBy, FromQry bytes.Buffer

	var status = make(map[string]interface{})

	var testColArrs [20]string
	testColArrs[0] = "CFQ.id"
	testColArrs[1] = "CFQ.question"
	testColArrs[2] = "CFQ.status"
	testColArrs[3] = "created_by"

	var testArrs []map[string]string
	testArrs = append(testArrs, map[string]string{
		"key":   "question",
		"value": "CFQ.question",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "status",
		"value": "CFQ.status",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "created_at",
		"value": "DATE(from_unixtime(CFQ.created_at))",
	})

	QryCnt.WriteString(" COUNT(CFQ.id) AS cnt ")
	QryFilter.WriteString(" COUNT(CFQ.id) AS cnt ")

	Qry.WriteString(" CFQ.id, CFQ.question, CONCAT(from_unixtime(CFQ.created_at),' ',CU.username) AS created_by, ST.status, ST.id AS status_id ")

	FromQry.WriteString(" FROM cf_feedback_questions AS CFQ ")
	FromQry.WriteString(" INNER JOIN cf_user AS CU ON CU.id = CFQ.created_by  ")
	FromQry.WriteString(" INNER JOIN status AS ST ON ST.id = CFQ.status ")
	FromQry.WriteString(" WHERE CFQ.status <> 3 ")
	Data, err := JQueryTable(r, reqMap, Qry, FromQry, QryCnt, QryFilter, GroupBy, testArrs, testColArrs, status, "")
	return Data, err
}
