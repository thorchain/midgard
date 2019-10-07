package resolvers

//go:generate go run github.com/99designs/gqlgen
import (
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"math/rand"

	"gitlab.com/thorchain/bepswap/chain-service/api/graphQL/v1/codegen"
	"gitlab.com/thorchain/bepswap/chain-service/api/graphQL/v1/models"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct {
	todos []*models.Todo
}

func (r *Resolver) Mutation() codegen.MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() codegen.QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateTodo(ctx context.Context, input models.NewTodo) (*models.Todo, error) {
	todo := &models.Todo{
		Text: input.Text,
		ID:   fmt.Sprintf("T%d", rand.Int()),
		User: &models.User{
			ID: input.UserID,
		},
	}
	spew.Dump(todo)
	r.todos = append(r.todos, todo)
	return todo, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Todos(ctx context.Context) ([]*models.Todo, error) {
	return r.todos, nil
}
