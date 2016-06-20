var f_guid = function() {
  'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
    var r = Math.random()*16|0, v = c == 'x' ? r : (r&0x3|0x8);
    return v.toString(16);
  });
};

var guid = f_guid();

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

