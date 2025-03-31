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
	app.Get("/check-report/:reportID", handlers.CheckReportStatus)

	adminRoutes := app.Group("/admin", middleware.JWTMiddleware("admin"))
	adminRoutes.Post("/create-user", handlers.CreateUserHandler)
	adminRoutes.Patch("/update-user/:id", handlers.UpdateUserHandler)
	adminRoutes.Delete("/delete-user/:id", handlers.DeleteUserHandler)
	adminRoutes.Get("/audit-log", handlers.GetAuditLogs)

	addEvidence := app.Group("/add-evidence", middleware.JWTMiddleware("admin", "investigator", "officer"))
	addEvidence.Post("/text", handlers.AddTextEvidenceHandler)
	addEvidence.Post("/image", handlers.AddImageEvidenceHandler)

	caseRoutes := app.Group("/update-case", middleware.JWTMiddleware("admin", "investigator", "officer"))
	caseRoutes.Put("/:caseid", middleware.JWTMiddleware("admin", "investigator"), handlers.UpdateCaseHandler)
	caseRoutes.Post("/:caseid/add-person", middleware.JWTMiddleware("admin", "investigator"), handlers.AddPersonHandler)
	caseRoutes.Post("/:caseid/add-officer", middleware.JWTMiddleware("admin", "investigator"), handlers.AddOfficerToCaseHandler)
	caseRoutes.Patch("/:caseid/status", handlers.UpdateCaseStatusHandler)

	viewCase := app.Group("/case", middleware.JWTMiddleware("admin", "investigator", "officer"))
	viewCase.Get("/partial/:caseid", handlers.GetPartialCaseDetailsHandler)
	viewCase.Get("/full/:caseid", handlers.GetFullCaseDetailsHandler)
	viewCase.Get("/pdf/:caseid", handlers.GenerateCasePDFHandler)

	viewEvidence := app.Group("/evidence", middleware.JWTMiddleware("admin", "investigator", "officer"))
	viewEvidence.Get("/top-ten", handlers.GetTopWordsInTextEvidence)
	viewEvidence.Get("/get-urls/:caseid", handlers.GetCaseURLs)
	viewEvidence.Get("/details/:evidenceid", handlers.GetEvidenceHandler)
	viewEvidence.Get("/get-image/:evidenceid", handlers.GetImageEvidenceHandler)
	viewEvidence.Get("/update/:evidenceid", middleware.JWTMiddleware("admin", "investigator"), handlers.UpdateEvidence)
	viewEvidence.Get("/soft-delete/:evidenceid", middleware.JWTMiddleware("admin", "investigator"), handlers.SoftDeleteEvidence)

	// Evidence delete confirmation stages
	viewEvidence.Post("/hard-delete/:evidenceid", middleware.JWTMiddleware("admin", "investigator"), handlers.HardDeleteEvidence)
	viewEvidence.Patch("/hard-delete/:evidenceid", middleware.JWTMiddleware("admin", "investigator"), handlers.HardDeleteEvidence)
	viewEvidence.Delete("/hard-delete/:evidenceid", middleware.JWTMiddleware("admin", "investigator"), handlers.HardDeleteEvidence)

	reports := app.Group("/reports", middleware.JWTMiddleware("admin", "investigator", "officer"))
	reports.Get("/all", handlers.GetAllReports)
	reports.Post("/case/:reportID", handlers.LinkReportToCase)

	// Long pooling
	viewEvidence.Get("/hard-delete-status/:evidenceid", middleware.JWTMiddleware("admin"), handlers.LongPollDeleteStatus)

}
