package server

import (
	"fmt"
	"net/http"
	"strconv"
)

type LedControlServer struct {
	EffectHandler func(int)
}

func (s *LedControlServer) RunServer() {
	http.HandleFunc("/", s.handleClient)
	http.HandleFunc("/effect", s.handleEffect)

	fmt.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

func (s *LedControlServer) handleClient(w http.ResponseWriter, r *http.Request) {
	html := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>DASBLIKENCONTROLLER</title>
		 <style>
            body {
                display: flex;
                justify-content: center;
                align-items: center;
                height: 100vh;
                margin: 0;
            }
            .grid-container {
                display: grid;
                grid-template-columns: repeat(3, 1fr);
                gap: 10px;
            }
            .grid-container button {
                width: 100px;
                height: 100px;
                font-size: 24px;
            }
        </style>
		<script>
            function callEffect(index) {
                fetch('/effect?index=' + index)
                    .then(response => response.text())
                    .then(data => console.log(data))
                    .catch(error => console.error('Error:', error));
            }
        </script>
	</head>
	<body>
		<div>
		<h1>DASBLIKENCONTROLLER</h1>
        <div class="grid-container">
            <button onclick="callEffect(1)">1</button>
            <button onclick="callEffect(2)">2</button>
            <button onclick="callEffect(3)">3</button>
            <button onclick="callEffect(4)">4</button>
            <button onclick="callEffect(5)">5</button>
            <button onclick="callEffect(6)">6</button>
            <button onclick="callEffect(7)">7</button>
            <button onclick="callEffect(8)">8</button>
            <button onclick="callEffect(9)">9</button>
		</div>
		</div>
	</body>
	</html>
	`
	w.Write([]byte(html))
}

func (s *LedControlServer) handleEffect(w http.ResponseWriter, r *http.Request) {
	// Extract the index from the query string
	indexStr := r.URL.Query().Get("index")
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		http.Error(w, "Invalid index", http.StatusBadRequest)
		return
	}
	fmt.Printf("Effect index: %d\n", index)

	s.EffectHandler(index - 1)

	// Return a 200 OK status
	w.WriteHeader(http.StatusOK)
}
