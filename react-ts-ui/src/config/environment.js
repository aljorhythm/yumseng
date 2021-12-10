const SERVER_PORT_HEROKU = process.env.PORT;
const SERVER_PORT_DEFAULT = process.env.DEFAULT_SERVER_PORT;
const SERVER_PORT_REACT_APP = process.env.REACT_APP_SERVER_PORT;

const SERVER_PORT =
  SERVER_PORT_HEROKU || SERVER_PORT_DEFAULT || SERVER_PORT_REACT_APP || 80;
export { SERVER_PORT };