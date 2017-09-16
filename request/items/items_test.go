package items

import (
	"testing"
	"net/http/httptest"
	"io/ioutil"
	"bytes"
	"github.com/haisum/simplepos/db"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	assert2 "github.com/stretchr/testify/assert"
	"regexp"
	"github.com/haisum/simplepos/db/models/items"
	"encoding/json"
)

func TestList(t *testing.T){
	mockdb,mock, _ := db.ConnectMock("sqlite3", t)
	defer mockdb.Db.Close()
	query, _,_ := db.Get().From("Item").ToSql()
	assert := assert2.New(t)
	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(sqlmock.NewRows([]string{"ID", "Name"}).AddRow(1, "First Item"))
	req := httptest.NewRequest("GET", "https://localhost:8443/items", bytes.NewBufferString("{}"))
	w := httptest.NewRecorder()
	List(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(200, resp.StatusCode)
	rspExpected := struct{
		Error error
		Items []items.Item
		Ok bool
	}{}
	json.Unmarshal(body, &rspExpected)
	assert.Equal(true, rspExpected.Ok )
	assert.Equal(int64(1), rspExpected.Items[0].ID)
	assert.Equal("First Item", rspExpected.Items[0].Name)
	assert.Equal(int64(0), rspExpected.Items[0].Stock)
}


func TestAdd(t *testing.T){
	mockdb,mock, _ := db.ConnectMock("sqlite3", t)
	defer mockdb.Db.Close()
	items := []items.Item{
		{
			ID: 1,
			Name : "Hello world",
		},
		{
			ID: 2,
			Name : "Hello world 2",
		},
	}
	query, _,_ := db.Get().From("Item").ToInsertSql(items)
	assert := assert2.New(t)
	mock.ExpectExec(regexp.QuoteMeta(query)).WillReturnResult(sqlmock.NewResult(0, 2))
	jsItems, _ := json.Marshal(&items)
	req := httptest.NewRequest("POST", "https://localhost:8443/items", bytes.NewBuffer(jsItems))
	w := httptest.NewRecorder()
	Add(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(200, resp.StatusCode)
	rspExpected := struct{
		Error error
		Affected int64
		Ok bool
	}{}
	json.Unmarshal(body, &rspExpected)
	assert.Equal(true, rspExpected.Ok )
	assert.Equal(int64(2), rspExpected.Affected)
}
