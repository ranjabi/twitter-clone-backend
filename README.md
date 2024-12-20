# Twitter Clone Backend

<p align="center">
<img src="images/thumbnail.png" alt="Twitter Clone" style="width: 40%; height: 40%; align: center"/>
</p>

## Table of Contents
* [About The Project](#about-the-project)
    + [Build With](#build-with)
* [Development Process](#development-process)
    + [Stage 1](#stage-1)
    + [Stage 2](#stage-2)
    + [Stage 3](#stage-3)
* [ERD](#erd)
* [Caching Mechanism](#caching-mechanism)
    + [Cache Expiration](#cache-expiration)
* [Lesson Learned](#lesson-learned)

## About The Project
A clone of twitter web application with features like tweets, user profile, and news feed. I tried to use go library as little as possible, maximizing the existing features of GO. Some of the things that I implemented can be seen at [Leeson Learned](#lesson-learned)

### Build With
- Go
- PostgreSQL
- Redis
- CI/CD: Docker, Github Actions

## Development Process
### Stage 1
Features:
- [x] Register
- [x] Login
- [x] Tweet (create, edit, delete)
- [x] Follow/unfollow
- [x] E2E Testing

### Stage 2
From this stage, the development was done in `/v2/...` endpoint prefix.

Improvements:
- [x] Refactored to layered architecture (handler, service, and repository)
- [x] Request validation

Features:
- [x] Like/unlike

### Stage 3
- [x] Implemented redis for user profile caching
 <!-- (including user last 10 tweets) -->
- [x] News feed

## ERD
![Entity Relationship Diagram](/images/erd.png)

## Caching Mechanism
User profile display the information about a user, such as username, follower/following count, and recent tweets. The structure of response looks like below. containing below information is getting cached.
```
userResponse := struct {
    Id                 int            `json:"id"`
    Username           string         `json:"username"`
    FollowerCount      int            `json:"followerCount"`
    FollowingCount     int            `json:"followingCount"`
    RecentTweetsLength int            `json:"recentTweetsLength"`
    RecentTweets       []models.Tweet `json:"recentTweets"`
}
```

The decision tree of caching can be seen below. Cache saved as `key=user.id:{id}, value = userResponse`.
```
is $ (root key) exist?
    yes
    is $.recentTweets exist?
        yes
        -> return cache

        no
        cache miss for $.recentTweets
        get recentTweets from db and store it to cache
        -> return data from db (recentTweets) + cache (userProfile)
    no
    get userProfile with recentTweets from db
    store it to cache
    -> return fully from db
```

### Cache Expiration
Everytime user profile is added to cache, the expiration time will set to 10 minutes.
User profile can only meet its expiration date until the end (and be deleted after that) if only these operations are performed:

- `JSON.SET` running partial update on `$.recentTweets`
- `JSON.DEL` running partial update on `$.recentTweets`

Both operation performed when a user create a new tweet. The expiration time will be reset again to 10 minutes when **user profile get accessed**.

## Lesson Learned
- Define my own HTTP appHandler for custom error handling. This made me able to structure the error response to be send to the client.
- Implemented logging by reading the incoming HTTP request. The result is error method, request url, and request body can be seen in the logs.

