Generating Configuration'dev'
3.0-Removing existing config filr for 'dev' (if exists)
rm: cannot remove 'iot-onboarding-infra-config-dev.json': No such file or directory
3.1-Generating new config file for env 'dev'
3.2 CReating admin User and getting rfresh token
1. Create user iot-onboarding-admin-2023-08-21-23-21-18@example.com for testing
{
    "User": {
        "Username": "b408e428-6071-7015-bcad-447949d63a7b",
        "Attributes": [
            {
                "Name": "sub",
                "Value": "b408e428-6071-7015-bcad-447949d63a7b"
            },
            {
                "Name": "email",
                "Value": "iot-onboarding-admin-2023-08-21-23-21-18@example.com"
            }
        ],
        "UserCreateDate": "2023-08-21T23:21:20.458000+05:30",
        "UserLastModifiedDate": "2023-08-21T23:21:20.458000+05:30",
        "Enabled": true,
        "UserStatus": "FORCE_CHANGE_PASSWORD"
    }
}
{
    "ChallengeParameters": {},
    "AuthenticationResult": {
        "AccessToken": "eyJraWQiOiJtNG51ZEo2WG5MV01Pekltd1BTQXgxeDJWWWxaWGxXNmtMaExRYUlxNmlFPSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiJiNDA4ZTQyOC02MDcxLTcwMTUtYmNhZC00NDc5NDlkNjNhN2IiLCJpc3MiOiJodHRwczpcL1wvY29nbml0by1pZHAudXMtZWFzdC0xLmFtYXpvbmF3cy5jb21cL3VzLWVhc3QtMV8xTFBMRUtQUDAiLCJjbGllbnRfaWQiOiJoc29qYW5nNTVsNXIxMjVncmdzNjk4ZTY4Iiwib3JpZ2luX2p0aSI6ImRiMWMyYjY2LWE2OTEtNDRkYi04NDA2LTg3NDdkZjViM2FiOCIsImV2ZW50X2lkIjoiYTAyNDA2YzgtMTkyZC00MGE2LTk2NTYtYTAxY2UwMTk1ODJjIiwidG9rZW5fdXNlIjoiYWNjZXNzIiwic2NvcGUiOiJhd3MuY29nbml0by5zaWduaW4udXNlci5hZG1pbiIsImF1dGhfdGltZSI6MTY5MjY0MDI4NCwiZXhwIjoxNjkyNjQzODg0LCJpYXQiOjE2OTI2NDAyODQsImp0aSI6ImMxYzViMzRiLTRkZDAtNDJkNC04NTM5LWE4OWJkMGJiMDU5MCIsInVzZXJuYW1lIjoiYjQwOGU0MjgtNjA3MS03MDE1LWJjYWQtNDQ3OTQ5ZDYzYTdiIn0.QCB55BQUYB7c8nFcBYGMqTZ2WxChDODHZUoN_QJV9vmATxWmL58uw6qq7LlkhY6d2UMhCQLpXxgR-pXQATCSAOume1UIuyTh1Wjmj5567zpfMd5d1RGymHr0QzekIRthflanpQnQvGr-PcXnFJcnNRXDNBKC569mSASj00Hnm8Y6gHh9gKka2QvBQIEyBpG6ylI7y3nWC4UFaK_tTEp7LGcQsz7EtEG9AW8EsDwl9OI0yPYwNU0wPLdBUp3AosqjAhw2FndJJ3iVHVmklxcvQTSQi9KTN2cH5H8TlbpDLWMf7ONHEUx6Dv5FSZsYmv9w37oKa7bbASmCZwu8ej_1vw",
        "ExpiresIn": 3600,
        "TokenType": "Bearer",
        "RefreshToken": "eyJjdHkiOiJKV1QiLCJlbmMiOiJBMjU2R0NNIiwiYWxnIjoiUlNBLU9BRVAifQ.hZl88yQUkxSSNiR-MJP3kHnFGlgncR_lMvbBjef9mnGP99K7uuXtz4AFZYlCPGkmn8kgTZOVpywvc89TuZoN6rxSnxLGResZ5Dq1dO-cHH-UNqIjCueoMsKq7CSeaypfnwehOirMJm1HchL0j6tCyF2dZyaXqj5t6lipRrwimNqnzKmd8wHmtg5KdZ2ebwqi-epDUKJz6M8WWzQc8AG8B0MoeNKyYzAQXq9-S1IBrNmKsF3Zh9bHtm97qph8Bey2G8gmbvoLIbFEB86Bt6h820BRUhMTk4FWTYx5BHLqnA787Z8kRcU6SYItDPNxxNiVGrJtGK_wWEE37h77A8TW4g.kUgwnWUoLqUPcYmg.R8eYH-2xwGStWqfpvtdX57IDJwpENez0oF0RwTCMryZ5tsBmOtvt8-Savn-lZEdHgRSu9oytb_Yn2hY3DIWixTsCb-9uXR1eKbkMESJkVWlLllpoSSdhQlBZivQcle9cMZqoDAzIgefF1i_kRbM8hfRJG6xvvtArpmZUlG-kwxieLTjgz3RwPlcL4hMFzZilCPDmen_UnuJ-f0Htn6UDTjAcL-AhBYmhU0He0rnwHhJdOSiHjTj5cTFQBGHoKgWsphu_A73O9fbZEhL_Hf04ZfnL_09AmceGo3l_4MnxbzoT1VFaBOZeH2cj2zNV9O_wquDWg_kldtWiUJpZ6lgRpw1oBhLbaZs92KUE8K9ejn000gl-uNsWgX-b3lA0hujF8tzjyEI2PCJcFP4DTDnrYguOEEkRQ7quBcVXkDHt_0HE90CU9Uk5zc8-091CNosVBa2y1f0qNAy5861KnAbetK6fVyl7yqyIjwaoLiG3d71mBoS0kIAyynzY7SKYlhj4Pn1gV4Fa62qh1DhwMUBZHRpwSUDoysnQfZfMoN2MVssIpBciiVq3f5UvQOBnWdGWRJ2WZzKmgUdKTA_P4gaDs39tWNBba7yl_HcqTCyM2XtBJilnzNhlfGm7mp66lKubEm1TNLZ7FN9E6mj4L2ZvqEGTKBXWz-0D-HvutBEeatS08l5oJLqfFReblitYJVxlO6fvI-MVDlHotOtF6-11wIekCh9ifX_1l4Af8vBJMEpU6bd7gToP1dGn1vY70BuXg8BSXoU49pDmnlzzDSbuJbthdItSBGt7Nrs0_7orzEkAN4C7MDXYKWalb0MZHbfUgoe4UA3nMu0QIpMn3VdrXUYh9-z3NrJ4qxwNHjy3nwrDkBy9s1gHl9u9NGCA_gP31O6oyY-ENrdYKfaJhhX4JE5HHUkYYfI2jORdVDnE2vSfSOorKN8F-5UkwEfc-ufHUJ1iqx1b0kJ-43Zl1SXlwLkOidhjs8EA0zEinKfCwtKosxDkARYmi9xWJqvZ_PcU485gMVlPd6e9qkYs5vz75Rji9asaLQvov0KDH8I1ZM4k5eNPzBjkZcOMbM3jud2B3SZAuwvyFU51GvBnwTmQcAdDQTHaU-vXo8hU7_pJE_5L6_HKVYbWrSKRY6JAFOsCeN6eaBjAeTd29C2Lq4Zh7iHlHsyog5jNb1WG0VDMAkW-aEd_9pqZhoyNwhlxcoHcZrh2V_RF9tpCKa94FZlyYGBEMyL_YxLT3aDr3gXyJLMjIZxVOjZ09GLzfYXGWgRCXIW2UAOf6uqQt9VA99jMGaaahW-zg0lQSPDQEjGlDQ4sHZMO6lPbym7a.QN0pMpJSeoE_q1fN1nQ4Wg",
        "IdToken": "eyJraWQiOiIxKzVwTHhSK0ljUXJpMTdZWlNYVTZiSTlIVXNqaXhadWVhVnA0WnE0ckN3PSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiJiNDA4ZTQyOC02MDcxLTcwMTUtYmNhZC00NDc5NDlkNjNhN2IiLCJpc3MiOiJodHRwczpcL1wvY29nbml0by1pZHAudXMtZWFzdC0xLmFtYXpvbmF3cy5jb21cL3VzLWVhc3QtMV8xTFBMRUtQUDAiLCJjb2duaXRvOnVzZXJuYW1lIjoiYjQwOGU0MjgtNjA3MS03MDE1LWJjYWQtNDQ3OTQ5ZDYzYTdiIiwib3JpZ2luX2p0aSI6ImRiMWMyYjY2LWE2OTEtNDRkYi04NDA2LTg3NDdkZjViM2FiOCIsImF1ZCI6Imhzb2phbmc1NWw1cjEyNWdyZ3M2OThlNjgiLCJldmVudF9pZCI6ImEwMjQwNmM4LTE5MmQtNDBhNi05NjU2LWEwMWNlMDE5NTgyYyIsInRva2VuX3VzZSI6ImlkIiwiYXV0aF90aW1lIjoxNjkyNjQwMjg0LCJleHAiOjE2OTI2NDM4ODQsImlhdCI6MTY5MjY0MDI4NCwianRpIjoiMGM1N2FlZTMtOTE4Yi00MjFiLThkZGUtNTM3NWQ5ZWJlNWQ3IiwiZW1haWwiOiJpb3Qtb25ib2FyZGluZy1hZG1pbi0yMDIzLTA4LTIxLTIzLTIxLTE4QGV4YW1wbGUuY29tIn0.HU17bxEFHkLpW-_0SlUk7gg-UqKYYrtpEs-wmarqtvQnEV-UiEbBBDqQbFRbNx640Yi_UbSHsrS47Z3OV_vpbn93tf0of7mX6Fcv_uZPxPRAkU6Hz11ZAAETQm6prLJYSzYgO3wzYelni2SIOWVJaf8Lvrp3VUennJy7y5Ha_O1YEGFdl-aS6eSgS6QwokqdTkvcBZ9lDyg5PhwuFurbFR7qBxs2mGkcRzqi67qPBLhWAb0lh5YVEQoOYFiMVkO5_RyfuhDxE-5hlCdWR9VVB09bK3LIxKButDliS2hPppeCLULYg2h4JVmXHYrcl8W8DBTBqKd_JbXpQk1IJeM2qA"
    }
}

    <h3>AWS IOT Connectivity QuickStart Output Values</h3></br>
    --------------------------------------------------------------------------------------------</br>
    | Cognito URL             |  https://iot-onboarding-quickstart-657907747545-dev.auth.us-east-1.amazoncognito.com/oauth2/token</br>
    --------------------------------------------------------------------------------------------</br>
    | API Gateway URL         |  https://q2uwvnxoud.execute-api.us-east-1.amazonaws.com/</br>
    --------------------------------------------------------------------------------------------</br>
    | Client ID               |  hsojang55l5r125grgs698e68</br>
    --------------------------------------------------------------------------------------------</br>
    | Refresh Token           |  eyJjdHkiOiJKV1QiLCJlbmMiOiJBMjU2R0NNIiwiYWxnIjoiUlNBLU9BRVAifQ.hZl88yQUkxSSNiR-MJP3kHnFGlgncR_lMvbBjef9mnGP99K7uuXtz4AFZYlCPGkmn8kgTZOVpywvc89TuZoN6rxSnxLGResZ5Dq1dO-cHH-UNqIjCueoMsKq7CSeaypfnwehOirMJm1HchL0j6tCyF2dZyaXqj5t6lipRrwimNqnzKmd8wHmtg5KdZ2ebwqi-epDUKJz6M8WWzQc8AG8B0MoeNKyYzAQXq9-S1IBrNmKsF3Zh9bHtm97qph8Bey2G8gmbvoLIbFEB86Bt6h820BRUhMTk4FWTYx5BHLqnA787Z8kRcU6SYItDPNxxNiVGrJtGK_wWEE37h77A8TW4g.kUgwnWUoLqUPcYmg.R8eYH-2xwGStWqfpvtdX57IDJwpENez0oF0RwTCMryZ5tsBmOtvt8-Savn-lZEdHgRSu9oytb_Yn2hY3DIWixTsCb-9uXR1eKbkMESJkVWlLllpoSSdhQlBZivQcle9cMZqoDAzIgefF1i_kRbM8hfRJG6xvvtArpmZUlG-kwxieLTjgz3RwPlcL4hMFzZilCPDmen_UnuJ-f0Htn6UDTjAcL-AhBYmhU0He0rnwHhJdOSiHjTj5cTFQBGHoKgWsphu_A73O9fbZEhL_Hf04ZfnL_09AmceGo3l_4MnxbzoT1VFaBOZeH2cj2zNV9O_wquDWg_kldtWiUJpZ6lgRpw1oBhLbaZs92KUE8K9ejn000gl-uNsWgX-b3lA0hujF8tzjyEI2PCJcFP4DTDnrYguOEEkRQ7quBcVXkDHt_0HE90CU9Uk5zc8-091CNosVBa2y1f0qNAy5861KnAbetK6fVyl7yqyIjwaoLiG3d71mBoS0kIAyynzY7SKYlhj4Pn1gV4Fa62qh1DhwMUBZHRpwSUDoysnQfZfMoN2MVssIpBciiVq3f5UvQOBnWdGWRJ2WZzKmgUdKTA_P4gaDs39tWNBba7yl_HcqTCyM2XtBJilnzNhlfGm7mp66lKubEm1TNLZ7FN9E6mj4L2ZvqEGTKBXWz-0D-HvutBEeatS08l5oJLqfFReblitYJVxlO6fvI-MVDlHotOtF6-11wIekCh9ifX_1l4Af8vBJMEpU6bd7gToP1dGn1vY70BuXg8BSXoU49pDmnlzzDSbuJbthdItSBGt7Nrs0_7orzEkAN4C7MDXYKWalb0MZHbfUgoe4UA3nMu0QIpMn3VdrXUYh9-z3NrJ4qxwNHjy3nwrDkBy9s1gHl9u9NGCA_gP31O6oyY-ENrdYKfaJhhX4JE5HHUkYYfI2jORdVDnE2vSfSOorKN8F-5UkwEfc-ufHUJ1iqx1b0kJ-43Zl1SXlwLkOidhjs8EA0zEinKfCwtKosxDkARYmi9xWJqvZ_PcU485gMVlPd6e9qkYs5vz75Rji9asaLQvov0KDH8I1ZM4k5eNPzBjkZcOMbM3jud2B3SZAuwvyFU51GvBnwTmQcAdDQTHaU-vXo8hU7_pJE_5L6_HKVYbWrSKRY6JAFOsCeN6eaBjAeTd29C2Lq4Zh7iHlHsyog5jNb1WG0VDMAkW-aEd_9pqZhoyNwhlxcoHcZrh2V_RF9tpCKa94FZlyYGBEMyL_YxLT3aDr3gXyJLMjIZxVOjZ09GLzfYXGWgRCXIW2UAOf6uqQt9VA99jMGaaahW-zg0lQSPDQEjGlDQ4sHZMO6lPbym7a.QN0pMpJSeoE_q1fN1nQ4Wg</br>
    --------------------------------------------------------------------------------------------</br>
    | Environment             |  dev</br>
    --------------------------------------------------------------------------------------------</br>
    | Region                  |  us-east-1</br>
    --------------------------------------------------------------------------------------------</br>
    | Cognito User Pool ID    |  us-east-1_1LPLEKPP0</br>
    --------------------------------------------------------------------------------------------</br>
    | Password                |  testPassword1$26867</br>
    --------------------------------------------------------------------------------------------</br>
    | Glue DB Name            |  iot-onboarding-sensors-data-dev</br>
    -------------------------------------------------------------------------------------------- </br>
    | Athena Table Name       |  iotonboardinginfrastackd_iotonboardingsensorsdata_169ivu6qiwkyf</br>
    --------------------------------------------------------------------------------------------</br>
    | Iot Sitewise Role       |  arn:aws:iam::657907747545:role/IOTOnboardingInfraStackde-iotOnboardingIotSitewise-1A7K1OW29C5FQ</br>
    --------------------------------------------------------------------------------------------</br>