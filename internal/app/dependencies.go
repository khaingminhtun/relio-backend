package app

import (
	"github.com/khaingminhtun/production-go-api/internal/features/auth"
	"github.com/khaingminhtun/production-go-api/internal/features/relationship"
	"github.com/khaingminhtun/production-go-api/internal/features/user"
	redisinfra "github.com/khaingminhtun/production-go-api/internal/infrastructure/redis"
	transaction "github.com/khaingminhtun/production-go-api/internal/shared/dbutils"
	"github.com/khaingminhtun/production-go-api/internal/shared/security"
	"gorm.io/gorm"
)

type Dependencies struct {
	UserHandler         *user.Handler
	AuthHandler         *auth.Handler
	RelationshipHandler *relationship.Handler

	JWTManager *security.JWTManager
}

func NewDependencies(db *gorm.DB,
	redisStore redisinfra.RedisStore,
	emailQueue redisinfra.EmailQueue,

	jwtManager *security.JWTManager,
) *Dependencies {

	//transaction Manager
	txManager := transaction.NewManager(db)

	//Repository
	userRepository := user.NewRepository(db)
	authRepository := auth.NewRepository(db)
	relationshipRepository := relationship.NewRelationshipRepository(db)
	memberRepository := relationship.NewRelationshipMemberRepository(db)

	//Service
	userService := user.NewService(userRepository)
	authService := auth.NewService(userRepository, authRepository, redisStore, emailQueue, jwtManager)
	relationshipService := relationship.NewService(relationshipRepository, memberRepository, txManager)

	//Handler
	userHandler := user.NewHandler(userService)
	authHandler := auth.NewHandler(authService)
	relationshipHandler := relationship.NewHandler(relationshipService)

	return &Dependencies{
		UserHandler:         userHandler,
		AuthHandler:         authHandler,
		RelationshipHandler: relationshipHandler,

		JWTManager: jwtManager,
	}
}
