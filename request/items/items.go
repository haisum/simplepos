package items

import (
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/haisum/simplepos/db/models/items"
	"github.com/haisum/simplepos/request"
	"github.com/haisum/simplepos/stringutils"
)

//List handler lists/filters items based on some criteria
var List = http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
	var criteria items.Criteria
	err := request.Decode(r, w, &criteria)
	if err != nil {
		return
	}
	itemList, err := items.List(criteria)
	if err != nil {
		request.RespondWithError(w, err, "Couldn't find records.", http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"Items" : itemList,
		"Ok" : err == nil,
		"Error" : err,
	}
	jsonResp, err := json.Marshal(response)
	if err != nil {
		request.RespondWithError(w, err, "Couldn't convert to json", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, stringutils.ByteToString(jsonResp))
})

//Update handler updates Items
var Update = http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
	var itemList []items.Item
	err := request.Decode(r, w, &itemList)
	if err != nil {
		return
	}
	affected, err := items.Update(itemList)
	if err != nil {
		request.RespondWithError(w, err, "Couldn't update records.", http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"Affected" : affected,
		"Ok" : err == nil,
		"Error" : err,
	}
	jsonResp, err := json.Marshal(response)
	if err != nil {
		request.RespondWithError(w, err, "Couldn't convert to json", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, stringutils.ByteToString(jsonResp))
})

//Add handler adds one or more Items
var Add = http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
	var itemList []items.Item
	err := request.Decode(r, w, &itemList)
	if err != nil {
		return
	}
	affected, err := items.Add(itemList)
	if err != nil {
		request.RespondWithError(w, err, "Couldn't add records.", http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"Affected" : affected,
		"Ok" : err == nil,
		"Error" : err,
	}
	jsonResp, err := json.Marshal(response)
	if err != nil {
		request.RespondWithError(w, err, "Couldn't convert to json", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, stringutils.ByteToString(jsonResp))
})

//Delete handler deletes one or more Items
var Delete = http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
	var ids []int64
	err := request.Decode(r, w, &ids)
	if err != nil {
		return
	}
	affected, err := items.Delete(ids)
	if err != nil {
		request.RespondWithError(w, err, "Couldn't delete records.", http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"Affected" : affected,
		"Ok" : err == nil,
		"Error" : err,
	}
	jsonResp, err := json.Marshal(response)
	if err != nil {
		request.RespondWithError(w, err, "Couldn't convert to json", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, stringutils.ByteToString(jsonResp))
})