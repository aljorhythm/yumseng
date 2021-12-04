interface RoomProp {
  name: string | null;
}

const Room = (prop: RoomProp) => {
  return (
    <div
      style={{
        display: "flex",
        verticalAlign: "middle",
        justifyContent: "center",
        alignSelf: "center",
      }}
    >
      {prop.name}
    </div>
  );
};

export default Room;
