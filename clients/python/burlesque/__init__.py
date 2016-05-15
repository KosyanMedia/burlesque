from logging import Logger

from tornado.gen import coroutine, sleep, maybe_future
from tornado.httpclient import AsyncHTTPClient, HTTPError
from tornado.queues import Queue
from tornado.ioloop import IOLoop

class Burlesque:

    def __init__(self, url, **kwargs):
        self.client = AsyncHTTPClient()
        self.logger = kwargs.get('logger', Logger('burlesque'))
        self._request_timeout = kwargs.get('request_timeout', 30)
        self._retry_count = kwargs.get('retry', 3)
        self._url = url
        self._queue = Queue(maxsize=64)

    @coroutine
    def send(self, queue_name, body):
        try_count = 0
        url = '%s/publish?queue=%s' % (self._url, queue_name)
        while True:
            try:
                resp = yield self.client.fetch(url, method='POST', body=body,
                    request_timeout=self._request_timeout)
                self.logger.debug("successfully sent data to queue %s", queue_name)
                return
            except HTTPError as e:
                if try_count < self._retry_count:
                    self.logger.warning("can't send data to queue: %s, err: %s", url, e)
                    yield sleep(1)
                else:
                    raise(e)
            finally:
                try_count += 1

    @coroutine
    def listen(self, queues, fn, workers_count=4):
        [self._worker(queues, fn) for x in range(workers_count)]
        yield self._queue.put(True)
        yield self._queue.join()

    @coroutine
    def _worker(self, queues, fn):
        url = '%s/subscribe?queues=%s' % (self._url, ','.join(queues))
        while True:
            try:
                resp = yield self.client.fetch(url, request_timeout=self._request_timeout)
                queue = resp.headers["Queue"]
                try:
                    yield maybe_future(fn(queue, resp.body))
                except Exception as e:
                    self.send(queue, resp.body)
                    self.logger.warning("msg sent back to queue %s", queue)
            except HTTPError as e:
                if e.code != 599:  # Do not annoy logs on timeouts
                    self.logger.warning("can't receive data from queue: %s, err: %s", url, e)
                yield sleep(1)

@coroutine
def main():
    @coroutine
    def fn(queue, body):
        print(queue, body)
    import sys
    import logging
    logger = logging.getLogger('burlesque')
    logger.setLevel('DEBUG')
    queue = Burlesque(sys.argv[1], logger=logger)
    yield queue.send(sys.argv[2], sys.argv[3])
    print('msg %s sent to %s' % (sys.argv[2], sys.argv[3]))
    yield queue.listen([sys.argv[2]], fn)

"""
for test run just exceute:
docker run -d --name burlesque aviasales/burlesque
docker run --rm --link burlesque -ti -v `pwd`:/app python:latest /bin/bash -c 'cd app && python3 setup.py install && python3 burlesque/__init__.py http://burlesque:4401 1-test gggg'
"""
if __name__ == '__main__':
    IOLoop.current().run_sync(main)
