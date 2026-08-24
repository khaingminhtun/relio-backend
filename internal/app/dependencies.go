package app

import (
	"github.com/khaingminhtun/relio-backend/internal/features/auth"
	"github.com/khaingminhtun/relio-backend/internal/features/relationship"
	"github.com/khaingminhtun/relio-backend/internal/features/user"
	redisinfra "github.com/khaingminhtun/relio-backend/internal/infrastructure/redis"
	transaction "github.com/khaingminhtun/relio-backend/internal/shared/dbutils"
	"github.com/khaingminhtun/relio-backend/internal/shared/security"
	"gorm.io/gorm"
)

type Dependencies struct {
	UserHandler              *user.Handler
	AuthHandler              *auth.Handler
	RelationshipHandler      *relationship.RelationshipHandler
	RelatonshipMemberHandler *relationship.RelationshipMemberHandler

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
	relationshipService := relationship.NewRelationshipService(relationshipRepository, memberRepository, txManager)
	relationshipMemberService := relationship.NewRelationshipMemberService(
		memberRepository,
	)

	//Handler
	userHandler := user.NewHandler(userService)
	authHandler := auth.NewHandler(authService)
	relationshipHandler := relationship.NewHandler(relationshipService)
	relationshipMemberHandler := relationship.NewRelationshipMemberHandler(
		relationshipMemberService,
	)

	return &Dependencies{
		UserHandler:              userHandler,
		AuthHandler:              authHandler,
		RelationshipHandler:      relationshipHandler,
		RelatonshipMemberHandler: relationshipMemberHandler,

		JWTManager: jwtManager,
	}
}
