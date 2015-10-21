function vals(name) {
    if (name == "Lines") {
        return {Lines: $("#Lines")[0].value.split('\n')};
    }

    var inputs = $('input[name='+name+']');
    var data = {};
    if (inputs.length == 1) {
        if (inputs[0].type=="checkbox") {
            data[name] = inputs[0].checked;
        } else {
            data[name] = inputs.val();
        }
    } else if (inputs.length > 1) {
        data[name] = inputs.map(function(k,v) { return v.value; }).get();
    }
    return data
}

window.addEventListener('input', function(ev) {
    if (ev.target.name == "frame") {
        selectFrame();
        return
    }
    if (ev.target.name == "Lines") {
        updateFrame();
        // fall through
    }

    $.post("set", vals(ev.target.name));
})

function checkedVal(name) {
    var data = {};
    data[name] = $('input[name='+name+']')[0].checked;
    return data;
}

$('input[type=checkbox]').click(function(ev){
    $.post("set", checkedVal(ev.target.name));
});

function initValue(key, val) {
    if (key == "Lines") {
        $("#Lines")[0].value = val.join('\n');
        return;
    }

    var input = $('input[name='+key+']');
    if (input.length == 0) {
        console.log('No such input: ', key)
        return;
    }

    if (val.length && input.length == val.length) {
        for (var i = 0; i < input.length; i++) {
            input[i].value = val[i];
        }
    } else {
        input[0].value = val;
        input[0].checked = val;
    }
}
$.get("settings",  function(data) { $.each(data, initValue); });

function buttonClick(event) {
    var myValue = '\\x0' + event.target.textContent;
    var textArea = $('#Lines')[0];
    var startPos = textArea.selectionStart;
    var endPos = textArea.selectionEnd;
    var scrollTop = textArea.scrollTop;
    textArea.value = textArea.value.substring(0, startPos) + myValue +
                     textArea.value.substring(endPos, textArea.value.length);
    textArea.focus();
    textArea.selectionStart = startPos + myValue.length;
    textArea.selectionEnd = startPos + myValue.length;
    textArea.scrollTop = scrollTop;

    onInput();
}

function onInput(event) {
    var textArea = $('#Lines')[0];
    var content = textArea.value;
    content = content.replace('\\x0', '');
    $.post("set", vals("Lines"));
    var message = $('#CharsLeft')[0].textContent =
        "Random number: " + (32 - content.length + 1);
}
