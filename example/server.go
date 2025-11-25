package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	dasblinken "barnstar.com/dasblinken"
)

type LedControlServer struct {
	EffectHandler func(string, dasblinken.Channel) error
	StopHandler   func(dasblinken.Channel)
	EffectFetcher func() []dasblinken.Effect
}

type EffectInfo struct {
	Name string `json:"name"`
}

func (s *LedControlServer) RunServer() {
	http.HandleFunc("/", s.handleClient)
	http.HandleFunc("/switch", s.handleSwitch)
	http.HandleFunc("/list", s.handleList)
	http.HandleFunc("/stop", s.handleStop)

	fmt.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

func (s *LedControlServer) handleClient(w http.ResponseWriter, r *http.Request) {
	head := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>DASBLINKENLIGHTS</title>
		 <style>
            body {
                display: flex;
                justify-content: center;
                align-items: center;
                height: 100vh;
                margin: 0;
            }
            .effect-button {
                width: 320px;
                height: 40px;
                font-size: 18px;
				margin: 0px 2px;
            }
			.stop-button {
                height: 50px;
				width: 320px;
                font-size: 24px;
                margin: 10px auto;
            }
        </style>
		<script>
            function callEffect(name) {
                fetch('/switch?name=' + name + '&channel=0')
                    .then(response => response.text())
                    .then(data => console.log(data))
                    .catch(error => console.error('Error:', error));
            }

			function stop() {
                fetch('/stop')
                    .then(response => response.text())
                    .then(data => console.log(data))
                    .catch(error => console.error('Error:', error));
            }
        </script>
	</head>
	<body>
		<div>
		<h1>DASBLIKENCONTROLLER</h1>
        <div>
	`
	w.Write([]byte(head))

	effects := s.EffectFetcher()
	for _, effect := range effects {
		name := effect.Opts().Name
		fmt.Fprintf(w, "<div><button class=\"effect-button\" onclick='callEffect(\"%s\")'>%s</button></div>\n", name, name)
	}

	foot := `
		</div>
        <button class="stop-button" onclick="stop()">STOP</button>
		</div>
	</body>
	</html>
	`
	w.Write([]byte(foot))

}

// handleEffectList returns a list of effect names in JSON format
func (s *LedControlServer) handleList(w http.ResponseWriter, r *http.Request) {
	effects := s.EffectFetcher()
	output := make([]EffectInfo, 0, len(effects))
	for _, effect := range effects {
		info := EffectInfo{Name: effect.Opts().Name}
		output = append(output, info)
	}

	jsonData, err := json.Marshal(output)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (s *LedControlServer) handleSwitch(w http.ResponseWriter, r *http.Request) {
	// Extract the index from the query string
	name := r.URL.Query().Get("name")
	channelStr := r.URL.Query().Get("channel")
	channel, err := strconv.Atoi(channelStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("Effect name: %d\n", name)

	err = s.EffectHandler(name, dasblinken.Channel(channel))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Return a 200 OK status
	w.WriteHeader(http.StatusOK)
}

func (s *LedControlServer) handleStop(w http.ResponseWriter, r *http.Request) {
	s.StopHandler(0)
	w.WriteHeader(http.StatusOK)
}
