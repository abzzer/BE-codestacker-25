package routes

import (
	"github.com/abzzer/BE-codestacker-25/internal/handlers"
	"github.com/abzzer/BE-codestacker-25/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("We have a working API with databases and RBAC!!")
	})

	app.Post("/login", handlers.LoginHandler)
	app.Post("/logout", handlers.LogoutHandler)
	app.Post("/submit-report", handlers.SubmitCrimeReportHandler)
	app.Post("/add-case", middleware.JWTMiddleware("admin", "investigator"), handlers.CreateCaseHandler)

	app.Patch("/update-case-status/:caseid", middleware.JWTMiddleware("officer"), handlers.UpdateCaseStatusHandler)

	adminRoutes := app.Group("/admin", middleware.JWTMiddleware("admin"))
	adminRoutes.Post("/create-user", handlers.CreateUserHandler)
	adminRoutes.Patch("/update-user/:id", handlers.UpdateUserHandler)
	adminRoutes.Delete("/delete-user/:id", handlers.DeleteUserHandler)

	evidence := app.Group("/add-evidence", middleware.JWTMiddleware("admin", "investigator", "officer"))
	evidence.Post("/text", handlers.AddTextEvidenceHandler)
	evidence.Post("/image", handlers.AddImageEvidenceHandler)

	caseRoutes := app.Group("/update-case", middleware.JWTMiddleware("admin", "investigator"))
	caseRoutes.Post("/:caseid/add-person", handlers.AddPersonHandler)
	caseRoutes.Put("/:caseid", handlers.UpdateCaseHandler)

	viewCase := app.Group("/case", middleware.JWTMiddleware("admin", "investigator", "officer"))
	viewCase.Get("/:caseid", handlers.GetCaseDetailsHandler)

	app.Get("/evidence/:evidenceid", middleware.JWTMiddleware("admin", "investigator", "officer"), handlers.GetEvidenceHandler)

}
