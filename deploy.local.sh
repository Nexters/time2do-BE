rm -rf mysql
git fetch --all
git pull origin main
docker-compose down
docker-compose build
docker-compose up -d
docker-compose stop time2do-go