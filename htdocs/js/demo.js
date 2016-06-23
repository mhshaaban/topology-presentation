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
var step = 1;
var btnStep = document.getElementById("step");
btnStep.onclick = function() {
  switch (step) {
    case 1:
      document.getElementById('step').innerHTML='Step ' + step + ' - 15 years ago, the pet!';
      step = 2;
      break;
    case 2:
      document.getElementById('step').innerHTML='Step ' + step + ' - 10 years ago, a menagerie';
      refreshSlide1();
      step = 3;
      break;
    case 3:
      document.getElementById('step').innerHTML='Step ' + step + ' - 5 years ago, to tame the animals';
      step = 4;
      break;
    case 4:
      document.getElementById('step').innerHTML='Step ' + step + ' - now, the cattle';
      step = 5;
      break;
    default:
      document.getElementById('step').innerHTML='Step ' + step + ' - 15 years ago, the pet!';
      step = 1;
      refreshSlide1();
      break;
  } 
};

var demourl = window.location.protocol +'//'+ window.location.hostname+':'+window.location.port;
document.getElementById("demourl").innerHTML = demourl;
document.getElementById("demourl").href = demourl;
new QRCode(document.getElementById("qrcode"), {
  text: demourl,
  //  width: 150,
  //  height: 150,
});
document.getElementById("codedemo").innerHTML  = sessionID;
//document.getElementById("codedemo2").innerHTML  = sessionID;
//Reveal.addEventListener( 'slideJoin', function() {
// (function() {
//JSONData = {
//  "nodes": [{}],
//  "links": [{}]
//};
JSONData = {
      "nodes": [
        {
          "name": "Presenter", 
          "id":0,
          "icon": "/img/laptop.png"
        },
      ],
      "links": [
        {"source":0,"target":0}
      ]
    };

//var color = d3.scale.category10();

var width = 960,
  height = 400;

var svg = d3.select("#demo").append("svg:svg")
  .attr("width", width)
  .attr("height", height);

var force = d3.layout.force()
  .charge(-500)
  .distance(300)
  .gravity(0.1)
  .size([width, height]);


var refreshSlide1 = function() {
var node = svg.selectAll(".node");
var link = svg.selectAll(".link");
  console.log('entering refreshSlide1');

  force
    .nodes(JSONData.nodes);
    //.links(JSONData.links)
    //.start();

  if (step > 1) {
    force.links(JSONData.links);
  }
    force.start();

    node.data(JSONData.nodes, function(d) { console.log('Refreshing '+ d.id); return d.id; })
    .enter().append("g")
    .attr("class", "node")
    .call(force.drag);

    link.data(JSONData.links)
    .enter().append("line")
    .attr("class", "link");

  var images = node.append("image")
    .attr("xlink:href", function(d) { return d.icon; })
    .attr("x", -40)
    .attr("y", -40)
    .attr("width", 64)
    .attr("height", 64);

  node.append("text")
    .attr("fill", function(d) { return d.color; })
    .style("font-size","20px")
    .attr("dx", 20)
    .attr("font-family","sans-serif")
    .attr("font-size","20px")
    .attr("dy", ".05em")
    .text(function(d) { return d.name; });

  // make the image grow a little on mouse over and add the text details on click
  var setEvents = images
  // Append hero text
    .on( 'click', function (d) {
      d3.select("#NodeName").html(d.name); 
      d3.select("#NodeStatus").html(d.status); 
    });

  /*
    .on( 'mouseenter', function() {
      // select element in current context
      d3.select( this )
        .transition()
        .attr("x", function(d) { return -60;})
        .attr("y", function(d) { return -60;})
        .attr("height", 100)
        .attr("width", 100);
    })
  // set back
    .on( 'mouseleave', function() {
      d3.select( this )
        .transition()
        .attr("x", function(d) { return -25;})
        .attr("y", function(d) { return -25;})
        .attr("height", 50)
        .attr("width", 50);
    });
*/
  /*
  force.on("tick", function() {
//    node.attr("cx", function(d) { return d.x; })
//      .attr("cy", function(d) { return d.y; });
    link.attr("x1", function(d) { return d.source.x; })
      .attr("y1", function(d) { return d.source.y; })
      .attr("x2", function(d) { return d.target.x; })
      .attr("y2", function(d) { return d.target.y; });

    node.attr("transform", function(d) { return "translate(" + d.x + "," + d.y + ")"; });
  });
  */
  force.on('tick', function(e) {
    node
      //.transition().ease('linear').duration(400)
      .attr('transform', function(d, i) {
      return 'translate('+ d.x +', '+ d.y +')';
    });

    link
      .attr('x1', function(d) { return d.source.x; })
      .attr('y1', function(d) { return d.source.y; })
      .attr('x2', function(d) { return d.target.x; })
      .attr('y2', function(d) { return d.target.y; });
  });

};

ws.onmessage = function(evt) {
  // append new data from the socket
  var elements = JSON.parse(evt.data);
  console.log(JSON.stringify(elements));
  if (elements.message === "ping") {
    var msg = {
      name: "presenter",
      id: 0,
      uuid: guid,
      icon: "/img/laptop.png",
      status: "pong"
    };
    ws.send(JSON.stringify(msg));
    console.log("pong sent");
    return;
  }
  JSONData.nodes = elements.nodes;
  JSONData.links =elements.links;
  refreshSlide1();
  refreshSlide1();
};

refreshSlide1();

//  })();
//} );


