{
    "AwsAccountId": "__ACCOUNT_ID",
    "DashboardId": "__DASHBOARD_ID",
    "Name": "Rigado Quickstart Dashboard __ENV",
    "Permissions": [
        {
            "Principal": "__ADMIN_USER_ARN",
            "Actions": [
                "quicksight:DescribeDashboard",
                "quicksight:ListDashboardVersions",
                "quicksight:UpdateDashboardPermissions",
                "quicksight:QueryDashboard",
                "quicksight:UpdateDashboard",
                "quicksight:DeleteDashboard",
                "quicksight:DescribeDashboardPermissions",
                "quicksight:UpdateDashboardPublishedVersion"
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
    ],
    "VersionDescription": "1",
    "DashboardPublishOptions": {
        "AdHocFilteringOption": {
            "AvailabilityStatus": "ENABLED"
        },
        "ExportToCSVOption": {
            "AvailabilityStatus": "ENABLED"
        },
        "SheetControlsOption": {
            "VisibilityState": "EXPANDED"
        }
    }
}