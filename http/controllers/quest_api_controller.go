package controllers

import (
	"github.com/Akkadius/spire/http/routes"
	"github.com/Akkadius/spire/questapi"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

type QuestApiController struct {
	logger *logrus.Logger
	parser *questapi.ParseService
	sourcer *questapi.QuestExamplesGithubSourcer
}

func NewQuestApiController(
	logger *logrus.Logger,
	parser *questapi.ParseService,
	sourcer *questapi.QuestExamplesGithubSourcer,
) *QuestApiController {
	return &QuestApiController{logger: logger, parser: parser, sourcer: sourcer}
}

func (d *QuestApiController) Routes() []*routes.Route {
	return []*routes.Route{
		routes.RegisterRoute(http.MethodGet, "quest-api/methods", d.methods, nil),
		routes.RegisterRoute(
			http.MethodPost,
			"quest-api/source-examples/org/:org/repo/:repo/branch/:branch",
			d.searchGithubExamples,
			nil,
		),
	}
}

func (d *QuestApiController) methods(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"data": d.parser.Parse(false)})
}

type SearchTermRequest struct {
	SearchTerms []string `json:"search_terms"`
	Language    string   `json:"language"`
}

// searches quest examples
func (d *QuestApiController) searchGithubExamples(c echo.Context) error {
	// body - bind
	p := new(SearchTermRequest)
	if err := c.Bind(p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// params
	org := c.Param("org")
	repo := c.Param("repo")
	branch := c.Param("branch")

	// result
	return c.JSON(
		http.StatusOK,
		echo.Map{
			"data": d.sourcer.Search(org, repo, branch, p.SearchTerms, p.Language),
		},
	)
}
