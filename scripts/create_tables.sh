#!/bin/bash

# crate tables defined in the directory "table"

set -e
HOSTNAME="app-container-db.cs5vq0ejfcck.rds.cn-north-1.amazonaws.com.cn"
USERNAME="root"
PASSWORD="HFBkCNuF2WcucFvC"
#DATABASE="app_container_db"
BATCH_FILE=`mktemp batchfile.XXXXXX`
echo "use $DATABASE;" > $BATCH_FILE && echo "use $DATABASE;" >> $BATCH_FILE

for i in `ls tables`
do
  if [ -f tables/$i ]
  then
    echo "processing table $i"
    cat tables/$i >> $BATCH_FILE
  fi
done

echo "start creating tables..."
mysql -h$HOSTNAME -u$USERNAME -p$PASSWORD < $BATCH_FILE
echo "start creating tables...Done!"

echo "remove file $BATCH_FILE"
rm $BATCH_FILE
