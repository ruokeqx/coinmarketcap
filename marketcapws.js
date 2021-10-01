var websocket;

function createWebSocketConnection() {
    connect("wss://stream.coinmarketcap.com/price/latest");
}

//Make a websocket connection with the server.
function connect(host) {
  if (websocket === undefined) {
    // port.postMessage("open");
    websocket = new WebSocket(host);
  }

  websocket.onopen = function () {
    websocket.send(
      '{"method":"subscribe","id":"price","data":{"cryptoIds":[1,1027,1839,52,5994],"index":null}}'
    );
  };

  websocket.onmessage = function (event) {
    console.log(event.data);
  };

  //If the websocket is closed but the session is still active, create new connection again
  websocket.onclose = function () {
    websocket.send(
      '{"method":"unsubscribe","id":"unsubscribePrice"}'
    )
    websocket = undefined;
  };
}

function close(){
  websocket.send(
    '{"method":"unsubscribe","id":"unsubscribePrice"}'
  )
}
// createWebSocketConnection();