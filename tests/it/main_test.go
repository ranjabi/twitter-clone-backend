package it

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/ranjabi/twitter-clone-backend/config"
	"github.com/ranjabi/twitter-clone-backend/db"
	"github.com/ranjabi/twitter-clone-backend/messagebroker"
	"github.com/ranjabi/twitter-clone-backend/types"
	"github.com/ranjabi/twitter-clone-backend/usecases/auth"
	"github.com/ranjabi/twitter-clone-backend/usecases/tweet"
	"github.com/ranjabi/twitter-clone-backend/usecases/user"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	rabbitMq       messagebroker.RabbitMq
	ctx            context.Context
	cfg            *config.Config
	migrationsPath string

	userRepository  user.Repository
	tweetRepository tweet.TweetRepository

	authService  auth.Service
	userService  user.Service
	tweetService tweet.Service

	validUser     *types.User
	validUser2    *types.User
	validTweet    *types.Tweet
	notExistUser  types.User
	notExistTweet types.Tweet
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

	err = db.SetupConnection(s.ctx, s.cfg)
	s.NoError(err)

	s.rabbitMq = messagebroker.NewRabbitMq(s.ctx, db.RbConn, db.RbCh)
	err = s.rabbitMq.DeclareFeedQeueu()
	s.NoError(err)

	s.userRepository = user.NewRepository(s.ctx, db.PgConn, db.RdConn)
	s.tweetRepository = tweet.NewRepository(s.ctx, db.PgConn, db.RdConn)

	s.authService = auth.NewService(s.ctx, s.cfg, s.userRepository)
	s.userService = user.NewService(s.ctx, s.userRepository)
	s.tweetService = tweet.NewService(s.tweetRepository, s.userRepository, s.rabbitMq)

	cwd, err := os.Getwd()
	s.NoError(err)

	s.migrationsPath = filepath.Join(cwd, "..", "..", "db", "migrations")
	actions := []string{"migrate.reset"}
	err = db.ApplyMigrationsAndSeed(s.ctx, s.cfg, actions, s.migrationsPath, "", true)
	s.NoError(err)

	err = s.userRepository.ClearRedisCache()
	s.NoError(err)
}

func (s *TestSuite) TearDownSuite() {
	db.PgConn.Close()
}

func (s *TestSuite) SetupTest() {
	actions := []string{"migrate.up"}
	err := db.ApplyMigrationsAndSeed(s.ctx, s.cfg, actions, s.migrationsPath, "", true)
	s.NoError(err)

	s.validUser, err = s.userRepository.Create(types.User{
		Id:       1,
		Email:    "test@example.com",
		Username: "testusername",
		FullName: "Test Full Name",
		Password: "$2a$14$ZqZ1FmMgZNYvO.Q2rSht3.fGTX4IBq6VJMBoJ7bRXMAaEQk3pAP9i",
	})
	s.NoError(err)
	s.NotNil(s.validUser)

	s.validUser2, err = s.userRepository.Create(types.User{
		Id:       2,
		Email:    "test2@example.com",
		Username: "testusername2",
		FullName: "Test Full Name 2",
		Password: "$2a$14$ZqZ1FmMgZNYvO.Q2rSht3.fGTX4IBq6VJMBoJ7bRXMAaEQk3pAP9i",
	})
	s.NoError(err)
	s.NotNil(s.validUser2)

	s.validTweet, err = s.tweetRepository.Create(types.Tweet{
		Id:      1,
		Content: "Content",
		User: types.User{
			Id: s.validUser.Id,
		},
	})
	s.NoError(err)
	s.NotNil(s.validTweet)

	s.notExistUser = types.User{
		Id: 9999,
	}

	s.notExistTweet = types.Tweet{
		Id: 9999,
	}
}

func (s *TestSuite) TearDownTest() {
	actions := []string{"migrate.reset"}
	err := db.ApplyMigrationsAndSeed(s.ctx, s.cfg, actions, s.migrationsPath, "", true)
	s.NoError(err)

	err = s.userRepository.ClearRedisCache()
	s.NoError(err)
}

var (
//	validUser2 = types.User{
//		Id:       2,
//		Email:    "test2@example.com",
//		Username: "test2",
//		FullName: "Test test 2",
//		Password: "password",
//	}
//
//	notExistUser = types.User{
//		Id: 100,
//	}
//
//	validTweet = types.Tweet{
//		Id:      1,
//		Content: "content",
//		UserId:  1,
//	}
//

)
