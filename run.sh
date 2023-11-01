cd /eebot/src/eebot
go build .
cd /eebot
./stop.sh
mv /eebot/src/eebot/eebot /eebot/
./start.sh