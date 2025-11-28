package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	dasblinken "barnstar.com/dasblinken"
	effects "barnstar.com/effects"
	"tailscale.com/tsnet"
)

type LedControlServer struct {
	Das         *dasblinken.Dasblinken
	Hostname    string
	AuthKey     string
	ConfigFile  string
	EffectsFile string
}

type EffectInfo struct {
	Name string `json:"name"`
}

type EffectsResponse struct {
	Effects []EffectInfo `json:"effects"`
}

func (s *LedControlServer) RunServer() {
	srv := &tsnet.Server{
		Hostname: s.Hostname,
		AuthKey:  s.AuthKey,
		Dir:      "/var/lib/tsnet-" + s.Hostname,
	}

	defer srv.Close()

	ln, err := srv.Listen("tcp", ":80")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleClient)
	mux.HandleFunc("/api/effects", s.handleEffects)
	mux.HandleFunc("/api/switch", s.handleSwitch)
	mux.HandleFunc("/api/stop", s.handleStop)
	mux.HandleFunc("/api/config", s.handleConfig)
	mux.HandleFunc("/api/update-config", s.handleUpdateConfig)

	fmt.Printf("Server starting on tailnet as %s...\\n", s.Hostname)
	if err := http.Serve(ln, mux); err != nil {
		log.Fatal(err)
	}
}

