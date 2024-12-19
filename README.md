# twitter-clone-backend

## How to run
```
air
```

## Development Process
### Stage 1
Features:
- [x] Sign up
- [x] Login
- [x] Tweet (create, edit, delete)
- [x] Follow/unfollow

## Stage 2
Implement layered architecture & request validation. Developed in '/v2/{...}' endpoint.

###
Improvements:
- [x] E2E Testing
- [x] Refactored to layered architecture

Features:
- [x] Like

## Stage 3
- [ ] Implement redis for cache user profile (including user last 10 tweets)

## ERD
Migration version: 20241218063017
![Entity Relationship Diagram](/images/erd.png)

## Caching Mechanism
User profile containing below information is getting cached.
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
todo: does check $ exist (json) return same response as check if key exist (check key immediately)?
```
is $ (root key) exist?
    yes
    is $.recentTweets exist?
        yes
        return cache

        no
        get recentTweets from db
        store recentTweets to cache
        return (db + cache)
    no
    get userProfile with recentTweets from db
    store it to cache
    return db
```

### Expiration
Everytime user profile added to cache, the expiration will set to 10 minutes.
User profile can only meet its expiration date until the end if only these operations are performed:
- json.set partial update on recentTweets
- json.del partial update on recentTweets

Both operation performed when user create a new tweet. The other operation will clear the expiration (key will no longer expire)

The expiration will be extended when **user profile get accessed**.
