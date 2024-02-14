package object

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"syscall/js"
)

type DocObject struct {
	Object js.Value `json:"object"`
	Type   string   `json:"type_html_object"` // normal, figma
}

type JSObject struct {
	Object         js.Value    `json:"object"`
	ObjectData     js.Value    `json:"object_data"` // object for linked data
	ObjectID       string      `json:"object_id"`
	ObjectType     string      `json:"object_type"`
	Document       *DocObject  `json:"document"`
	ObjectExtender interface{} `json:"object_extender"`
	Childs         []*JSObject `json:"childs list"`
	Parent         *JSObject   `json:"parent object"`
}

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
	ID   string `json:"id"`
	Data string `json:"data"`
	//    OnClick func(js.Value, []js.Value, interface{}) interface{}
}

type TableRowItem struct {
	ID          string      `json:"id"`
	TableItems  []TableItem `json:"table_items"`
	TableIDDict map[string]TableItem
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
}

type TreeItem struct {
	ID        string     `json:"id"`
	Data      string     `json:"data"`
	TreeItems []TreeItem `json:"tree_items"`
	////    OnClick func(js.Value, []js.Value, interface{}) interface{}
}

type Tree struct {
	TreeItems    []TreeItem    `json:"tree_items"`
	ChangeCBName string        `json:"change_callback_name"`
	ChangeCB     reflect.Value `json:"change_callback"`
	ClickCBName  string        `json:"click_callback_name"`
	ClickCB      reflect.Value `json:"click_callback"`
}

func NewDocument() *DocObject {
	doc := js.Global().Get("document")
	do := DocObject{}
	do.Object = doc
	return &do
}

func GetObjectByName(doc *DocObject, id string) js.Value {
	obj := doc.Object.Call("getElementById", id)
	if obj.IsNull() {
		obj_a := doc.Object.Call("getElementsByClassName", id)
		if obj_a.IsNull() {
			return obj_a
		}
		obj_l := DOMObjectToList(obj_a)
		obj = obj_l[0]
	}
	return obj
}

func NewJSObject(doc *DocObject, id string, o_type string, oe interface{}) *JSObject {
	obj := GetObjectByName(doc, id)
	childs := []*JSObject{}
	parent := &JSObject{}

	jso := JSObject{obj, obj, id, o_type, doc, oe, childs, parent}
	return &jso
}

type ObjCallBack func(this js.Value, args []js.Value) interface{}

func (jso *JSObject) SetCallBack(name string, ocb ObjCallBack) {
	jso.Object.Call("addEventListener", name, js.FuncOf(ocb), false)
}

func (jso *JSObject) AddChild(jso_c *JSObject) {
	jso.Childs = append(jso.Childs, jso_c)
}

func (jso *JSObject) SetParent(jso_p *JSObject) {
	jso.Parent = jso_p
}

func (doc *DocObject) NewBlock(id string) *JSObject {
	b := Block{}
	jso := NewJSObject(doc, id, "block", b)
	return jso
}

func (jso *JSObject) SetValue(value string) {
	switch jso.ObjectType {
	case "block":
		jso.Object.Set("innerHTML", value)
	case "label":
		jso.Object.Set("innerHTML", value)
	case "text":
		doc := jso.Document
		if doc.Type == "figma" {
			jso.ObjectData.Set("value", value)
		} else {
			jso.Object.Set("value", value)
		}
	case "button":
		jso.Object.Set("innerHTML", value)
	}
}

func (jso *JSObject) GetValue() string {
	switch jso.ObjectType {
	case "block":
		return jso.Object.Get("innerHTML").String()
	case "label":
		return jso.Object.Get("innerHTML").String()
	case "text":
		doc := jso.Document
		if doc.Type == "figma" {
			return jso.ObjectData.Get("value").String()
		} else {
			return jso.Object.Get("value").String()
		}
	}
	return ""
}

