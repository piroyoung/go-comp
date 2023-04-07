function completion(prompt) {

    let xhr = new XMLHttpRequest();
    xhr.open('POST', '/v1/stream', true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4 && xhr.status === 200) {
            let chunks = xhr.responseText.split("event:chunk\ndata:");
            console.log(chunks);
            for (let i = 1; i < chunks.length; i++) {
                let chunk = chunks[i];
                let completion = processChunk(chunk);
                $("#completion").append(completion);
            }
        }
    };

    xhr.send(JSON.stringify({
        "value": prompt
    }));
}

// チャンクを処理する
function processChunk(chunk) {
    if (chunk.length === 0) {
        return '';
    }
    let data = JSON.parse(chunk);
    return data.value;
}

// 2秒間入力がなければ補完を実行
let timer;
document.addEventListener('keydown', () => {
    clearTimeout(timer);
    timer = setTimeout(() => {
        $("#completion").empty();
        let p = $("#prompt").val();
        completion(p);
    }, 2000);
});