import React from "react";
import ReactDOM from "react-dom";

// font
import "../css/root.css";

const EntryJSX = () => <div>Hi</div>;
// react
// ReactDOM.render(<EntryJSX />, document.getElementById("ReactEntry"));
ReactDOM.render((<><EntryJSX /><EntryJSX /></>), document.getElementById("react-entry"));
