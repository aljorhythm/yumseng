import React from "react";
import Room from "./components/Room";
import { connectionSIO } from "./connections/sio";
const getDummyRoom = (): string => {
  return "dummyRoom12345";
};
const App = () => {
  const [profile, setProfile] = React.useState<Object | null>({
    user: "dummyUser",
    sessionId: "uniqueDummySession",
  });
  const conn = connectionSIO();
  const [roomId, setRoomId] = React.useState<string | null>(null);
  React.useEffect(() => {
    const roomId = getDummyRoom();
    setRoomId(roomId);
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
