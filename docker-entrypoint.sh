echo "wait db server"
./wait-for-it.sh tcp://time2do-db:3306 -t 20

echo "start go server"
./server