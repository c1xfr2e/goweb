#!/usr/bin/env bash
curl -v -X POST http://localhost:1323/REST/ping \
  -H 'content-type: application/json' \
  -d '{"token":"iRWSPR96doUY0r2T1W9Pk6TOjtsI49A4HW96UiSqsXLKdOQhdrtiGxt5HWkEr467", "fingerprint":{"a":"A","b":"B","c":"C","d":"D","info":"MY_FINGERPRINT_A"}}'
