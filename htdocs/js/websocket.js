var sessionID = Math.floor(1000 + Math.random() * 9000);

var ws = new WebSocket("wss://"+window.location.host+"/serveWs/"+sessionID);

ws.onopen = function(evt) {
  var msg = {
    name: "presenter",
    id: 0,
    uuid: guid,
    icon: "/img/laptop.png"
  };
  ws.send(JSON.stringify(msg));
};

var refreshSlide1 = function() {
  console.log('slide1');
};

