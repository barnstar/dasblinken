package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

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
	mux.HandleFunc("/effects-config", s.handleEffectsConfigPage)
	mux.HandleFunc("/api/effects", s.handleEffects)
	mux.HandleFunc("/api/switch", s.handleSwitch)
	mux.HandleFunc("/api/stop", s.handleStop)
	mux.HandleFunc("/api/config", s.handleConfig)
	mux.HandleFunc("/api/update-config", s.handleUpdateConfig)
	mux.HandleFunc("/api/effect-config", s.handleEffectConfig)
	mux.HandleFunc("/api/update-effect-config", s.handleUpdateEffectConfig)

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
        .nav-button {
            display: inline-block;
            padding: 12px 24px;
            margin: 10px 5px;
            border: none;
            border-radius: 8px;
            background: linear-gradient(135deg, #ffa726 0%, #fb8c00 100%);
            color: white;
            font-size: 1rem;
            font-weight: 600;
            cursor: pointer;
            text-decoration: none;
            transition: all 0.3s ease;
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
        }
        .nav-button:hover {
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
            <a href="/effects-config" class="nav-button">🔧 Configure Effects</a>
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
        <div class="effects-container">
		    <h2>Available Effects</h2>
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

func (s *LedControlServer) handleEffectsConfigPage(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Effects Configuration - DasBlinkenLights</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }
        .container { max-width: 1400px; margin: 0 auto; }
        .header { text-align: center; color: white; margin-bottom: 40px; }
        .header h1 { font-size: 2.5rem; font-weight: 700; margin-bottom: 10px; text-shadow: 2px 2px 4px rgba(0,0,0,0.3); }
        .header p { font-size: 1.1rem; opacity: 0.9; }
        .nav-button {
            display: inline-block;
            padding: 10px 20px;
            margin: 10px 5px;
            border: none;
            border-radius: 8px;
            background: rgba(255,255,255,0.2);
            color: white;
            font-size: 0.95rem;
            font-weight: 600;
            cursor: pointer;
            text-decoration: none;
            transition: all 0.3s ease;
        }
        .nav-button:hover { background: rgba(255,255,255,0.3); }
        .main-content { display: flex; gap: 20px; }
        .sidebar {
            flex: 0 0 280px;
            background: white;
            border-radius: 12px;
            padding: 20px;
            box-shadow: 0 8px 16px rgba(0,0,0,0.2);
            max-height: calc(100vh - 200px);
            overflow-y: auto;
        }
        .sidebar h2 { color: #333; margin-bottom: 15px; font-size: 1.3rem; }
        .effect-type {
            padding: 12px;
            margin: 5px 0;
            border-radius: 6px;
            cursor: pointer;
            background: #f5f5f5;
            transition: all 0.3s ease;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .effect-type:hover { background: #e0e0e0; }
        .effect-type.active { background: #667eea; color: white; font-weight: 600; }
        .effect-count {
            background: rgba(0,0,0,0.1);
            padding: 2px 8px;
            border-radius: 12px;
            font-size: 0.85rem;
        }
        .content-area {
            flex: 1;
            background: white;
            border-radius: 12px;
            padding: 25px;
            box-shadow: 0 8px 16px rgba(0,0,0,0.2);
            max-height: calc(100vh - 200px);
            overflow-y: auto;
        }
        .content-area h2 { color: #333; margin-bottom: 20px; font-size: 1.5rem; }
        .effect-list { display: grid; gap: 15px; }
        .effect-item {
            border: 2px solid #e0e0e0;
            border-radius: 8px;
            padding: 15px;
            background: #fafafa;
            transition: all 0.3s ease;
        }
        .effect-item:hover { border-color: #667eea; background: #f8f9ff; }
        .effect-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 10px;
        }
        .effect-name { font-size: 1.1rem; font-weight: 600; color: #333; }
        .effect-actions { display: flex; gap: 8px; }
        .btn {
            padding: 6px 12px;
            border: none;
            border-radius: 6px;
            cursor: pointer;
            font-size: 0.85rem;
            font-weight: 600;
            transition: all 0.3s ease;
        }
        .btn-edit { background: #ffa726; color: white; }
        .btn-edit:hover { background: #fb8c00; }
        .btn-delete { background: #ef5350; color: white; }
        .btn-delete:hover { background: #e53935; }
        .btn-add { 
            width: 100%;
            padding: 12px;
            background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
            color: white;
            font-size: 1rem;
            margin-bottom: 20px;
        }
        .btn-add:hover { transform: translateY(-2px); box-shadow: 0 4px 8px rgba(0,0,0,0.2); }
        .effect-props {
            display: flex;
            flex-direction: column;
            gap: 8px;
            font-size: 0.9rem;
            color: #666;
        }
        .prop { 
            display: flex; 
            justify-content: space-between;
            padding: 8px 12px;
            background: #f8f9fa;
            border-radius: 6px;
        }
        .prop-name { font-weight: 600; }
        .prop-value { color: #667eea; }
        .modal {
            display: none;
            position: fixed;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            background: rgba(0,0,0,0.5);
            z-index: 1000;
            align-items: center;
            justify-content: center;
        }
        .modal.show { display: flex; }
        .modal-content {
            background: white;
            border-radius: 12px;
            padding: 30px;
            max-width: 600px;
            width: 90%;
            max-height: 80vh;
            overflow-y: auto;
        }
        .modal-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 20px;
        }
        .modal-header h3 { color: #333; font-size: 1.5rem; }
        .close-btn {
            background: none;
            border: none;
            font-size: 1.5rem;
            cursor: pointer;
            color: #999;
        }
        .close-btn:hover { color: #333; }
        .form-group {
            margin-bottom: 15px;
        }
        .form-group label {
            display: block;
            color: #666;
            font-size: 0.9rem;
            margin-bottom: 5px;
            font-weight: 600;
        }
        .form-group input, .form-group select {
            width: 100%;
            padding: 10px;
            border: 2px solid #e0e0e0;
            border-radius: 6px;
            font-size: 1rem;
        }
        .form-group input:focus, .form-group select:focus {
            outline: none;
            border-color: #667eea;
        }
        .btn-save {
            width: 100%;
            padding: 12px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            font-size: 1.1rem;
        }
        .btn-save:hover { transform: translateY(-2px); box-shadow: 0 4px 8px rgba(0,0,0,0.2); }
        .status-message {
            position: fixed;
            top: 20px;
            right: 20px;
            padding: 16px 24px;
            background: white;
            border-radius: 8px;
            box-shadow: 0 4px 12px rgba(0,0,0,0.15);
            display: none;
            z-index: 2000;
        }
        .status-message.show { display: block; animation: slideIn 0.3s ease; }
        @keyframes slideIn {
            from { transform: translateX(400px); opacity: 0; }
            to { transform: translateX(0); opacity: 1; }
        }
        .status-message.success { border-left: 4px solid #4caf50; color: #2e7d32; }
        .status-message.error { border-left: 4px solid #f44336; color: #c62828; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔧 Effects Configuration</h1>
            <p>Manage LED Effect Parameters</p>
            <a href="/" class="nav-button">← Back to Controller</a>
        </div>
        <div class="main-content">
            <div class="sidebar">
                <h2>Effect Types</h2>
                <div id="effectTypes"></div>
            </div>
            <div class="content-area">
                <h2 id="contentTitle">Select an effect type</h2>
                <button class="btn btn-add" id="addBtn" style="display:none;" onclick="showAddModal()">+ Add New Effect</button>
                <div id="effectList"></div>
            </div>
        </div>
    </div>
    
    <div class="modal" id="editModal">
        <div class="modal-content">
            <div class="modal-header">
                <h3 id="modalTitle">Edit Effect</h3>
                <button class="close-btn" onclick="closeModal()">×</button>
            </div>
            <form id="effectForm">
                <div id="formFields"></div>
                <button type="submit" class="btn btn-save">Save Changes</button>
            </form>
        </div>
    </div>
    
    <div class="status-message" id="statusMessage"></div>
    
    <script>
        let effectsConfig = {};
        let currentType = null;
        let editingIndex = null;
        
        const effectSchemas = {
            balls: {
                name: { type: 'text', label: 'Effect Name', required: true },
                topology: { type: 'select', label: 'Topology', options: ['linear', 'matrix', 'any'], required: true },
                numBalls: { type: 'number', label: 'Number of Balls', min: 1, max: 100, required: true },
                trailLen: { type: 'number', label: 'Trail Length', min: 0, max: 20, required: true },
                palette: { type: 'select', label: 'Color Palette', options: ['rainbow', 'heat', 'cold', 'green', 'ice', 'festive', 'fullfire'], required: true }
            },
            race: {
                name: { type: 'text', label: 'Effect Name', required: true },
                topology: { type: 'select', label: 'Topology', options: ['linear', 'matrix', 'any'], required: true },
                length: { type: 'number', label: 'Length', min: 1, max: 100, required: true },
                mirrored: { type: 'checkbox', label: 'Mirrored', required: false },
                numRacers: { type: 'number', label: 'Number of Racers', min: 1, max: 10, required: true }
            },
            wave: {
                name: { type: 'text', label: 'Effect Name', required: true },
                topology: { type: 'select', label: 'Topology', options: ['linear', 'matrix', 'any'], required: true },
                palette: { type: 'select', label: 'Color Palette', options: ['rainbow', 'heat', 'cold', 'green', 'ice', 'festive', 'fullfire'], required: true }
            },
            chase: {
                name: { type: 'text', label: 'Effect Name', required: true },
                topology: { type: 'select', label: 'Topology', options: ['linear', 'matrix', 'any'], required: true },
                speed: { type: 'number', label: 'Speed', min: 0.1, max: 5.0, step: 0.1, required: true },
                palette: { type: 'select', label: 'Color Palette', options: ['rainbow', 'heat', 'cold', 'green', 'ice', 'festive', 'fullfire'], required: true }
            },
            fire: {
                name: { type: 'text', label: 'Effect Name', required: true },
                topology: { type: 'select', label: 'Topology', options: ['linear', 'matrix', 'any'], required: true },
                sparking: { type: 'number', label: 'Sparking', min: 0, max: 1, step: 0.1, required: true },
                cooling: { type: 'number', label: 'Cooling', min: 0, max: 0.1, step: 0.01, required: true },
                doubleEnded: { type: 'checkbox', label: 'Double Ended', required: false },
                palette: { type: 'select', label: 'Color Palette', options: ['heat', 'cold', 'rainbow', 'green', 'ice', 'festive', 'fullfire'], required: true }
            },
            fireMatrix: {
                name: { type: 'text', label: 'Effect Name', required: true },
                topology: { type: 'select', label: 'Topology', options: ['matrix'], required: true },
                sparking: { type: 'number', label: 'Sparking', min: 0, max: 1, step: 0.1, required: true },
                cooling: { type: 'number', label: 'Cooling', min: 0, max: 0.1, step: 0.01, required: true },
                palette: { type: 'select', label: 'Color Palette', options: ['heat', 'cold', 'rainbow', 'green', 'ice', 'festive', 'fullfire'], required: true }
            },
            snow: {
                name: { type: 'text', label: 'Effect Name', required: true },
                topology: { type: 'select', label: 'Topology', options: ['linear', 'matrix', 'any'], required: true },
                numFlakes: { type: 'number', label: 'Number of Flakes', min: 1, max: 100, required: true },
                speed: { type: 'number', label: 'Speed', min: 0.1, max: 5.0, step: 0.1, required: true }
            },
            solid: {
                name: { type: 'text', label: 'Effect Name', required: true },
                topology: { type: 'select', label: 'Topology', options: ['linear', 'matrix', 'any'], required: true },
                frameDelay: { type: 'number', label: 'Frame Delay', min: 0, max: 1000, step: 1, required: true },
                palette: { type: 'select', label: 'Color Palette', options: ['', 'heat', 'cold', 'rainbow', 'green', 'ice', 'festive', 'fullfire'], required: false },
                mode: { type: 'select', label: 'Color Mode', options: ['rotate', 'random'], required: true }
            },
            textScroll: {
                name: { type: 'text', label: 'Effect Name', required: true },
                topology: { type: 'select', label: 'Topology', options: ['matrix'], required: true },
                text: { type: 'text', label: 'Text to Display', required: true },
                palette: { type: 'select', label: 'Color Palette', options: ['rainbow', 'heat', 'cold', 'green', 'ice', 'festive', 'fullfire'], required: true },
                colorTransform: { type: 'select', label: 'Color Transform', options: ['rotate', 'random'], required: true }
            },
            fontTest: {
                name: { type: 'text', label: 'Effect Name', required: true },
                topology: { type: 'select', label: 'Topology', options: ['matrix'], required: true }
            },
            static: {
                name: { type: 'text', label: 'Effect Name', required: true },
                topology: { type: 'select', label: 'Topology', options: ['linear', 'matrix', 'any'], required: true },
                imageFile: { type: 'text', label: 'Image File Path', required: true }
            },
            clock: {
                name: { type: 'text', label: 'Effect Name', required: true },
                topology: { type: 'select', label: 'Topology', options: ['matrix'], required: true }
            },
            flicker: {
                name: { type: 'text', label: 'Effect Name', required: true },
                topology: { type: 'select', label: 'Topology', options: ['linear', 'matrix', 'any'], required: true },
                numGroups: { type: 'number', label: 'Number of Groups', min: 1, max: 50, required: true },
                minRadius: { type: 'number', label: 'Min Radius', min: 1, max: 100, required: true },
                maxRadius: { type: 'number', label: 'Max Radius', min: 1, max: 100, required: true },
                flickerSpeed: { type: 'number', label: 'Flicker Speed', min: 0.1, max: 10, step: 0.1, required: true },
                flickerAmt: { type: 'number', label: 'Flicker Amount', min: 0, max: 1, step: 0.05, required: true },
                palette: { type: 'select', label: 'Color Palette', options: ['', 'heat', 'cold', 'rainbow', 'green', 'ice', 'festive', 'fullfire'], required: false }
            }
        };
        
        async function loadEffectsConfig() {
            try {
                const response = await fetch('/api/effect-config');
                effectsConfig = await response.json();
                renderEffectTypes();
            } catch (error) {
                showStatus('Failed to load effects configuration', 'error');
            }
        }
        
        function renderEffectTypes() {
            const container = document.getElementById('effectTypes');
            container.innerHTML = '';
            
            for (const [type, effects] of Object.entries(effectsConfig)) {
                const div = document.createElement('div');
                div.className = 'effect-type';
                div.innerHTML = '<span>' + type + '</span><span class="effect-count">' + effects.length + '</span>';
                div.onclick = () => selectEffectType(type);
                container.appendChild(div);
            }
        }
        
        function selectEffectType(type) {
            currentType = type;
            document.querySelectorAll('.effect-type').forEach(el => el.classList.remove('active'));
            event.target.closest('.effect-type').classList.add('active');
            document.getElementById('contentTitle').textContent = type + ' Effects';
            document.getElementById('addBtn').style.display = 'block';
            renderEffectList();
        }
        
        function renderEffectList() {
            const container = document.getElementById('effectList');
            const effects = effectsConfig[currentType] || [];
            
            if (effects.length === 0) {
                container.innerHTML = '<p style="text-align:center; color:#999; padding:40px;">No effects configured. Click "Add New Effect" to create one.</p>';
                return;
            }
            
            container.innerHTML = '';
            effects.forEach((effect, index) => {
                const div = document.createElement('div');
                div.className = 'effect-item';
                
                const props = Object.entries(effect)
                    .filter(([key]) => key !== 'name')
                    .sort(([a], [b]) => a.localeCompare(b))
                    .map(([key, value]) => '<div class="prop"><span class="prop-name">' + key + ':</span><span class="prop-value">' + (typeof value === 'boolean' ? (value ? 'Yes' : 'No') : value) + '</span></div>').join('');
                
                const header = '<div class="effect-header"><div class="effect-name">' + effect.name + '</div>' +
                    '<div class="effect-actions">' +
                    '<button class="btn btn-edit" onclick="editEffect(' + index + ')">Edit</button>' +
                    '<button class="btn btn-delete" onclick="deleteEffect(' + index + ')">Delete</button>' +
                    '</div></div>';
                
                div.innerHTML = header + '<div class="effect-props">' + props + '</div>';
                container.appendChild(div);
            });
        }
        
        function showAddModal() {
            editingIndex = null;
            document.getElementById('modalTitle').textContent = 'Add New ' + currentType + ' Effect';
            renderForm({});
            document.getElementById('editModal').classList.add('show');
        }
        
        function editEffect(index) {
            editingIndex = index;
            document.getElementById('modalTitle').textContent = 'Edit Effect';
            renderForm(effectsConfig[currentType][index]);
            document.getElementById('editModal').classList.add('show');
        }
        
        function renderForm(data) {
            const schema = effectSchemas[currentType];
            const container = document.getElementById('formFields');
            container.innerHTML = '';
            
            for (const [field, config] of Object.entries(schema)) {
                const div = document.createElement('div');
                div.className = 'form-group';
                
                let input;
                if (config.type === 'select') {
                    input = '<select id="field_' + field + '" name="' + field + '">' +
                        config.options.map(opt => 
                            '<option value="' + opt + '" ' + (data[field] === opt ? 'selected' : '') + '>' + opt + '</option>'
                        ).join('') +
                    '</select>';
                } else if (config.type === 'checkbox') {
                    input = '<input type="checkbox" id="field_' + field + '" name="' + field + '" ' + (data[field] ? 'checked' : '') + '>';
                } else {
                    const attrs = [];
                    if (config.min !== undefined) attrs.push('min="' + config.min + '"');
                    if (config.max !== undefined) attrs.push('max="' + config.max + '"');
                    if (config.step !== undefined) attrs.push('step="' + config.step + '"');
                    input = '<input type="' + config.type + '" id="field_' + field + '" name="' + field + '" value="' + (data[field] || '') + '" ' + attrs.join(' ') + '>';
                }
                
                div.innerHTML = '<label for="field_' + field + '">' + config.label + (config.required ? ' *' : '') + '</label>' + input;
                container.appendChild(div);
            }
        }
        
        function closeModal() {
            document.getElementById('editModal').classList.remove('show');
        }
        
        document.getElementById('effectForm').addEventListener('submit', async (e) => {
            e.preventDefault();
            const formData = new FormData(e.target);
            const data = {};
            const schema = effectSchemas[currentType];
            
            for (const [field, config] of Object.entries(schema)) {
                const value = formData.get(field);
                if (config.type === 'number') {
                    data[field] = parseFloat(value);
                } else if (config.type === 'checkbox') {
                    data[field] = document.getElementById('field_' + field).checked;
                } else {
                    data[field] = value;
                }
            }
            
            if (editingIndex !== null) {
                effectsConfig[currentType][editingIndex] = data;
            } else {
                if (!effectsConfig[currentType]) {
                    effectsConfig[currentType] = [];
                }
                effectsConfig[currentType].push(data);
            }
            
            await saveEffectsConfig();
            closeModal();
            renderEffectList();
        });
        
        async function deleteEffect(index) {
            if (!confirm('Are you sure you want to delete this effect?')) return;
            
            effectsConfig[currentType].splice(index, 1);
            await saveEffectsConfig();
            renderEffectList();
        }
        
        async function saveEffectsConfig() {
            try {
                const response = await fetch('/api/update-effect-config', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(effectsConfig)
                });
                
                if (response.ok) {
                    showStatus('Effects configuration saved successfully!', 'success');
                } else {
                    showStatus('Failed to save configuration', 'error');
                }
            } catch (error) {
                showStatus('Error saving configuration', 'error');
            }
        }
        
        function showStatus(message, type) {
            const statusDiv = document.getElementById('statusMessage');
            statusDiv.textContent = message;
            statusDiv.className = 'status-message show ' + type;
            setTimeout(() => statusDiv.classList.remove('show'), 3000);
        }
        
        loadEffectsConfig();
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

func (s *LedControlServer) handleEffectConfig(w http.ResponseWriter, r *http.Request) {
	// Read the effects.json file and return it
	data, err := os.ReadFile(s.EffectsFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to read effects config: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (s *LedControlServer) handleUpdateEffectConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Read the new effects configuration from request body
	var effectsConfig map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&effectsConfig); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request body: %v", err)
		return
	}

	// Write to effects file
	data, err := json.MarshalIndent(effectsConfig, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to marshal effects config: %v", err)
		return
	}

	if err := os.WriteFile(s.EffectsFile, data, 0644); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to save effects config: %v", err)
		return
	}

	// Reload effects
	s.Das.StopAll()
	s.Das.ClearEffects()
	config := s.Das.Config()
	effects.LoadEffectsFromFile(s.EffectsFile, s.Das.RegisterEffect, config)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Effects configuration updated successfully")
}
