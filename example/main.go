package main

import (
	"fmt"
	"syscall/js"

	object "github.com/wanderer69/js_object"
)

type AppT struct {
	Jsoa []*object.JSObject
	Jsod map[string]*object.JSObject
}

func (at AppT) BlockClickCallBack(this js.Value, args []js.Value, jso *object.JSObject) interface{} {
	fmt.Printf("BlockClickCallBack\r\n")
	return nil
}

func (at AppT) BlockChangeCallBack(this js.Value, args []js.Value, jso *object.JSObject) interface{} {
	fmt.Printf("BlockChangeCallBack\r\n")
	return nil
}

func (at AppT) TextChangeCallBack(this js.Value, args []js.Value, jso *object.JSObject) interface{} {
	fmt.Printf("TextChangeCallBack\r\n")
	fmt.Printf("%#v\r\n", args)
	return nil
}

func (at AppT) ListChangeCallBack(this js.Value, args []js.Value, jso *object.JSObject) interface{} {
	//do something
	fmt.Printf("ListChangeCallBack\r\n")
	fmt.Printf("%#v\r\n", this)
	for i := range args {
		fmt.Printf("--> %v\r\n", args[i].Type().String())
		switch args[i].Type().String() {
		case "string":
			fmt.Printf("%v\r\n", args[i].String())
		case "number":
			fmt.Printf("%v\r\n", args[i].Int())
		case "float":
			fmt.Printf("%v\r\n", args[i].Float())
		case "bool":
			fmt.Printf("%v\r\n", args[i].Bool())
		}
	}
	return nil
}

func (at AppT) SelectorChangeCallBack(this js.Value, args []js.Value, jso *object.JSObject) interface{} {
	fmt.Printf("SelectorChangeCallBack\r\n")
	fmt.Printf("%#v\r\n", args)
	return nil
}

func (at AppT) ButtonClickCallBack(this js.Value, args []js.Value, jso *object.JSObject) interface{} {
	fmt.Printf("ButtonClickCallBack\r\n")
	text1 := at.Jsod["text_1"]
	if text1 != nil {
		v := text1.GetValue()
		if v == "33333" {
			text1.SetValue("22222")
		} else {
			text1.SetValue("33333")
		}
	}
	return nil
}

func (at AppT) LabelClickCallBack(this js.Value, args []js.Value, jso *object.JSObject) interface{} {
	fmt.Printf("LabelClickCallBack\r\n")
	return nil
}

func (at AppT) ListClickCallBack(this js.Value, args []js.Value, jso *object.JSObject) interface{} {
	fmt.Printf("ListClickCallBack\r\n")
	fmt.Printf("%#v\r\n", this)
	for i := range args {
		fmt.Printf("-- %#v\r\n", args[i])
	}
	return nil
}

func (at AppT) TableClickCallBack(this js.Value, args []js.Value, jso *object.JSObject) interface{} {
	fmt.Printf("TableClickCallBack\r\n")
	fmt.Printf("%#v\r\n", args)
	table1 := at.Jsod["table_1"]
	if table1 != nil {
	}
	return nil
}

func (at AppT) TableRowClickCallBack(this js.Value, args []js.Value, jso *object.JSObject) interface{} {
	fmt.Printf("TableRowClickCallBack\r\n")
	for i := range args {
		fmt.Printf("-- %v\r\n", args[i].Type().String())
		switch args[i].Type().String() {
		case "string":
			fmt.Printf("%v\r\n", args[i].String())
		case "number":
			fmt.Printf("%v\r\n", args[i].Int())
		case "float":
			fmt.Printf("%v\r\n", args[i].Float())
		case "bool":
			fmt.Printf("%v\r\n", args[i].Bool())
		}
	}
	table1 := at.Jsod["table_1"]
	if table1 != nil {
		fmt.Printf("table1 %v\r\n", table1)
	}
	return nil
}

func (at AppT) TableChangeCallBack(this js.Value, args []js.Value, jso *object.JSObject) interface{} {
	fmt.Printf("TableChangeCallBack\r\n")
	fmt.Printf("%#v\r\n", args)
	return nil
}

func (at AppT) TreeClickCallBack(this js.Value, args []js.Value, jso *object.JSObject) interface{} {
	fmt.Printf("TreeClickCallBack\r\n")
	fmt.Printf("%#v\r\n", args)
	return nil
}

func (at AppT) TreeChangeCallBack(this js.Value, args []js.Value, jso *object.JSObject) interface{} {
	fmt.Printf("TreeChangeCallBack\r\n")
	fmt.Printf("%#v\r\n", args)
	return nil
}

func main() {
	fmt.Println("Test 1")
	done := make(chan bool)
	jsoa, jsod := object.CreateDocConstructor()
	at := AppT{}
	at.Jsoa = jsoa
	at.Jsod = jsod
	fmt.Printf("jsoa %v jsod %v\r\n", jsoa, jsod)

	doc := js.Global().Get("document")
	doc.Call("querySelectorAll", ".table-clicable tbody tr")
	object.BindCallBack(at, jsoa, &at)

	block_1 := jsod["block_1"]
	text_1 := jsod["text_1"]
	button_1 := jsod["button_1"]
	block_1.SetValue("11111")
	label_1 := jsod["label_1"]
	label_1.SetLabelText("This label text")
	/*
		list_1 := jsod["list_1"]
		id := []string{"Col 1", "Col 2", "Col 3", "Col 4"}
		data := []string{"Item10", "Item20", "Item30", "Item40"}
		list_1.UpdateList(at, id, data)
	*/

	selector_1 := jsod["selector_1"]
	ol := []string{"Item1", "Item2", "Item3", "Item4"}
	selector_1.AddOption(ol)

	table_1 := jsod["table_1"]

	hl := []string{"Col 1", "Col 2", "Col 3", "Col 4"}
	ol1 := []string{"Item10", "Item20", "Item30", "Item40"}
	ol2 := []string{"Item11", "Item21", "Item31", "Item41"}
	oll := [][]string{ol1, ol2}
	table_1.SetTable(at, hl, oll)

	bb_cb := func(this js.Value, args []js.Value) interface{} {
		text_1.SetValue("333333")
		return nil
	}
	button_1.SetCallBack("click", bb_cb)

	<-done
}
