package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
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
var presets map[uint8]Preset

func main() {
	err := LoadConfiguration(&config)
	if err != nil {
		panic("Cannot load config")
	}
	if loadPresets() != nil {
		panic("Cannot load presets from config")
	}
	getIndexHTML(os.Stdout)

	defer log.Println("Exit.")

	r := mux.NewRouter()
	r.HandleFunc("/preset/{presetNum}", presetHandler).Methods(http.MethodGet)
	r.HandleFunc("/source/{sourceName}", sourceHandler).Methods(http.MethodGet)
	r.HandleFunc("/pureDirect/{onOrOff}", pureDirectHandler).Methods(http.MethodGet)
	r.HandleFunc("/exit", exitHandler).Methods(http.MethodGet)
	r.HandleFunc("/", indexHandler).Methods(http.MethodGet)           // index from template
	r.HandleFunc("/index.html", indexHandler).Methods(http.MethodGet) // index from template
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./")))        // rest from fs

	r.Use(mux.CORSMethodMiddleware(r))

	srvAddr := "0.0.0.0:" + strconv.Itoa(int(config.Server.ListenPort))
	log.Println("Starging local Web Listener on Addr: " + srvAddr)
	srv := &http.Server{
		Handler: r,
		Addr:    srvAddr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Println("Using Yamaha RX-V771 @ " + config.YamahaReceiverHost)
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

func LoadConfiguration(config *Config) error {
	configFile, err := os.Open(ConfigFileName)
	defer configFile.Close()
	if err != nil {
		return err
	}
	jsonParser := json.NewDecoder(configFile)
	return jsonParser.Decode(config)
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

func loadPresets() error {
	presets = make(map[uint8]Preset)
	for _, preset := range config.Presets {
		if presets[preset.Index].Index > 0 {
			log.Println("Duplicate entry for Index '" + strconv.Itoa(int(preset.Index)) + "'. Overwriting it with current one.")
		}
		presets[preset.Index] = preset
	}

	return nil
}

// =============================================================================================

func getIndexHTML(wr io.Writer) {
	indexTemplate, err := template.ParseFiles("index.html.template")
	if err != nil {
		log.Println(err)
		return
	}

	// Execute needs some sort of io.Writer
	//err = indexTemplate.Execute(wr, config)
	err = indexTemplate.Execute(wr, presets)
	if err != nil {
		log.Println(err)
		return
	}
}

// =============================================================================================

func exitHandler(w http.ResponseWriter, r *http.Request) {
	//defer signal.Notify(shutdownSig, os.Interrupt)	// doesnt work
	w.Header().Set("Access-Control-Allow-Origin", "*")

	log.Println("Exit requested")
	w.Write([]byte("OK"))

	shutdownSig <- os.Interrupt
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	getIndexHTML(w)
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

func sourceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sourceName := vars["sourceName"]
	log.Println("### Source: " + sourceName)

	w.Header().Set("Access-Control-Allow-Origin", "*")

	var lala = bytes.Buffer{}
	lala.WriteString("<YAMAHA_AV cmd=\"PUT\"><Main_Zone><Input><Input_Sel>" + sourceName + "</Input_Sel></Input></Main_Zone></YAMAHA_AV>")

	_, err := http.Post("http://"+config.YamahaReceiverHost+"/YamahaRemoteControl/ctrl", "text/xml", &lala)
	if err != nil {
		log.Println("FAIL")
	} else {
		log.Println("OK")
	}

	w.Write([]byte("OK"))
}

func pureDirectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	onOrOff := vars["onOrOff"]
	log.Println("### State: " + onOrOff)
	if onOrOff != "On" && onOrOff != "Off" {
		log.Println("### Unknown state!")
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	var lala = bytes.Buffer{}
	lala.WriteString("<YAMAHA_AV cmd=\"PUT\"><Main_Zone><Sound_Video><Pure_Direct><Mode>" + onOrOff + "</Mode></Pure_Direct></Sound_Video></Main_Zone></YAMAHA_AV>")
	_, err := http.Post("http://"+config.YamahaReceiverHost+"/YamahaRemoteControl/ctrl", "text/xml", &lala)
	if err != nil {
		log.Println("FAIL")
	} else {
		log.Println("OK")
	}

	w.Write([]byte("OK"))
}
