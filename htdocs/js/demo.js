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

var demourl = window.location.protocol +'//'+ window.location.hostname+':'+window.location.port;
document.getElementById("demourl").innerHTML = demourl;
document.getElementById("demourl").href = demourl;
new QRCode(document.getElementById("qrcode"), {
  text: demourl,
  //  width: 150,
  //  height: 150,
});
document.getElementById("codedemo").innerHTML  = sessionID;
document.getElementById("codedemo2").innerHTML  = sessionID;
//Reveal.addEventListener( 'slideJoin', function() {
// (function() {
JSONData = {
  "nodes": [{}],
  "links": [{}]
};
//JSONData = {
//      "nodes": [
//        {
//          "name": "Presenter", 
//          "id":0,
//          "icon": "/img/laptop.png"
//        },
//      ],
//      "links": [
//        {"source":0,"target":0}
//      ]
//    };

//var color = d3.scale.category10();

var width = 960,
  height = 400;

var svg = d3.select("#demo").append("svg:svg")
  .attr("width", width)
  .attr("height", height);

var force = d3.layout.force()
  .charge(-200)
  .distance(300)
  .gravity(0.01)
  .size([width, height]);

var refreshSlide1 = function() {
  console.log('entering refreshSlide1');

  force
    .nodes(JSONData.nodes)
    .links(JSONData.links)
    .start();

  var node = svg.selectAll(".node")
    .data(JSONData.nodes)//, function(d) { return d.uuid; })
    .enter().append("g")
    .attr("class", "node")
    .call(force.drag);

  var link = svg.selectAll(".link")
    .data(JSONData.links)
    .enter().append("line")
    .attr("class", "link");

  node.append("image")
    .attr("xlink:href", function(d) { return d.icon; })
    .attr("x", -8)
    .attr("y", -8)
    .attr("width", 64)
    .attr("height", 64);

  node.append("text")
    .attr("dx", 54)
    .attr("dy", ".05em")
    .text(function(d) { return d.name; });


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
      node.attr('transform', function(d, i) {
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
  //JSONData.nodes = JSONData.nodes.concat(elements.nodes);
  //JSONData.links = JSONData.links.concat(elements.links);
  JSONData2 = {
    "nodes": [{}],
    "links": [{"source": 0, "target": 0}]
  };
  for (var j=0;j< elements.nodes.length; j++) {
    console.log('looping for element'+j);
    var found = false;
    for (var i=0;i<JSONData.nodes.length;i++) {
      // If match
      if (elements.nodes[j].uuid === JSONData.nodes[i].uuid) {
        console.log('Element found in the array '+elements.nodes[j].uuid+" "+JSONData.nodes[i].uuid);
        JSONData.nodes[i].status = elements.nodes[j].status;
        found = true;
        break;
      }
    }
    if (found === false) {
      console.log('element not found, appending it');
      console.log('Length of JSONData2.nodes: '+JSONData2.nodes.length);
      if (JSONData2.nodes.length > 1) {
        JSONData2.nodes = JSONData2.nodes.concat(elements.nodes[j]);
      } else {
        JSONData2.nodes[0] = elements.nodes[j];
      }
    }
  }
  for (var t=0;t< elements.links.length; t++) {
    console.log('looping for element'+t);
    var found2 = false;
    for (var i2=0;i2<JSONData.links.length;i2++) {
      // If match
      if (elements.links[t].source === JSONData.links[i2].source && elements.links[t].target === JSONData.links[i2].target) {
        console.log('Element found in the array '+elements.links[t].uuid+" "+JSONData.links[i2].uuid);
        found2 = true;
        break;
      }
    }
    if (found2 === false) {
      console.log('element not found, appending it');
      JSONData2.links = JSONData2.links.concat(elements.links[t]);
    }
  }
//  JSONData.nodes = elements.nodes;
    JSONData.nodes = JSONData2.nodes;
   JSONData.links = elements.links;
//  JSONData.links = elements.links;
  console.log("to be appened"+JSON.stringify(JSONData));
  refreshSlide1();
};

refreshSlide1();

//  })();
//} );