func (jso *JSObject) Enable() {
	switch jso.ObjectType {
	case "block":
		jso.Object.Get("style").Set("display", "block")
	case "text":
		jso.Object.Get("style").Set("display", "block")
	case "image":
		jso.Object.Get("style").Set("display", "block")
	case "label":
		jso.Object.Get("style").Set("display", "block")
	case "button":
		jso.Object.Get("style").Set("display", "block")
	}
}

func (jso *JSObject) Disable() {
	switch jso.ObjectType {
	case "block":
		jso.Object.Get("style").Set("display", "none")
	case "text":
		jso.Object.Get("style").Set("display", "none")
	case "image":
		jso.Object.Get("style").Set("display", "none")
	case "label":
		jso.Object.Get("style").Set("display", "none")
	case "button":
		jso.Object.Get("style").Set("display", "none")
	}
}

func (jso *JSObject) SetStyleProperty(property string, value string) {
	switch jso.ObjectType {
	case "block":
		jso.Object.Get("style").Set(property, value)
	case "text":
		jso.Object.Get("style").Set(property, value)
	case "image":
		jso.Object.Get("style").Set(property, value)
	case "label":
		jso.Object.Get("style").Set(property, value)
	case "button":
		jso.Object.Get("style").Set(property, value)
	}
}

func (doc *DocObject) NewText(id string, text_by_default string) *JSObject {
	t := Text{}
	t.TextByDefault = text_by_default
	jso := NewJSObject(doc, id, "text", t)
	if len(t.TextByDefault) > 0 {
		jso.SetValue(t.TextByDefault)
	}
	return jso
}

func (jso *JSObject) CorrectText(baseclass string) *JSObject {
	doc := jso.Document
	switch jso.ObjectType {
	case "text":
		if doc.Type == "figma" {
			txt := `<input id="%v_id" type="text" class="%v">`
			ss := fmt.Sprintf(txt, baseclass, baseclass)
			jso.Object.Set("innerHTML", ss)
			jso.ObjectData = GetObjectByName(doc, fmt.Sprintf("%v_id", baseclass))
		}
	case "label":
		if doc.Type == "figma" {
		}
	}

	return jso
}

func (doc *DocObject) NewLabel(id string, text_by_default string) *JSObject {
	lb := Label{}
	lb.TextByDefault = text_by_default
	jso := NewJSObject(doc, id, "label", lb)
	if len(lb.TextByDefault) > 0 {
		jso.SetValue(lb.TextByDefault)
	}
	return jso
}

func (jso *JSObject) SetLabelText(value string) {
	switch jso.ObjectType {
	case "label":
		jso.ObjectData.Set("textContent", value)
		//jso.Object.Get("style").Set("display", "none")
	}
}

func (doc *DocObject) NewImage(id string) *JSObject {
	im := Image{}
	jso := NewJSObject(doc, id, "image", im)
	return jso
}

func (doc *DocObject) NewOrdinary(id string) *JSObject {
	o := Ordinary{}
	jso := NewJSObject(doc, id, "ordinary", o)
	if doc.Type == "figma" {
		if false {
			jso.Object.Get("style").Set("display", "none")
		}
	}
	return jso
}

func (doc *DocObject) NewButton(id string, button_name string) *JSObject {
	b := Button{}
	b.ButtonName = button_name
	jso := NewJSObject(doc, id, "button", b)
	if len(b.ButtonName) > 0 {
		jso.SetValue(b.ButtonName)
	}
	return jso
}

func DOMObjectToList(o js.Value) []js.Value {
	var out []js.Value
	length := o.Get("length").Int()
	for i := 0; i < length; i++ {
		out = append(out, o.Call("item", i))
	}
	return out
}

func (doc *DocObject) NewList(id string, list_id []string, list_data []string) *JSObject {
	if len(list_id) != len(list_data) {
		return nil
	}
	l := List{}
	l.ListIDDict = make(map[string]ListItem)
	tags_lst := []string{}
	for i, item := range list_id {
		li := ListItem{}
		li.ID = item
		li.Data = list_data[i]
		tags_lst = append(tags_lst, fmt.Sprintf("<span class=\"custom-option\">%v</span><br>", item))
		//        OnClick func(*js.Object, []*js.Object) interface{}
		l.ListItems = append(l.ListItems, li)
		l.ListIDDict[item] = li
	}
	ss := strings.Join(tags_lst, "\r\n")
	jso := NewJSObject(doc, id, "list", l)
	jso.Object.Set("innerHTML", ss)

	return jso
}

