package server

const dashboardTmpl = `
{{define "dashboard"}}

<!DOCTYPE html>
<html lang="en">
<head>
    <title>Queues @ {{.hostname}}</title>
    <meta charset="utf8">

<style type="text/css">

* { box-sizing: border-box; }
.title, td, th {
    font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif;
    font-size: 1.3em;
    padding: 0.6em;
    font-weight: 300;
}

.title, table {
    position: absolute;
    left: 50%;
    width: 650px;
    margin-left: -325px;
}
.title {
    top: 10px;
    text-align: center;
    font-size: 1.8em;
}
table {
    top: 100px;
    border-spacing: 0;
    border-collapse: collapse;
}
th {
    font-weight: 400;
}
thead tr {
    border-bottom: #666 1px solid;
}
tbody tr:nth-child(even) {
    background-color: #f5f5f5;
}
.name {
    width: 350px;
    max-width: 350px;
    text-align: left;
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
}
.messages, .subscriptions {
    width: 150px;
    max-width: 150px;
    text-align: right;
}
.zero {
    color: #aaa;
}
.fat {
    font-weight: 600;
}
.hot {
    font-weight: 600;
    color: #f20;
}
#loading {
    display: none;
    position: absolute;
    bottom: 10px;
    right: 10px;
    font-size: 0.5em;
    width: auto;
}
#placeholder td {
    text-align: center;
}

</style>

</head>
<body>

<h1 class="title">Burlesque v{{.version}} at {{.hostname}}</h1>

<table class="stats">
    <thead>
        <tr>
            <th class="name">Queue</th>
            <th class="messages">Messages</th>
            <th class="subscriptions">Subscriptions</th>
        </tr>
    </thead>
    <tbody id="queues">
        <tr id="placeholder">
            <td colspan="3">Loading queues...</td>
        </tr>
    </tbody>
    <div id="loading">Loading...</div>
</table>

<script type="text/javascript">

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

        messagesCol.innerHTML = meta.messages;
        subscriptionsCol.innerHTML = meta.subscriptions;

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

</script>

</body>
</html>

{{end}}
`
