var demourl = window.location.protocol +'//'+ window.location.hostname+':'+window.location.port;
document.getElementById("demourl").innerHTML = demourl;
document.getElementById("demourl").href = demourl;
new QRCode(document.getElementById("qrcode"), {
  text: demourl,
  width: 128,
  height: 128,
});

Reveal.addEventListener( 'slideJoin', function() {
  (function() {
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
      .charge(-100)
      .distance(100)
      .gravity(0.05)
      .size([width, height]);

    var refreshGraph = function() {
      force
        .nodes(JSONData.nodes)
        .links(JSONData.links)
        .start();

      var node = svg.selectAll(".node")
        .data(JSONData.nodes)
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

      force.on("tick", function() {
        link.attr("x1", function(d) { return d.source.x; })
          .attr("y1", function(d) { return d.source.y; })
          .attr("x2", function(d) { return d.target.x; })
          .attr("y2", function(d) { return d.target.y; });

        node.attr("transform", function(d) { return "translate(" + d.x + "," + d.y + ")"; });
      });
    };

    var ws = new WebSocket("wss://"+window.location.host+"/slideJoin");

    //var data = [];

    ws.onmessage = function(evt) {
      // append new data from the socket
      var elements = JSON.parse(evt.data);
      console.log(JSON.stringify(elements));
      JSONData.nodes = JSONData.nodes.concat(elements.nodes);
      console.log(JSON.stringify(JSONData));
      refreshGraph();
    };

    refreshGraph();

  })();
} );


