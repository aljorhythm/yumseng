import { SERVER_PORT } from "../config/environment";

const websocketScheme =
  window.location.protocol === "http:" ? "ws://" : "wss://";

const serverAddr = document.location.hostname + `:${SERVER_PORT}`;
const connectionWS = (() => {
  console.log("establishing ws");
  console.log("server address: " + serverAddr);
  const conn: WebSocket = new WebSocket(
    websocketScheme + serverAddr + "/rooms/events"
  );
  console.log("initialised ws", conn);
  return () => conn;
})();

const verboseWsReadyState = (state: number): string => {
  switch (state) {
    case 0:
      return "CONNECTING";
    case 1:
      return "OPEN";
    case 2:
      return "CLOSING";
    case 3:
      return "CLOSED";
  }
  return "UNKNOWN";
};
const isWebSocketAlive = () => {
  console.log(
    "websocket ready state: " + verboseWsReadyState(connectionWS().readyState)
  );
  return connectionWS().readyState === WebSocket.OPEN;
};

export { connectionWS, isWebSocketAlive };
