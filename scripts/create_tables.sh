#!/bin/bash

# crate tables defined in the directory "table"

HOSTNAME=""
USERNAME=""
PASSWORD=""
echo "Usage:"
echo "    HOSTNAME=XXX USERNAME=XXX PASSWORD=XXX ./create_table.sh"

set -e
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
