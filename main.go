package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

const ConfigFileName = "config.json"

type Preset struct {
	Index   uint8
	Name    string
	IconURL string
}

type Config struct {
	YamahaReceiverHost string `json:"yamahaReceiverHost"` // IP or FQDN
	Presets            []Preset
	Server             struct {
		ListenPort uint16
	}
}

var shutdownSig chan os.Signal
var config Config

func main() {
	LoadConfiguration(&config)
	defer log.Println("Exit.")

	r := mux.NewRouter()
	r.HandleFunc("/preset/{presetNum}", presetHandler).Methods(http.MethodGet)
	r.HandleFunc("/exit", exitHandler).Methods(http.MethodGet)
	r.HandleFunc("/pureDirectOff", pureDirectOffHandler).Methods(http.MethodGet)
	r.HandleFunc("/pureDirectOn", pureDirectOnHandler).Methods(http.MethodGet)
	r.HandleFunc("/AUDIO2", audio2Handler).Methods(http.MethodGet)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./")))

	r.Use(mux.CORSMethodMiddleware(r))

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:9000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Println("Starting Server using Yamaha RX-V771 @ " + config.YamahaReceiverHost)
	go srv.ListenAndServe()

	shutdownSig = make(chan os.Signal)
	signal.Notify(shutdownSig, os.Interrupt)
	signal.Notify(shutdownSig, syscall.SIGKILL)

	// wait for our death
	<-shutdownSig
	// and persist our state
	PersistConfiguration(&config)

}

// =============================================================================================

func LoadConfiguration(config *Config) {
	configFile, err := os.Open(ConfigFileName)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(config)
}

func PersistConfiguration(config *Config) {
	configFile, err := os.OpenFile(ConfigFileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer configFile.Close()

	err = os.Truncate(ConfigFileName, 0)
	if err != nil {
		log.Fatal(err)
		return
	}

	_buf, err := json.MarshalIndent(config, "", "    ")
	if err == nil {
		_, err = configFile.Write(_buf)
		if err == nil {
			log.Println("Config written.")
			return
		}
	}
	log.Fatal("Couldn't write config:")
	log.Fatal(err)
}

func isFile(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()

}

// =============================================================================================

func exitHandler(w http.ResponseWriter, r *http.Request) {
	//defer signal.Notify(shutdownSig, os.Interrupt)	// doesnt work
	w.Header().Set("Access-Control-Allow-Origin", "*")

	log.Println("Exit requested")
	w.Write([]byte("OK"))

	shutdownSig <- os.Interrupt
}

func presetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	presetNum := vars["presetNum"]
	log.Println("### Preset: " + presetNum)

	w.Header().Set("Access-Control-Allow-Origin", "*")

	var lala = bytes.Buffer{}
	lala.WriteString("<YAMAHA_AV cmd=\"PUT\"><NET_RADIO><Play_Control><Preset><Preset_Sel>" + presetNum + "</Preset_Sel></Preset></Play_Control></NET_RADIO></YAMAHA_AV>")

	_, err := http.Post("http://"+config.YamahaReceiverHost+"/YamahaRemoteControl/ctrl", "text/xml", &lala)
	if err != nil {
		log.Println("FAIL")
	} else {
		log.Println("OK")
	}

	w.Write([]byte("OK"))
}

func pureDirectOnHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var lala = bytes.Buffer{}
	lala.WriteString("<YAMAHA_AV cmd=\"PUT\"><Main_Zone><Sound_Video><Pure_Direct><Mode>On</Mode></Pure_Direct></Sound_Video></Main_Zone></YAMAHA_AV>")
	_, err := http.Post("http://"+config.YamahaReceiverHost+"/YamahaRemoteControl/ctrl", "text/xml", &lala)
	if err != nil {
		log.Println("FAIL")
	} else {
		log.Println("OK")
	}

	w.Write([]byte("OK"))
}
func pureDirectOffHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var lala = bytes.Buffer{}
	//lala.WriteString("<YAMAHA_AV cmd=\"PUT\"><System><Pure_Direct><Mode>Off</Mode></Pure_Direct></System></YAMAHA_AV>")
	//lala.WriteString("<YAMAHA_AV cmd=\"PUT\"><Main_Zone><Pure_Direct><Mode>Off</Mode></Pure_Direct></Main_Zone></YAMAHA_AV>")
	//lala.WriteString("<YAMAHA_AV cmd=\"PUT\"><Main><Pure_Direct><Mode>Off</Mode></Pure_Direct></Main></YAMAHA_AV>")
	lala.WriteString("<YAMAHA_AV cmd=\"PUT\"><Main_Zone><Sound_Video><Pure_Direct><Mode>Off</Mode></Pure_Direct></Sound_Video></Main_Zone></YAMAHA_AV>")

	_, err := http.Post("http://"+config.YamahaReceiverHost+"/YamahaRemoteControl/ctrl", "text/xml", &lala)
	if err != nil {
		log.Println("FAIL")
	} else {
		log.Println("OK")
	}

	w.Write([]byte("OK"))
}
func audio2Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var lala = bytes.Buffer{}
	//lala.WriteString("<YAMAHA_AV cmd=\"PUT\"><System><Pure_Direct><Mode>Off</Mode></Pure_Direct></System></YAMAHA_AV>")
	//lala.WriteString("<YAMAHA_AV cmd=\"PUT\"><Main_Zone><Pure_Direct><Mode>Off</Mode></Pure_Direct></Main_Zone></YAMAHA_AV>")
	//lala.WriteString("<YAMAHA_AV cmd=\"PUT\"><Main><Pure_Direct><Mode>Off</Mode></Pure_Direct></Main></YAMAHA_AV>")
	lala.WriteString("<YAMAHA_AV cmd=\"PUT\"><Main_Zone><Input><Input_Sel>AUDIO2</Input_Sel></Input></Main_Zone></YAMAHA_AV>")

	_, err := http.Post("http://"+config.YamahaReceiverHost+"/YamahaRemoteControl/ctrl", "text/xml", &lala)
	if err != nil {
		log.Println("FAIL")
	} else {
		log.Println("OK")
	}

	w.Write([]byte("OK"))
}
func NETRADIOHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var lala = bytes.Buffer{}
	//lala.WriteString("<YAMAHA_AV cmd=\"PUT\"><System><Pure_Direct><Mode>Off</Mode></Pure_Direct></System></YAMAHA_AV>")
	//lala.WriteString("<YAMAHA_AV cmd=\"PUT\"><Main_Zone><Pure_Direct><Mode>Off</Mode></Pure_Direct></Main_Zone></YAMAHA_AV>")
	//lala.WriteString("<YAMAHA_AV cmd=\"PUT\"><Main><Pure_Direct><Mode>Off</Mode></Pure_Direct></Main></YAMAHA_AV>")
	lala.WriteString("<YAMAHA_AV cmd=\"PUT\"><Main_Zone><Input><Input_Sel>NET RADIO</Input_Sel></Input></Main_Zone></YAMAHA_AV>")

	_, err := http.Post("http://"+config.YamahaReceiverHost+"/YamahaRemoteControl/ctrl", "text/xml", &lala)
	if err != nil {
		log.Println("FAIL")
	} else {
		log.Println("OK")
	}

	w.Write([]byte("OK"))
}
