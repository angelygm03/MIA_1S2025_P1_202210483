import { useState } from 'react';
import './App.css';

function App() {
  const [input, setInput] = useState('');
  const [output, setOutput] = useState('');

  const handleExecute = async () => {
    try {
      const commands = input.trim().split("\n"); // Commands are separated by lines
      let results = [];
    
      for (const command of commands) {
        const params = command.split(" ");
        let requestBody = {};
        let endpoint = "";
  
        if (command.startsWith("mkdisk")) {
          // Get parameters for mkdisk
          let size = 0, unit = "k", fit = "", path = "";
          params.forEach(param => {
            if (param.startsWith("-size=")) size = parseInt(param.split("=")[1]); // Convert to number
            if (param.startsWith("-unit=")) unit = param.split("=")[1].toLowerCase();
            if (param.startsWith("-fit=")) fit = param.split("=")[1].toLowerCase();
            if (param.startsWith("-path=")) path = param.split("=")[1].replace(/"/g, ''); // Remove ""
          });
          
          // Set request body and endpoint for mkdisk
          requestBody = { size, unit, fit, path };
          endpoint = "mkdisk";
  
        } else if (command.startsWith("rmdisk")) {
          // Get parameters for rmdisk
          let path = "";
          params.forEach(param => {
            if (param.startsWith("-path=")) path = param.split("=")[1].replace(/"/g, '');
          });
          
          // Set request body and endpoint for rmdisk
          requestBody = { path };
          endpoint = "rmdisk";

        } else if (command.startsWith("fdisk")) {
          let size = 0, unit = "k", path = "", type = "p", fit= "wf", name = "";
          params.forEach(param => {
            if (param.startsWith("-size=")) size = parseInt(param.split("=")[1]);
            if (param.startsWith("-path=")) path = param.split("=")[1].replace(/"/g, '');
            if (param.startsWith("-name=")) name = param.split("=")[1].replace(/"/g, '');
            if (param.startsWith("-unit=")) unit = param.split("=")[1].toLowerCase();
            if (param.startsWith("-type=")) type = param.split("=")[1].toLowerCase();
            if (param.startsWith("-fit="))  fit = param.split("=")[1].toLowerCase();
          });
  
          requestBody = { size, unit, path, type, fit, name};
          endpoint = "fdisk";
        
        } else  if (command.startsWith("mount")) {
            let path = "", name = "";
            params.forEach(param => {
              if (param.startsWith("-path=")) path = param.split("=")[1].replace(/"/g, '');
              if (param.startsWith("-name=")) name = param.split("=")[1].replace(/"/g, '');
            });
          
            requestBody = { path, name };
            endpoint = "mount";
          
        } else  if (command.startsWith("login")) {
          let user = "", password = "", id="";
          params.forEach(param => {
            if (param.startsWith("-user=")) user = param.split("=")[1];
            if (param.startsWith("-password=")) password = param.split("=")[1];
            if (param.startsWith("-id=")) id = param.split("=")[1].toLowerCase();
          });
        
          requestBody = { user, password, id };
          endpoint = "login";
          
        } else  if (command.startsWith("rep")) {
          let path = "", name = "", id = ""; 
          params.forEach(param => {
            if (param.startsWith("-path=")) path = param.split("=")[1].replace(/"/g, '');
            if (param.startsWith("-name=")) name = param.split("=")[1].replace(/"/g, '');
            if (param.startsWith("-id=")) id = param.split("=")[1].toLowerCase();
          });
        
          requestBody = { path, name, id };
          endpoint = "report";
        
        } else  if (command.startsWith("mkfs")) {
          let id = "", type = "full"; 
          params.forEach(param => {
            if (param.startsWith("-id=")) id = param.split("=")[1].toLowerCase();
            if (param.startsWith("-type=")) type = param.split("=")[1].toLowerCase();
          });
        
          requestBody = { id, type};
          endpoint = "mkfs";
        
        } else if (command.startsWith("logout")) {
          requestBody = {}; 
          endpoint = "logout";      

        } else {
          results.push(`==================================\nComando no reconocido: ${command}\n==================================\n`);
          continue;
        }
  
        // Send request to the server
        const response = await fetch(`http://localhost:8080/${endpoint}`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(requestBody),
        });
  
        const text = await response.text();
        results.push(`===============================================\nComando: ${command}\nRespuesta: ${text}\n===============================================\n`);
      }
  
      // Show results in the output
      setOutput(results.join("\n"));
  
    } catch (error) {
      setOutput(`Error al ejecutar comandos: ${error.message}`);
    }
  };
  
  
  const handleFileUpload = (event) => {
    const file = event.target.files[0];
    if (file) {
      const reader = new FileReader();
      reader.onload = (e) => setInput(e.target.result);
      reader.readAsText(file);
    }
  };

  return (
    <div className="container">
      <h1>Sistema de Archivos EXT2</h1>
      <div className="textarea-container">
        <textarea
          className="input-area"
          placeholder="Ingrese comandos aquí..."
          value={input}
          onChange={(e) => setInput(e.target.value)}
        ></textarea>
        <textarea
          className="output-area"
          placeholder="Salida..."
          value={output}
          readOnly
        ></textarea>
      </div>
      <div className="buttons">
        <input type="file" accept=".smia" onChange={handleFileUpload} />
        <button onClick={handleExecute}>Ejecutar</button>
      </div>
    </div>
  );
}

export default App;