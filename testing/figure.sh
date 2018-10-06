#!/usr/bin/env bash

http GET http://localhost:1323/REST/dataset/figure \
ids=='YRD.LoanFacilitation.TotalLoansFacilitated' \
dateType==1 \
token==s45BX7l7vzqQRuJI5MpcgnV17RLT7xf7NM1D8epbkLmGU1fGETms7JnWz6931sK8 \
filters=='[{"key":"id","values":["1", "2", "3"]}]' \
fp==BCGD
