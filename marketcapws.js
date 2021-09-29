var websocket;

async function getCurrentTab() {
  let queryOptions = { active: true, currentWindow: true };
  let [tab] = await chrome.tabs.query(queryOptions);
  return tab;
}
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

    chrome.tabs.query({ active: true, currentWindow: true }, function (tabs) {
      if (tabs.length == 0) {
        return;
      }
      chrome.tabs.sendMessage(tabs[0].id, event.data, function (response) {});
    });
  };

  //If the websocket is closed but the session is still active, create new connection again
  websocket.onclose = function () {
    // port.postMessage("close");
    websocket = undefined;
    chrome.storage.local.get(["demo_session"], function (data) {
      if (data.demo_session) {
        createWebSocketConnection();
      }
    });
  };
}

//Close the websocket connection
function closeWebSocketConnection(username) {
  if (websocket != null || websocket != undefined) {
    websocket.close();
    websocket = undefined;
  }
  e;
}

// createWebSocketConnection();