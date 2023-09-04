import { useEffect } from "react";
import SlotBook from "./slot-book/slot-book";
import useWebSocket from 'react-use-websocket';

const WS_URL = "ws://localhost:8080/ws"; // connect to ws go server

const App = () => {
  const { lastMessage } = useWebSocket(WS_URL, {
    onOpen: () => {
      console.log('WebSocket connection established.');
    },
    onClose: () => {
      console.log('WebSocket connection closed.');
    }
  });

  useEffect(() => {
    console.log(lastMessage?.data)
  }, [lastMessage])

  return (
    <div className="App">
      <SlotBook />
    </div>
  );
}

export default App;
