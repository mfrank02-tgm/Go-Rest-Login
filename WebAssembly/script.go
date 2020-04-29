package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"syscall/js"
)

type Person struct {
	ID       string
	Username string
	Password string
}

func main() {
	println("Hello")
	go func() {
		time.Sleep(time.Second * 60)
		//js.Global.Get("document").Call("write", "Hello world!")
		//js.Global.Get("location").Call("reload")
	}()
	println("World")

	url := "./users"
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		print(err.Error())
	}
	resp, err := client.Do(req)
	if err != nil {
		print(err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err.Error())
	}
	//println(body)

	erg := make([]Person, 0)
	json.Unmarshal(body, &erg)

	js.Global().Get("document").Call("write", "<table style='border: 1px solid black; width:100%; text-align: center'><thead><tr><th>ID</th><th>Username</th><th>Password</th></tr></thead><tbody>")

	for _, element := range erg {
		js.Global().Get("document").Call("write", "<tr><td>"+element.ID+"</td><td>"+element.Username+"</td><td>"+element.Password+"</td></tr>")
	}

	js.Global().Get("document").Call("write", "</tbody></table>")

	//println(erg)
}
