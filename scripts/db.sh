#!/bin/bash

# helper script to create MySQL , Postgres , Redis containers
# Usage: ./db.sh <command> <kind> [mode]
# command : create, drop, start, stop
# kind : mysql, postgres, redis, mysql-admin, postgres-admin
# mode : dev, test, prod (default: dev)

# read yaml config from `.env/mode.yaml`

function help() {
cat <<EOF
  Usage: ./db.sh <command> <kind> [mode]

    command : create, drop, restart, stop, delete
    kind    : mysql, postgres, redis, mysql-admin, postgres-admin
    mode    : dev, test, prod
EOF
}

function main() {
    parse_args $@
    echo -e "command: $command\nkind: $db_kind\nmode: $mode"
    eval $(parse_yaml .env/$mode.yaml)
    run_command
}

function parse_args() {
    if [ $# -lt 2 ]; then
        help
        exit 1
    fi
    command=$1
    db_kind=$2
    mode=${3:-dev}

    # check if mode is valid
    if [ ! -f .env/$mode.yaml ]; then
        echo "Invalid mode: $mode"
        echo "./env/$mode.yaml not found"
        exit 1
    fi
}

# via : https://stackoverflow.com/questions/5014632/how-can-i-parse-a-yaml-file-from-a-linux-shell-script
function parse_yaml() {
   local prefix=$2
   local s='[[:space:]]*' w='[a-zA-Z0-9_]*' fs=$(echo @|tr @ '\034')
   sed -ne "s|^\($s\):|\1|" \
        -e "s|^\($s\)\($w\)$s:$s[\"']\(.*\)[\"']$s\$|\1$fs\2$fs\3|p" \
        -e "s|^\($s\)\($w\)$s:$s\(.*\)$s\$|\1$fs\2$fs\3|p"  $1 |
   awk -F$fs '{
      indent = length($1)/2;
      vname[indent] = $2;
      for (i in vname) {if (i > indent) {delete vname[i]}}
      if (length($3) > 0) {
         vn=""; for (i=0; i<indent; i++) {vn=(vn)(vname[i])("_")}
         printf("%s%s%s=\"%s\"\n", "'$prefix'",vn, $2, $3);
      }
   }'
}

function run_command() {
    echo "Running command: $command"
    container_name=$db_kind-$mode-$database_name
    case $command in
        create)
            create
            ;;
        drop)
            drop
            ;;
        restart)
            restart
            ;;
        stop)
            stop
            ;;
        delete)
            delete
            ;;
        *)
            help
            exit 1
            ;;
    esac
}

function create() {
    echo "Creating $db_kind container"
    echo "container_name: $container_name"

    case $db_kind in
        mysql)
            docker run -d \
                --name $container_name \
                -e MYSQL_ROOT_PASSWORD=$database_password \
                -e MYSQL_DATABASE=$database_name \
                -p $database_port:3306 \
                mysql
            ;;
        postgres)
            docker run -d \
                --name $container_name \
                -e POSTGRES_USER=$database_user \
                -e POSTGRES_PASSWORD=$database_password \
                -e POSTGRES_DB=$database_name \
                -p $database_port:5432 \
                postgres
            ;;
        redis)
            docker run -d \
                --name $container_name \
                --requirepass $database_password \
                -p $database_port:6379 \
                redis
            ;;
        mysql-admin)
            docker run -d \
                --name $container_name \
                --link mysql-$mode-$database_name:db \
                -e PMA_HOST=db \
                -e PMA_PORT=3306 \
                -p 33306:80 \
                phpmyadmin
            ;;
        postgres-admin)
            docker run -d \
                --name $container_name \
                --link postgres-$mode-$database_name:db \
                -e PGADMIN_DEFAULT_EMAIL=$database_user@gmail.com \
                -e PGADMIN_DEFAULT_PASSWORD=$database_password \
                -p 54321:80 \
                dpage/pgadmin4
            ;;
        *)
            help
            exit 1
            ;;
    esac
}

function drop() {
    echo "Dropping $db_kind container"
    echo "container_name: $container_name"

    case $db_kind in
        mysql)
            docker exec -it $container_name mysql -u root -p$database_password -e "DROP DATABASE $database_name"
            ;;
        postgres)
            docker exec -it $container_name psql -U $database_user -c "DROP DATABASE $database_name"
            ;;
        redis)
            docker exec -it $container_name redis-cli -a $database_password "FLUSHALL"
            ;;
        mysql-admin|postgres-admin)
            echo "Cannot drop admin container"
            echo "Use delete command instead"
            ;;
        *)
            help
            exit 1
            ;;
    esac
}

function restart() {
    echo "Starting $db_kind container"
    echo "container_name: $container_name"
    docker start $container_name
}

function stop() {
    echo "Stopping $db_kind container"
    echo "container_name: $container_name"
    docker stop $container_name
}

function delete() {
    echo "Deleting $db_kind container"
    echo "container_name: $container_name"
    docker rm $container_name -f
}

if ! (return 0 2> /dev/null); then
    main "$@"
fi