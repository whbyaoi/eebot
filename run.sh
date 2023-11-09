cd /home/eebot
go build .
cd /eebot
./stop.sh
mv /home/eebot/eebot /eebot/
./start.sh