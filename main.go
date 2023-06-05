package main

import (
	"fmt"
	"github.com/pickme-lk/quick-doc/qdoc"
	"github.com/pickme-lk/quick-doc/ui"
	"net/http"
)

func main() {
	Doc()
}

func Doc() {
	doc := qdoc.NewDoc(qdoc.Config{
		Title:       "Quick Doc Demo",
		Description: "Quick Doc demo API documentation example",
		Version:     "1.0.0",
		Servers: qdoc.Servers(
			"http://localhost:8080",
			"http://dev.quickdoc.com",
		),
		SpecPath: "/doc/json",

		UiConfig: qdoc.UiConfig{
			Enabled:      true,
			Path:         "/doc/ui",
			DefaultTheme: ui.SWAGGER_UI,
			ThemeByQuery: false,
		},
	})

	// Post request example
	doc.Post(&qdoc.Endpoint{
		Path: "/api/user",
		Desc: "Create a new user",
		ReqBody: qdoc.ReqJson(doc.Schema(User{
			Username: "testuser1",
			Password: "123456",
			Age:      24,
			Team:     "testteam1",
		})),
		RespSet: qdoc.RespSet{
			Success: qdoc.ResJson("User creation success", nil),
			BadReq:  qdoc.ResJson("Invalid user data", nil),
			ISE:     qdoc.ResJson("Internal server error", nil),
		},
	}).Tag("User")

	// Get request example with path parameter
	doc.Get(&qdoc.Endpoint{
		Path: "/api/user/{userId}",
		Desc: "Get user by user id",
		PathParams: qdoc.PathParams(
			qdoc.RequiredParam("user id", doc.Schema(0)), // 0 is int type and example value will be 0
		),
		RespSet: qdoc.RespSet{
			Success:  qdoc.ResJson("User found", doc.Schema(User{})), // schema and example will be generated
			NotFound: qdoc.ResJson("User not found", nil),
			ISE:      qdoc.ResJson("Internal server error", nil),
		},
	}).Tag("User").WithBearerAuth() // Add bearer token authentication requirement

	// Get request example with query parameter
	doc.Get(&qdoc.Endpoint{
		Path: "/api/user",
		Desc: "Get users",
		QueryParams: qdoc.QueryParams(
			qdoc.OptionalParam("username", doc.Schema("testuser1")), // example value will be "testuser1"
			qdoc.OptionalParam("age", doc.Schema(11)),               // example value will be 11
			qdoc.OptionalParam("team", doc.Schema("testteam1")),     // example value will be "testteam1"
		),
		RespSet: qdoc.RespSet{
			Success:  qdoc.ResJson("User found", doc.Schema(User{})), // schema and example will be generated
			NotFound: qdoc.ResJson("User not found", nil),
			ISE:      qdoc.ResJson("Internal server error", nil),
		},
	}).Tag("User").WithBearerAuth() // Add bearer token authentication requirement

	// Post request example with required request header
	doc.Post(&qdoc.Endpoint{
		Path: "/api/team",
		Desc: "Create a new team",
		ReqBody: qdoc.ReqJson(doc.Schema(Team{
			Name:        "testteam1",
			Description: "test team 1",
		})),
		Headers: qdoc.Headers(
			qdoc.RequiredParam("origin", doc.Schema("mobile-app")), // example value will be "mobile-app"
		),
		RespSet: qdoc.RespSet{
			Success: qdoc.ResJson("Team creation success", nil),
			BadReq:  qdoc.ResJson("Invalid team data", nil),
			ISE:     qdoc.ResJson("Internal server error", nil),
		},
	}).Tag("Team").WithBearerAuth() // Add bearer token authentication requirement

	// Get request with complex schema
	doc.Get(&qdoc.Endpoint{
		Path: "/api/team/{teamId}",
		Desc: "Get team with users",
		PathParams: qdoc.PathParams(
			qdoc.RequiredParam("team id", doc.Schema(0)), // 0 is int type and example value will be 0
		),
		RespSet: qdoc.RespSet{
			Success:  qdoc.ResJson("Team found", doc.Schema(test{})), // schema and example will be generated
			NotFound: qdoc.ResJson("Team not found", nil),
			ISE:      qdoc.ResJson("Internal server error", nil),
		},
	}).Tag("Team").WithBearerAuth() // Add bearer token authentication requirement

	doc.Get(&qdoc.Endpoint{
		Summary: "Get a Option",
		Desc:    "Get a Option Endpoint",
		Path:    "/v1.0/sku/option/{option}",
		PathParams: qdoc.PathParams(
			qdoc.RequiredParam("option", doc.Schema(nil)),
		),
		Headers: qdoc.Headers(
			qdoc.OptionalParam("type", doc.Schema(nil)),
		),
		RespSet: qdoc.RespSet{
			Success: qdoc.ResJson("Success", doc.Schema(OptionGetResponse{
				Payload: struct {
					Id         int64        `json:"id"`
					MerchantId int64        `json:"merchantId"`
					Data       OptionDetail `json:"data"`
				}{
					MerchantId: 1,
					Id:         1,
					Data: OptionDetail{
						Id:         1,
						Name:       "name",
						Label:      "label",
						Note:       "note",
						MerchantId: 1,
						Required:   true,
						Selection:  1,
						Minimum:    1,
						Maximum:    1,
						Quantity:   1,
						Skus: struct {
							Added     []Sku `json:"added"`
							Available []Sku `json:"available"`
						}{Added: []Sku{{
							Id:              1,
							Name:            "name",
							CurrencyCode:    "currencyCode",
							Price:           2.36,
							IsOriginalPrice: true,
							MerchantStatus:  5,
							Order:           5,
							Options:         []OptionInSku{},
						}}, Available: []Sku{}},
					},
				},
			})),
			ISE: qdoc.ResJson("Internal Server Error", nil),
		},
	}).WithBearerAuth().Tag("SKUs")
	// Compile the doc config
	cd, err := doc.Compile()
	if err != nil {
		panic(err)
	}

	s := cd.ServeMux()

	fmt.Println("Swagger UI: http://localhost:8080/doc/ui")
	err = http.ListenAndServe(":8080", s)
	if err != nil {
		return
	}
}

type test struct {
	//name string `json:"name"`
	names []test2 `json:"names"`
}

type test2 struct {
	name2 string `json:"name_2"`
}

type OptionGetResponse struct {
	Payload struct {
		Id         int64        `json:"id"`
		MerchantId int64        `json:"merchantId"`
		Data       OptionDetail `json:"data"`
	} `json:"payload"`
}

type OptionDetail struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Label      string `json:"label"`
	Note       string `json:"note"`
	MerchantId int64  `json:"merchantId"`
	Required   bool   `json:"required"`
	Selection  int8   `json:"selection"`
	Minimum    int64  `json:"minimum"`
	Maximum    int64  `json:"maximum"`
	Quantity   int64  `json:"quantity"`
	Skus       struct {
		Added     []Sku `json:"added"`
		Available []Sku `json:"available"`
	} `json:"skus"`
}
type Sku struct {
	Id              int64         `json:"id"`
	Name            string        `json:"name"`
	CurrencyCode    string        `json:"currencyCode"`
	Price           float64       `json:"price"`
	IsOriginalPrice bool          `json:"is_original_price"`
	MerchantStatus  int           `json:"merchantStatus"`
	Order           int           `json:"order"`
	Options         []OptionInSku `json:"options"`
}

type OptionInSku struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Order int    `json:"order"`
	Skus  []Sku  `json:"skus"`
}

type Team struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Age      int    `json:"age"`
	Team     string `json:"team"`
}
