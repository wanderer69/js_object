package object

import (
	"reflect"
	"syscall/js"
)

type Block struct {
	ChangeCBName string        `json:"change_callback_name"`
	ClickCBName  string        `json:"click_callback_name"`
	ChangeCB     reflect.Value `json:"change_callback"`
	ClickCB      reflect.Value `json:"click_callback"`
}

type Text struct {
	TextByDefault string        `json:"text_by_default"`
	ChangeCBName  string        `json:"change_callback_name"`
	ChangeCB      reflect.Value `json:"change_callback"`
}

type Label struct {
	TextByDefault string        `json:"text_by_default"`
	ClickCBName   string        `json:"click_callback_name"`
	ClickCB       reflect.Value `json:"click_callback"`
}

type Image struct {
	ClickCBName string        `json:"click_callback_name"`
	ClickCB     reflect.Value `json:"click_callback"`
}

type Ordinary struct {
	ChangeCBName string        `json:"change_callback_name"`
	ChangeCB     reflect.Value `json:"change_callback"`
}

type Button struct {
	ButtonName  string        `json:"button_name"`
	ClickCBName string        `json:"click_callback_name"`
	ClickCB     reflect.Value `json:"click_callback"`
}

type ListItem struct {
	ID      string `json:"id"`
	Data    string `json:"data"`
	OnClick func(js.Value, []js.Value, interface{}) interface{}
}

type List struct {
	ListItems       []ListItem `json:"list_items"`
	ListIDDict      map[string]ListItem
	ChangeCBName    string        `json:"change_callback_name"`
	ChangeCB        reflect.Value `json:"change_callback"`
	Tags_lst        []string
	CurrentListItem *ListItem
}

type SelectorItem struct {
	ID      string `json:"id"`
	Data    string `json:"data"`
	OnClick func(js.Value, []js.Value, interface{}) interface{}
}

type Selector struct {
	SelectorItems       []SelectorItem `json:"selector_items"`
	SelectorIDDict      map[string]SelectorItem
	ChangeCBName        string        `json:"change_callback_name"`
	ChangeCB            reflect.Value `json:"change_callback"`
	Tags_lst            []string
	CurrentSelectorItem *SelectorItem
}

type TableItem struct {
	ID       string        `json:"id"`
	Data     string        `json:"data"`
	ClickCB  reflect.Value `json:"click_callback"`
	ChangeCB reflect.Value `json:"change_callback"`
}

type TableRowItem struct {
	ID          string      `json:"id"`
	TableItems  []TableItem `json:"table_items"`
	TableIDDict map[string]TableItem
	ClickCB     reflect.Value `json:"click_callback"`
}

type Table struct {
	HeaderItems         []TableItem    `json:"header_items"`
	TableRowItems       []TableRowItem `json:"table_row_items"`
	TableIDDict         map[string]TableRowItem
	ChangeCBName        string        `json:"change_callback_name"`
	ChangeCB            reflect.Value `json:"change_callback"`
	ClickCBName         string        `json:"click_callback_name"`
	ClickCB             reflect.Value `json:"click_callback"`
	Tags_lst            []string
	CurrentTableRowItem *TableRowItem
	ItemChangeCBName    string        `json:"item_change_callback_name"`
	ItemChangeCB        reflect.Value `json:"item_change_callback"`
	ItemClickCBName     string        `json:"item_click_callback_name"`
	ItemClickCB         reflect.Value `json:"item_click_callback"`
	RowClickCBName      string        `json:"row_click_callback_name"`
	RowClickCB          reflect.Value `json:"row_click_callback"`
	ID                  string
}

type TreeItem struct {
	ID        string     `json:"id"`
	Data      string     `json:"data"`
	TreeItems []TreeItem `json:"tree_items"`
}

type Tree struct {
	TreeItems    []TreeItem    `json:"tree_items"`
	ChangeCBName string        `json:"change_callback_name"`
	ChangeCB     reflect.Value `json:"change_callback"`
	ClickCBName  string        `json:"click_callback_name"`
	ClickCB      reflect.Value `json:"click_callback"`
}

type TreeData struct {
	Data      string     `json:"data"`
	TreeDatas []TreeData `json:"tree_data"`
}
