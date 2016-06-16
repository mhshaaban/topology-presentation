var ws = new WebSocket('wss://' + window.location.host + '/progress');


ws.onopen = function() {
    var msg = {
        topic:  "null",
        sender: "null",
        message: "Connection init",
        date: Date.now()
    };
     
    ws.send(JSON.stringify(msg));
}

// Write message on receive
ws.onmessage = function(e) {
    console.log("Received:",e);
    //document.getElementById('output').innerHTML += "Received: " + e.data + "<br>";
    var obj = JSON.parse(e.data);
    dataval = obj.score;
    $('.progress .amount').css("width", 100 - dataval + "%");
    $('.progress').attr("data-amount", dataval);
};

function senddata(like) {
    // Construct a msg object containing the data the server needs to process the message from the chat client.
    var msg = {
        topic:  document.getElementById("topic").value,
        like: like,
        date: Date.now()
    };
     
    ws.send(JSON.stringify(msg));
    console.log("Sending:",msg);
    //document.getElementById('output').innerHTML += "Sent: " + JSON.stringify(msg) + "<br>";
}
