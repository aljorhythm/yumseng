const websocketScheme =
    window.location.protocol === "http:" ? "ws://" : "wss://";

const serverAddr = document.location.hostname + ":80";
const connectionWS = (() => {
    console.log("establishing ws")
    console.log(serverAddr)
    const conn: WebSocket = new WebSocket(
        websocketScheme + serverAddr + "/rooms/events"
    )
    console.log("connected ws", conn)
    return () => conn;
})();

export {connectionWS};