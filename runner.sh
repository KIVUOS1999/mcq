current_dir=$(pwd)
path_to_db="$current_dir/bkend-db"
path_to_redis="$current_dir/bkend-redis"
path_to_bkend="$current_dir/bkend"
path_to_build="$current_dir/build/"

cd $path_to_db
GOOS=linux GOARCH=amd64 go build
mv bkend-db $path_to_build

cd $path_to_redis
GOOS=linux GOARCH=amd64 go build
mv bkend-redis $path_to_build

cd $path_to_bkend
GOOS=linux GOARCH=amd64 go build -o bkend
mv bkend $path_to_build

cd ..

docker compose -f docker-compose.yaml up