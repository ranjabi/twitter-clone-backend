package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/ranjabi/twitter-clone-backend/config"
	"github.com/ranjabi/twitter-clone-backend/db"
	"github.com/ranjabi/twitter-clone-backend/messagebroker"
	"github.com/ranjabi/twitter-clone-backend/types"
	"github.com/ranjabi/twitter-clone-backend/usecases/auth"
	"github.com/ranjabi/twitter-clone-backend/usecases/tweet"
	"github.com/ranjabi/twitter-clone-backend/usecases/user"
)

/*
	Create 11 users
	Create 3 tweets by every user, interval now, -1 hour, -2 hours

	User 1 follows other 10 users

	User 1 feed will be filled by tweets from user2 - user11
*/

func main() {
	fmt.Println("Starting seed programatically ...")

	// Setup
	env := ".env.dev.local"
	err := godotenv.Load(".env.dev.local")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("Loaded ENV:", env)

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	if err = db.SetupConnection(ctx, cfg); err != nil {
		log.Fatal(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	migrationsPath := filepath.Join(cwd, "db", "migrations")
	actions := []string{"migrate.reset", "migrate.up"}
	err = db.ApplyMigrationsAndSeed(ctx, cfg, actions, migrationsPath, "", false)
	if err != nil {
		log.Fatal(err)
	}

	rabbitMq := messagebroker.NewRabbitMq(ctx, db.RbConn, db.RbCh)
	if err = rabbitMq.DeclareFeedQeueu(); err != nil {
		log.Fatal(err)
	}

	userRepository := user.NewRepository(ctx, db.PgConn, db.RdConn)
	tweetRepository := tweet.NewRepository(ctx, db.PgConn, db.RdConn)

	authService := auth.NewService(ctx, cfg, userRepository)

	// Create 10 users
	fmt.Println("Starting create 10 users")
	var users []types.User
	for i := 1; i <= 10; i++ {
		_user := types.User{
			FullName: fmt.Sprintf("User %d", i),
			Username: fmt.Sprintf("user%d", i),
			Email:    fmt.Sprintf("user%d@gmail.com", i),
			Password: "password",
		}
		newUser, err := authService.Register(_user)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, *newUser)
	}
	fmt.Println("Finished create 10 users")

	// Create 3 tweets by every user, interval now, -1 hour, -2 hours
	fmt.Println("Starting create 3 tweets for each user")
	for _, _user := range users {
		for i := 1; i <= 3; i++ {
			_tweet := types.Tweet{
				Content: fmt.Sprintf("Tweet %d", i),
				User:    _user,
			}
			newTweet, err := tweetRepository.Create(_tweet)
			if err != nil {
				log.Fatal(err)
			}

			query := `
				UPDATE
					tweets
				SET
					created_at = @created_at,
					modified_at = @modified_at
				WHERE
					id = @id
			`

			var createdAt time.Time
			var modifiedAt time.Time
			switch i {
			case 1:
				createdAt = time.Now()
				modifiedAt = time.Now()
			case 2:
				createdAt = time.Now().Add(-1 * time.Hour)
				modifiedAt = time.Now().Add(-1 * time.Hour)
			case 3:
				createdAt = time.Now().Add(-2 * time.Hour)
				modifiedAt = time.Now().Add(-2 * time.Hour)
			}

			args := pgx.NamedArgs{
				"created_at":  createdAt,
				"modified_at": modifiedAt,
				"id":          newTweet.Id,
			}
			if _, err := db.PgConn.Exec(ctx, query, args); err != nil {
				log.Fatal(err)
			}
		}
	}
	fmt.Println("Finished create 3 tweets for each user")

	// User 1 follows other 10 users
	fmt.Println("Starting user1 follow user2 - user10")
	for i := 2; i <= 10; i++ {
		if err := userRepository.FollowOtherUser(1, i); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Finished user1 follow user2 - user10")
}
