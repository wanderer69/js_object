package object

import (
	"fmt"
	"reflect"
	"syscall/js"
)

type JSObject struct {
	Object         js.Value    `json:"object"`
	ObjectData     js.Value    `json:"object_data"` // object for linked data
	ObjectID       string      `json:"object_id"`
	ObjectType     string      `json:"object_type"`
	Document       *DocObject  `json:"document"`
	ObjectExtender interface{} `json:"object_extender"`
	Childs         []*JSObject `json:"childs_list"`
	Parent         *JSObject   `json:"parent_object"`
	Header         *JSObject   `json:"header_object"`
}

func GetObjectByName(doc *DocObject, id string) js.Value {
	obj := doc.Object.Call("getElementById", id)
	if obj.IsNull() {
		objA := doc.Object.Call("getElementsByClassName", id)
		if objA.IsNull() {
			return objA
		}
		objL := DOMObjectToList(objA)
		obj = objL[0]
	}
	return obj
}

func NewJSObject(doc *DocObject, id string, oType string, oe interface{}) *JSObject {
	obj := GetObjectByName(doc, id)
	childs := []*JSObject{}
	parent := &JSObject{}

	jso := JSObject{
		Object:         obj,
		ObjectData:     obj,
		ObjectID:       id,
		ObjectType:     oType,
		Document:       doc,
		ObjectExtender: oe,
		Childs:         childs,
		Parent:         parent,
	}
	return &jso
}

type ObjCallBack func(this js.Value, args []js.Value) interface{}

func (jso *JSObject) SetCallBack(name string, ocb ObjCallBack) {
	jso.Object.Call("addEventListener", name, js.FuncOf(ocb), false)
}

func (jso *JSObject) AddChild(jsoC *JSObject) {
	jso.Childs = append(jso.Childs, jsoC)
}

func (jso *JSObject) SetParent(jsoP *JSObject) {
	jso.Parent = jsoP
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

func (jso *JSObject) CorrectText(baseClass string) *JSObject {
	doc := jso.Document
	switch jso.ObjectType {
	case "text":
		if doc.Type == "figma" {
			txt := `<input id="%v_id" type="text" class="%v">`
			ss := fmt.Sprintf(txt, baseClass, baseClass)
			jso.Object.Set("innerHTML", ss)
			jso.ObjectData = GetObjectByName(doc, fmt.Sprintf("%v_id", baseClass))
		}
	case "label":
		if doc.Type == "figma" {
		}
	}

	return jso
}

func (jso *JSObject) SetLabelText(value string) {
	switch jso.ObjectType {
	case "label":
		jso.ObjectData.Set("textContent", value)
	}
}

func DOMObjectToList(o js.Value) []js.Value {
	var out []js.Value
	length := o.Get("length").Int()
	for i := 0; i < length; i++ {
		out = append(out, o.Call("item", i))
	}
	return out
}

func (jso *JSObject) SetList(cb reflect.Value) *JSObject {
	// cb reflect.Value
	l := jso.ObjectExtender.(List)
	ol := jso.Document.Object.Call("querySelector", ".custom-select")
	cl := ol.Get("classList")
	cl.Call("toggle", "open")

	poslValue := jso.Document.Object.Call("querySelectorAll", ".custom-option")
	posL := DOMObjectToList(poslValue)
	for i := range posL {
		spanItem := posL[i]
		spanItem.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			pn := this.Get("parentNode")
			ol := pn.Call("querySelector", ".custom-option.selected")
			if !ol.IsNull() {
				ol.Get("classList").Call("remove", "selected")
			}
			this.Get("classList").Call("add", "selected")
			val := this.Get("textContent").String()
			li := l.ListIDDict[val]
			l.CurrentListItem = &li
			if li.OnClick != nil {
				li.OnClick(this, args, l)
			}
			fmt.Printf("--> %v this %v args %v\r\n", cb, this, args)
			v := js.ValueOf(i)
			args = append(args, v)
			InvokeCall(cb, this, args, jso)
			return js.Null()
		}), false)
	}

	return jso
}

