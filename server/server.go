package server

import (
	"fmt"
	"net/http"
)

type LedControlServer struct {
}

func (s *LedControlServer) RunServer() {
	http.HandleFunc("/", s.handleClient)
	http.HandleFunc("/on", s.handleOn)
	http.HandleFunc("/off", s.handleOff)

	fmt.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

func (s *LedControlServer) handleClient(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("I'm alive...\n"))
}

func (s *LedControlServer) handleOn(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("On!\n"))
}

func (s *LedControlServer) handleOff(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Off!\n"))
}
