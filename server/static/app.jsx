var Chart = React.createClass({
    getInitialState: function() {
        return {pointsIn: [], pointsOut: []};
    },

    componentDidMount: function() {
        this.buildPoints(this.props);
    },

    componentWillReceiveProps: function(nextProps) {
        this.buildPoints(nextProps);
    },

    buildPoints: function(props) {
        var maxDouble = this.calcMaxDouble(props.valuesIn, props.valuesOut) || 1,
            pointsIn = [],
            pointsOut = [];

        for (var i = 0; i < props.valuesIn.length; i++) {
            var normIn = Math.ceil(props.valuesIn[i] / maxDouble * props.height),
                normOut = Math.ceil(props.valuesOut[i] / maxDouble * props.height),
                pointIn = props.height/2 - normIn,
                pointOut = props.height/2 + normOut;

            pointsIn.push(pointIn);
            pointsOut.push(pointOut);
        }

        this.setState({pointsIn: pointsIn, pointsOut: pointsOut});
    },

    calcMaxDouble: function(a, b) {
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
    },

    buildPathD: function(points) {
        var d = ['M0,'+ this.props.height/2],
            missing = this.props.width - points.length;

        for (var i = 0; i < missing; i++) {
            d.push('L'+ i +','+ this.props.height/2);
        }
        for (var i = 0; i < points.length; i++) {
            d.push('L'+ missing+i +','+ points[i]);
        }
        d.push('L'+ this.props.width +','+ this.props.height/2, 'Z');

        return d.join(' ');
    },

    render: function() {
        var viewBox = [0, 0, this.props.width, this.props.height].join(' ');
        return (
            <svg className="chart" viewBox={viewBox}>
                <path className="in" d={this.buildPathD(this.state.pointsIn)} />
                <path className="out" d={this.buildPathD(this.state.pointsOut)} />
            </svg>
        );
    }
});

var QueuesList = React.createClass({
    render: function(){
        if (!this.props.isDataRecieved) {
            return (
                <tbody id="queues">
                    <tr>
                        <td colSpan="3" className="placeholder">Loading...</td>
                    </tr>
                </tbody>
            )
        }

        if (Object.keys(this.props.queues) === 0) {
            return (
                <tbody id="queues">
                    <tr>
                        <td colSpan="3" className="placeholder">This server has no queues</td>
                    </tr>
                </tbody>
            )
        }

        var queues = this.props.queues;
        var createQueue = function(name) {
            var meta = queues[name],
                titleClasses = ['title'],
                messagesClasses = ['messages'],
                subscriptionsClasses = ['subscriptions'];

            if (meta.messages > 1000) {
                titleClasses.push('hot');
                messagesClasses.push('hot');
            } else if (meta.messages > 100) {
                titleClasses.push('fat');
                messagesClasses.push('fat');
            } else if (meta.messages === 0) {
                messagesClasses.push('zero');
            }
            if (meta.subscriptions === 0) {
                subscriptionsClasses.push('zero');
            }
            return (
                <tr key={name}>
                    <td className="title">
                        <div className={titleClasses.join(' ')}>{name}</div>
                        <Chart
                            valuesIn={meta.in_rate_history}
                            valuesOut={meta.out_rate_history}
                            width={300}
                            height={40} />
                    </td>
                    <td className={messagesClasses.join(' ')}>
                        <div className="num">{meta.messages}</div>
                    </td>
                    <td className={subscriptionsClasses.join(' ')}>
                        <div className="num">{meta.subscriptions}</div>
                    </td>
                </tr>
            );
        };

        return (
            <tbody id="queues">
                {Object.keys(queues).map(createQueue)}
            </tbody>
        );
    }
});

var Dashboard = React.createClass({
    getInitialState: function() {
        return {queues: {}, isDataRecieved: false};
    },

    componentDidMount: function() {
        this.loop(this.props.interval, this.refresh);
    },

    componentWillUnmount: function() {
        clearTimeout(this.timeout);
    },

    loop: function (timeout, func) {
        var loop = this.loop;
        func();
        this.timeout = setTimeout(function(){
            loop(timeout, func);
        }, timeout);
    },

    refresh: function () {
        var xhr = new XMLHttpRequest()
            self = this;
        xhr.open('GET', '/status?rates=please', true);
        xhr.onreadystatechange = function() {
            if (xhr.readyState === 4) {
                if (xhr.status === 200) {
                    self.setState({
                        queues: JSON.parse(xhr.responseText),
                        isDataRecieved: true
                    });
                }
            }
        };
        xhr.send(null);
    },

    render: function() {
        return (
            <table className="stats">
                <thead>
                    <tr>
                        <th className="name">Queue</th>
                        <th className="messages">Messages</th>
                        <th className="subscriptions">Subscriptions</th>
                    </tr>
                </thead>
                <QueuesList queues={this.state.queues} isDataRecieved={this.state.isDataRecieved} />
            </table>
        );
    }
});
