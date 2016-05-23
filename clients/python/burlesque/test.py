from tornado.testing import gen_test, AsyncTestCase, main

import Burlesque

class BurlesqueTest(AsyncTestCase):

    def _get_burlesque(self):
        return Burlesque(os.environ.get('BURLESQUE_URL', 'http://localhost:4401'))

    @gen_test
    def test_send(self):
        burlesque = self._get_burlesque()
        msg = 'test msg'
        queue = '1-test'
        yield burlesque.send(queue, msg)

if __name__ == '__main__':
    main()
