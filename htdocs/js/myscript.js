var rawData = null;
var referer = null;

function template2html(url) {
    $.getJSON(url,function(rawData){
        var viewHelpers = {
            template2html: function(url) {
                template2html(url);
            },
            displayMarkdown: function(description){
                var converter = new showdown.Converter({
                    'github_flavouring': true,
                    'tables': true
                });
                return converter.makeHtml(description);
            }
        }

        var data = { target:rawData };
        _.extend(data,viewHelpers);
        _.extend(data,{referer: referer});
        referer = url;
        var template = _.template( $("#tpl-html").text() );
        $("#output").html( template(data) );
    });
}

var $searchableTree = $('#tree');

var search = function(e) {
    var pattern = $('#input-search').val();
    var options = {
        ignoreCase: true,
        exactMatch: false,
        revealResults: false
    };
    var results = $searchableTree.treeview('search', [ pattern, options  ]);

    $('#tree').treeview('collapseAll', { silent: false  });
    if (results.length == 1) {
        var result = results[0];
        $('#tree').treeview('revealNode', [result ,{levels: 5, silent: true}] );
        if (result.href != null){
            //alert(data.href);
            template2html(result.href);
        }
    } else {
        var output = '<p>' + results.length + ' matches found</p>';
        $.each(results, function (index, result) {
            $('#tree').treeview('revealNode', [result ,{levels: 5, silent: true}] );
            output += '<p>- <a href="#" onclick="template2html(\''+result.href+'\');">' + result.text + '</a></p>';
        });
        $('#output').html(output);
    };
};

$('#input-search').on('keyup', search);

//https://github.com/jonmiles/bootstrap-treeview/issues/179
$.getJSON('json/tree.json',function(data){
    console.dir(JSON.stringify(data)); //data received as String 
    $('#tree').treeview({
        data: "["+JSON.stringify(data)+"]"
    });;
    $('#tree').on('nodeSelected', function(event, data) {
        if (data.href != null){
            alert(data.href);
            template2html(data.href);
        }
    });
});;

