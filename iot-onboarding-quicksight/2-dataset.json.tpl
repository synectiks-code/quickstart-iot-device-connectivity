{
    "AwsAccountId": "__ACCOUNT_ID",
    "DataSetId": "__DATA_SET_ID",
    "Name": "Rigado Dashboard Dataset",
    "PhysicalTableMap": {
            "51d140d2-c31f-4636-94f2-efe06f5e69c6": {
                "CustomSql": {
                    "DataSourceArn": "arn:aws:quicksight:__AWS_REGION:__ACCOUNT_ID:datasource/__DATASOURCE_ID",
                    "Name": "rigado_sensors_data",
                    "SqlQuery": "SELECT *\nFROM \"__GLUE_DB\".\"__ATHENA_TABLE_NAME\"\nWHERE (\"year\" = date_format(current_date,'%Y')\n        AND \"month\"=date_format(current_date,'%m')\n        AND \"day\"=date_format(current_date,'%d')) OR (\"year\" = date_format(current_date - interval '1' day,'%Y')\n        AND \"month\"=date_format(current_date - interval '1' day,'%m')\n        AND \"day\"=date_format(current_date - interval '1' day,'%d')) ",
                    "Columns": [
                        {
                            "Name": "device.dtmi",
                            "Type": "STRING"
                        },
                        {
                            "Name": "device.gatewayid",
                            "Type": "STRING"
                        },
                        {
                            "Name": "device.deviceid",
                            "Type": "STRING"
                        },
                        {
                            "Name": "device.capabilitymodelid",
                            "Type": "STRING"
                        },
                        {
                            "Name": "measurements.batterylevel.string",
                            "Type": "STRING"
                        },
                        {
                            "Name": "measurements.batterylevel.int",
                            "Type": "INTEGER"
                        },
                        {
                            "Name": "measurements.batterylevel.double",
                            "Type": "DECIMAL"
                        },
                        {
                            "Name": "measurements.temperature.string",
                            "Type": "STRING"
                        },
                        {
                            "Name": "measurements.temperature.int",
                            "Type": "INTEGER"
                        },
                        {
                            "Name": "measurements.temperature.double",
                            "Type": "DECIMAL"
                        },
                        {
                            "Name": "measurements.humidity.int",
                            "Type": "INTEGER"
                        },
                        {
                            "Name": "measurements.humidity.double",
                            "Type": "DECIMAL"
                        },
                        {
                            "Name": "measurements.rssi.string",
                            "Type": "STRING"
                        },
                        {
                            "Name": "measurements.rssi.int",
                            "Type": "INTEGER"
                        },
                        {
                            "Name": "ts",
                            "Type": "INTEGER"
                        },
                        {
                            "Name": "measurements.lastchange",
                            "Type": "STRING"
                        },
                        {
                            "Name": "measurements.count",
                            "Type": "INTEGER"
                        },
                        {
                            "Name": "measurements.occupied",
                            "Type": "BOOLEAN"
                        },
                        {
                            "Name": "measurements.interval",
                            "Type": "STRING"
                        },
                        {
                            "Name": "measurements.lastseen",
                            "Type": "STRING"
                        },
                        {
                            "Name": "measurements.history",
                            "Type": "STRING"
                        },
                        {
                            "Name": "measurements.statechange",
                            "Type": "BOOLEAN"
                        },
                        {
                            "Name": "measurements.accelerometeraverage.avgy",
                            "Type": "DECIMAL"
                        },
                        {
                            "Name": "measurements.accelerometeraverage.avgx",
                            "Type": "DECIMAL"
                        },
                        {
                            "Name": "measurements.accelerometeraverage.avgz",
                            "Type": "DECIMAL"
                        },
                        {
                            "Name": "measurements.accel.y",
                            "Type": "DECIMAL"
                        },
                        {
                            "Name": "measurements.accel.x",
                            "Type": "DECIMAL"
                        },
                        {
                            "Name": "measurements.accel.z",
                            "Type": "DECIMAL"
                        },
                        {
                            "Name": "measurements.launchcount",
                            "Type": "STRING"
                        },
                        {
                            "Name": "measurements.alarms",
                            "Type": "STRING"
                        },
                        {
                            "Name": "measurements.screentype",
                            "Type": "STRING"
                        },
                        {
                            "Name": "measurements.size",
                            "Type": "INTEGER"
                        },
                        {
                            "Name": "measurements.cap",
                            "Type": "INTEGER"
                        },
                        {
                            "Name": "measurements.displayedpictureid",
                            "Type": "STRING"
                        },
                        {
                            "Name": "measurements.temp",
                            "Type": "INTEGER"
                        },
                        {
                            "Name": "measurements.version",
                            "Type": "STRING"
                        },
                        {
                            "Name": "measurements.maxnodes",
                            "Type": "INTEGER"
                        },
                        {
                            "Name": "measurements.advinterval",
                            "Type": "INTEGER"
                        },
                        {
                            "Name": "measurements.txpower",
                            "Type": "INTEGER"
                        },
                        {
                            "Name": "partition_0",
                            "Type": "STRING"
                        },
                        {
                            "Name": "year",
                            "Type": "STRING"
                        },
                        {
                            "Name": "month",
                            "Type": "STRING"
                        },
                        {
                            "Name": "day",
                            "Type": "STRING"
                        },
                        {
                            "Name": "hour",
                            "Type": "STRING"
                        }
                    ]
                }
            }
        },
    "Permissions": [
        {
            "Principal": "__ADMIN_USER_ARN",
            "Actions": [
                "quicksight:DescribeDataSet",
                "quicksight:DescribeDataSetPermissions",
                "quicksight:PassDataSet",
                "quicksight:DescribeIngestion","quicksight:ListIngestions",
                "quicksight:UpdateDataSet","quicksight:DeleteDataSet",
                "quicksight:CreateIngestion",
                "quicksight:CancelIngestion",
                "quicksight:UpdateDataSetPermissions"
            ]
        }
    ],
        "LogicalTableMap": {
            "51d140d2-c31f-4636-94f2-efe06f5e69c6": {
                "Alias": "__DATA_SET_ID_logical_map",
                "DataTransforms": [
                    {
                        "CastColumnTypeOperation": {
                            "ColumnName": "measurements.temperature.string",
                            "NewColumnType": "DECIMAL"
                        }
                    },
                    {
                        "ProjectOperation": {
                            "ProjectedColumns": [
                                "device.dtmi",
                                "device.gatewayid",
                                "device.deviceid",
                                "device.capabilitymodelid",
                                "measurements.batterylevel.string",
                                "measurements.batterylevel.int",
                                "measurements.batterylevel.double",
                                "measurements.temperature.string",
                                "measurements.temperature.int",
                                "measurements.temperature.double",
                                "measurements.humidity.int",
                                "measurements.humidity.double",
                                "measurements.rssi.string",
                                "measurements.rssi.int",
                                "ts",
                                "measurements.lastchange",
                                "measurements.count",
                                "measurements.occupied",
                                "measurements.interval",
                                "measurements.lastseen",
                                "measurements.history",
                                "measurements.statechange",
                                "measurements.accelerometeraverage.avgy",
                                "measurements.accelerometeraverage.avgx",
                                "measurements.accelerometeraverage.avgz",
                                "measurements.accel.y",
                                "measurements.accel.x",
                                "measurements.accel.z",
                                "measurements.launchcount",
                                "measurements.alarms",
                                "measurements.screentype",
                                "measurements.size",
                                "measurements.cap",
                                "measurements.displayedpictureid",
                                "measurements.temp",
                                "measurements.version",
                                "measurements.maxnodes",
                                "measurements.advinterval",
                                "measurements.txpower",
                                "partition_0",
                                "year",
                                "month",
                                "day",
                                "hour"
                            ]
                        }
                    }
                ],
                "Source": {
                    "PhysicalTableId": "51d140d2-c31f-4636-94f2-efe06f5e69c6"
                }
            }
        },
        "ImportMode": "DIRECT_QUERY"
}