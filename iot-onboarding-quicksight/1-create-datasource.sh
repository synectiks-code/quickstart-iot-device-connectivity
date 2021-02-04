env=$1
account=$2
datasourceId=$3

aws quicksight create-data-source --aws-account-id $account  --cli-input-json file://out/1-datasource.json > out/datasource.json
    rc=$?
    if [ $rc -ne 0 ]; then
      echo "An Error Occured. Existing with status $rc" >&2
      exit $rc
    fi
aws quicksight describe-data-source --aws-account-id $account --data-source-id $datasourceId