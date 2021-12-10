import { useEffect, useRef, useState } from "react";

import ButtonSendCheer from "./ButtonSendCheer";
import CheerVolumeBar from "./CheerVolumeBar";

import imgSrc from "../assets/face_1.jpg";

interface RoomProp {
  name: string | null;
  redSymbol: number;
  setRedSymbol: React.Dispatch<((n: number) => number) | number>;
  conn: WebSocket;
}

const randomOffset = (parentLength: number, elementLength: number): number => {
  const posOffsetBetweenBoundary =
    Math.random() * (parentLength - elementLength);
  console.log(posOffsetBetweenBoundary);
  return posOffsetBetweenBoundary;
};
interface PositionOffset {
  x: number;
  y: number;
}

const imgWidth = 100;
const imgHeight = 120;
const Room = (props: RoomProp) => {
  const { redSymbol, setRedSymbol } = props;

  const [intensity, setIntensity] = useState<number>(0);

  const [imgPos, setImgPos] = useState<PositionOffset>({ x: 0, y: 0 });
  const stageCanvasRef = useRef<HTMLDivElement>(null);

  const { conn } = props;

  const onMessageEvents = (thisConn: WebSocket) => {
    thisConn.onmessage = (msgEvent) => {
      const { data } = msgEvent;
      console.log("onmessage event data", data);
      const event = JSON.parse(data);
      const eventName = event["event_name"];
      if (eventName === "EVENT_CHEER_ADDED") {
        console.log("CHEER ADDED");
        setRedSymbol((prev) => {
          return (prev + 20) % 255;
        });
        setImgPos((prevPos) => {
          if (stageCanvasRef.current) {
            const height = stageCanvasRef.current.offsetHeight;
            const width = stageCanvasRef.current.offsetWidth;
            return {
              x: randomOffset(width, imgWidth),
              y: randomOffset(height, imgHeight),
            };
          }
          return { ...prevPos };
        });
      } else if (eventName === "EVENT_LAST_SECONDS_COUNT") {
        const { count } = event;
        setIntensity(count);
      }
    };
  };
  useEffect(() => {
    onMessageEvents(conn);
  }, [conn]);
  return (
    <>
      <div
        className="room"
        ref={stageCanvasRef}
        style={{
          display: "grid",
          border: "solid",
          height: "100%",
          width: "60%",
          alignItems: "center",
          flexFlow: "column wrap",
          justifyContent: "center",
          textAlign: "center",
          position: "relative",
        }}
      >
        <img
          src={imgSrc}
          alt=""
          style={{
            position: "absolute",
            width: `${imgWidth}px`,
            height: `${imgHeight}px`,
            top: `${imgPos.y}px`,
            left: `${imgPos.x}px`,
          }}
        ></img>

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
          <ButtonSendCheer conn={conn}></ButtonSendCheer>
        </div>
      </div>
    </>
  );
};

export default Room;
