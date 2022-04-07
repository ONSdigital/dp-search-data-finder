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

    Scenario: App health returning a WARNING (429) status as search reindex api health is in warning state
        Given the Search Reindex API is unhealthy with state warning
        And I wait "1500" milliseconds for the healthcheck to update
        When I GET "/health"
        Then the HTTP status code should be "429"
        And the response header "Content-Type" should be "application/json; charset=utf-8"
        And I should receive the following health JSON response:
        """
            {
                "status":"WARNING",
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
                        "status":"WARNING",
                        "status_code":429,
                        "message":"dp-search-reindex-api is degraded, but at least partially functioning"
                    }
                ]
            }
        """

    Scenario: App health returning WARNING (429), then CRITICAL (500) as search reindex api health is in critical state
        Given the Search Reindex API is unhealthy with state critical
        And I wait "1500" milliseconds for the healthcheck to update
        When I GET "/health"
        Then the HTTP status code should be "429"
        And the response header "Content-Type" should be "application/json; charset=utf-8"
        And I should receive the following health JSON response:
        """
            {
                "status":"WARNING",
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
                        "status":"CRITICAL",
                        "status_code":500,
                        "message":"dp-search-reindex-api functionality is unavailable or non-functioning"
                    }
                ]
            }
        """
        And I wait "2000" milliseconds for the healthcheck to update after surpassing the critical timeout
        When I GET "/health"
        Then the HTTP status code should be "500"
        And the response header "Content-Type" should be "application/json; charset=utf-8"
        And I should receive the following health JSON response:
        """
            {
                "status":"CRITICAL",
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
                        "status":"CRITICAL",
                        "status_code":500,
                        "message":"dp-search-reindex-api functionality is unavailable or non-functioning"
                    }
                ]
            }
        """