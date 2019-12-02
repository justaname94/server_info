cluster_init() {
  sleep 5
  echo Starting one time cluster initialization
  ./cockroach init --insecure

  echo Running schema scripts
  sleep 5

  PARAMS="--insecure -e"
  SQL_FILE="/db_setup/scripts.sql"

  ./cockroach sql $PARAMS "$(cat $SQL_FILE)"
}

cluster_init &
./cockroach start --insecure --join=roach1
