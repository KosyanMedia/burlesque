"""
How to use this:
1. Copy this to /opt/datadog-agent/agent/checks.d/burlesque.py
2. Add config to /etc/dd-agent/conf.d/burlesque.yaml
    init_config:

    instances:
        - host: localhost
          port: 4401
3. Make /etc/init.d/datadog-agent restart
4. Run /etc/init.d/datadog-agent info and you will see this
    burlesque
    ---------
      - instance #0 [OK]
      - Collected 11 metrics, 0 events & 1 service check
"""

# stdlib
import urllib2

# 3rd party
import simplejson as json

# project
from checks import AgentCheck

class Burlesque(AgentCheck):
    """Tracks burlesque metrics via the stats monitoring port
    Connects to burlesque via the configured stats port.
    $ curl localhost:4401
    {
        "0-hotellook_search_deeplinks": {
            "messages": 0,
            "subscriptions": 9
        },
        "1-hotellook_search_results": {
            "messages": 0,
            "subscriptions": 12
        }
    }
    """
    def check(self, instance):
        if 'host' not in instance:
            raise Exception('Burlesque instance missing "host" value.')
        tags = instance.get('tags', [])

        response = self._get_data(instance)
        self.log.debug(u"Burlesque `response`: {0}".format(response))

        if not response:
            self.log.warning(u"No response received from Burlesque.")
            return

        metrics = Burlesque.parse_json(response, tags)

        for row in metrics:
            try:
                name, value, tags = row
                self.gauge(name, value, tags)
            except Exception, e:
                self.log.error(
                    u'Could not submit metric: %s: %s',
                    repr(row), str(e)
                )

    def _get_data(self, instance):
        host = instance.get('host')
        port = int(instance.get('port', 4401)) # 4401 is default

        url = "http://%s:%s/status" % (host, port)
        try:
            response = urllib2.urlopen(url).read()
        except Exception as e:
            self.log.warning("unable to get to %s, err: %s", url, e)
            return None

        return response

    @classmethod
    def parse_json(cls, raw, tags=None):
        if tags is None:
            tags = []
        parsed = json.loads(raw)
        metric_base = 'burlesque'
        output = []

        for key, val in parsed.iteritems():
            if isinstance(val, dict):
                metric_name = '%s.%s' % (metric_base, 'queue')
                ctags = tags + ['name:%s' % key]
                output.append((metric_name, val.get('messages', 0), ctags))
        return output