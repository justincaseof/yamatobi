package main

import (
    "bytes"
    "log"
	"time"
    "net/http"
	"github.com/gorilla/mux"
)

func main() {
	presetHandler := func(w http.ResponseWriter, r *http.Request) {
		log.Println("###")
		
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == http.MethodOptions {
			return
		}
		
		var lala = bytes.Buffer{}
		lala.WriteString(`<YAMAHA_AV cmd="PUT"><NET_RADIO><Play_Control><Preset><Preset_Sel>1</Preset_Sel></Preset></Play_Control></NET_RADIO></YAMAHA_AV>`)
		
		_, err := http.Post("http://192.168.178.7/YamahaRemoteControl/ctrl", "text/xml", &lala)
		if err != nil {
			log.Println("FAIL")
		} else {
			log.Println("OK")
		}
		
		
		w.Write([]byte("OK"))
	}

    r := mux.NewRouter()

    // This will serve files under http://localhost:9000/static/<filename>
    r.HandleFunc("/foo", presetHandler).Methods(http.MethodGet)
	r.PathPrefix("/").Handler( http.FileServer( http.Dir("./") ))
	r.HandleFunc("/preset", presetHandler)
	
	r.Use(mux.CORSMethodMiddleware(r))
    
	srv := &http.Server{
        Handler:      r,
        Addr:         "127.0.0.1:9000",
        // Good practice: enforce timeouts for servers you create!
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }

    log.Fatal(srv.ListenAndServe())
}