func (jso *JSObject) SetList(cb reflect.Value) *JSObject {
	l := jso.ObjectExtender.(List)
	ol := jso.Document.Object.Call("querySelector", ".custom-select")
	cl := ol.Get("classList")
	cl.Call("toggle", "open")
	posl := jso.Document.Object.Call("querySelectorAll", ".custom-option")
	posl_ := DOMObjectToList(posl)
	for i := range posl_ {
		span_item := posl_[i]
		span_item.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			pn := this.Get("parentNode")
			ol := pn.Call("querySelector", ".custom-option.selected")
			if !ol.IsNull() { // js.Null()
				ol.Get("classList").Call("remove", "selected")
			}
			this.Get("classList").Call("add", "selected")
			val := this.Get("textContent").String()
			li := l.ListIDDict[val]
			l.CurrentListItem = &li
			if li.OnClick != nil {
				li.OnClick(this, args, l)
			}
			/*
			   fmt.Printf("--> %v this %v args %v\r\n", l.ChangeCB, this, args)
			   InvokeCall(l.ChangeCB, this, args, jso)
			*/
			fmt.Printf("--> %v this %v args %v\r\n", cb, this, args)
			InvokeCall(cb, this, args, jso)
			/*
			   if G_debug_i > 5 {
			       fmt.Printf("span click %v %v %v\r\n", this, this.Get("classList"), val)
			   }
			*/
			return js.Null()
		}), false)
	}

	return jso
}

func (doc *DocObject) NewSelector(id string, selector_id []string, selector_data []string) *JSObject {
	if len(selector_id) != len(selector_data) {
		return nil
	}
	l := Selector{}
	/*
	   <select id="yourDropDownElementId"><select/>

	   // Get the select element
	   var select = document.getElementById("yourDropDownElementId");
	   // Create a new option element
	   var el = document.createElement("option");
	   // Add our value to the option
	   el.textContent = "Example Value";
	   el.value = "Example Value";
	   // Set the option to selected
	   el.selected = true;
	   // Add the new option element to the select element
	   select.appendChild(el);
	*/
	l.SelectorIDDict = make(map[string]SelectorItem)
	//    tags_lst := []string{}
	for i, item := range selector_id {
		li := SelectorItem{}
		li.ID = item
		li.Data = selector_data[i]
		//	ss := fmt.Sprintf("<option value=\"%v\">%v</option>", i+1, item)
		//	tags_lst = append(tags_lst, ss)
		l.SelectorItems = append(l.SelectorItems, li)
		l.SelectorIDDict[item] = li
	}
	//    fmt.Printf("selector id %v\r\n", id)
	jso := NewJSObject(doc, id, "selector", l)
	for i, item := range selector_id {
		el := doc.Object.Call("createElement", "option")
		el.Set("textContent", item)
		el.Set("value", fmt.Sprintf("%v", i+1))
		//        el.Set("selected","true")
		jso.Object.Call("appendChild", el)
	}
	return jso
}

func (jso *JSObject) AddOption(selector_id []string) {
	switch jso.ObjectType {
	case "selector":
		for i, item := range selector_id {
			el := jso.Document.Object.Call("createElement", "option")
			el.Set("textContent", item)
			el.Set("value", fmt.Sprintf("%v", i+1))
			//              el.Set("selected","true")
			jso.Object.Call("appendChild", el)
		}
	}
}

