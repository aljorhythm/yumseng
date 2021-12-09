import ButtonSendCheer from "./ButtonSendCheer";
import { useEffect } from "react";
interface RoomProp {
  name: string | null;
  conn: WebSocket | null;
  cheersSent: number;
  redSymbol: number;
  setRedSymbol: React.Dispatch<((n: number) => number) | number>;
  setCheersSent: React.Dispatch<((n: number) => number) | number>;
}

const Room = (props: RoomProp) => {
  const { cheersSent, setCheersSent, redSymbol, setRedSymbol } = props;
  useEffect(() => {
    const { conn: thisConn } = props;
    if (thisConn) {
      thisConn.onmessage = (msgEvent) => {
        const { data } = msgEvent;
        console.log("onmessage event data", data);
        const event = JSON.parse(data);
        const eventName = event["event_name"];
        if (eventName == "EVENT_CHEER_ADDED") {
          console.log("CHEER ADDED");
          console.log("redSymbol", redSymbol);
          setRedSymbol((prev) => {
            console.log("prev", prev);
            return (prev + 20) % 255;
          });
          setCheersSent((prev) => prev + 1);
        }
      };
    }
  }, [redSymbol]);
  return (
    <>
      <div
        style={{
          display: "grid",
          border: "solid",
          height: "100%",
          width: "100%",
          alignItems: "center",
          flexFlow: "column wrap",
          justifyContent: "center",
          textAlign: "center",
        }}
      >
        <div
          style={{
            display: "inherit",
            verticalAlign: "top",
            flex: "1",
            color: `rgb(${redSymbol}, ${redSymbol}, 0)`,
            backgroundColor: `rgb(${255 - redSymbol}, ${255 - redSymbol}, 255)`,
            transition: "background-color 1000ms linear, color 1s linear",
            opacity: "1",
          }}
        >
          {props.name}
        </div>

        <div
          style={{
            display: "inherit",
            height: "100%",
            verticalAlign: "middle",
          }}
        >
          <ButtonSendCheer
            cheersSent={cheersSent}
            conn={props.conn}
          ></ButtonSendCheer>
        </div>
      </div>
    </>
  );
};

export default Room;
