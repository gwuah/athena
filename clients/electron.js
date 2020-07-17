const WebSocket = require("ws");

ws = new WebSocket("ws://localhost:8080/electron-ws/250");

ws.on("open", (e) => {
  setInterval(() => {
    ws.send("Abeka");
  }, 1000);
});

ws.on("message", function (data) {
  console.log(data);
});