func (doc *DocObject) NewSelector1(id string, selector_id []string, selector_data []string) *JSObject {
	if len(selector_id) != len(selector_data) {
		return nil
	}
	l := Selector{}
	l.SelectorIDDict = make(map[string]SelectorItem)
	tags_lst := []string{}
	tags_lst = append(tags_lst, "<select id=\""+id+"\">")
	fmt.Printf("selector_id %v\r\n", selector_id)
	for i, item := range selector_id {
		li := SelectorItem{}
		li.ID = item
		li.Data = selector_data[i]
		ss := fmt.Sprintf("<option value=\"%v\">%v</option>", i+1, item)
		tags_lst = append(tags_lst, ss)
		l.SelectorItems = append(l.SelectorItems, li)
		l.SelectorIDDict[item] = li
	}
	tags_lst = append(tags_lst, "</select>")
	fmt.Printf("selector id %v\r\n", id)
	ss := strings.Join(tags_lst, "\r\n")
	div_name := fmt.Sprintf("%v_", id)
	obj_div := doc.Object.Call("getElementById", div_name)
	obj_div.Set("innerHTML", ss)
	jso := NewJSObject(doc, id, "selector", l)
	return jso
}

func (doc *DocObject) NewTable(id string, table_id []string, table_data [][]string) *JSObject {
	if len(table_id) != len(table_data) {
		return nil
	}
	/*
	   var table = document.getElementById("myTable");
	   // Create an empty <thead> element and add it to the table:
	   var header = table.createTHead();
	   // Create an empty <tr> element and add it to the first position of <thead>:
	   var row = header.insertRow(0);
	   // Insert a new cell (<td>) at the first position of the "new" <tr> element:
	   var cell = row.insertCell(0);
	*/
	l := Table{}
	l.TableIDDict = make(map[string]TableRowItem)
	//tags_lst := []string{}
	//tags_lst = append(tags_lst, "<select id=\""+ id +"\">")
	fmt.Printf("table_id %v\r\n", table_id)
	for i, item := range table_id {
		li := TableRowItem{}
		li.ID = item
		td := table_data[i]
		for j := range td {
			lii := TableItem{}
			lii.Data = td[j]
			//ss := fmt.Sprintf("<option value=\"%v\">%v</option>", i+1, item)
			//tags_lst = append(tags_lst, ss)
			li.TableItems = append(li.TableItems, lii)
			l.TableIDDict[item] = li
		}
	}
	//tags_lst = append(tags_lst, "</select>")
	fmt.Printf("table id %v\r\n", id)

	jso := NewJSObject(doc, id, "table", l)
	if len(table_id) > 0 {
		header := jso.Object.Call("createTHead")
		// Create an empty <tr> element and add it to the first position of <thead>:
		row := header.Call("insertRow", 0)
		for i, item := range table_id {
			// Insert a new cell (<td>) at the first position of the "new" <tr> element:
			cell := row.Call("insertCell", i)
			cell.Set("innerHTML", item)
		}
		if len(table_data) > 0 {
			for i := range table_data {
				td := table_data[i]
				row := jso.Object.Call("insertRow", i)
				for j := range td {
					cell := row.Call("insertCell", j)
					cell.Set("innerHTML", td[j])
				}
			}
		}
	}
	return jso
}

func (jso *JSObject) SetTable(table_id []string, table_data [][]string) {
	switch jso.ObjectType {
	case "table":
		header := jso.Object.Call("createTHead")
		// Create an empty <tr> element and add it to the first position of <thead>:
		row := header.Call("insertRow", 0)
		for i, item := range table_id {
			// Insert a new cell (<td>) at the first position of the "new" <tr> element:
			cell := row.Call("insertCell", i)
			cell.Set("innerHTML", item)
		}
		for i := range table_data {
			td := table_data[i]
			row := jso.Object.Call("insertRow", i+1)
			for j := range td {
				cell := row.Call("insertCell", j)
				cell.Set("innerHTML", td[j])
			}
		}
	}
}

type TreeData struct {
	Data      string     `json:"data"`
	TreeDatas []TreeData `json:"tree_data"`
}

