/*
 * Dashboard
 */

function loadStatus(callback) {
    var xhr = new XMLHttpRequest(),
        loading = document.getElementById('loading');

    loading.setAttribute('style', 'display: block;');
    xhr.open('GET', '/status?rates=please', true);
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

            var titleCol = document.createElement('td'),
                messagesCol = document.createElement('td'),
                subscriptionsCol = document.createElement('td'),
                nameDiv = document.createElement('div'),
                svg = document.createElementNS('http://www.w3.org/2000/svg', 'svg'),
                pathIn = document.createElementNS('http://www.w3.org/2000/svg', 'path'),
                pathOut = document.createElementNS('http://www.w3.org/2000/svg', 'path'),
                messagesDiv = document.createElement('div'),
                subscriptionsDiv = document.createElement('div');

            pathIn.setAttribute('class', 'in');
            pathOut.setAttribute('class', 'out');
            svg.setAttributeNS(null, 'viewbox', '0 0 300 40');
            svg.appendChild(pathIn);
            svg.appendChild(pathOut);

            nameDiv.setAttribute('class', 'name');
            svg.setAttribute('class', 'chart');
            titleCol.setAttribute('class', 'title');
            messagesCol.setAttribute('class', 'messages');
            subscriptionsCol.setAttribute('class', 'subscriptions');

            nameDiv.appendChild(document.createTextNode(queue));
            titleCol.appendChild(svg);
            titleCol.appendChild(nameDiv);
            messagesCol.appendChild(messagesDiv);
            subscriptionsCol.appendChild(subscriptionsDiv);
            tr.appendChild(titleCol);
            tr.appendChild(messagesCol);
            tr.appendChild(subscriptionsCol);
            queuesList.appendChild(tr);
        }

        var titleCol = tr.getElementsByClassName('title')[0],
            nameDiv = titleCol.getElementsByClassName('name')[0],
            svg = titleCol.getElementsByClassName('chart')[0],
            messagesCol = tr.getElementsByClassName('messages')[0],
            subscriptionsCol = tr.getElementsByClassName('subscriptions')[0],
            messagesDiv = messagesCol.getElementsByTagName('div')[0],
            subscriptionsDiv = subscriptionsCol.getElementsByTagName('div')[0],
            messages = Number(meta.messages).toLocaleString(),
            subscriptions = Number(meta.subscriptions).toLocaleString();

        messagesDiv.innerHTML = messages;
        subscriptionsDiv.innerHTML = subscriptions;

        if (meta.messages > hotThreshold) {
            nameDiv.setAttribute('class', 'name hot');
            messagesDiv.setAttribute('class', 'num messages hot');
        } else if (meta.messages > fatThreshold) {
            nameDiv.setAttribute('class', 'name fat');
            messagesDiv.setAttribute('class', 'num messages fat');
        } else if (meta.messages === 0) {
            messagesDiv.setAttribute('class', 'num messages zero');
        } else {
            nameDiv.setAttribute('class', 'name');
            messagesDiv.setAttribute('class', 'num messages');
        }

        if (meta.subscriptions === 0) {
            subscriptionsDiv.setAttribute('class', 'num subscriptions zero');
        } else {
            subscriptionsDiv.setAttribute('class', 'num subscriptions');
        }

        svg.setAttributeNS(null, 'viewbox', '0 0 300 '+ titleCol.offsetTop);
        drawChart(svg, titleCol.offsetTop, meta.in_rate_history, meta.out_rate_history);
    }
}

/*
 * Charts
 */

function drawChart(svg, maxHeight, valuesIn, valuesOut) {
    var pathIn = svg.getElementsByClassName('in')[0],
        pathOut = svg.getElementsByClassName('out')[0],
        // valuesIn = generateValues(300),
        // valuesOut = generateValues(300),
        maxDouble = calcMaxDouble(valuesIn, valuesOut),
        pointsIn = [],
        pointsOut = [];

    for (var i = 0; i < valuesIn.length; i++) {
        var normIn = Math.ceil(valuesIn[i] / maxDouble * maxHeight),
            normOut = Math.ceil(valuesOut[i] / maxDouble * maxHeight),
            pointIn = maxHeight/2 - normIn,
            pointOut = maxHeight/2 + normOut;

        pointsIn.push(pointIn);
        pointsOut.push(pointOut);
    }

    pathIn.setAttributeNS(null, 'd', buildPathD(pointsIn, maxHeight));
    pathIn.setAttributeNS(null, 'class', 'in');
    pathOut.setAttributeNS(null, 'd', buildPathD(pointsOut, maxHeight));
    pathOut.setAttributeNS(null, 'class', 'out');
}

function generateValues(num) {
    var values = [];
    for (var i = 0; i < num; i++) {
        var value = Math.ceil(Math.random() * 60) + 30;
        values.push(value);
    }

    return values;
}

function calcMaxDouble(a, b) {
    var doubleValue = 0;
    for (var i = 0; i < a.length; i++) {
        if (a[i] * 2 > doubleValue) {
            doubleValue = a[i] * 2;
        }
        if (b[i] * 2 > doubleValue) {
            doubleValue = b[i] * 2;
        }
    }

    return doubleValue * 1.2;
}

function buildPathD(points, maxHeight) {
    var d = ['M0,'+ maxHeight/2];
    for (var i = 0; i < points.length; i++) {
        d.push('L'+ i +','+ points[i]);
    }
    d.push('L300,'+ maxHeight/2, 'Z');

    return d.join(' ');
}

 /*
 * Starting up
 */

function loop(timeout, func) {
    func();
    window.setTimeout(function(){
        loop(timeout, func);
    }, timeout);
}

loop(1000, function(){
    loadStatus(updateDashboard);
});
