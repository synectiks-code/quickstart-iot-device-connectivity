env=$1
account=$2

aws quicksight create-analysis --cli-input-json file://out/3-DashboardFromTemplate.json
    rc=$?
    if [ $rc -ne 0 ]; then
      echo "An Error Occured. Existing with status $rc" >&2
      exit $rc
    fi
