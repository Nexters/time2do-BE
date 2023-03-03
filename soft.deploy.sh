git fetch --all
git pull origin main
git checkout main
git reset --hard origin/main
# docker-compose down
docker-compose build
docker-compose up -d