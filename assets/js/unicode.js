// Unicode 编码函数
function unicodeEncode(str) {
    return str.split('').map(function(char) {
        return '\\u' + ('0000' + char.charCodeAt(0).toString(16)).slice(-4);
    }).join('');
}

// Unicode 解码函数
function unicodeDecode(str) {
    return str.replace(/\\u([0-9a-fA-F]{4})/g, function(match, hex) {
        return String.fromCharCode(parseInt(hex, 16));
    });
}