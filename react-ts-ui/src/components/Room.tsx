import ButtonSendCheer from "./ButtonSendCheer";
interface RoomProp {
  name: string | null;
}

const Room = (prop: RoomProp) => {
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
        <div style={{ display: "inherit", verticalAlign: "top", flex: "1" }}>
          {prop.name}
        </div>

        <div
          style={{
            display: "inherit",
            height: "100%",
            verticalAlign: "middle",
          }}
        >
          <ButtonSendCheer></ButtonSendCheer>
        </div>
      </div>
    </>
  );
};

export default Room;
