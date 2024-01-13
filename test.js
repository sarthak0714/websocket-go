var socket = new WebSocket("ws://localhost:3000/chat");

socket.onerror = function(event) {
  console.error("WebSocket error observed:", event);
};

socket.onopen = function(event) {
  console.log("WebSocket is open now.");
};

socket.onclose = function(event) {
  console.log("WebSocket is closed now.");
};
