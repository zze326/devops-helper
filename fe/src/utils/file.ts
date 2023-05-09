export function wsSendFile(file: File, ws: WebSocket) {
    const reader = new FileReader();

    // Handle file load event
    reader.onload = (event) => {
        if (event.target?.result) {
            ws.send(event.target.result as ArrayBuffer);
        }
    };

    // Read file as binary data
    reader.readAsArrayBuffer(file);
}