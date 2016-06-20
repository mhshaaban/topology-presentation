var sessionID = Math.floor(1000 + Math.random() * 9000);

var ws = new WebSocket("wss://"+window.location.host+"/slideJoin");

ws.onopen = function(evt) {
  var msg = {
    id: sessionID
  };
  ws.send(JSON.stringify(msg));
};

var refreshSlide1 = function() {
  console.log('slide1');
};

