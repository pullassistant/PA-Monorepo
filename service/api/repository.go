package api

import (
	"encoding/json"
	"goji.io/pat"
	"net/http"
	"strconv"
)

type RepositoryHandler struct {
}

func (h RepositoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	installationIdParam := pat.Param(r, "installationId")
	if len(installationIdParam) == 0 {
		panic("No installationId")
	}

	installationId, err := strconv.ParseInt(installationIdParam, 10, 64)
	if err != nil {
		panic("Invalid installationId")
	}

	client := clientFromCtx(r.Context())

	repositories, _, err := client.Apps.ListUserRepos(r.Context(), int64(installationId), nil)
	if err != nil {
		panic(err)
	}

	var repositoriesResponse []Repository = nil
	for _, v := range repositories {
		repositoriesResponse = append(repositoriesResponse, Repository{
			ID:      *v.ID,
			HTMLURL: *v.HTMLURL,
			Name:    *v.Name,
		})
	}

	err = json.NewEncoder(w).Encode(Repositories{Repositories: repositoriesResponse})
	if err != nil {
		panic(err)
	}
}

type Repositories struct {
	Repositories []Repository `json:"repositories"`
}
type Repository struct {
	ID      int64  `json:"id"`
	HTMLURL string `json:"html_url"`
	Name    string `json:"name"`
}