func (jso *JSObject) UpdateList(dd interface{}, listID []string, listData []string) *JSObject {
	// cb reflect.Value
	l := jso.ObjectExtender.(List)

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

	bf, f := Invoke(dd, l.ChangeCBName)
	if bf {
		l.ChangeCB = f
		jso.SetList(f)
	}

	return jso
}

func (jso *JSObject) AddOption(selectorID []string) {
	switch jso.ObjectType {
	case "selector":
		for i, item := range selectorID {
			el := jso.Document.Object.Call("createElement", "option")
			el.Set("textContent", item)
			el.Set("value", fmt.Sprintf("%v", i+1))
			//              el.Set("selected","true")
			jso.Object.Call("appendChild", el)
		}
	}
}

func (jso *JSObject) SetTable(dd interface{}, tableID []string, tableData [][]string) {
	switch jso.ObjectType {
	case "table":
		t := jso.ObjectExtender.(Table)

		header := jso.Object.Call("createTHead")
		// Create an empty <tr> element and add it to the first position of <thead>:
		row := header.Call("insertRow", 0)
		for i, item := range tableID {
			// Insert a new cell (<td>) at the first position of the "new" <tr> element:
			cell := row.Call("insertCell", i)
			cell.Set("innerHTML", item)
		}
		rowID := t.ID + "_header"
		row.Set("id", rowID)
		rowObj := NewJSObject(jso.Document, rowID, "table_header", row)
		rowObj.Parent = jso
		jso.Header = rowObj
		body := jso.Object.Call("createTBody")
		for i := range tableData {
			td := tableData[i]
			row := body.Call("insertRow", i)
			rowID := fmt.Sprintf("%v_row_%v", t.ID, i)
			row.Set("id", rowID)
			for j := range td {
				cell := row.Call("insertCell", j)
				cell.Set("innerHTML", td[j])
			}
			rowObj := NewJSObject(jso.Document, rowID, "table_row", row)
			rowObj.Parent = jso
			jso.Childs = append(jso.Childs, rowObj)
		}

		bf, f := Invoke(dd, t.ChangeCBName)
		if bf {
			t.ChangeCB = f
			jso.Object.Call("addEventListener", "change", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				InvokeCall(t.ChangeCB, this, args, jso)
				return js.Null()
			}), false)
		}
		bf, f = Invoke(dd, t.ClickCBName)
		if bf {
			t.ClickCB = f
			jso.Object.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				InvokeCall(t.ClickCB, this, args, jso)
				return js.Null()
			}), false)
		}

		bf, f = Invoke(dd, t.RowClickCBName)
		for i := range jso.Childs {
			if bf {
				t.RowClickCB = f
				rowObj := jso.Childs[i]
				rowObj.Object.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
					args = append(args, rowObj.ObjectExtender.(js.Value))
					v := js.ValueOf(i)
					args = append(args, v)
					InvokeCall(t.RowClickCB, this, args, rowObj)
					return js.Null()
				}), false)
			}
		}
	}
}

func BindCallBack(dd interface{}, jsoa []*JSObject, at interface{}) {
	for i := range jsoa {
		jso := jsoa[i]
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
					InvokeCall(b.ClickCB, this, args, jso /*, at*/)
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

			bf, f = Invoke(dd, b.RowClickCBName)
			for i := range b.TableRowItems {
				if bf {
					b.RowClickCB = f
					rowObj := jso.Childs[i]
					rowObj.Object.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
						InvokeCall(b.RowClickCB, this, args, rowObj)
						return js.Null()
					}), false)
				}
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
				lst := jso.Document.Object.Call("getElementsByClassName", class_id)
				lst_len := lst.Length()

				for i := 0; i < lst_len; i++ {
					//fmt.Printf("%v %v\r\n", i, lst.Index(i))
					jso.Object.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
						pe := this.Get("parentElement")
						pe.Call("querySelector", ".nested").Get("classList").Call("toggle", "active")
						this.Get("classList").Call("toggle", "check-box")
						return js.Null()
					}), false)
				}
			}
		}
	}
}

func Invoke(any interface{}, name string) (bool, reflect.Value) {
	v := reflect.ValueOf(any)
	f := v.MethodByName(name)
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
