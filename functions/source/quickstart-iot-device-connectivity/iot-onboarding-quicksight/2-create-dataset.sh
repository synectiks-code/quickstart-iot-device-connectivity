env=$1
account=$2
datasetId=$3

aws quicksight create-data-set --cli-input-json file://out/2-dataset.json
    rc=$?
    if [ $rc -ne 0 ]; then
      echo "An Error Occured. Existing with status $rc" >&2
      exit $rc
    fi
aws quicksight describe-data-set --aws-account-id $account --data-set-id $datasetId
