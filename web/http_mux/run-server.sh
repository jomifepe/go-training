#!/bin/bash

SCRIPT_NAME=$(basename "$0") 
BIN_NAME="go_server"

run_server() {
   echo "$SCRIPT_NAME | Running server..."
   eval "/app/$BIN_NAME &"
}

rerun_server () {
   PID=$(pidof "$BIN_NAME")
   if [ ! -z "$PID" ]
   then
      echo "$SCRIPT_NAME | Killing old server PID: $PID"
      kill $PID
   fi

   echo "$SCRIPT_NAME | Building server..."
   go build -o "$BIN_NAME" main.go
   run_server
}

lock_build() {
   [ -f /tmp/server.lock ] && inotifywait -e DELETE /tmp/server.lock
   touch /tmp/server.lock
}

unlock_build() {
   rm -f /tmp/server.lock
}

# run the server for the first time
rerun_server

inotifywait -e MODIFY -r -m /app |
   while read path action file; do
      lock_build
      ext=${file: -3}
      if [[ "$ext" == ".go" ]]; then
         echo "$SCRIPT_NAME | File changed: $file"
         rerun_server
      fi
      unlock_build
   done