package it

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"twitter-clone-backend/config"
	"twitter-clone-backend/db"
	"twitter-clone-backend/models"
	"twitter-clone-backend/usecases/tweet"
	"twitter-clone-backend/usecases/user"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	pgConn         *pgxpool.Pool
	rdConn         *redis.Client
	ctx            context.Context
	cfg            *config.Config
	migrationsPath string

	userRepository  user.UserRepository
	tweetRepository tweet.TweetRepository

	userService  user.Service
	tweetService tweet.Service

	validUser     *models.User
	validTweet    *models.Tweet
	notExistTweet models.Tweet
}

func TestRunWithSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) SetupSuite() {
	var err error
	s.ctx = context.Background()

	err = godotenv.Load("../../.env.test.local")
	s.NoError(err)

	s.cfg, err = config.Load()
	s.NoError(err)

	s.pgConn, s.rdConn, err = db.SetupConnection(s.ctx, s.cfg)
	s.NoError(err)

	s.userRepository = user.NewRepository(s.ctx, s.pgConn, s.rdConn)
	s.tweetRepository = tweet.NewRepository(s.ctx, s.pgConn, s.rdConn)

	s.userService = user.NewService(s.ctx, s.cfg, s.userRepository)
	s.tweetService = tweet.NewService(s.tweetRepository, s.userRepository)

	cwd, err := os.Getwd()
	s.NoError(err)

	s.migrationsPath = filepath.Join(cwd, "..", "..", "db", "migrations")
	actions := []string{"migrate.reset"}
	err = db.ApplyMigrationsAndSeed(s.ctx, s.cfg, actions, s.migrationsPath, "", true)
	s.NoError(err)
}

func (s *TestSuite) TearDownSuite() {
	s.pgConn.Close()
}

func (s *TestSuite) SetupTest() {
	actions := []string{"migrate.up"}
	err := db.ApplyMigrationsAndSeed(s.ctx, s.cfg, actions, s.migrationsPath, "", true)
	s.NoError(err)

	s.validUser, err = s.userRepository.CreateUser(models.User{
		Id:           1,
		Email:        "test@example.com",
		Username:     "testusername",
		FullName:     "Test Full Name",
		Password:     "$2a$14$ZqZ1FmMgZNYvO.Q2rSht3.fGTX4IBq6VJMBoJ7bRXMAaEQk3pAP9i",
		ProfileImage: "https://twitter-clone-tzjvdg.s3.ap-southeast-1.amazonaws.com/purple-1.png",
	})
	s.NoError(err)
	s.NotNil(s.validUser)

	s.validTweet, err = s.tweetRepository.CreateTweet(models.Tweet{
		Id:      1,
		Content: "Content",
		UserId:  s.validUser.Id,
	})
	s.NoError(err)
	s.NotNil(s.validTweet)

	s.notExistTweet = models.Tweet{
		Id: 99999,
	}
}

func (s *TestSuite) TearDownTest() {
	// actions := []string{"migrate.reset"}
	// err := db.ApplyMigrationsAndSeed(s.ctx, s.cfg, actions, s.migrationsPath, "", true)
	// s.NoError(err)
}

var (
//	validUser2 = models.User{
//		Id:       2,
//		Email:    "test2@example.com",
//		Username: "test2",
//		FullName: "Test test 2",
//		Password: "password",
//	}
//
//	notExistUser = models.User{
//		Id: 100,
//	}
//
//	validTweet = models.Tweet{
//		Id:      1,
//		Content: "content",
//		UserId:  1,
//	}
//

)
