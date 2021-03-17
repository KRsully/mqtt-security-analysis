#!/bin/bash

PAYLD_FILE=.payload_lengths

if [[ $# -lt 2 ]];then
	echo 'Usage: ./mosq_payloads <timeout length> <host address> [port number]'
	exit -1
fi

if [[ $# -eq 2 ]];then
	timeout $1 mosquitto_sub -h $2 -t '#' -F %l -u payload-tester > $PAYLD_FILE
else
	timeout $1 mosquitto_sub -h $2 -t '#' -F %l -u payload-tester -p $3 > $PAYLD_FILE
fi
sum=$(awk '{sum += $1}  END {print sum}' $PAYLD_FILE)
num=$(wc -l < $PAYLD_FILE)
max=$(awk 'BEGIN {max=-inf}{if ($1 > max){max=$1}} END {print max}' $PAYLD_FILE)
printf "total packets: %d, maximum payload (bytes): %d\n Avg. Payload (bytes): " $num $max
echo "scale=2;$sum/$num" | bc
exit 0
