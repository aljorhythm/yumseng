#!/bin/bash

port=80

echo kill process using port $port
kill $(lsof -t -i:$port)

echo running program from source
PORT=$port make run-source &
progpid=$!

echo sleeping for server to start
sleep 5

echo run integration tests
PORT=$port make integration-test
testresult=$?

echo killing server $progpid
kill $progpid

echo kill process using port $port
kill $(lsof -t -i:$port)

if [ $testresult -eq 0 ]; then
  echo "ok"
else
  echo "not ok"
fi

exit $testresult
