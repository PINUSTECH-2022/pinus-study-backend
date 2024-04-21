# API Documentation

## Follow

### Endpoint 1: /follow/thread/:userid

#### Description
Get the list of threads posted by users followed by a certain user.

#### Request
- Method: GET
- Param:
    - userid: The id of the user for whom we aim to retrieve a list of threads from users whom they follow (integer)

#### Response
- Status Code: 200 [OK]
- Body:
    - status: Whether the API call is successful ("success" or "failure")
    - threads: List of the id of threads which posted by the users followed by the user specified in the request (array of integer)
    - Example:
        1. In the absence of threads posted by users followed by the user specified in the request
            ```json
            {
                "status": "success",
                "threads": []
            }
            ```
        2. When there are threads (thread 1 and thread 2) posted by users followed by the user specified in the request
            ```json
            {
                "status": "success",
                "threads": [
                    1,
                    2
                ]
            }
            ```

## Likes

### Endpoint 1: /likes/thread/:threadid/likes

#### Description:
Get the list of users who like a certain thread.

#### Request:
- Method: GET
- Param: 
    - threadid: The id of the thread we aim to retrieve the list of users who like it (integer)

#### Response:
- Status Code: 200 [OK]
- Body:
    - likes : List of the user id of user who like the thread (array of integer)
    - status: Whether the API call is successful ("success" or "failure") 
    - Example:
        1. In the absence of any user liking the thread
            ```json
            {
                "likes": [],
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
Get the list of users who dislike a certain thread.

#### Request:
- Method: GET
- Param: 
    - threadid: The id of the thread we aim to retrieve the list of users who dislike it (integer)

#### Response:
- Status Code: 200 [OK]
- Body:
    - dislikes : List of the user id of user who dislike the thread (array of integer)
    - status: Whether the API call is successful ("success" or "failure")
- Example: 
    1. In the absence of any user disliking the thread
        ```json
        {
            "dislikes": [],
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
Get the list of users who like a certain comment.

#### Request:
- Method: GET
- Param: 
    - commentid: The id of the comment we aim to retrieve the list of users who like it (integer)

#### Response:
- Status Code: 200 [OK]
- Body:
    - likes : List of the user id of user who like the comment (array of integer)
    - status: Whether the API call is successful ("success" or "failure") 
    - Example:
        1. In the absence of any user liking the comment
            ```json
            {
                "likes": [],
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
Get the list of users who dislike a certain comment.

#### Request:
- Method: GET
- Param: 
    - commentid: The id of the comment we aim to retrieve the list of users who dislike it (integer)

#### Response:
- Status Code: 200 [OK]
- Body:
    - dislikes : List of the user id of user who dislike the comment (array of integer)
    - status: Whether the API call is successful ("success" or "failure")
- Example: 
    1. In the absence of any user disliking the comment
        ```json
        {
            "dislikes": [],
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
