cd /eebot/src/eebot
go build .
./eebot/stop.sh
mv /eebot /eebot/
./eebot/start.sh