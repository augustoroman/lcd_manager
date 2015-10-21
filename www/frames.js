
function saveFrame() {
    var id = (new Date()).getTime();
    var val = $('#Lines').val();
    $('#frame')
        .append($("<option></option>")
            .attr("value", id)
            .text(val));
    $('#frame').val(id);
    uploadFrames();
}

function delFrame() {
    var sel = $('#frame');
    var id = sel.val();
    if (!id) {
        return;
    }
    var option = $('option[value='+id+']');
    var next = option.next();
    if (!next) {
        next = option.prev();
    }
    if (next) {
        sel.val(next.val());
    }
    option.remove();
    uploadFrames();
}

function selectFrame() {
    var id = $('#frame').val();
    if (!id) {
        return
    }
    var val = $('option[value='+id+']').text();
    $("#Lines").val(val);
    $.post("set", vals("Lines"));
}

function updateFrame() {
    var id = $('#frame').val();
    if (!id) {
        return
    }
    var val = $("#Lines").val();
    $('option[value='+id+']').text(val);
    uploadFrames();
}

$(document).keydown(function(event) {
        // If Control or Command key is pressed and the S key is pressed
        // run save function. 83 is the key code for S.
        if((event.ctrlKey || event.metaKey) && event.which == 83) {

            saveFrame();

            event.preventDefault();
            return false;
        };
    }
);

function getFrames() {
    var opts = $('#frame option');
    var frames = {};
    for (var i = 0; i < opts.length; i++) {
        var opt = $(opts[i]);
        var id = opt.val();
        var val = opt.text();
        frames[id] = val;
    }
    return frames;
}
function uploadFrames() {
    $.ajax('frames', {
        type: 'put',
        data: JSON.stringify(getFrames()),
    });
}
function setFrames(frames) {
    for (id in frames) {
        var val = frames[id];
        $('#frame')
            .append($("<option></option>")
                .attr("value", id)
                .text(val));
    }
}

$.get('frames', setFrames);
