package object

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type JSObjectO struct {
	ObjectID       string      `json:"object_id"`
	ObjectType     string      `json:"object_type"`
	ObjectExtender interface{} `json:"object_extender"`
}

type BlockO struct {
	ChangeCB string `json:"change_callback"`
	ClickCB  string `json:"click_callback"`
}

type TextO struct {
	TextByDefault string `json:"text_by_default"`
	ChangeCB      string `json:"change_callback"`
}

type ButtonO struct {
	ButtonName string `json:"button_name"`
	ClickCB    string `json:"click_callback"`
}

type ListItemO struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

type ListO struct {
	ListItems []ListItemO `json:"list_items"`
	ChangeCB  string      `json:"change_callback"`
}

type SelectorItemO struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

type SelectorO struct {
	SelectorItems []SelectorItemO `json:"selector_items"`
	ChangeCB      string          `json:"change_callback"`
}

type TableItemO struct {
	Data string `json:"data"`
}

type TableRowItemO struct {
	ID        string       `json:"id"`
	TableItem []TableItemO `json:"table_row_items"`
}

type TableO struct {
	HeaderItems   []TableItemO    `json:"header_items"`
	TableRowItems []TableRowItemO `json:"table_row_items"`
	ChangeCB      string          `json:"change_callback"`
	ClickCB       string          `json:"click_callback"`
	ItemChangeCB  string          `json:"item_change_callback"`
	ItemClickCB   string          `json:"item_click_callback"`
	RowClickCB    string          `json:"row_click_callback"`
}

type TreeItemO struct {
	Data      string      `json:"data"`
	TreeItems []TreeItemO `json:"tree_items"`
}

type TreeO struct {
	TreeItems []TreeItemO `json:"tree_items"`
	ChangeCB  string      `json:"change_callback"`
	ClickCB   string      `json:"click_callback"`
}

type DocConstructorItem struct {
	Object *JSObjectO `json:"object"`
}

type DocConstructor struct {
	DocConstructors []DocConstructorItem `json:"doc_constructors"`
}

func LoadDocConstructorFromString(template string) (*DocConstructor, error) {
	var dc DocConstructor

	err := json.Unmarshal([]byte(template), &dc)
	if err != nil {
		fmt.Println("error:", err)
		return nil, err
	}
	return &dc, nil
}

func transcode(in, out interface{}) {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(in)
	json.NewDecoder(buf).Decode(out)
}

func CreateDocConstructor() ([]*JSObject, map[string]*JSObject) {
	doc := NewDocument()
	DocConstructorTextSection := doc.Object.Call("getElementById", "DocConstructor")

	if DocConstructorTextSection.IsNull() {
		return nil, nil
	}

	DocConstructorText := DocConstructorTextSection.Get("innerHTML").String()
	DocConstructor, err := LoadDocConstructorFromString(DocConstructorText)
	if err != nil {
		return nil, nil
	}
	jsoa := []*JSObject{}
	jsoDict := make(map[string]*JSObject)

	for i := range DocConstructor.DocConstructors {
		dci := DocConstructor.DocConstructors[i]
		//fmt.Printf("dci %v dci.Object %#v\r\n", dci, dci.Object)
		jso := &JSObject{}
		switch dci.Object.ObjectType {
		case "block":
			jso = doc.NewBlock(dci.Object.ObjectID)
			b := Block{}
			bo := BlockO{}
			transcode(dci.Object.ObjectExtender, &bo)
			b.ChangeCBName = bo.ChangeCB
			b.ClickCBName = bo.ClickCB
			jso.ObjectExtender = b
		case "text":
			to := TextO{}
			transcode(dci.Object.ObjectExtender, &to)
			jso = doc.NewText(dci.Object.ObjectID, to.TextByDefault)
			t := Text{}
			t.ChangeCBName = to.ChangeCB
			jso.ObjectExtender = t
		case "button":
			bo := ButtonO{}
			transcode(dci.Object.ObjectExtender, &bo)
			jso = doc.NewButton(dci.Object.ObjectID, bo.ButtonName)
			b := Button{}
			b.ClickCBName = bo.ClickCB
			jso.ObjectExtender = b
		case "list":
			lo := ListO{}
			transcode(dci.Object.ObjectExtender, &lo)
			listID := []string{}
			listData := []string{}
			for i := range lo.ListItems {
				listID = append(listID, lo.ListItems[i].ID)
				listData = append(listData, lo.ListItems[i].Data)
			}
			jso = doc.NewList(dci.Object.ObjectID, listID, listData)
			l := List{}
			l.ChangeCBName = lo.ChangeCB
			jso.ObjectExtender = l
		case "selector":
			lo := SelectorO{}
			transcode(dci.Object.ObjectExtender, &lo)
			listID := []string{}
			listData := []string{}
			for i := range lo.SelectorItems {
				listID = append(listID, lo.SelectorItems[i].ID)
				listData = append(listData, lo.SelectorItems[i].Data)
			}
			jso = doc.NewSelector(dci.Object.ObjectID, listID, listData)
			l := Selector{}
			l.ChangeCBName = lo.ChangeCB
			jso.ObjectExtender = l
		case "table":
			lo := TableO{}
			transcode(dci.Object.ObjectExtender, &lo)
			listID := []string{}
			listData := [][]string{}
			for i := range lo.HeaderItems {
				listID = append(listID, lo.HeaderItems[i].Data)
			}
			for i := range lo.TableRowItems {
				ld := []string{}
				for j := range lo.TableRowItems {
					ld = append(ld, lo.TableRowItems[i].TableItem[j].Data)
				}
				listData = append(listData, ld)
			}
			jso = doc.NewTable(dci.Object.ObjectID, listID, listData)
			l := Table{}
			l.ClickCBName = lo.ClickCB
			l.ChangeCBName = lo.ChangeCB
			l.RowClickCBName = lo.RowClickCB
			l.ItemClickCBName = lo.ItemClickCB
			l.ItemChangeCBName = lo.ItemChangeCB
			jso.ObjectExtender = l
		case "tree":
			lo := TreeO{}
			transcode(dci.Object.ObjectExtender, &lo)

			type ConvertFunc func(tree_data []TreeItemO) []TreeData
			var cf ConvertFunc
			cf = func(treeDataIn []TreeItemO) []TreeData {
				tia := []TreeData{}
				for i := range treeDataIn {
					ti := TreeData{}
					ti.Data = treeDataIn[i].Data
					ti.TreeDatas = cf(treeDataIn[i].TreeItems)
					tia = append(tia, ti)
				}
				return tia
			}
			treeData := cf(lo.TreeItems)

			listID := "tree"
			jso = doc.NewTree(dci.Object.ObjectID, listID, treeData)
			l := Tree{}
			l.ClickCBName = lo.ClickCB
			l.ChangeCBName = lo.ChangeCB
			jso.ObjectExtender = l
		}
		jsoa = append(jsoa, jso)
		jsoDict[dci.Object.ObjectID] = jso
	}
	return jsoa, jsoDict
}
