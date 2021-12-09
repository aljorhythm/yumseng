import React from "react";

import Button from "react-bootstrap/Button";
import { connectionWS } from "../connections/websocket";

interface CheerButtonProp {
  cheersSent: number;
}
const thisConn = connectionWS();

const ButtonSendCheer = (props: CheerButtonProp) => {
  const { cheersSent } = props;
  return (
    <Button
      style={{ height: "100px" }}
      onClick={() => {
        const cheersMessage = JSON.stringify({
          value: "yum",
          client_created_at: new Date().toJSON(),
        });

        thisConn.send(cheersMessage);
      }}
    >
      click {cheersSent}
    </Button>
  );
};

export default ButtonSendCheer;
