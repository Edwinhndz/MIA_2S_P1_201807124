import axios from 'axios';

const instance = axios.create({
    baseURL: 'http://localhost:3000/',
    timeout: 15000,
    headers: {
        'Content-Type': 'application/json',
    },
});

export const analisis = async (value) => {
    console.log("Valor de comando antes de enviar:", value); // Verifica que esto sea una cadena de texto simple
    const payload = { comando: value.trim() };
    console.log("Payload enviado al backend:", payload);
    const { data } = await instance.post('command', payload);
    return data;
}
