#!/usr/bin/env bash
http --form POST http://localhost:1323/REST/user/signin \
email='zh' \
password="000" \
fingerprint='{"fingerPrint":{"b":"B","c":"C","g":"G","d":"D","info":"MY_FINGERPRINT_A"}}'