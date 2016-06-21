var ws;
var state = 'initial';
function generateUUID(){
  var d = new Date().getTime();
  if(window.performance && typeof window.performance.now === "function"){
    d += performance.now(); //use high-precision timer if available
  }
  var uuid = 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
    var r = (d + Math.random()*16)%16 | 0;
    d = Math.floor(d/16);
    return (c=='x' ? r : (r&0x3|0x8)).toString(16);
  });
  return uuid;
}
var guid = generateUUID();
console.log(guid);

var status = "initial";

var id = 0;


window.onload = function() {
  var btnConnect = document.getElementById("btnConnect");
  var btnDisConnect = document.getElementById("btnDisConnect");
  var btnCreate = document.getElementById("btnCreate");
  var btnDelete = document.getElementById("btnDelete");
  var btnStop = document.getElementById("btnStop");
  var btnStart = document.getElementById("btnStart");
  var btnConfigure = document.getElementById("btnConfigure");

  btnCreate.onclick = function() {
    $("#heart").removeClass("hidden");
    status = "created";
    senddata(status);
  };
  btnConfigure.onclick = function() {
    if (status == "created") {
      $("#heartimg").removeClass("sepia");
      status = "configured";
      senddata(status);
    }
  };
  btnStart.onclick = function() {
    if (status == "configured") {
      $("#heartimg").addClass("bottom");
      status = "started";
      senddata(status);
    }
  };
  btnStop.onclick = function() {
    if (status == "started") {
      $("#heartimg").removeClass("bottom");
      status = "stopped";
      senddata(status);
    }
  };
  btnDelete.onclick = function() {
    if (status === "stopped") {
      $("#heart").addClass("hidden");
      $("#heartimg").addClass("sepia");
      status= "initial"; 
      senddata(status);
    }
  };

  btnDisConnect.onclick = function () {
      ws.onclose = function () {}; // disable onclose handler first
      ws.close();
      $('#login').removeClass('hidden');
      $('#main').addClass('hidden');
      $('#footer').addClass('hidden');
  };

  btnConnect.onclick = function () {
    var name = document.getElementById('name').value;
    var sessionID = document.getElementById('demoCode').value;
    if (name === null || name === "") {
      alert("Please enter your name");
      return;
    }
    console.log('opening websocket');
    ws = new WebSocket('wss://' + window.location.host + '/serveWs/'+ sessionID);
    ws.onopen = function() {

      $('#login').addClass('hidden');
      $('#main').removeClass('hidden');
      $('#footer').removeClass('hidden');
      console.log('hiding login');
      var msg = {
        name: name,  
        id: 0,
        uuid: guid,
        device: navigator.userAgent,
        status: status
      };
      console.log("Sending"+JSON.stringify(msg));
      ws.send(JSON.stringify(msg));
    };

    // Write message on receive
    ws.onmessage = function(e) {
      console.log("Received:",e);
      //document.getElementById('output').innerHTML += "Received: " + e.data + "<br>";
      var obj = JSON.parse(e.data);
      console.log("Message: "+obj.message);
      switch (obj.message) {
        case "ping":
          senddate("pong");
          break;
        case "error":
          $('#main').addClass('hidden');
          $('#footer').addClass('hidden');
          $('body').css('background-color', 'black');
          break;
        default:
          $('#main').removeClass('hidden');
          $('#footer').removeClass('hidden');
          $('body').css('background-color', 'grey');
          break;
      }
    };
    ws.onclose = function() {
      document.getElementById('btnDisConnect').innerHTML='Disconnect, click to reconnect';
      $('#btnDisConnect').removeClass('btn-success');
      $('#btnDisConnect').addClass('btn-warning');
      
    };
  };
};

function readDeviceOrientation() {

  if (Math.abs(window.orientation) === 90) {
    // Landscape
    senddata("error"); 
  } else {
    // Portrait
    senddata("initial"); 

  }

}

window.onorientationchange = readDeviceOrientation;
function senddata(status) {
  // Construct a msg object containing the data the server needs to process the message from the chat client.
  if (ws !== null) {
    var msg = {
      name: document.getElementById('name').value,
      uuid: guid,
      status: status,
    };

    ws.send(JSON.stringify(msg));
    console.log("Sending:",msg);
    //document.getElementById('output').innerHTML += "Sent: " + JSON.stringify(msg) + "<br>";
  }
}
