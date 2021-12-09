import React from "react";

import Button from "react-bootstrap/Button";

interface CheerButtonProp {
  conn: WebSocket | null;
  cheersSent: number;
}
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
        console.log("sending message via connection", props.conn);
        props.conn?.send(cheersMessage);
      }}
    >
      click {cheersSent}
    </Button>
  );
};

export default ButtonSendCheer;
