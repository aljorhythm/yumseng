import React, { useState } from "react";

import Button from "react-bootstrap/Button";
import { connectionWS } from "../connections/websocket";

const thisConn = connectionWS();

const ButtonSendCheer = () => {
  const [cheersSent, setCheersSent] = useState<number>(0);

  return (
    <Button
      style={{ height: "100px" }}
      onClick={() => {
        const cheersMessage = JSON.stringify({
          value: "yum",
          client_created_at: new Date().toJSON(),
        });
        console.log(
          "websocket state when button clicked" + thisConn.readyState
        );
        thisConn.send(cheersMessage);
        setCheersSent((prev) => prev + 1);
      }}
    >
      click {cheersSent}
    </Button>
  );
};

export default ButtonSendCheer;
