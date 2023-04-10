function completion(prompt) {
    let xhr = new XMLHttpRequest();
    xhr.open('POST', '/v1/stream');
    xhr.setRequestHeader('Content-Type', 'application/json');
    let cursor = 0;
    xhr.onreadystatechange = () => {
        if (cursor < xhr.responseText.length) {
            let diff = xhr.responseText.slice(cursor);
            cursor = xhr.responseText.length;
            diff.split('event:chunk\ndata:').forEach((chunk) => {
                if (chunk.length > 0) {
                    let value = processChunk(chunk);
                    $("#completion").append(value);
                }
            });
        }
    }

    xhr.send(JSON.stringify({value: prompt}));
}

function processChunk(chunk) {
    if (chunk.length === 0) {
        return '';
    }
    let data = JSON.parse(chunk);
    return data.value;
}

let timer;
document.addEventListener('keydown', () => {
    clearTimeout(timer);
    timer = setTimeout(() => {
        $("#completion").empty();
        let p = $("#prompt").val();
        console.log("post prompt: " + p);
        completion(p);
    }, 2000);
});