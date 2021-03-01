{
    "AwsAccountId": "__ACCOUNT_ID",
    "DataSourceId": "__DATA_SOURCE_ID",
    "Name": "Rigado Qickstart Datasource",
    "Type": "ATHENA",
    "DataSourceParameters": {
        "AthenaParameters": {
            "WorkGroup": "primary"
        }
    },
    "Permissions": [
        {
            "Principal": "__ADMIN_USER_ARN",
            "Actions": [
                "quicksight:UpdateDataSourcePermissions",
                "quicksight:DescribeDataSource",
                "quicksight:DescribeDataSourcePermissions",
                "quicksight:PassDataSource",
                "quicksight:UpdateDataSource",
                "quicksight:DeleteDataSource"
            ]
        }
    ],
    "Tags": [
        {
            "Key": "Name",
            "Value": "API-AthenaDataSource"
        }
    ]
}