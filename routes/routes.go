package routes

import (
	"catatan-keuangan/database"
	"catatan-keuangan/handler"
	"catatan-keuangan/middleware"
	"catatan-keuangan/modules/accounts"
	"catatan-keuangan/modules/attachments"
	"catatan-keuangan/modules/auth"
	"catatan-keuangan/modules/categories"
	"catatan-keuangan/modules/groups"
	"catatan-keuangan/modules/invitations"
	"catatan-keuangan/modules/transactions"
	"catatan-keuangan/modules/users"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine) {
	userRepository := users.NewRepository(database.DB)
	userService := users.NewService(userRepository)
	userHandler := handler.NewUserHandler(userService)
	accountRepository := accounts.NewRepository(database.DB)
	accountService := accounts.NewService(accountRepository)
	accountHandler := handler.NewAccountHandler(accountService)
	groupRepository := groups.NewRepository(database.DB)
	groupService := groups.NewService(groupRepository, database.DB)
	groupHandler := handler.NewGroupHandler(groupService)
	invitationRepository := invitations.NewRepository(database.DB)
	invitationService := invitations.NewService(invitationRepository)
	invitationHandler := handler.NewInvitationHandler(invitationService)
	categoryRepository := categories.NewRepository(database.DB)
	categoryService := categories.NewService(categoryRepository, database.DB)
	categoryHandler := handler.NewCategory(categoryService)
	transactionRepository := transactions.NewRepository(database.DB)
	attachmentRepository := attachments.NewRepository(database.DB)
	attachmentService := attachments.NewService(attachmentRepository, "uploads/attachments", "/uploads/attachments")
	transactionService := transactions.NewService(transactionRepository, database.DB, attachmentService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	refreshRepository := auth.NewRefreshRepository(database.DB)
	authService := auth.NewService(userRepository, refreshRepository)
	authHandler := handler.NewAuthHandler(authService)
	v1 := r.Group("/api/v1")
	{
		v1.POST("/auth/login", authHandler.Login)
		v1.POST("/auth/register", authHandler.Register)
		v1.POST("/auth/refresh", authHandler.RefreshToken)
		v1.POST("/auth/logout", authHandler.Logout)

		v1.Use(middleware.AuthMiddleware())
		v1.GET("/auth/me", authHandler.Me)
		v1.GET("/users", userHandler.GetAllUsers)
		v1.POST("/users", userHandler.PostUserHandler)
		v1.GET("/users/:id", userHandler.GetUserByID)
		v1.PUT("/users/:id", userHandler.UpdateUserHandler)
		v1.DELETE("/users/:id", userHandler.DeleteUserHandler)

		v1.POST("/accounts", accountHandler.CreateAccountHandler)
		v1.GET("/accounts/:id", accountHandler.GetAccountByIDHandler)
		v1.PUT("/accounts/:id", accountHandler.UpdateAccountHandler)
		v1.DELETE("/accounts/:id", accountHandler.DeleteAccountHandler)
		v1.GET("/accounts/user/:id", accountHandler.GetAccountByUserIDHandler)

		v1.POST("/groups", groupHandler.CreateGroupHandler)
		v1.GET("/groups/:id", groupHandler.GetGroupByIDHandler)
		v1.GET("/groups/user/:id", groupHandler.GetGroupByUserID)

		v1.POST("/invitations", invitationHandler.SendInvitationGroupHandler)
		v1.POST("/invitations/accept", invitationHandler.AcceptInvitationGroupHandler)

		v1.POST("/category", categoryHandler.CreateCategoryHandler)

		v1.POST("/transactions", transactionHandler.CreateTransaction)
		v1.GET("/transactions", transactionHandler.GetAllTransactions)
		v1.GET("/transactions/:id", transactionHandler.GetTransactionByID)
	}
}
