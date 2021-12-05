import { io, Socket } from "socket.io-client";
// This module handles socket io connection.

const defaultServerPort = 80;

const connectionSIO = (() => {
  let conn: Socket | null = null;
  return () => conn || io("http://localhost:" + defaultServerPort);
})();

export { connectionSIO };
