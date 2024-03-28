# API Documentation

## Likes

### Endpoint 1: /likes/thread/:threadid/likes

#### Description:
Get the list of users who like a certain thread

#### Request:
- Method: GET

#### Response:
- Status Code: 200 [OK]
- Body:
    - likes : List of the user id of user who like the thread (array of integer)
    - status: Whether the API call is successful ("success" or "failure") 
    - Example:
        1. In the absence of any user liking the thread
            ```json
            {
                "likes": null,
                "status": "success"
            }
            ```

        2. When there are users (user 1 and user 2) who like the thread
            ```json
            {
                "likes": [
                    {
                        "Id": 1,
                        "Username": "firstUser"
                    },
                    {
                        "Id": 2,
                        "Username": "secondUser"
                    }
                ],
                "status": "success"
            }
            ```

### Endpoint 2: /likes/thread/:threadid/dislikes

#### Description:
Get the list of users who dislike a certain thread

#### Request:
- Method: GET

#### Response:
- Status Code: 200 [OK]
- Body:
    - dislikes : List of the user id of user who dislike the thread (array of integer)
    - status: Whether the API call is successful ("success" or "failure")
- Example: 
    1. In the absence of any user disliking the thread
        ```json
        {
            "dislikes": null,
            "status": "success"
        }
        ```
    2. When there are users (user 1 and user 2) who dislike the thread
        ```json
        {
                "dislikes": [
                    {
                        "Id": 1,
                        "Username": "firstUser"
                    },
                    {
                        "Id": 2,
                        "Username": "secondUser"
                    }
                ],
                "status": "success"
            }
        ```

### Endpoint 3: /likes/comment/:commentid/likes

#### Description:
Get the list of users who like a certain comment

#### Request:
- Method: GET

#### Response:
- Status Code: 200 [OK]
- Body:
    - likes : List of the user id of user who like the comment (array of integer)
    - status: Whether the API call is successful ("success" or "failure") 
    - Example:
        1. In the absence of any user liking the comment
            ```json
            {
                "likes": null,
                "status": "success"
            }
            ```

        2. When there are users (user 1 and user 2) who like the comment
            ```json
            {
                "likes": [
                    {
                        "Id": 1,
                        "Username": "firstUser"
                    },
                    {
                        "Id": 2,
                        "Username": "secondUser"
                    }
                ],
                "status": "success"
            }
            ```

### Endpoint 4: /likes/comment/:commentid/dislikes

#### Description:
Get the list of users who dislike a certain comment

#### Request:
- Method: GET

#### Response:
- Status Code: 200 [OK]
- Body:
    - dislikes : List of the user id of user who dislike the comment (array of integer)
    - status: Whether the API call is successful ("success" or "failure")
- Example: 
    1. In the absence of any user disliking the comment
        ```json
        {
            "dislikes": null,
            "status": "success"
        }
        ```
    2. When there are users (user 1 and user 2) who dislike the comment
        ```json
        {
            "dislikes": [
                {
                    "Id": 1,
                    "Username": "firstUser"
                },
                {
                    "Id": 2,
                    "Username": "secondUser"
                }
            ],
            "status": "success"
        }
        ```
