package object

import (
	"fmt"
	"strings"
	"syscall/js"
)

type DocObject struct {
	Object js.Value `json:"object"`
	Type   string   `json:"type_html_object"` // normal, figma
}

func NewDocument() *DocObject {
	doc := js.Global().Get("document")
	do := DocObject{}
	do.Object = doc
	return &do
}

func (doc *DocObject) NewRow(id string, parent *JSObject, r *TableRowItem) *JSObject {
	jso := NewJSObject(doc, id, "row", r)
	jso.Parent = parent
	return jso
}

func (doc *DocObject) NewText(id string, textByDefault string) *JSObject {
	t := Text{}
	t.TextByDefault = textByDefault
	jso := NewJSObject(doc, id, "text", t)
	if len(t.TextByDefault) > 0 {
		jso.SetValue(t.TextByDefault)
	}
	return jso
}

func (doc *DocObject) NewLabel(id string, textByDefault string) *JSObject {
	lb := Label{}
	lb.TextByDefault = textByDefault
	jso := NewJSObject(doc, id, "label", lb)
	if len(lb.TextByDefault) > 0 {
		jso.SetValue(lb.TextByDefault)
	}
	return jso
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

func (doc *DocObject) NewButton(id string, buttonName string) *JSObject {
	b := Button{}
	b.ButtonName = buttonName
	jso := NewJSObject(doc, id, "button", b)
	if len(b.ButtonName) > 0 {
		jso.SetValue(b.ButtonName)
	}
	return jso
}

func (doc *DocObject) NewList(id string, listID []string, listData []string) *JSObject {
	if len(listID) != len(listData) {
		return nil
	}
	l := List{}
	l.ListIDDict = make(map[string]ListItem)
	tags_lst := []string{}
	for i, item := range listID {
		li := ListItem{}
		li.ID = item
		li.Data = listData[i]
		tags_lst = append(tags_lst, fmt.Sprintf("<span class=\"custom-option\">%v</span><br>", item))
		l.ListItems = append(l.ListItems, li)
		l.ListIDDict[item] = li
	}
	ss := strings.Join(tags_lst, "\r\n")
	jso := NewJSObject(doc, id, "list", l)
	jso.Object.Set("innerHTML", ss)

	return jso
}

func (doc *DocObject) NewSelector(id string, selectorID []string, selectorData []string) *JSObject {
	if len(selectorID) != len(selectorData) {
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
	for i, item := range selectorID {
		li := SelectorItem{}
		li.ID = item
		li.Data = selectorData[i]
		l.SelectorItems = append(l.SelectorItems, li)
		l.SelectorIDDict[item] = li
	}
	jso := NewJSObject(doc, id, "selector", l)
	for i, item := range selectorID {
		el := doc.Object.Call("createElement", "option")
		el.Set("textContent", item)
		el.Set("value", fmt.Sprintf("%v", i+1))
		jso.Object.Call("appendChild", el)
	}
	return jso
}

func (doc *DocObject) NewSelector1(id string, selectorID []string, selectorData []string) *JSObject {
	if len(selectorID) != len(selectorData) {
		return nil
	}
	l := Selector{}
	l.SelectorIDDict = make(map[string]SelectorItem)
	tagsLst := []string{}
	tagsLst = append(tagsLst, "<select id=\""+id+"\">")
	fmt.Printf("selector_id %v\r\n", selectorID)
	for i, item := range selectorID {
		li := SelectorItem{}
		li.ID = item
		li.Data = selectorData[i]
		ss := fmt.Sprintf("<option value=\"%v\">%v</option>", i+1, item)
		tagsLst = append(tagsLst, ss)
		l.SelectorItems = append(l.SelectorItems, li)
		l.SelectorIDDict[item] = li
	}
	tagsLst = append(tagsLst, "</select>")
	fmt.Printf("selector id %v\r\n", id)
	ss := strings.Join(tagsLst, "\r\n")
	divName := fmt.Sprintf("%v_", id)
	objDiv := doc.Object.Call("getElementById", divName)
	objDiv.Set("innerHTML", ss)
	jso := NewJSObject(doc, id, "selector", l)
	return jso
}

func (doc *DocObject) NewTable(id string, tableID []string, tableData [][]string) *JSObject {
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
	l.ID = id
	l.TableIDDict = make(map[string]TableRowItem)
	for i, item := range tableData {
		li := TableRowItem{}
		li.ID = fmt.Sprintf("%v", i)
		td := item
		for j := range td {
			lii := TableItem{}
			lii.Data = td[j]
			lii.ID = tableID[j]
			li.TableItems = append(li.TableItems, lii)
			l.TableIDDict[li.ID] = li
		}
	}

	jso := NewJSObject(doc, id, "table", l)
	if len(tableID) > 0 {
		header := jso.Object.Call("createTHead")
		// Create an empty <tr> element and add it to the first position of <thead>:
		row := header.Call("insertRow", 0)
		for i, item := range tableID {
			// Insert a new cell (<td>) at the first position of the "new" <tr> element:
			cell := row.Call("insertCell", i)
			cell.Set("innerHTML", item)
		}
		rowID := id + "_header"
		row.Set("id", rowID)
		value := "/content-page"
		row.Set("data-href", value)
		rowObj := NewJSObject(doc, rowID, "table_row", row)
		rowObj.Parent = jso
		jso.Header = rowObj
		if len(tableData) > 0 {
			body := jso.Object.Call("createTBody")
			for i := range tableData {
				td := tableData[i]
				row := body.Call("insertRow", i)
				for j := range td {
					cell := row.Call("insertCell", j)
					cell.Set("innerHTML", td[j])
				}
				rowID := fmt.Sprintf("%v_row_%v", id, i)
				row.Set("id", rowID)
				rowObj := NewJSObject(doc, rowID, "table_row", row)
				rowObj.Parent = jso
				jso.Childs = append(jso.Childs, rowObj)
			}
		}
	}
	return jso
}

func (doc *DocObject) NewBlock(id string) *JSObject {
	b := Block{}
	jso := NewJSObject(doc, id, "block", b)
	return jso
}

func (doc *DocObject) NewTree(id string, treeID string, treeData []TreeData) *JSObject {
	l := Tree{}
	type ConvertFunc func(treeDataIn []TreeData) []TreeItem
	var cf ConvertFunc
	gNum := 1
	cf = func(treeDataIn []TreeData) []TreeItem {
		tia := []TreeItem{}
		for i := range treeDataIn {
			ti := TreeItem{}
			ti.ID = fmt.Sprintf("%v", gNum)
			gNum = gNum + 1
			ti.Data = treeDataIn[i].Data
			ti.TreeItems = cf(treeDataIn[i].TreeDatas)
			tia = append(tia, ti)
		}
		return tia
	}
	l.TreeItems = cf(treeData)

	jso := NewJSObject(doc, id, "tree", l)
	if len(treeData) > 0 {
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
		classID := "class_" + id
		type SetFunc func(treeItems []TreeItem) string
		var sf SetFunc
		sf = func(treeItems []TreeItem) string {
			ss := ""
			for i := range treeItems {
				ti := treeItems[i]
				level := ""
				if len(ti.TreeItems) > 0 {
					down := sf(ti.TreeItems)
					level = fmt.Sprintf("<li><span class=\"%v\">%v</span><ul class=\"nested\">%v</ul></li>", classID, ti.Data, down)
				} else {
					level = fmt.Sprintf("<li>%v</li>", ti.Data)
				}
				ss = ss + level
			}
			return ss
		}
		data := sf(l.TreeItems)
		header := fmt.Sprintf("<ul id=\"%v\">%v</ul>", treeID, data)
		objDiv := doc.Object.Call("getElementById", id)
		objDiv.Set("innerHTML", header)
	}
	return jso
}
