package delivery

import (
	"context"
	"go-clean-api/domain"
	"go-clean-api/entity"
	"log"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/labstack/echo/v4"
)

type GraphQLHandler struct {
	usecase domain.UserUsecase
	schema  graphql.Schema
}

func NewGraphQLHandler(usecase domain.UserUsecase) *GraphQLHandler {
	// Define the user type
	userType := graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"first_name": &graphql.Field{
				Type: graphql.String,
			},
			"last_name": &graphql.Field{
				Type: graphql.String,
			},
			"email": &graphql.Field{
				Type: graphql.String,
			},
			"province": &graphql.Field{
				Type: graphql.String,
			},
			"district": &graphql.Field{
				Type: graphql.String,
			},
			"subdistrict": &graphql.Field{
				Type: graphql.String,
			},
			"zip_code": &graphql.Field{
				Type: graphql.String,
			},
			"detail_address": &graphql.Field{
				Type: graphql.String,
			},
			"phone": &graphql.Field{
				Type: graphql.String,
			},
			"password": &graphql.Field{
				Type: graphql.String,
			},
			"created_at": &graphql.Field{
				Type: graphql.String,
			},
			"updated_at": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	// Define the root query
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"user": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"].(int)
					if !ok {
						return nil, nil
					}
					usecase := p.Context.Value("usecase").(domain.UserUsecase)
					return usecase.GetUserByID(uint(id))
				},
			},
			"users": &graphql.Field{
				Type: graphql.NewList(userType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					usecase := p.Context.Value("usecase").(domain.UserUsecase)
					return usecase.GetAllUsers()
				},
			},
		},
	})

	// Define the root mutation
	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createUser": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"first_name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"last_name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"email": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"province": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"district": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"subdistrict": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"zip_code": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"detail_address": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"phone": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"password": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					usecase := p.Context.Value("usecase").(domain.UserUsecase)
					
					user := &entity.User{
						First_name:     p.Args["first_name"].(string),
						Last_name:      p.Args["last_name"].(string),
						Email:          p.Args["email"].(string),
						Province:       p.Args["province"].(string),
						District:       p.Args["district"].(string),
						Subdistrict:    p.Args["subdistrict"].(string),
						Zip_code:       p.Args["zip_code"].(string),
						Detail_address: p.Args["detail_address"].(string),
						Phone:          p.Args["phone"].(string),
						Password:       p.Args["password"].(string),
					}
					
					err := usecase.CreateUser(user)
					if err != nil {
						return nil, err
					}
					
					return user, nil
				},
			},
			"updateUser": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"first_name": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"last_name": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"email": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"province": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"district": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"subdistrict": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"zip_code": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"detail_address": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"phone": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"password": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					usecase := p.Context.Value("usecase").(domain.UserUsecase)
					
					id, ok := p.Args["id"].(int)
					if !ok {
						return nil, nil
					}
					
					user, err := usecase.GetUserByID(uint(id))
					if err != nil {
						return nil, err
					}
					if user == nil {
						return nil, nil
					}
					
					// Update only provided fields
					if firstName, ok := p.Args["first_name"].(string); ok {
						user.First_name = firstName
					}
					if lastName, ok := p.Args["last_name"].(string); ok {
						user.Last_name = lastName
					}
					if email, ok := p.Args["email"].(string); ok {
						user.Email = email
					}
					if province, ok := p.Args["province"].(string); ok {
						user.Province = province
					}
					if district, ok := p.Args["district"].(string); ok {
						user.District = district
					}
					if subdistrict, ok := p.Args["subdistrict"].(string); ok {
						user.Subdistrict = subdistrict
					}
					if zipCode, ok := p.Args["zip_code"].(string); ok {
						user.Zip_code = zipCode
					}
					if detailAddress, ok := p.Args["detail_address"].(string); ok {
						user.Detail_address = detailAddress
					}
					if phone, ok := p.Args["phone"].(string); ok {
						user.Phone = phone
					}
					if password, ok := p.Args["password"].(string); ok {
						user.Password = password
					}
					
					err = usecase.UpdateUser(user)
					if err != nil {
						return nil, err
					}
					
					return user, nil
				},
			},
			"deleteUser": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					usecase := p.Context.Value("usecase").(domain.UserUsecase)
					
					id, ok := p.Args["id"].(int)
					if !ok {
						return false, nil
					}
					
					err := usecase.DeleteUser(uint(id))
					if err != nil {
						return false, err
					}
					
					return true, nil
				},
			},
		},
	})

	// Create the schema
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})
	
	if err != nil {
		log.Fatalf("Failed to create GraphQL schema: %v", err)
	}

	return &GraphQLHandler{
		usecase: usecase,
		schema:  schema,
	}
}

func (h *GraphQLHandler) GraphQLHandler(c echo.Context) error {
	// Create a new context with the usecase
	ctx := context.WithValue(c.Request().Context(), "usecase", h.usecase)
	
	// Create a new request with the context
	req := c.Request().WithContext(ctx)
	
	// Create a response writer
	w := c.Response()
	
	// Create a new handler for each request
	handler := handler.New(&handler.Config{
		Schema:   &h.schema,
		Pretty:   true,
		GraphiQL: false,
	})
	
	// Handle the request
	handler.ServeHTTP(w, req)
	return nil
}

func (h *GraphQLHandler) PlaygroundHandler(c echo.Context) error {
	// Create a new context with the usecase
	ctx := context.WithValue(c.Request().Context(), "usecase", h.usecase)
	
	// Create a new request with the context
	req := c.Request().WithContext(ctx)
	
	// Create a response writer
	w := c.Response()
	
	// Create a new handler for the playground
	handler := handler.New(&handler.Config{
		Schema:   &h.schema,
		Pretty:   true,
		GraphiQL: true,
	})
	
	// Handle the request
	handler.ServeHTTP(w, req)
	return nil
}