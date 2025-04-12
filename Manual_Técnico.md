# Manual Técnico File System EXT2
## Manejo e Implementación de Archivos
### Primer Semestre 2025
```js
Universidad San Carlos de Guatemala
Facultad de Ingeniería
Nombre: Angely Lucrecia García Martínez
Carne: 202210483
```
---
## Descripción del Proyecto
El sistema de archivos EXT2 está diseñado como una aplicación web que utiliza una arquitectura cliente-servidor. Está compuesto por dos módulos principales: Frontend  y Backend, que trabajan en conjunto para procesar comandos y simular el comportamiento de un sistema de archivos EXT2.

## Objetivo
El manual técnico proporciona una guía detallada para desarrolladores, administradores de sistemas y otros usuarios técnicos sobre cómo instalar, configurar y utilizar el Sistema de Archivos EXT2. Contiene información sobre la estructura del código, las dependencias del proyecto, las instrucciones de implementación, y la explicación de los diferentes componentes y su funcionamiento.

El manual técnico está diseñado para ofrecer una referencia completa para aquellos que deseen comprender el funcionamiento interno del sistema, realizar personalizaciones o contribuir al desarrollo continuo del proyecto.

## Características Principales
* **Simulación del Sistema de Archivos EXT2**
    * Implementación de estructuras como MBR, EBR, Superbloques, Inodos y Bloques.
    * Soporte para operaciones como creación de discos, particiones, montaje, formateo y manejo de archivos.

* **Interfaz Gráfica Intuitiva**
    * Permite a los usuarios ingresar comandos, cargar archivos .smia y visualizar los resultados.

* **Gestión de Usuarios y Grupos**
    * Comandos para crear, eliminar y gestionar usuarios y grupos.
    * Soporte para inicio y cierre de sesión.

* **Generación de Reportes**
    * Creación de reportes visuales en formato Graphviz para estructuras como MBR, particiones y sistema de archivos.

* **Arquitectura Cliente-Servidor**
    * Comunicación entre el frontend y backend mediante solicitudes HTTP.

## Tecnologías utilizadas
* **Frontend**
    * React + Vite: Para la creación de la interfaz gráfica.
    * HTML/CSS: Para el diseño y estilo de la aplicación.

* **Backend**
    * Golang: Para la implementación de la lógica del sistema de archivos y manejo de estructuras internas.
    * HTTP Server: Para exponer los endpoints que procesan los comandos

* **Visualización de Reportes**
    * Graphviz: Para generar reportes visuales en formato .dot y convertirlos a imágenes.

* **Almacenamiento**
    * Archivos binarios .mia para simular discos y particiones.

## Requisitos del sistema
1.  **Hardware**
    * Procesador: Intel Core i3 (de preferencia superior)
    * Memoria RAM: 8GB o más
    * Espacio en disco: Al menos 100 MB para la aplicación y almacenamiento de discos simulados.

1.  **Software**
    * Sistema operativo: Linux (Ubuntu recomendado)
    * Node.js: Para ejecutar el frontend
    * Golang: Para compilar y ejecutar el backend
    * Graphviz: Para la generación de reportes

