import React, { useEffect, useState } from "react";
import Room from "./components/Room";
import {
  connectionWS,
  isWebSocketAlive,
  checkSocketReadyStateInterval,
} from "./connections/websocket";

const getDummyRoom = (): string => {
  return "dummyRoom12345";
};

const getDummyUser = (): string => "dummyUser";

const thisConn = connectionWS();
checkSocketReadyStateInterval();
const App = () => {
  const userId = getDummyUser();
  const roomId = getDummyRoom();
  const [redSymbol, setRedSymbol] = useState<number>(0);
  const [isConnAlive, setIsConnAlive] = useState<boolean>(false);

  const wsSocketsAndEvents = () => {
    thisConn.onerror = (event) => {
      console.log("Error " + event);
    };
    thisConn.onopen = (_) => {
      console.log("connection opened");
      setIsConnAlive(isWebSocketAlive());
      const userDetails = JSON.stringify({
        room_name: roomId,
        user_id: userId,
      });
      console.log("send first message (user and room info)", userDetails);
      thisConn.send(userDetails);
    };
  };
  useEffect(wsSocketsAndEvents, [roomId, userId]);

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
