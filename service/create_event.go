package service

import (
	"github.com/angdev/chocolat/model"
	"github.com/angdev/chocolat/support/repo"
)

type CreateEventParams struct {
	CollectionName string
	Events         map[string][]repo.Doc
}

func CreateEvent(p *model.Project, params *CreateEventParams) (repo.Doc, error) {
	r := repo.NewRepository(p.RepoName())
	defer r.Close()

	docs := params.Events[params.CollectionName]

	result := insertEvents(r, params.CollectionName, docs...)
	if result[0]["success"] == true {
		return repo.Doc{"created": true}, nil
	} else {
		return repo.Doc{"created": false}, nil
	}
}

func CreateMultipleEvents(p *model.Project, params *CreateEventParams) (map[string][]repo.Doc, error) {
	r := repo.NewRepository(p.RepoName())
	defer r.Close()

	result := map[string][]repo.Doc{}
	for event, docs := range params.Events {
		result[event] = insertEvents(r, event, docs...)
	}

	return result, nil
}

func insertEvents(r *repo.Repository, event string, docs ...repo.Doc) []repo.Doc {
	result := []repo.Doc{}
	for _, doc := range docs {
		if err := r.Insert(event, &doc); err != nil {
			result = append(result, repo.Doc{"success": false})
		} else {
			result = append(result, repo.Doc{"success": true})
		}
	}

	return result
}
