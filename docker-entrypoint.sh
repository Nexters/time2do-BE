echo "wait db server"
./wait-for-it.sh tcp://time2do-mysql:3306 -t 20

echo "start go server"
./server