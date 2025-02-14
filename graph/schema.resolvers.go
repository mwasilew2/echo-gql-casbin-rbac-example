package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.55

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	echo "github.com/labstack/echo/v4"

	"github.com/mwasilew2/echo-gqlgen-casbin-rbac-example/graph/model"
	"github.com/mwasilew2/echo-gqlgen-casbin-rbac-example/util"
)

// CreateAccount is the resolver for the createAccount field.
func (r *mutationResolver) CreateAccount(ctx context.Context, input model.NewAccount) (*model.Account, error) {
	panic(fmt.Errorf("not implemented: CreateAccount - createAccount"))
}

// CreateNamespace is the resolver for the createNamespace field.
func (r *mutationResolver) CreateNamespace(ctx context.Context, input model.NewNamespace) (*model.Namespace, error) {
	panic(fmt.Errorf("not implemented: CreateNamespace - createNamespace"))
}

const (
	AuthorizationResourceStack = "stack"

	AuthorizationActionCreate = "create"
	AuthorizationActionRead   = "read"
)

// CreateStack is the resolver for the createStack field.
func (r *mutationResolver) CreateStack(ctx context.Context, input model.NewStack) (*model.Stack, error) {
	// extract echo context
	ec, err := util.ExtractEchoContext(ctx)
	if err != nil {
		r.logger.Error("Error getting echo context", "error", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	// extract account from session
	sess, err := session.Get(util.CookieKeySessionName, ec)
	if err != nil {
		r.logger.Error("Error getting session", "error", err)
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Error getting session")
	}
	userID, userExists := sess.Values[util.SessionKeyUserID]
	if !userExists {
		r.logger.Debug("No user ID in session")
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Not logged in")
	}
	accountID, accountExists := sess.Values[util.SessionKeyAccountID]
	if !accountExists {
		r.logger.Debug("No account ID in session")
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Not logged in")
	}

	// get user from database
	user := &model.User{}
	err = r.db.Where("id = ?", userID).First(user).Error
	if err != nil {
		r.logger.Error("Error getting user", "error", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	// get account from database
	account := &model.Account{}
	err = r.db.Where("ulid = ?", accountID).First(account).Error
	if err != nil {
		r.logger.Error("Error getting account", "error", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	// authorization
	hasAccess, err := r.authorizationService.IsAuthorized(user.Username, account.Name, AuthorizationResourceStack, AuthorizationActionCreate)
	if err != nil {
		r.logger.Error("Error checking authorization", "error", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
	if !hasAccess {
		r.logger.Debug("Not authorized")
		return nil, echo.NewHTTPError(http.StatusForbidden, "Not authorized")
	}

	// create stack
	stack := &model.Stack{
		Name:    input.Name,
		Account: account,
	}
	err = r.db.Create(stack).Error
	if err != nil {
		r.logger.Error("Error creating stack", "error", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	return stack, nil
}

// Account is the resolver for the account field.
func (r *queryResolver) Account(ctx context.Context) (*model.Account, error) {
	// extract echo context
	ec, err := util.ExtractEchoContext(ctx)
	if err != nil {
		r.logger.Error("Error getting echo context", "error", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	// extract user and accouant from session
	sess, err := session.Get(util.CookieKeySessionName, ec)
	if err != nil {
		r.logger.Error("Error getting session", "error", err)
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Error getting session")
	}
	userID, userExists := sess.Values[util.SessionKeyUserID]
	if !userExists {
		r.logger.Debug("No user ID in session")
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Not logged in")
	}
	r.logger.Debug("Got user ID from session", "userID", userID)
	accountID, accountExists := sess.Values[util.SessionKeyAccountID]
	if !accountExists {
		r.logger.Debug("No account ID in session")
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Not logged in")
	}

	return &model.Account{
		Ulid: accountID.(string),
	}, nil
}

// Namespaces is the resolver for the namespaces field.
func (r *queryResolver) Namespaces(ctx context.Context) ([]*model.Namespace, error) {
	panic(fmt.Errorf("not implemented: Namespaces - namespaces"))
}

// Stacks is the resolver for the stacks field.
func (r *queryResolver) Stacks(ctx context.Context) ([]*model.Stack, error) {
	// extract echo context
	ec, err := util.ExtractEchoContext(ctx)
	if err != nil {
		r.logger.Error("Error getting echo context", "error", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	// extract account from session
	sess, err := session.Get(util.CookieKeySessionName, ec)
	if err != nil {
		r.logger.Error("Error getting session", "error", err)
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Error getting session")
	}
	userID, userExists := sess.Values[util.SessionKeyUserID]
	if !userExists {
		r.logger.Debug("No user ID in session")
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Not logged in")
	}
	accountID, accountExists := sess.Values[util.SessionKeyAccountID]
	if !accountExists {
		r.logger.Debug("No account ID in session")
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Not logged in")
	}

	// get user from database
	user := &model.User{}
	err = r.db.Where("id = ?", userID).First(user).Error
	if err != nil {
		r.logger.Error("Error getting user", "error", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	// get account from database
	account := &model.Account{}
	err = r.db.Where("ulid = ?", accountID).First(account).Error
	if err != nil {
		r.logger.Error("Error getting account", "error", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	// authorization
	hasAccess, err := r.authorizationService.IsAuthorized(user.Username, account.Name, AuthorizationResourceStack, AuthorizationActionRead)
	if err != nil {
		r.logger.Error("Error checking authorization", "error", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
	if !hasAccess {
		r.logger.Debug("Not authorized", "username", user.Username, "account", account.Name, "resource", AuthorizationResourceStack, "action", AuthorizationActionRead)
		return nil, echo.NewHTTPError(http.StatusForbidden, "Not authorized")
	}
	r.logger.Debug("Authorized to get stacks")

	// get stack
	stacks := []*model.Stack{}
	err = r.db.Where("account_id = ?", account.ID).Find(&stacks).Error
	if err != nil {
		r.logger.Error("Error getting stack", "error", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	return stacks, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
