import React from "react";
import NavBar from "../components/NavBar";
import Consola from "../components/Editor";
import Service from "../Services/Service";
import 'bootstrap/dist/css/bootstrap.min.css';

function Index() {
    const [value, setValue] = React.useState("");
    const [response, setResponse] = React.useState("");

    const changeText = (text) => {
        console.log("Texto recibido en changeText:", text); // Verifica si se está actualizando correctamente
    
        setValue(text);
    };

    const handlerClick = async () => {
        //console.log("Valor de value al hacer clic en enviar:", value); // Agregar este log

        if (value === "") {
            alert("NO PUEDES ENVIAR UN COMANDO VACIO");
            return;
        }
        Service.analisis(value)
            .then((res) => {
                setResponse(res.respuesta);
            })
            .catch((err) => {
                console.error(err);
            });
    };

    const handlerLimpiar = () => {
        if (value === "") {
            alert("NO PUEDES LIMPIAR UN CAMPO VACIO");
            return;
        }
        changeText("");
        setResponse("");
    };

    const handleLoadClick = () => {
        const input = document.createElement('input');
        input.type = 'file';
        input.addEventListener('change', handleFileChange);
        input.click();
    };

    const handleFileChange = (e) => {
        const file = e.target.files[0];
        const reader = new FileReader();
        reader.onload = (e) => {
            const text = e.target.result;
            console.log("Contenido del archivo cargado:", text);
            changeText(text);
        };
        reader.readAsText(file);
    };

    return (
        <>
            <NavBar />
            <div className="container mt-4">
                <h1 className="text-center mb-4">PROYECTO 1</h1>
                <div className="card p-4 mb-4">
                    <Consola 
                        text="CONSOLA DE ENTRADA" 
                        handlerChange={changeText} 
                        value={value} 
                        className="form-control mb-3"
                    />
                </div>
                <div className="d-flex justify-content-around my-3">
                    <button type="button" className="btn btn-primary" onClick={handlerClick}>Enviar</button>
                    <button type="button" className="btn btn-secondary" onClick={handlerLimpiar}>Limpiar</button>
                    <button type="button" className="btn btn-success" onClick={handleLoadClick}>Cargar</button>
                </div>
                <div className="card p-4">
                    <Consola 
                        text="CONSOLA DE SALIDA" 
                        handlerChange={changeText} 
                        value={response} 
                        readOnly={false} 
                        className="form-control"
                    />
                </div>
            </div>
        </>
    );
}

export default Index;
