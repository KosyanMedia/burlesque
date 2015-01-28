function loadStatus(callback) {
    var xhr = new XMLHttpRequest(),
        loading = document.getElementById('loading');

    loading.setAttribute('style', 'display: block;');
    xhr.open('GET', '/status', true);
    xhr.onreadystatechange = function() {
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                var queues = JSON.parse(xhr.responseText);
                loading.setAttribute('style', 'display: none;');
                callback(queues);
            }
        }
    };
    xhr.send(null);
}

function updateDashboard(queues) {
    var queuesList = document.getElementById('queues'),
        placeholder = document.getElementById('placeholder'),
        fatThreshold = 100,
        hotThreshold = 1000;

    if (Object.keys(queues).length === 0) {
        var td = placeholder.getElementsByTagName('td')[0];
        td.innerHTML = 'Empty';
    } else if (placeholder) {
        queuesList.removeChild(placeholder);
    }

    for (queue in queues) {
        var meta = queues[queue],
            id = 'queue_' + queue,
            tr = document.getElementById(id);

        if (!tr) {
            tr = document.createElement('tr');
            tr.setAttribute('id', id);

            var nameCol = document.createElement('td');
            nameCol.appendChild(document.createTextNode(queue));
            tr.appendChild(nameCol);
            tr.appendChild(document.createElement('td'));
            tr.appendChild(document.createElement('td'));

            queuesList.appendChild(tr);
        }

        var cols = tr.getElementsByTagName('td'),
            nameCol = cols[0],
            messagesCol = cols[1],
            subscriptionsCol = cols[2];

        messagesCol.innerHTML = Number(meta.messages).toLocaleString();
        subscriptionsCol.innerHTML = Number(meta.subscriptions).toLocaleString();

        if (meta.messages > hotThreshold) {
            nameCol.setAttribute('class', 'name hot');
            messagesCol.setAttribute('class', 'messages hot');
        } else if (meta.messages > fatThreshold) {
            nameCol.setAttribute('class', 'name fat');
            messagesCol.setAttribute('class', 'messages fat');
        } else if (meta.messages === 0) {
            messagesCol.setAttribute('class', 'messages zero');
        } else {
            nameCol.setAttribute('class', 'name');
            messagesCol.setAttribute('class', 'messages');
        }

        if (meta.subscriptions === 0) {
            subscriptionsCol.setAttribute('class', 'subscriptions zero');
        } else {
            subscriptionsCol.setAttribute('class', 'subscriptions');
        }
    }
}

function loop(timeout, func) {
    func();
    window.setTimeout(function(){
        loop(timeout, func);
    }, timeout);
}

loop(1000, function(){
    loadStatus(updateDashboard);
});
