package assets

import (
	"fmt"
	"github.com/Akkadius/spire/internal/github"
	appmiddleware "github.com/Akkadius/spire/internal/http/middleware"
	"github.com/labstack/echo/v4"
	gocache "github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type SpireAssets struct {
	logger     *logrus.Logger
	cache      *gocache.Cache
	downloader *github.GithubSourceDownloader
}

const (
	organization = "Akkadius"
	repository   = "eq-asset-preview"
	version      = 1
)

func NewSpireAssets(
	logger *logrus.Logger,
	cache *gocache.Cache,
	downloader *github.GithubSourceDownloader,
) *SpireAssets {
	return &SpireAssets{
		logger:     logger,
		cache:      cache,
		downloader: downloader,
	}
}

func (a SpireAssets) ServeStatic() echo.MiddlewareFunc {
	// cleanup old versions
	for i := 1; i < version; i++ {
		oldPath := filepath.Join(a.downloader.GetSourceRoot(), fmt.Sprintf("%v-v%v", repository, i))
		fmt.Printf("Deleting old assets [%v]\n", oldPath)
		err := os.RemoveAll(oldPath)
		if err != nil {
			a.logger.Fatal(err)
		}
	}

	zippedPath := a.getZippedPath()
	if !a.AssetsExists() && len(zippedPath) == 0 {
		downloader := github.NewGithubSourceDownloader(a.logger, a.cache)
		branch := fmt.Sprintf("v%v", version)
		downloader.SourceToUserCacheDir(true)
		downloader.SetReadFiles(false)
		r := downloader.Source(organization, repository, branch, false)
		if len(r.ZippedPath) > 0 {
			zippedPath = r.ZippedPath
		}
	}

	// serve
	return appmiddleware.StaticAsset(appmiddleware.StaticConfig{
		Root:        "/",
		StripPrefix: string(filepath.Separator) + "eq-asset-preview-master",
		Filesystem:  http.Dir(zippedPath),
	})
}

func (a SpireAssets) AssetsExists() bool {
	if _, err := os.Stat(a.getRepoDir()); err == nil {
		return true
	}

	return false
}

func (a SpireAssets) getZippedPath() string {
	repoDir := a.getRepoDir()
	zippedPath := ""
	files, _ := ioutil.ReadDir(repoDir)
	if len(files) > 0 {
		zippedPath = filepath.Join(repoDir, files[0].Name())
	}

	return zippedPath
}

func (a SpireAssets) getRepoDir() string {
	return filepath.Join(a.downloader.GetSourceRoot(), fmt.Sprintf("%v-v%v", repository, version))
}
