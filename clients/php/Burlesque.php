<?php

namespace KosyanMedia;

class Burlesque {
    private $url;

    public function __construct($url = 'http://127.0.0.1:4401')
    {
        $this->url = $url;
    }

    public function get(array $queues, $timeout = 5)
    {
        $context = stream_context_create([
            'http' => [
                'method' => 'GET',
                'timeout' => $timeout,
            ],
        ]);
        $data = @file_get_contents($this->url . '/subscribe?queues=' . implode(',', $queues), false, $context);
        if ($http_response_header === []) {
            return null;
        }
        $queue = null;
        foreach ($http_response_header as $header) {
            if (substr($header, 0, 6) == 'Queue:') {
                $queue = trim(substr($header, 7));
            }
        }
        return [$queue, $data];
    }

    public function put($queue, $data)
    {
        $contextOptions = [
            'http' => [
                'method' => 'POST',
                'content' => $data,
                'header' => 'Content-type: text/plain'
            ]
        ];
        $context = stream_context_create($contextOptions);
        file_get_contents($this->url . '/publish?queue=' . $queue, false, $context);
    }

    public function length($queue)
    {
        $data = $this->status();
        if (isset($data[$queue])) {
            return $data[$queue]['messages'];
        } else {
            return 0;
        }
    }

    public function status()
    {
        return json_decode(file_get_contents($this->url . '/status'), true);
    }
}
