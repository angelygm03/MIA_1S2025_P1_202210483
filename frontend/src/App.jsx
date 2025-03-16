import { useState } from 'react';
import './App.css';

function App() {
  const [input, setInput] = useState('');
  const [output, setOutput] = useState('');

  const handleExecute = () => {
    // Simulación de respuesta del backend
    setOutput(`Ejecutando: \n${input}`);
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