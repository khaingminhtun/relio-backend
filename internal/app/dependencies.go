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
	ProfileHandler           *user.UserProfileHandler
	AuthHandler              *auth.Handler
	RelationshipHandler      *relationship.RelationshipHandler
	RelatonshipMemberHandler *relationship.RelationshipMemberHandler
	InvitationHandler        *relationship.InvitationHandler

	JWTManager *security.JWTManager
}

func NewDependencies(db *gorm.DB,
	redisStore redisinfra.RedisStore,
	emailQueue redisinfra.EmailQueue,
	appBaseURL string,
	jwtManager *security.JWTManager,
) *Dependencies {

	//transaction Manager
	txManager := transaction.NewManager(db)

	//Repository
	userRepository := user.NewRepository(db)
	userProfileRepository := user.NewUserProfileRepository(db)
	authRepository := auth.NewRepository(db)
	relationshipRepository := relationship.NewRelationshipRepository(db)
	memberRepository := relationship.NewRelationshipMemberRepository(db)
	invitationRepository := relationship.NewInvitationRepository(db)

	//Service
	userService := user.NewService(userRepository)
	userProfileService := user.NewUserProfileService(userProfileRepository)
	authService := auth.NewService(userRepository, userProfileRepository, authRepository, redisStore, emailQueue, jwtManager, txManager)
	relationshipService := relationship.NewRelationshipService(relationshipRepository, memberRepository, txManager)
	relationshipMemberService := relationship.NewRelationshipMemberService(
		memberRepository,
	)
	invitationservice := relationship.NewInvitationService(invitationRepository, emailQueue, appBaseURL)

	//Handler
	userHandler := user.NewHandler(userService)
	userProfileHandler := user.NewUserProfileHandler(userProfileService)
	authHandler := auth.NewHandler(authService)
	relationshipHandler := relationship.NewHandler(relationshipService)
	relationshipMemberHandler := relationship.NewRelationshipMemberHandler(
		relationshipMemberService,
	)
	invitationHandler := relationship.NewInvitationHandler(invitationservice)

	return &Dependencies{
		UserHandler:              userHandler,
		ProfileHandler:           userProfileHandler,
		AuthHandler:              authHandler,
		RelationshipHandler:      relationshipHandler,
		RelatonshipMemberHandler: relationshipMemberHandler,
		InvitationHandler:        invitationHandler,
		JWTManager:               jwtManager,
	}
}
