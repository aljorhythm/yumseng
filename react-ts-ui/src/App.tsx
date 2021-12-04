import React from "react";
import "./App.css";
import Room from "./components/Room";

const getDummyRoom = (): string => {
  return "dummyRoom12345";
};
const App = () => {
  const [profile, setProfile] = React.useState({
    user: "dummyUser",
    sessionId: "uniqueDummySession",
  });
  const [roomId, setRoomId] = React.useState<string | null>(null);
  React.useEffect(() => {
    const roomId = getDummyRoom();
    setRoomId(roomId);
  }, []);
  return (
    <div
      className="vp-cheers-tgt"
      style={{
        height: "inherit",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
      }}
    >
      <Room key={roomId} name={roomId}></Room>
    </div>
  );
};

export default App;