func (doc *DocObject) NewTree(id string, tree_id string, tree_data []TreeData) *JSObject {
	l := Tree{}
	type ConvertFunc func(tree_data_ []TreeData) []TreeItem
	var cf ConvertFunc
	g_num := 1
	cf = func(tree_data_ []TreeData) []TreeItem {
		tia := []TreeItem{}
		for i := range tree_data_ {
			ti := TreeItem{}
			ti.ID = fmt.Sprintf("%v", g_num)
			g_num = g_num + 1
			ti.Data = tree_data_[i].Data
			ti.TreeItems = cf(tree_data_[i].TreeDatas)
			tia = append(tia, ti)
		}
		return tia
	}
	l.TreeItems = cf(tree_data)

	jso := NewJSObject(doc, id, "tree", l)
	if len(tree_data) > 0 {
		/*
		   <ul id="myUL">
		     <li><span class="box">Beverages</span>
		       <ul class="nested">
		         <li>Water</li>
		         <li>Coffee</li>
		         <li><span class="box">Tea</span>
		           <ul class="nested">
		             <li>Black Tea</li>
		             <li>White Tea</li>
		             <li><span class="box">Green Tea</span>
		               <ul class="nested">
		                 <li>Sencha</li>
		                 <li>Gyokuro</li>
		                 <li>Matcha</li>
		                 <li>Pi Lo Chun</li>
		               </ul>
		             </li>
		           </ul>
		         </li>
		       </ul>
		     </li>
		   </ul>
		*/
		class_id := "class_" + id
		type SetFunc func(tree_items []TreeItem) string
		var sf SetFunc
		sf = func(tree_items []TreeItem) string {
			ss := ""
			for i := range tree_items {
				ti := tree_items[i]
				level := ""
				if len(ti.TreeItems) > 0 {
					down := sf(ti.TreeItems)
					level = fmt.Sprintf("<li><span class=\"%v\">%v</span><ul class=\"nested\">%v</ul></li>", class_id, ti.Data, down)
				} else {
					level = fmt.Sprintf("<li>%v</li>", ti.Data)
				}
				ss = ss + level
			}
			return ss
		}
		data := sf(l.TreeItems)
		header := fmt.Sprintf("<ul id=\"%v\">%v</ul>", tree_id, data)

		//ol := jso.Document.Object.Call("querySelector", ".custom-select")

		obj_div := doc.Object.Call("getElementById", id)
		obj_div.Set("innerHTML", header)
	}
	return jso
}

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
	//    ID string                  `json:"id"`
	Data string `json:"data"`
}

type TableRowItemO struct {
	//    ID string                          `json:"id"`
	TableItem []TableItemO `json:"table_row_items"`
}

