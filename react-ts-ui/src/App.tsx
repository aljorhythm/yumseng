import React, { useEffect, useState } from "react";
import Room from "./components/Room";
import { newWebSocket } from "./connections/websocket";

const getDummyRoom = (): string => {
  return "dummyRoom12345";
};

const getDummyUser = (): string => "dummyUser";

const App = () => {
  const userId = getDummyUser();
  const roomId = getDummyRoom();
  const [redSymbol, setRedSymbol] = useState<number>(0);
  const [isConnAlive, setIsConnAlive] = useState<boolean>(false);
  const [conn, setConn] = useState<WebSocket>(newWebSocket());
  const wsSocketsAndEvents = (thisConn: WebSocket) => {
    thisConn.onerror = (event) => {
      console.log("Error " + event);
      setIsConnAlive(false);
      setConn(newWebSocket());
    };
    thisConn.onopen = (_) => {
      console.log("connection opened");
      setIsConnAlive(true);
      const userDetails = JSON.stringify({
        room_name: roomId,
        user_id: userId + Number(Math.random().toFixed(3)) * 100,
      });
      console.log("send first message (user and room info)", userDetails);
      thisConn.send(userDetails);
    };
  };
  useEffect(() => {
    wsSocketsAndEvents(conn);
  }, [conn]);

  return isConnAlive ? (
    <>
      <div>
        <div
          style={{
            display: "flex",
            height: "200px",
            justifyContent: "center",
            alignItems: "center",
            textAlign: "center",
            verticalAlign: "middle",
          }}
        >
          <header
            style={{
              display: "flex",
              justifyContent: "center",
              alignItems: "center",
              height: "inherit",
              textAlign: "center",
            }}
          >
            TODO: YUMSENG HEADER: NAV / LOGO
          </header>
        </div>
      </div>
      <div
        style={{
          height: "inherit",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
        }}
      >
        <Room
          conn={conn}
          key={roomId}
          setRedSymbol={setRedSymbol}
          redSymbol={redSymbol}
          name={roomId}
        ></Room>
      </div>
    </>
  ) : (
    <div
      style={{
        height: "inherit",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
      }}
    >
      No connection to server
    </div>
  );
};

export default App;
