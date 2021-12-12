import { SERVER_PORT } from "../config/environment";

const websocketScheme =
  window.location.protocol === "http:" ? "ws://" : "wss://";

const insecureScheme = "ws://";

const serverAddr = document.location.hostname + `:${SERVER_PORT}`;

const socketRegEndpoint = insecureScheme + serverAddr + "/rooms/events";

const verboseWsReadyState = (ws: WebSocket): string => {
  switch (ws.readyState) {
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

const newWebSocket = () => {
  console.log("Websocket target endpoint:" + socketRegEndpoint)
  const conn = new WebSocket(socketRegEndpoint);
  console.log("Websocket connection initialized: " + verboseWsReadyState(conn));
  return conn;
};
export { newWebSocket, verboseWsReadyState };