type TableO struct {
	HeaderItems   []TableItemO    `json:"header_items"`
	TableRowItems []TableRowItemO `json:"table_row_items"`
	ChangeCB      string          `json:"change_callback"`
	ClickCB       string          `json:"click_callback"`
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
	// json data
	var dc DocConstructor

	// unmarshall it
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
	//    fmt.Printf("DocConstructorText %v\r\n", DocConstructorText)
	DocConstructor, err := LoadDocConstructorFromString(DocConstructorText)
	if err != nil {
		return nil, nil
	}
	jsoa := []*JSObject{}
	jso_dict := make(map[string]*JSObject)

	//    fmt.Printf("> %v\r\n", DocConstructor)
	for i := range DocConstructor.DocConstructors {
		dci := DocConstructor.DocConstructors[i]
		fmt.Printf("dci %v dci.Object %#v\r\n", dci, dci.Object)
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
			list_id := []string{}
			list_data := []string{}
			for i := range lo.ListItems {
				list_id = append(list_id, lo.ListItems[i].ID)
				list_data = append(list_data, lo.ListItems[i].Data)
			}
			jso = doc.NewList(dci.Object.ObjectID, list_id, list_data)
			l := List{}
			l.ChangeCBName = lo.ChangeCB
			jso.ObjectExtender = l
		case "selector":
			lo := SelectorO{}
			transcode(dci.Object.ObjectExtender, &lo)
			list_id := []string{}
			list_data := []string{}
			for i := range lo.SelectorItems {
				list_id = append(list_id, lo.SelectorItems[i].ID)
				list_data = append(list_data, lo.SelectorItems[i].Data)
			}
			jso = doc.NewSelector(dci.Object.ObjectID, list_id, list_data)
			l := Selector{}
			l.ChangeCBName = lo.ChangeCB
			jso.ObjectExtender = l
		case "table":
			lo := TableO{}
			transcode(dci.Object.ObjectExtender, &lo)
			list_id := []string{}
			list_data := [][]string{}
			for i := range lo.HeaderItems {
				list_id = append(list_id, lo.HeaderItems[i].Data)
			}
			for i := range lo.TableRowItems {
				ld := []string{}
				for j := range lo.TableRowItems {
					ld = append(ld, lo.TableRowItems[i].TableItem[j].Data)
				}
				list_data = append(list_data, ld)
			}
			jso = doc.NewTable(dci.Object.ObjectID, list_id, list_data)
			l := Table{}
			l.ClickCBName = lo.ClickCB
			l.ChangeCBName = lo.ChangeCB
			jso.ObjectExtender = l
		case "tree":
			lo := TreeO{}
			transcode(dci.Object.ObjectExtender, &lo)

			type ConvertFunc func(tree_data []TreeItemO) []TreeData
			var cf ConvertFunc
			//g_num := 1
			cf = func(tree_data []TreeItemO) []TreeData {
				tia := []TreeData{}
				for i := range tree_data {
					ti := TreeData{}
					ti.Data = tree_data[i].Data
					ti.TreeDatas = cf(tree_data[i].TreeItems)
					tia = append(tia, ti)
				}
				return tia
			}
			tree_data := cf(lo.TreeItems)

			list_id := "tree"
			jso = doc.NewTree(dci.Object.ObjectID, list_id, tree_data)
			l := Tree{}
			l.ClickCBName = lo.ClickCB
			l.ChangeCBName = lo.ChangeCB
			jso.ObjectExtender = l
		}
		jsoa = append(jsoa, jso)
		jso_dict[dci.Object.ObjectID] = jso
	}
	return jsoa, jso_dict
}