## Instalación
Para poder ejecutar la aplicación, se debe clonar o descargar el repositorio del proyecto desde [Github](https://github.com/angelygm03/MIA_1S2025_P1_202210483.git) para luego abrir el proyecto en el IDE de su preferencia (se recomienda Visual Studio Code) y correr la aplicación.

## Bibliotecas utilizadas

#### **Frontend**:
- **React.js**: Para la creación de componentes y manejo del estado.
- **Axios**: Para realizar solicitudes HTTP al backend.

#### **Backend**:
- **net/http**: Para crear el servidor HTTP y manejar solicitudes.
- **encoding/json**: Para procesar datos en formato JSON.
- **os**: Para operaciones con archivos y directorios.
- **fmt**: Para formateo y salida de texto.
- **strings**: Para manipulación de cadenas.
- **flag**: Para el manejo de parámetros en comandos.

#### **Generación de Reportes**:
- **Graphviz**: Para la creación de archivos `.dot` y generación de imágenes.


## Uso del código
El archivo main.go contiene el script principal que debe ejecutarse para ejecutar el backend. En una terminal se ejecuta el comando go run main.go, para verificar su funcionamiento, se puede acceder al navegador con http://localhost:8080. Para levantar el frontend se debe tener instalado Node.js y npm, se navega hasta el archivo App.jsx y con el comando npm run dev se inicia el servidor, para verificarlo se accede a http://localhost:3000.

[![image-2025-04-02-132444799.png](https://i.postimg.cc/66htVRCh/image-2025-04-02-132444799.png)](https://postimg.cc/Bttr4jY8)

## Estructuras y Comandos
El sistema utiliza varias estructuras de datos fundamentales para simular el sistema de archivos EXT2. Estas estructuras se almacenan y gestionan dentro de un archivo binario `.mia`.

#### **MBR (Master Boot Record)**
- **Descripción:** Contiene información sobre el disco y sus particiones.
- **Campos:**
  - `MbrSize`: Tamaño total del disco.
  - `CreationDate`: Fecha de creación del disco.
  - `Signature`: Identificador único del disco.
  - `Fit`: Algoritmo de ajuste para particiones (`bf`, `ff`, `wf`).
  - `Partitions`: Arreglo de hasta 4 particiones (primarias o extendidas).
- **Función:** Define la estructura general del disco y organiza las particiones.
- **Código de implementación:** 

[![image-2025-04-02-152835379.png](https://i.postimg.cc/QCGjdJL2/image-2025-04-02-152835379.png)](https://postimg.cc/2bwsHvSG)


#### **EBR (Extended Boot Record)**
- **Descripción:** Utilizado para gestionar particiones lógicas dentro de una partición extendida.
- **Campos:**
  - `PartStart`: Inicio de la partición lógica.
  - `PartSize`: Tamaño de la partición lógica.
  - `PartNext`: Dirección del siguiente EBR (o `-1` si es el último).
  - `PartName`: Nombre de la partición lógica.
- **Función:** Permite la creación de múltiples particiones lógicas dentro de una partición extendida.
- **Código de implementación:** 

[![image-2025-04-02-152942621.png](https://i.postimg.cc/BbBfQWBF/image-2025-04-02-152942621.png)](https://postimg.cc/sGxNwN8f)

#### **Superblock**
- **Descripción:** Contiene información sobre el sistema de archivos.
- **Campos:**
  - `S_inodes_count`: Número total de inodos.
  - `S_blocks_count`: Número total de bloques.
  - `S_free_blocks_count`: Bloques libres.
  - `S_free_inodes_count`: Inodos libres.
  - `S_inode_start`: Inicio de la tabla de inodos.
  - `S_block_start`: Inicio de la tabla de bloques.
- **Función:** Gestiona el estado del sistema de archivos.
- **Código de implementación:** 

[![image-2025-04-02-153126549.png](https://i.postimg.cc/K8r3S4Tc/image-2025-04-02-153126549.png)](https://postimg.cc/9D0Qym1s)

#### **Inodo**
- **Descripción:** Representa un archivo o carpeta.
- **Campos:**
  - `I_uid`: ID del usuario propietario.
  - `I_gid`: ID del grupo propietario.
  - `I_size`: Tamaño del archivo.
  - `I_block`: Punteros a los bloques de datos.
  - `I_type`: Tipo (`0` para carpeta, `1` para archivo).
- **Función:** Almacena metadatos y punteros a los bloques de datos.
- **Código de implementación:** 

[![image-2025-04-02-153150607.png](https://i.postimg.cc/250qyJTP/image-2025-04-02-153150607.png)](https://postimg.cc/8FrknZ2H)

#### **Bloques**
- **Tipos:**
  - **Folderblock:** Almacena referencias a otros inodos (carpetas o archivos).
  - **Fileblock:** Almacena contenido de archivos.
  - **Pointerblock:** Almacena punteros a otros bloques.
- **Función:** Gestionan el almacenamiento de datos y la estructura jerárquica del sistema de archivos.
[![image-2025-04-02-153322257.png](https://i.postimg.cc/3wNmHmHj/image-2025-04-02-153322257.png)](https://postimg.cc/CZyzGZ5d)
---

### **Descripción de los Comandos Implementados**

A continuación, se describen los comandos disponibles en el sistema, junto con ejemplos de uso, código de implementación y sus efectos o resultados.

#### **1. MKDISK**
- **Descripción:** Crea un disco virtual.
- **Parámetros:**
  - `-size`: Tamaño del disco (en KB o MB).
  - `-unit`: Unidad (`k` o `m`).
  - `-fit`: Algoritmo de ajuste (`bf`, `ff`, `wf`).
  - `-path`: Ruta donde se creará el disco.
- **Ejemplo:**
  ```bash
  mkdisk -size=1024 -unit=m -fit=ff -path="/home/user/disk1.mia"
  ```
- **Efecto:** Crea un archivo binario con un MBR inicializado.
- **Código de implementación:** 

[![image-2025-04-02-151746201.png](https://i.postimg.cc/x8c1X1qm/image-2025-04-02-151746201.png)](https://postimg.cc/gLbW5dVc)

#### **2. RMDISK**
- **Descripción:** Elimina un disco virtual.
- **Parámetros:**
  - `-path`: Ruta del disco a eliminar.
- **Ejemplo:**
  ```bash
  rmdisk -path="/home/user/disk1.mia"
  ```
- **Efecto:** Elimina el archivo binario correspondiente al disco.

#### **3. FDISK**
- **Descripción:** Crea, elimina o modifica particiones.
- **Parámetros:**
  - `-size`: Tamaño de la partición.
  - `-unit`: Unidad (`k`, `m` o `b`).
  - `-type`: Tipo de partición (`p`, `e`, `l`).
  - `-fit`: Algoritmo de ajuste (`bf`, `ff`, `wf`).
  - `-path`: Ruta del disco.
  - `-name`: Nombre de la partición.
- **Ejemplo:**
  ```bash
  fdisk -size=500 -unit=m -type=p -fit=bf -path="/home/user/disk1.mia" -name="part1"
  ```
- **Efecto:** Crea una partición en el disco especificado.
- **Código de implementación:** 

[![image-2025-04-02-152111222.png](https://i.postimg.cc/s2mjKhhS/image-2025-04-02-152111222.png)](https://postimg.cc/SXXbxjVN)

#### **4. MOUNT**
- **Descripción:** Monta una partición.
- **Parámetros:**
  - `-path`: Ruta del disco.
  - `-name`: Nombre de la partición.
- **Ejemplo:**
  ```bash
  mount -path="/home/user/disk1.mia" -name="part1"
  ```
- **Efecto:** Asigna un ID único a la partición y la marca como montada.
- **Código de implementación:** 

[![image-2025-04-02-152243937.png](https://i.postimg.cc/Bbr4t5rW/image-2025-04-02-152243937.png)](https://postimg.cc/q6x95KBQ)

#### **5. MKFS**
- **Descripción:** Formatea una partición montada.
- **Parámetros:**
  - `-id`: ID de la partición.
  - `-type`: Tipo de formato (`full` o `fast`).
- **Ejemplo:**
  ```bash
  mkfs -id=4831a -type=full
  ```
- **Efecto:** Inicializa el sistema de archivos EXT2 en la partición.
- **Código de implementación:** 

[![image-2025-04-02-152415443.png](https://i.postimg.cc/xdfYk14d/image-2025-04-02-152415443.png)](https://postimg.cc/CRQ9XSYy)

#### **6. LOGIN**
- **Descripción:** Inicia sesión en el sistema.
- **Parámetros:**
  - `-user`: Nombre de usuario.
  - `-pass`: Contraseña.
  - `-id`: ID de la partición.
- **Ejemplo:**
  ```bash
  login -user=root -pass=123 -id=4831a
  ```
- **Efecto:** Marca la partición como activa y permite ejecutar comandos relacionados con usuarios.
- **Código de implementación:** 

[![image-2025-04-02-152624109.png](https://i.postimg.cc/jjS5pcKv/image-2025-04-02-152624109.png)](https://postimg.cc/rD7MWSDt)

#### **7. REP**
- **Descripción:** Genera reportes del sistema.
- **Parámetros:**
  - `-name`: Tipo de reporte (`mbr`, `disk`, etc.).
  - `-path`: Ruta donde se generará el reporte.
  - `-id`: ID de la partición.
- **Ejemplo:**
  ```bash
  rep -name=mbr -path="/home/user/mbr_report.jpg" -id=4831a
  ```
- **Efecto:** Genera un archivo gráfico con la información solicitada.

#### **8. MKGRP (Crear Grupo)**

- **Descripción**: Crea un nuevo grupo en el sistema.
- **Parámetros**:
  - `-name`: Nombre del grupo a crear.
- **Ejemplo**:
  ```bash
  mkgrp -name=developers
  ```
- **Efecto**:
  - Se agrega un nuevo grupo al archivo `users.txt` en el sistema de archivos.
  - El grupo se registra con un ID único y su nombre.
- **Notas**:
  - Solo el usuario `root` puede ejecutar este comando.
  - Si el grupo ya existe, se devuelve un error.
- **Código de implementación:** 

[![image-2025-04-02-154116596.png](https://i.postimg.cc/v84y6BGP/image-2025-04-02-154116596.png)](https://postimg.cc/xqVhw9JH)


#### **9. CAT (Leer Archivo)**

- **Descripción**: Muestra el contenido de un archivo en el sistema de archivos.
- **Parámetros**:
  - `-file`: Ruta del archivo a leer.
- **Ejemplo**:
  ```bash
  cat -file=/home/docs/manual.txt
  ```
- **Efecto**:
  - Se busca el archivo en el sistema de archivos utilizando su ruta.
  - Si el archivo existe, se muestra su contenido en la salida.
- **Notas**:
  - El usuario debe tener permisos de lectura sobre el archivo.
  - Si el archivo no existe, se devuelve un error.
- **Código de implementación:** 

[![image-2025-04-02-154228650.png](https://i.postimg.cc/05KgyBCF/image-2025-04-02-154228650.png)](https://postimg.cc/qgrDmQSx)


#### **10. RMGRP (Eliminar Grupo)**

- **Descripción**: Elimina un grupo existente del sistema.
- **Parámetros**:
  - `-name`: Nombre del grupo a eliminar.
- **Ejemplo**:
  ```bash
  rmgrp -name=developers
  ```
- **Efecto**:
  - Se elimina el grupo del archivo `users.txt`.
  - Los usuarios asociados al grupo no se eliminan, pero su grupo queda desasignado.
- **Notas**:
  - Solo el usuario `root` puede ejecutar este comando.
  - No se puede eliminar el grupo `root`.

#### **11. MKUSR (Crear Usuario)**

- **Descripción**: Crea un nuevo usuario en el sistema.
- **Parámetros**:
  - `-user`: Nombre del usuario.
  - `-pass`: Contraseña del usuario.
  - `-grp`: Nombre del grupo al que pertenece el usuario.
- **Ejemplo**:
  ```bash
  mkusr -user=user1 -pass=12345 -grp=developers
  ```
- **Efecto**:
  - Se agrega un nuevo usuario al archivo `users.txt` con su nombre, contraseña y grupo asociado.
- **Notas**:
  - Solo el usuario `root` puede ejecutar este comando.
  - El grupo especificado debe existir previamente.
  - Si el usuario ya existe, se devuelve un error.
- **Código de implementación:** 

[![image-2025-04-02-154450038.png](https://i.postimg.cc/wMLbcnfB/image-2025-04-02-154450038.png)](https://postimg.cc/V0sRzZcy)


#### **12. CHGRP (Cambiar Grupo de Usuario)**

- **Descripción**: Cambia el grupo al que pertenece un usuario.
- **Parámetros**:
  - `-user`: Nombre del usuario.
  - `-grp`: Nombre del nuevo grupo.
- **Ejemplo**:
  ```bash
  chgrp -user=user1 -grp=admins
  ```
- **Efecto**:
  - Se actualiza el grupo del usuario en el archivo `users.txt`.
- **Notas**:
  - Solo el usuario `root` puede ejecutar este comando.
  - El grupo especificado debe existir previamente.
  - Si el usuario no existe, se devuelve un error.
- **Código de implementación:** 

[![image-2025-04-02-154537506.png](https://i.postimg.cc/h4dwB7Dq/image-2025-04-02-154537506.png)](https://postimg.cc/JtRK30HT)


#### **13. RMUSR (Eliminar Usuario)**

- **Descripción**: Elimina un usuario del sistema.
- **Parámetros**:
  - `-user`: Nombre del usuario a eliminar.
- **Ejemplo**:
  ```bash
  rmusr -user=john
  ```
- **Efecto**:
  - Se elimina el usuario del archivo `users.txt`.
  - Los archivos y carpetas creados por el usuario permanecen en el sistema, pero su propietario queda desasignado.
- **Notas**:
  - Solo el usuario `root` puede ejecutar este comando.
  - No se puede eliminar el usuario `root`.


## Manejo de Errores

En esta sección se describen los posibles errores que pueden surgir durante la instalación y ejecución del sistema, así como las soluciones recomendadas. También se aborda el manejo de excepciones en el código.

---

### **1. Instalación de Bibliotecas**

##### **Error: Falta de Dependencias en el Frontend**
- **Problema**: Al ejecutar `npm install`, pueden surgir errores relacionados con la falta de dependencias o versiones incompatibles.
- **Solución**:
  1. Asegurarse de tener una versión actualizada de **Node.js** y **npm**:
     ```bash
     node -v
     npm -v
     ```
     Si no están actualizados, instalar la última versión:
     ```bash
     sudo apt update
     sudo apt install nodejs npm
     ```
  2. Eliminar la carpeta `node_modules` y el archivo `package-lock.json`:
     ```bash
     rm -rf node_modules package-lock.json
     ```
  3. Reinstalar las dependencias:
     ```bash
     npm install
     ```

---

### **2. Error de Importación de Graphviz**

##### **Error: `dot` Command Not Found**
- **Problema**: Al generar reportes, el sistema puede devolver un error indicando que el comando `dot` no está disponible.
- **Causa**: Graphviz no está instalado en el sistema.
- **Solución**:
  1. Instalar Graphviz en el sistema:
     ```bash
     sudo apt-get install graphviz
     ```
  2. Verificar que el comando `dot` esté disponible:
     ```bash
     dot -V
     ```
     Esto debería devolver la versión instalada de Graphviz.

---

### **3. Manejo de Excepciones en el Código**

El sistema implementa manejo de errores en diferentes niveles para garantizar la estabilidad y la correcta ejecución de las operaciones.

##### **3.1. Backend (Golang)**

1. **Errores en la Lectura de Archivos**:
   - **Problema**: El archivo binario no existe o no se puede abrir.
   - **Solución en el Código**:
     ```go
     file, err := os.Open("disk.mia")
     if err != nil {
         fmt.Println("Error al abrir el archivo:", err)
         return
     }
     defer file.Close()
     ```

2. **Errores en la Decodificación de JSON**:
   - **Problema**: El cliente envía datos mal formateados.
   - **Solución en el Código**:
     ```go
     var req LoginRequest
     if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
         http.Error(w, "Solicitud inválida: formato JSON incorrecto", http.StatusBadRequest)
         return
     }
     ```

3. **Errores en Comandos del Sistema**:
   - **Problema**: Fallo al ejecutar un comando externo como `dot`.
   - **Solución en el Código**:
     ```go
     cmd := exec.Command("dot", "-Tpng", "-o", outputPath, inputPath)
     err := cmd.Run()
     if err != nil {
         fmt.Println("Error al ejecutar Graphviz:", err)
         return
     }
     ```

4. **Errores en Operaciones de Disco**:
   - **Problema**: Intento de acceder a una partición no montada.
   - **Solución en el Código**:
     ```go
     if !DiskControl.IsPartitionMounted(partitionID) {
         return "Error: La partición no está montada"
     }
     ```

##### **3.2. Frontend (React.js)**

1. **Errores en Solicitudes HTTP**:
   - **Problema**: El backend no responde o devuelve un error.
   - **Solución en el Código**:
     ```jsx
     const handleExecute = async () => {
       try {
         const response = await axios.post("http://localhost:8080/execute", { command: input });
         setOutput(response.data);
       } catch (error) {
         setOutput(`Error al ejecutar el comando: ${error.message}`);
       }
     };
     ```

2. **Errores en la Carga de Archivos**:
   - **Problema**: El usuario intenta cargar un archivo no válido.
   - **Solución en el Código**:
     ```jsx
     const handleFileUpload = (event) => {
       const file = event.target.files[0];
       if (!file.name.endsWith(".smia")) {
         setOutput("Error: Solo se permiten archivos con extensión .smia");
         return;
       }
       // Procesar el archivo...
     };
     ```

---

### **4. Recomendaciones Generales**

1. **Validación de Entradas**:
   - Asegurar la validación de todas las entradas del usuario tanto en el frontend como en el backend para evitar errores inesperados.

2. **Logs de Errores**:
   - Implementa un sistema de logs para registrar errores en el backend. Esto facilita la depuración:
     ```go
     log.Printf("Error: %v\n", err)
     ```






