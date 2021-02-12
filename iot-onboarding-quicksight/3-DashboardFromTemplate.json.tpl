{
    "AwsAccountId": "__ACCOUNT_ID",
    "AnalysisId": "__DASHBOARD_ID",
    "Name": "Rigado Quickstart Dashboard __ENV",
    "Permissions": [
        {
            "Principal": "__ADMIN_USER_ARN",
            "Actions": [
                "quicksight:RestoreAnalysis", 
                "quicksight:UpdateAnalysisPermissions", 
                "quicksight:DeleteAnalysis", 
                "quicksight:QueryAnalysis", 
                "quicksight:DescribeAnalysisPermissions", 
                "quicksight:DescribeAnalysis", 
                "quicksight:UpdateAnalysis"
            ]
        }
    ],
    "SourceEntity": {
        "SourceTemplate": {
            "DataSetReferences": [
                {
                    "DataSetPlaceholder": "placeholder",
                    "DataSetArn": "arn:aws:quicksight:__AWS_REGION:__ACCOUNT_ID:dataset/__DATA_SET_ID"
                }
            ],
            "Arn": "__SOURCE_TEMPLATE_ARN"
        }
    },
    "Tags": [
        {
            "Key": "Name",
            "Value": "API-DemoDashboard"
        }
    ]
}