#!/bin/bash
while true; do
	for a in abcd1234 abcd4321 xyz9232; do
		for m in requests files messages; do
			NOW=$(date -u +"%Y-%m-%dT%H:%M:%S.000Z")
			VALUE=$(( $RANDOM % 50 + 1 ))
			echo -e "time:$NOW\ta:$a\tm:$m\tv:$VALUE"
		done
	done
	sleep 1
done
