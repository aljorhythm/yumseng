import React from "react";
import Room from "./components/Room";
import { connectionWS } from "./connections/websocket";

const getDummyRoom = (): string => {
  return "dummyRoom12345";
};

const getDummyUser = (): string => "dummyUser"


const App = () => {
  const [userId, setUserId] = React.useState<string | null>(null);
    const [roomId, setRoomId] = React.useState<string | null>(null);
    const conn = connectionWS();

    conn.onopen = (_) => {
        console.log("connection opened");
        const userDetails = JSON.stringify({
            room_name: roomId,
            user_id: userId,
        });
        console.log("send first message (user info)", userDetails);
        conn.send(userDetails);
    };


  React.useEffect(() => {
    const roomId = getDummyRoom();
    setRoomId(roomId);
    const userId = getDummyUser();
    setUserId(userId);
  }, []);
  return (
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
        <Room key={roomId} name={roomId}></Room>
      </div>
    </>
  );
};

export default App;
