import React, { useEffect, useState } from "react";

import Button from "react-bootstrap/Button";
import { verboseWsReadyState } from "../connections/websocket";
interface ButtonSendCheerProps {
  conn: WebSocket;
}
const ButtonSendCheer = (props: ButtonSendCheerProps) => {
  const [cheersSent, setCheersSent] = useState<number>(0);
  const { conn } = props;

  return (
    <Button
      style={{ height: "100px" }}
      onClick={() => {
        if (conn) {
          const cheersMessage = JSON.stringify({
            value: "yum",
            client_created_at: new Date().toJSON(),
          });
          console.log(
            "websocket state when button clicked: " + verboseWsReadyState(conn)
          );
          conn.send(cheersMessage);
          setCheersSent((prev) => prev + 1);
        }
      }}
    >
      click {cheersSent}
    </Button>
  );
};

export default ButtonSendCheer;
