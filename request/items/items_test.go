package items

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/haisum/simplepos/db"
	"github.com/haisum/simplepos/db/models/items"
	assert2 "github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"gopkg.in/doug-martin/goqu.v4"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestList(t *testing.T) {
	mockdb, mock, _ := db.ConnectMock("sqlite3", t)
	defer mockdb.Db.Close()
	query, _, _ := db.Get().From("Item").ToSql()
	assert := assert2.New(t)
	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(sqlmock.NewRows([]string{"ID", "Name"}).AddRow(1, "First Item"))
	req := httptest.NewRequest("GET", "https://localhost:8443/items", bytes.NewBufferString("{}"))
	w := httptest.NewRecorder()
	List(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(http.StatusOK, resp.StatusCode)
	rspExpected := struct {
		Error error
		Items []items.Item
		Ok    bool
	}{}
	json.Unmarshal(body, &rspExpected)
	assert.Equal(true, rspExpected.Ok)
	assert.Equal(int64(1), rspExpected.Items[0].ID)
	assert.Equal("First Item", rspExpected.Items[0].Name)
	assert.Equal(int64(0), rspExpected.Items[0].Stock)
	assert.Nil(mock.ExpectationsWereMet())
}

func TestAdd(t *testing.T) {
	mockdb, mock, _ := db.ConnectMock("sqlite3", t)
	defer mockdb.Db.Close()
	itemList := []items.Item{
		{
			ID:   1,
			Name: "Hello world",
		},
		{
			ID:   2,
			Name: "Hello world 2",
		},
	}
	query, _, _ := db.Get().From("Item").ToInsertSql(itemList)
	assert := assert2.New(t)
	mock.ExpectExec(regexp.QuoteMeta(query)).WillReturnResult(sqlmock.NewResult(0, 2))
	jsItems, _ := json.Marshal(&itemList)
	req := httptest.NewRequest("POST", "https://localhost:8443/itemList", bytes.NewBuffer(jsItems))
	w := httptest.NewRecorder()
	Add(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(http.StatusOK, resp.StatusCode)
	rspExpected := struct {
		Error    error
		Affected int64
		Ok       bool
	}{}
	json.Unmarshal(body, &rspExpected)
	assert.Equal(true, rspExpected.Ok)
	assert.Equal(int64(2), rspExpected.Affected)
	assert.Nil(mock.ExpectationsWereMet())
}

func TestUpdate(t *testing.T) {
	mockdb, mock, _ := db.ConnectMock("sqlite3", t)
	defer mockdb.Db.Close()
	assert := assert2.New(t)

	itemList := []items.Item{
		{
			ID:   1,
			Name: "Hello world",
		},
		{
			ID:   2,
			Name: "Hello world 2",
		},
	}
	jsItems, _ := json.Marshal(&itemList)

	//scenario 1, happy all commit
	req := httptest.NewRequest("PUT", "https://localhost:8443/itemList", bytes.NewBuffer(jsItems))
	w := httptest.NewRecorder()
	mock.ExpectBegin()
	query, _, _ := db.Get().From("Item").Where(goqu.I("ID").Eq(itemList[0].ID)).ToUpdateSql(itemList[0])
	mock.ExpectExec(regexp.QuoteMeta(query)).WillReturnResult(sqlmock.NewResult(0, 1))
	query, _, _ = db.Get().From("Item").Where(goqu.I("ID").Eq(itemList[1].ID)).ToUpdateSql(itemList[1])
	mock.ExpectExec(regexp.QuoteMeta(query)).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	Update(w, req)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(http.StatusOK, resp.StatusCode)
	rspExpected := struct {
		Error    error
		Affected int64
		Ok       bool
	}{}
	json.Unmarshal(body, &rspExpected)
	assert.Equal(true, rspExpected.Ok)
	assert.Equal(int64(2), rspExpected.Affected)

	assert.Nil(mock.ExpectationsWereMet())
}

func TestUpdate_RollBack(t *testing.T) {
	mockdb, mock, _ := db.ConnectMock("sqlite3", t)
	defer mockdb.Db.Close()
	assert := assert2.New(t)

	itemList := []items.Item{
		{
			ID:   1,
			Name: "Hello world",
		},
		{
			ID:   2,
			Name: "Hello world 2",
		},
	}
	jsItems, _ := json.Marshal(&itemList)
	req := httptest.NewRequest("PUT", "https://localhost:8443/itemList", bytes.NewBuffer(jsItems))
	w := httptest.NewRecorder()

	mock.ExpectBegin()
	query, _, _ := db.Get().From("Item").Where(goqu.I("ID").Eq(itemList[0].ID)).ToUpdateSql(itemList[0])
	mock.ExpectExec(regexp.QuoteMeta(query)).WillReturnResult(sqlmock.NewResult(0, 1))
	query, _, _ = db.Get().From("Item").Where(goqu.I("ID").Eq(itemList[1].ID)).ToUpdateSql(itemList[1])
	mock.ExpectExec(regexp.QuoteMeta(query)).WillReturnError(errors.New("error in second update"))
	mock.ExpectRollback()
	Update(w, req)
	resp := w.Result()
	assert.Equal(http.StatusInternalServerError, resp.StatusCode)

	assert.Nil(mock.ExpectationsWereMet())
}

func TestDelete(t *testing.T) {
	mockdb, mock, _ := db.ConnectMock("sqlite3", t)
	defer mockdb.Db.Close()
	assert := assert2.New(t)

	ids := []int{1, 2, 3, 4, 5}
	mock.ExpectBegin()
	query, _, _ := db.Get().From("Item").Where(goqu.I("ID").In(ids)).ToDeleteSql()
	mock.ExpectExec(regexp.QuoteMeta(query)).WillReturnResult(sqlmock.NewResult(0, 5))
	mock.ExpectCommit()
	request := httptest.NewRequest("DELETE", "/items", bytes.NewBuffer([]byte("[1,2,3,4,5]")))
	w := httptest.NewRecorder()
	Delete(w, request)

	resp := w.Result()

	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(http.StatusOK, resp.StatusCode)
	rspExpected := struct {
		Error    error
		Affected int64
		Ok       bool
	}{}
	json.Unmarshal(body, &rspExpected)
	assert.Equal(true, rspExpected.Ok)
	assert.Equal(int64(5), rspExpected.Affected)

	assert.Nil(mock.ExpectationsWereMet())
}

func TestDelete_Rollback(t *testing.T) {
	mockdb, mock, _ := db.ConnectMock("sqlite3", t)
	defer mockdb.Db.Close()
	assert := assert2.New(t)

	ids := []int{1, 2, 3, 4, 5}
	mock.ExpectBegin()
	query, _, _ := db.Get().From("Item").Where(goqu.I("ID").In(ids)).ToDeleteSql()
	mock.ExpectExec(regexp.QuoteMeta(query)).WillReturnError(errors.New(""))
	mock.ExpectRollback()
	request := httptest.NewRequest("DELETE", "/items", bytes.NewBuffer([]byte("[1,2,3,4,5]")))
	w := httptest.NewRecorder()
	Delete(w, request)

	resp := w.Result()

	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(http.StatusInternalServerError, resp.StatusCode)
	t.Logf("%s", body)

	assert.Nil(mock.ExpectationsWereMet())
}