func (s *LedControlServer) handleClient(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>DasBlinkenLights</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }
        .container { max-width: 1200px; margin: 0 auto; }
        .header { text-align: center; color: white; margin-bottom: 40px; }
        .header h1 { font-size: 3rem; font-weight: 700; margin-bottom: 10px; text-shadow: 2px 2px 4px rgba(0,0,0,0.3); }
        .header p { font-size: 1.2rem; opacity: 0.9; }
        .effects-container {
            background: white;
            border-radius: 12px;
            padding: 20px;
            box-shadow: 0 8px 16px rgba(0,0,0,0.2);
        }
        .effects-container h2 { color: #333; margin-bottom: 20px; font-size: 1.5rem; }
        .effects-grid {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
            gap: 12px;
            margin-bottom: 20px;
        }
        .effect-button {
            padding: 16px;
            border: none;
            border-radius: 8px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            font-size: 1rem;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.3s ease;
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
        }
        .effect-button:hover {
            transform: translateY(-2px);
            box-shadow: 0 6px 12px rgba(0,0,0,0.2);
        }
        .stop-button {
            width: 100%;
            padding: 20px;
            border: none;
            border-radius: 8px;
            background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
            color: white;
            font-size: 1.2rem;
            font-weight: 700;
            cursor: pointer;
            transition: all 0.3s ease;
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
        }
        .stop-button:hover {
            transform: translateY(-2px);
            box-shadow: 0 6px 12px rgba(0,0,0,0.2);
        }
        .loading { text-align: center; padding: 40px; color: #666; }
        .status-message {
            position: fixed;
            top: 20px;
            right: 20px;
            padding: 16px 24px;
            background: white;
            border-radius: 8px;
            box-shadow: 0 4px 12px rgba(0,0,0,0.15);
            display: none;
            z-index: 1000;
        }
        .status-message.show { display: block; animation: slideIn 0.3s ease; }
        @keyframes slideIn {
            from { transform: translateX(400px); opacity: 0; }
            to { transform: translateX(0); opacity: 1; }
        }
        .status-message.success { border-left: 4px solid #4caf50; color: #2e7d32; }
        .status-message.error { border-left: 4px solid #f44336; color: #c62828; }
        .config-container {
            background: white;
            border-radius: 12px;
            padding: 20px;
            margin-bottom: 20px;
            box-shadow: 0 8px 16px rgba(0,0,0,0.2);
        }
        .config-container h2 { color: #333; margin-bottom: 20px; font-size: 1.5rem; }
        .config-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 15px;
            margin-bottom: 20px;
        }
        .config-field label {
            display: block;
            color: #666;
            font-size: 0.9rem;
            margin-bottom: 5px;
            font-weight: 600;
        }
        .config-field input {
            width: 100%;
            padding: 10px;
            border: 2px solid #e0e0e0;
            border-radius: 6px;
            font-size: 1rem;
            transition: border-color 0.3s ease;
        }
        .config-field input:focus {
            outline: none;
            border-color: #667eea;
        }
        .config-info {
            padding: 15px;
            background: #f8f9fa;
            border-radius: 8px;
            color: #666;
            font-size: 0.9rem;
            margin-bottom: 15px;
        }
        .update-button {
            width: 100%;
            padding: 15px;
            border: none;
            border-radius: 8px;
            background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
            color: white;
            font-size: 1.1rem;
            font-weight: 700;
            cursor: pointer;
            transition: all 0.3s ease;
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
        }
        .update-button:hover {
            transform: translateY(-2px);
            box-shadow: 0 6px 12px rgba(0,0,0,0.2);
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🌟 DasBlinkenLights 🌟</h1>
            <p>LED Effect Controller</p>
        </div>
        <div class="config-container">
            <h2>⚙️ Strip Configuration</h2>
            <div class="config-info" id="configInfo">Loading configuration...</div>
            <div class="config-grid">
                <div class="config-field">
                    <label for="width">Width (LEDs)</label>
                    <input type="number" id="width" min="1" max="1000" />
                </div>
                <div class="config-field">
                    <label for="height">Height (LEDs)</label>
                    <input type="number" id="height" min="1" max="100" />
                </div>
                <div class="config-field">
                    <label for="pin">GPIO Pin</label>
                    <input type="number" id="pin" min="0" max="40" />
                </div>
                <div class="config-field">
                    <label for="brightness">Brightness (0-255)</label>
                    <input type="number" id="brightness" min="0" max="255" />
                </div>
            </div>
            <button class="update-button" onclick="updateConfig()">💾 Save & Apply Configuration</button>
        </div>
        <div class="effects-container">\n            <h2>Available Effects</h2>
            <div class="effects-grid" id="effectsGrid">
                <div class="loading">Loading effects...</div>
            </div>
            <button class="stop-button" onclick="stopEffect()">⏹ STOP</button>
        </div>
    </div>
    <div class="status-message" id="statusMessage"></div>
    <script>
        let currentConfig = null;
        
        async function loadConfig() {
            try {
                const response = await fetch('/api/config');
                currentConfig = await response.json();
                document.getElementById('width').value = currentConfig.width;
                document.getElementById('height').value = currentConfig.height;
                document.getElementById('pin').value = currentConfig.pin;
                document.getElementById('brightness').value = currentConfig.brightness;
                
                const topology = currentConfig.height > 1 ? 'Matrix' : 'Linear';
                document.getElementById('configInfo').innerHTML = 
                    '<strong>Current:</strong> ' + currentConfig.width + 'x' + currentConfig.height + 
                    ' (' + topology + ') | GPIO Pin: ' + currentConfig.pin + 
                    ' | Brightness: ' + currentConfig.brightness + ' | FPS: ' + currentConfig.fps;
            } catch (error) {
                showStatus('Failed to load config', 'error');
            }
        }
        
        async function updateConfig() {
            const config = {
                pin: parseInt(document.getElementById('pin').value),
                brightness: parseInt(document.getElementById('brightness').value),
                width: parseInt(document.getElementById('width').value),
                height: parseInt(document.getElementById('height').value),
                fps: currentConfig ? currentConfig.fps : 30,
                invert: currentConfig ? currentConfig.invert : false
            };
            
            try {
                const response = await fetch('/api/update-config', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(config)
                });
                
                if (response.ok) {
                    showStatus('Configuration updated! Reloading effects...', 'success');
                    await loadConfig();
                    await loadEffects();
                } else {
                    showStatus('Failed to update configuration', 'error');
                }
            } catch (error) {
                showStatus('Error updating configuration', 'error');
            }
        }
        
        async function loadEffects() {
            const grid = document.getElementById('effectsGrid');
            grid.innerHTML = '<div class="loading">Loading effects...</div>';
            try {
                const response = await fetch('/api/effects');
                const data = await response.json();
                grid.innerHTML = '';
                if (data.effects.length === 0) {
                    grid.innerHTML = '<div class="loading">No effects available</div>';
                    return;
                }
                data.effects.forEach(effect => {
                    const button = document.createElement('button');
                    button.className = 'effect-button';
                    button.textContent = effect.name;
                    button.onclick = () => switchEffect(effect.name);
                    grid.appendChild(button);
                });
            } catch (error) {
                showStatus('Failed to load effects', 'error');
            }
        }
        async function switchEffect(name) {
            try {
                const response = await fetch('/api/switch?name=' + encodeURIComponent(name));
                if (response.ok) {
                    showStatus('Switched to ' + name, 'success');
                } else {
                    showStatus('Failed to switch effect', 'error');
                }
            } catch (error) {
                showStatus('Error switching effect', 'error');
            }
        }
        async function stopEffect() {
            try {
                const response = await fetch('/api/stop');
                if (response.ok) {
                    showStatus('Stopped', 'success');
                } else {
                    showStatus('Failed to stop', 'error');
                }
            } catch (error) {
                showStatus('Error stopping effect', 'error');
            }
        }
        function showStatus(message, type) {
            const statusDiv = document.getElementById('statusMessage');
            statusDiv.textContent = message;
            statusDiv.className = 'status-message show ' + type;
            setTimeout(() => statusDiv.classList.remove('show'), 3000);
        }
        
        loadConfig();
        loadEffects();
    </script>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func (s *LedControlServer) handleEffects(w http.ResponseWriter, r *http.Request) {
	effects := s.Das.Effects()
	effectInfos := make([]EffectInfo, 0, len(effects))
	for _, effect := range effects {
		effectInfos = append(effectInfos, EffectInfo{Name: effect.Opts().Name})
	}
	response := EffectsResponse{Effects: effectInfos}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *LedControlServer) handleSwitch(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	err := s.Das.SwitchToEffect(name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Switched to %s", name)
}

func (s *LedControlServer) handleStop(w http.ResponseWriter, r *http.Request) {
	s.Das.Stop()
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Stopped")
}

func (s *LedControlServer) handleConfig(w http.ResponseWriter, r *http.Request) {
	config := s.Das.Config()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}

func (s *LedControlServer) handleUpdateConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var config dasblinken.StripConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body: %v", err)
		return
	}

	// Stop all effects and update configuration
	s.Das.UpdateConfig(config)

	// Save configuration to file
	if err := config.SaveTo(s.ConfigFile); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to save config: %v", err)
		return
	}

	// Reload effects with new topology
	s.Das.ClearEffects()
	effects.LoadEffectsFromFile(s.EffectsFile, s.Das.RegisterEffect, config)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Configuration updated successfully")
}
