const DummyRequest = () => {};

const dummyCheersSent = (() => {
  let count = 0;
  return () => {
    count = count += 1;
    return count;
  };
})();
const sendCheerEvent = () => {
  return dummyCheersSent();
};
export { DummyRequest, sendCheerEvent };
