import React from "react";

import { sendCheerEvent } from "../requests/DummyRequest";

import Button from "react-bootstrap/Button";

const ButtonSendCheer = (prop: any) => {
  const [cheersSent, setCheersSent] = React.useState<number>(0);
  return (
    <Button
      style={{ height: "100px" }}
      onClick={() => {
        const cheersSentNow = sendCheerEvent();
        setCheersSent(cheersSentNow);
      }}
    >
      click {cheersSent}
    </Button>
  );
};

export default ButtonSendCheer;
