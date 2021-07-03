package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strconv"
)

const (
	Frequency = 302941000

	FanId = "1011"

	FanLow      = "0010000"
	FanMedium   = "0100000"
	FanHigh     = "1000000"
	FanOff      = "0000100"
	LightToggle = "0000010"
)

var cmdMap = map[string]string{"LOW": FanLow, "MED": FanMedium, "HIGH": FanHigh, "POWER": FanOff, "LIGHT": LightToggle}

func encode(cmdBitString string) string {
	str := ""
	for _, bit := range cmdBitString {
		str += "10" + string(bit)
	}
	return str
}

func makeCmdBitString(fanID string, cmd string) string {
	return encode(fmt.Sprintf("1%s1%s0", fanID, cmd))
}

func cmdHandler(writer http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		contents, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Fatal(err)
		}
		fanCmd := cmdMap[string(contents)]
		cmdBitString := makeCmdBitString(FanId, fanCmd)
		//sendook -f 302941000 -0 333 -1 333 -r 5 -p 10000 1011011001011011011001001011001001001
		cmd := exec.Command("/root/rpitx/sendook", "-f", strconv.Itoa(Frequency), "-0", "333", "-1", "333", "-r", "5", "-p", "10000", cmdBitString)
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
		writer.Write([]byte("success"))
	} else {
		writer.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func indexHandler(writer http.ResponseWriter, _ *http.Request) {
	tmpl := template.Must(template.ParseFiles("./index.html"))
	tmpl.Execute(writer, nil)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/cmd", cmdHandler)
	if err := http.ListenAndServe("0.0.0.0:80", mux); err != nil {
		log.Fatalln(err)
	}
}
