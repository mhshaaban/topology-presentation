var ws;
var state = 'initial';
window.onload = function() {
  var btnConnect = document.getElementById("btnConnect");
  btnConnect.onclick = function () {
    var name = document.getElementById('name').value;
    if (name === null || name === "") {
      alert("Please enter your name");
      return;
    }
    ws = new WebSocket('wss://' + window.location.host + '/phone');
    ws.onopen = function() {
      var msg = {
        name: document.getElementById('name').value,
        device: navigator.userAgent,
        state: "initial",
        date: Date.now()
      };
      ws.send(JSON.stringify(msg));
    };

    // Write message on receive
    ws.onmessage = function(e) {
      console.log("Received:",e);
      //document.getElementById('output').innerHTML += "Received: " + e.data + "<br>";
      var obj = JSON.parse(e.data);
      switch (obj.state) {
        case "runnable":
          document.body.style.background = 'green';
          state = 'runnable';
          break;
        case "notrunnable":
          document.body.style.background = 'red';
          state = 'notrunnable';
          break;
        case "connected":
          document.body.style.background = 'blue';
          break;
        case "running":
          document.getElementById('heart').style.visibility='visible';
          break;
        case "stopped":
          document.getElementById('heart').style.visibility='hidden';
          break;
        case "autonomous":
          document.getElementById('btn-start').style.visibility='visible';
          document.getElementById('btn-stop').style.visibility='visible';
          break;
        case "conducted":
          document.getElementById('btn-start').style.visibility='hidden';
          document.getElementById('btn-stop').style.visibility='hidden';
          break;
      }
    };
    ws.onclose = function() {
      document.body.style.background = 'white';
    };
  };
};

function changeRequest(req) {
  switch (req) {
    case "start":
      if (state == "runnable") {
        document.getElementById('heart').style.visibility='visible';
      }
      break;
    case "stop":
      document.getElementById('heart').style.visibility='hidden';
  } 
  senddata(state);
}
function readDeviceOrientation() {

  if (Math.abs(window.orientation) === 90) {
    // Landscape
    senddata("stop"); 
  } else {
    // Portrait
    senddata("start"); 

  }

}

window.onorientationchange = readDeviceOrientation;
function senddata(state) {
  // Construct a msg object containing the data the server needs to process the message from the chat client.
  if (ws !== null) {
    var msg = {
      name: document.getElementById('name').value,
      state: state,
      date: Date.now()
    };

    ws.send(JSON.stringify(msg));
    console.log("Sending:",msg);
    //document.getElementById('output').innerHTML += "Sent: " + JSON.stringify(msg) + "<br>";
  }
}
