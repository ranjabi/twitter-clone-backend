# Twitter Clone Backend

<p align="center">
<img src="images/thumbnail.png" alt="Twitter Clone" style="width: 60%; height: 60%; align: center"/>
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
A clone of the Twitter web application with features such as tweets, user profiles, and a news feed. I aimed to minimize the use of external Go libraries, leveraging Go's native features as much as possible. Some of the things I implemented can be found under [Leeson Learned](#lesson-learned).

The frontend repository can be found at: https://github.com/ranjabi/twitter-clone-frontend

### Build With
- Go
- PostgreSQL
- Redis
- CI/CD: Docker, Github Actions

## ERD
![Entity Relationship Diagram](/images/erd.png)

## Caching Mechanism
The user profile displays information such as the username, follower/following count, and recent tweets. The response structure is shown below, with the following information being cached.
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

The decision tree of caching can be seen below. The cache is saved as `key=user.id:{id}, value = userResponse` for user profile information and `key=user.id:{id}:recentTweets, value = recentTweets` for the recent tweets.
```
is user.id:{id} exist?
    yes
    is user.id:{id}:recentTweets exist?
        yes
        -> return cache

        no
        cache miss for user.id:{id}:recentTweets
        get recentTweets from db and store it to cache
        -> return the combined data from db (recentTweets) and cache (userProfile)
    no
    get userProfile with recentTweets from db
    store it to cache
    -> return fully from db
```

### Cache Expiration
Each time a user profile is added to the cache, the expiration time is set to 10 minutes.
The following operations will delete the user's recent tweets cache:
- Creating, updating, or deleting a tweet
- Liking or unliking a tweet by a follower

As a result, the cache will be recreated when a follower requests the user's profile. Additionally, the expiration time resets to 10 minutes each time **the user's profile is accessed**.

## Development Process
### Stage 1
Features:
- [x] Register
- [x] Login
- [x] Tweet (create, edit, delete)
- [x] Follow/unfollow
- [x] E2E Testing

### Stage 2
Improvements:
- [x] Refactored to layered architecture (handler, service, and repository)
- [x] Request validation

Features:
- [x] Like/unlike

### Stage 3
Features:
- [x] Implemented redis for user profile caching
- [x] News feed

## Lesson Learned
- Defined a custom HTTP appHandler for error handling, allowing structured error responses to be sent to the client.
- Implemented logging by capturing incoming HTTP requests. This logs the error method, request URL, and request body.
