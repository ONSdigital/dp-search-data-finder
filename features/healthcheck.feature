Feature: Healthcheck endpoint should inform the health of service

    Scenario: App health returning a OK (200) status when all downstream apps are healthy 
        Given all of the downstream services are healthy
        And I wait "1500" milliseconds for the healthcheck to update
        When I GET "/health"
        Then the HTTP status code should be "200"
        And the response header "Content-Type" should be "application/json; charset=utf-8"
        And I should receive the following health JSON response:
        """
            {
                "status":"OK",
                "version":{
                    "git_commit":"3t7e5s1t4272646ef477f8ed755",
                    "language":"go",
                    "language_version":"go1.16.5",
                    "version":"v1.2.3"
                },
                "checks":[
                    {
                        "name":"Kafka consumer",
                        "status":"OK",
                        "status_code":0,
                        "message":"OK"
                    },
                    {
                        "name":"Zebedee",
                        "status":"OK",
                        "status_code":200,
                        "message":"zebedee is ok"
                    },
                    {
                        "name":"Search Reindex API",
                        "status":"OK",
                        "status_code":200,
                        "message":"dp-search-reindex-api is ok"
                    }
                ]
            }
        """