func BindCallBack(dd interface{}, jsoa []*JSObject, at interface{}) {
	// fmt.Printf("> %v\r\n", dd)
	for i := range jsoa {
		jso := jsoa[i]
		// fmt.Printf("> %v\r\n", jso.ObjectType)
		switch jso.ObjectType {
		case "block":
			b := jso.ObjectExtender.(Block)
			bf, f := Invoke(dd, b.ChangeCBName)
			if bf {
				b.ChangeCB = f
				jso.Object.Call("addEventListener", "change", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
					InvokeCall(b.ChangeCB, this, args, jso, at)
					return js.Null()
				}), false)
			}
			bf, f = Invoke(dd, b.ClickCBName)
			if bf {
				b.ClickCB = f
				jso.Object.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
					InvokeCall(b.ClickCB, this, args, jso, at)
					return js.Null()
				}), false)
			}
		case "text":
			t := jso.ObjectExtender.(Text)
			bf, f := Invoke(dd, t.ChangeCBName)
			if bf {
				t.ChangeCB = f
				jso.Object.Call("addEventListener", "change", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
					// fmt.Printf("!OnChange\r\n")
					InvokeCall(t.ChangeCB, this, args, jso)
					return js.Null()
				}), false)
			}
		case "button":
			b := jso.ObjectExtender.(Button)
			bf, f := Invoke(dd, b.ClickCBName)
			if bf {
				b.ClickCB = f
				jso.Object.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
					InvokeCall(b.ClickCB, this, args, jso, at)
					return js.Null()
				}), false)
			}
		case "image":
			im := jso.ObjectExtender.(Image)
			bf, f := Invoke(dd, im.ClickCBName)
			if bf {
				im.ClickCB = f
				jso.Object.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
					InvokeCall(im.ClickCB, this, args, jso, at)
					return js.Null()
				}), false)
			}
		case "label":
			lb := jso.ObjectExtender.(Label)
			bf, f := Invoke(dd, lb.ClickCBName)
			if bf {
				lb.ClickCB = f
				jso.Object.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
					InvokeCall(lb.ClickCB, this, args, jso, at)
					return js.Null()
				}), false)
			}
		case "list":
			l := jso.ObjectExtender.(List)
			bf, f := Invoke(dd, l.ChangeCBName)
			if bf {
				l.ChangeCB = f
				jso.SetList(f)
			}
		case "selector":
			b := jso.ObjectExtender.(Selector)
			bf, f := Invoke(dd, b.ChangeCBName)
			if bf {
				b.ChangeCB = f
				jso.Object.Call("addEventListener", "change", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
					InvokeCall(b.ChangeCB, this, args, jso)
					return js.Null()
				}), false)
			}
		case "table":
			b := jso.ObjectExtender.(Table)
			bf, f := Invoke(dd, b.ChangeCBName)
			if bf {
				b.ChangeCB = f
				jso.Object.Call("addEventListener", "change", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
					InvokeCall(b.ChangeCB, this, args, jso)
					return js.Null()
				}), false)
			}
			bf, f = Invoke(dd, b.ClickCBName)
			if bf {
				b.ClickCB = f
				jso.Object.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
					InvokeCall(b.ClickCB, this, args, jso)
					return js.Null()
				}), false)
			}
		case "tree":
			fmt.Printf("jso %#v\r\n", jso)
			b := jso.ObjectExtender.(Tree)
			bf, f := Invoke(dd, b.ChangeCBName)
			if bf {
				b.ChangeCB = f
				jso.Object.Call("addEventListener", "change", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
					InvokeCall(b.ChangeCB, this, args, jso)
					return js.Null()
				}), false)
			}
			bf, f = Invoke(dd, b.ClickCBName)
			if bf {
				b.ClickCB = f
				jso.Object.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
					InvokeCall(b.ClickCB, this, args, jso)
					return js.Null()
				}), false)
				class_id := "class_" + jso.ObjectID
				/*
				   var toggler = document.getElementsByClassName("box");
				   var i;

				   for (i = 0; i < toggler.length; i++) {
				     toggler[i].addEventListener("click", function() {
				       this.parentElement.querySelector(".nested").classList.toggle("active");
				       this.classList.toggle("check-box");
				     });
				   }
				*/
				lst := jso.Document.Object.Call("getElementsByClassName", class_id)
				lst_len := lst.Length() //lst.Call("length")

				//lst_array := lst.Interface().([]interface{})
				for i := 0; i < lst_len; i++ {
					fmt.Printf("%v %v\r\n", i, lst.Index(i))
					jso.Object.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
						pe := this.Get("parentElement")
						pe.Call("querySelector", ".nested").Get("classList").Call("toggle", "active")
						this.Get("classList").Call("toggle", "check-box")
						return js.Null()
					}), false)
				}
				/*
					type ConvertFunc func(tree_data []TreeItemO) []TreeData
					var cf ConvertFunc
					//g_num := 1
					cf = func(tree_data []TreeItemO) []TreeData {
						//tia := []TreeData{}
						for i, _ := range tree_data {
							ti := TreeData{}
							ti.Data = tree_data[i].Data
							ti.TreeDatas = cf(tree_data[i].TreeItems)
							tia = append(tia, ti)
						}
						return tia
					}
					tree_data := cf(b.TreeItems)
				*/
			}
		}
	}
}

func Invoke(any interface{}, name string) (bool, reflect.Value) {
	v := reflect.ValueOf(any)
	//    fmt.Printf("%#v\r\n", v)
	//    fmt.Printf("name %v\r\n", name)
	f := v.MethodByName(name)
	// fmt.Printf("f %#v name %v\r\n", f, name)
	if f.IsValid() {
		return true, f
	}
	return false, f
}

func InvokeCall(f reflect.Value, args ...interface{}) {
	inputs := make([]reflect.Value, len(args))
	for i := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	if f.IsValid() {
		f.Call(inputs)
	}
}
