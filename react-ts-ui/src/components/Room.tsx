import { useEffect, useState } from "react";
import { connectionWS } from "../connections/websocket";

import ButtonSendCheer from "./ButtonSendCheer";
import CheerVolumeBar from "./CheerVolumeBar";

interface RoomProp {
  name: string | null;
  cheersSent: number;
  redSymbol: number;
  setRedSymbol: React.Dispatch<((n: number) => number) | number>;
  setCheersSent: React.Dispatch<((n: number) => number) | number>;
}

const thisConn = connectionWS();

const Room = (props: RoomProp) => {
  const { cheersSent, setCheersSent, redSymbol, setRedSymbol } = props;

  const [intensity, setIntensity] = useState<number>(0);
  useEffect(() => {
    if (thisConn) {
      thisConn.onmessage = (msgEvent) => {
        const { data } = msgEvent;
        console.log("onmessage event data", data);
        const event = JSON.parse(data);
        const eventName = event["event_name"];
        if (eventName == "EVENT_CHEER_ADDED") {
          console.log("CHEER ADDED");
          setRedSymbol((prev) => {
            return (prev + 20) % 255;
          });
          setCheersSent((prev) => prev + 1);
        } else if (eventName == "EVENT_LAST_SECONDS_COUNT") {
          const { count } = event;
          setIntensity(count);
        }
      };
    }
  }, []);
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
        <CheerVolumeBar intensity={intensity}></CheerVolumeBar>
        <div
          style={{
            display: "inherit",
            height: "100%",
            verticalAlign: "middle",
          }}
        >
          <ButtonSendCheer cheersSent={cheersSent}></ButtonSendCheer>
        </div>
      </div>
    </>
  );
};

export default Room;
