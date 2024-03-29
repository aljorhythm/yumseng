//todo LOW PRIORITY remove references to concrete environment `local` and `heroku`
const SERVER_PORT_HEROKU = process.env.PORT;
const SERVER_PORT_DEFAULT = process.env.DEFAULT_SERVER_PORT || 80;
const LOCAL_SERVER_WS_PORT = process.env.LOCAL_SERVER_WS_PORT || process.env.REACT_APP_LOCAL_SERVER_WS_PORT;

const SERVER_PORT =
  SERVER_PORT_HEROKU || LOCAL_SERVER_WS_PORT || SERVER_PORT_DEFAULT;
export { SERVER_PORT };
