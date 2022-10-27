# Introduction

Quick Doc is an easy-to-use Go module which can be used to generate API documentations using OpenAPI tech stack [1]. Quick Doc follows OpenAPI 3.0.3 guidelines.

The purpose of developing this is to,

-   Use OpenAPI tech stack to create API documentation
    
-   Easy integration with existing APIs with minimal code change
    
-   Auto generate OpenAPI schema(s) from Go objects
    
-   Use with other frameworks and modules. Ex: gin, mux-router, http
    
-   Provide easy to use, OpenAPI web viewer
    
-   Add new endpoint quickly with fewer changes
    

Unlike lexical parsing methods, this method works at the runtime. There for there is a few of advantages and disadvantages can be addressed,

| **Pros** | **Cons** |
|--|--|
|Easy to use with existing auto-completion  | There is an overhead at API start time. Can be reduced by serializing JSON spec string. |
|No boilerplate comments  | No automated way to identify the endpoints which haven’t documented yet. |
|It can be easily integrated with existing APIs with a minimal code change. | 
|Very flexible. | 
|No need of additional script or pipeline| 


# Usage

## Quick Start Example

Install,

```
go get -u https://github.com/pickme-lk/quick-doc@v1.0.0
```

Import qdoc and ui packages from quick-doc module,

```
import (
	"https://github.com/pickme-lk/quick-doc/qdoc"
	"https://github.com/pickme-lk/quick-doc/ui"
)
```

Create a new OpenAPI document,

```
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
		LogoUrl:      "<logoURL>",
	},
})
```

Add an endpoint details,

```
doc.Post(&qdoc.Endpoint{
	Path: "/doc/user",
	Desc: "Create a new user",
	ReqBody: qdoc.ReqJson(doc.Schema(ReqUserAdd{
		Name: "Student 1",
		Age:  16,
		Project: &Project{
			Name:        "Volunteer Project",
			Description: "This is a volunteer project",
		},
	})),
	RespSet: qdoc.RespSet{
		Success: qdoc.ResJson("User creation success", nil),
	},
}).Tag("User").WithBearerAuth()
```

Compile the quick-doc object,

```
cd, err := doc.Compile()
if err != nil {
	panic(err)
}
```

Get http.ServeMux handler,

```
fmt.Println("Server is running on port 8080")
fmt.Println("Swagger UI: http://localhost:8080/doc/ui")
err = http.ListenAndServe(":8080", s)
if err != nil {
	return
}
```

Complete Example (main.go),

```
package main

import (
	"fmt"
	"https://github.com/pickme-lk/quick-doc/qdoc"
	"https://github.com/pickme-lk/quick-doc/ui"
	"net/http"
)

type Project struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Description2 string `json:"description2"`
}

type ReqUserAdd struct {
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Project *Project `json:"project"`
}

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

	doc.Post(&qdoc.Endpoint{
		Path: "/doc/user",
		Desc: "Create a new user",
		ReqBody: qdoc.ReqJson(doc.Schema(ReqUserAdd{
			Name: "Student 1",
			Age:  16,
			Project: &Project{
				Name:        "Volunteer Project",
				Description: "This is a volunteer project",
			},
		})),
		RespSet: qdoc.RespSet{
			Success: qdoc.ResJson("User creation success", nil),
		},
	}).Tag("User").WithBearerAuth()

	cd, err := doc.Compile() // generate open api json target
	if err != nil {
		panic(err)
	}

	s := cd.ServeMux()

	fmt.Println("Server is running on port 8080")
	fmt.Println("Swagger UI: http://localhost:8080/doc/ui")
	err = http.ListenAndServe(":8080", s)
	if err != nil {
		return
	}
}
```

## Creating API document Step-by-Step guide

### 1) Creating document configuration object

qdoc package provides a function called `qdoc.NewDoc(...)` to create a new document configuration `object. qdoc.NewDoc` function accepts a `qdoc.Config{}` object.

#### Configuration Options 
### `qdoc.Config`

| **Field** |  **Type**  | **Description**|
| -- |  --  | --|
|Title|`string`|Open API documentation title. This will be used as both OpenAPI spec title and UI title.|
|Description | `string` |(**Optional**) Open API documentation description. This support markdown.|
Version|`string`|(**Optional**) Version information for OpenAPI specification. Example: `1.0.0` |
Servers|`[]string`|List of API host servers. There is a helper function to increase readability and constancy. <br/> <br/>Example:<br/><pre>qdoc.Servers(<br/>"http://localhost:8080",<br/>"http://dev.quickdoc.com",<br/>),</pre>|
AuthConf|`qdoc.AuthConf`|(**Optional**) Define authentication methods for API. There is a helper function to define this field. This field can be ignored, then automatically decide according to endpoint authentication details. <br/>Example: `qdoc.NewAuthConf().WithBearer()`|
SpecPath|`string`|(**Optional**) URL path to serve OpenAPI JSON. Default value is set to `/doc/openapi.json`
UiConfig|`qdoc.UiConfig`|(**Optional**) See below for more details


### `qdoc.UiConfig`

|**Field**|**Type**|**Description**|
|--|--|--|
|Enabled|`boolean`|Enable or disable built-in OpenAPI web viewer. Default value is `false`.|
Path|`string`|(**Optional**) URL path to serve OpenAPI web viewer. Default value will be set to `/doc/ui`. <br/>*Make sure that* `SpecPath `*and* `UiConfig.Path` *has same prefix string.*
DefaultTheme|`ui.Theme`|(**Optional**) `ui.SWAGGER_UI` or `ui.RAPI_DOC`, default value is `ui.SWAGGER_UI`
ThemeByQuery|`boolean`|(**Optional**) When this is set to true. Web viewer accepts optional query parameter called `theme=swagger-ui` or `theme=rapi-doc`
LogoUrl|`string`|CDN image URL to show in Swagger UI.
<br/>

Example of creating document configuration object,

```
doc := qdoc.NewDoc(qdoc.Config{
	Title:       "Quick Doc Example",
	Description: "Quick Doc Example API documentation example",
	Version:     "1.0.0",
	Servers: qdoc.Servers(
		"http://localhost:8080",
		"http://dev.quickdoc.com",
		"http://quickdoc.com",
	),
	SpecPath: "/doc/json",
    AuthConf: qdoc.NewAuthConf().WithBearer(),
	UiConfig: qdoc.UiConfig{
		Enabled:      true,
		Path:         "/doc/ui",
		DefaultTheme: ui.SWAGGER_UI,
		ThemeByQuery: true,
	},
})
```

### 2) Add endpoint to configuration

> Examples can be found at the end of this step.

`Doc` object provide `Get`, `Post`,` Put`, `Delete` methods which can be used to add endpoints to configuration object. Each of method accept a pointer to a `qdoc.Endpoint` object.

### `qdoc.Endpoint`

|**Field**|**Type**|**Description**|
--|--|--|
Path|`string`|URL path of the endpoint <br/>Example: `/api/user`
Summary|`string`|(**Optional**) Brief summary about endpoint <br/> Example: `get current user details`
Description|`string`|(**Optional**) Descriptive details about endpoint. This field has Markdown support
ReqBody|`qdoc.RequestBody`|(**Optional**) Request body schema and other details. Quick doc provides a couple of helper functions for creating these.<br/>`qdoc.ReqJson` - create JSON request.<br/>`qdoc.ReqJson` - create JSON request.<br/>`qdoc.ReqJson` - create JSON request.<br/>`qdoc.ReqForm` - create URL encoded form data request.<br/>qdoc.ReqBody - create custom request body with custom content types.<br/>All of these functions accept a pointer to a qdoc.SchemaConfig which provide details to generate OpenAPI schema. For more details about qdoc.SchemaConfig can be found below.<br/>Examples can be found below.
QueryParams|`qdoc.Parameters`|(**Optional**) Define query parameters in the request. Quick doc provides a couple of helper functions for creating these.<br/>`qdoc.QueryParams` - create `qdoc.Parameters`<br/>`qdoc.QueryParams` - create qdoc.Parameters<br/>`qdoc.OptionalParam` - create optional parameter<br/>`qdoc.RequiredParam` - create required parameter<br/>Both of these functions accepts two arguments,<br/>`name: string` - parameter name<br/>`sc: *qdoc.SchemaConfig - pointer to schema config (optional)<br/>Examples can be found below.
PathParams|`qdoc.Parameters`|(**Optional**) Define path parameters in the request. Same helper functions specified in QueryParams applied here.<br/>`qdoc.PathParams` - create qdoc.Parameters
Headers|`qdoc.Parameters`|(**Optional**) Define header parameters in the request. Same helper functions specified in QueryParams applied here.<br/>`qdoc.Headers` - create qdoc.Parameters
RespSet|`qdoc.RespSet`|Define set of response for the endpoint. Quick doc provide helper functions,<br/><pre>type RespSet struct {<br/>	Success   *Response<br/>	BadReq    *Response<br/>	UnAuth    *Response<br/>	Forbidden *Response<br/>	NotFound  *Response<br/>	ISE       *Response<br/>	others    map[HttpStatus]*Response<br/>}</pre><br/>`qdoc.ResJson` - define a JSON response.<br/>Examples can be found below.


### `qdoc.SchemaConfg`

Quick Doc provides a helper method inside document configuration object to create `SchemaConfig` objects. <br/> Example can be found below,
```
doc := qdoc.NewDoc(...)
doc.Schema(object: interface{}) // returns SchemaConfig pointer
```

### `qdoc.Parameter`

Quick Doc provides two helper function,

```
qdoc.OptionalParam(
    name: string // parameter name
    sc: *qdoc.SchemaConfig // schema config to define paramter schema
)

qdoc.RequiredParam(
    name: string // parameter name
    sc: *qdoc.SchemaConfig // schema config to define paramter schema
)
```

### `qdoc.RequestBody`

Quick Doc provides 3 helper function,
```
qdoc.ReqJson(
    sc: *qdoc.SchemaConfig // schema config to define req schema
)

qdoc.ReqForm(
    sc: *qdoc.SchemaConfig // schema config to define req schema
)

qdoc.ReqBody(
    sc: *qdoc.SchemaConfig // schema config to define req schema
)(
  []string // content types
)
```

### `qdoc.RespSet`
```
type RespSet struct {
	Success   *Response
	BadReq    *Response
	UnAuth    *Response
	Forbidden *Response
	NotFound  *Response
	ISE       *Response
	others    map[HttpStatus]*Response
}
```

Quick Doc provides helper functions to create `Response` objects

```
qdoc.ResJson(
    desc: string        // breif description about response
    sc: *SchemaConfig   // schema config to define response schema
)
```

#### Examples

Types and structs used in below examples,
```
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
```

1.  Post request example
    
    ```
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
    ```
    
2.  Get endpoint example with Path parameter
    
    ```
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
    ```
    
3.  Get request example with query parameter
    
    ```
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
    ```
    
4.  Post request example with required request header
    
    ```
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
    ```
    
5.  Get request with complex schema
    
    ```
    doc.Get(&qdoc.Endpoint{
    	Path: "/api/team/{teamId}",
    	Desc: "Get team with users",
    	PathParams: qdoc.PathParams(
    		qdoc.RequiredParam("team id", doc.Schema(0)), // 0 is int type and example value will be 0
    	),
    	RespSet: qdoc.RespSet{
    		Success: qdoc.ResJson("Team found", doc.Schema(struct {
    			Team  Team   `json:"team"`
    			Users []User `json:"users"`
    		}{
    			Team: Team{
    				Name:        "testteam1",
    				Description: "test team 1",
    			},
    			Users: []User{
    				{
    					Username: "testuser1",
    					Password: "123456",
    					Age:      24,
    					Team:     "testteam1",
    				},
    			},
    		})), // schema and example will be generated
    		NotFound: qdoc.ResJson("Team not found", nil),
    		ISE:      qdoc.ResJson("Internal server error", nil),
    	},
    }).Tag("Team").WithBearerAuth() // Add bearer token authentication requirement
    ```
    

### 3) Compiling and Serving OpenAPI document

**Compiling**
```
cd, err := doc.Compile() // returns *CompiledDoc object
if err != nil {
	panic(err) // when compilation is failed
}
```

**Serving**
`CompiledDoc` object has `cd.ServeMux` method which returns a `*http.ServeMux` http request multiplexer. Which can be used to serve both web UI and JSON spec string.

```
// SpecPath = /doc/json
// UIPath   = /doc/ui
s := cd.ServeMux()

// Using builtin http router
err = http.ListenAndServe(":8080", s)

// Using Gorilla mux router
r := mux.NewRouter()
r.PathPrefix("/doc/").Handler(OpenApiRouter())
fmt.Println("listening on :8080")
http.ListenAndServe(":8080", r)

// Using Gin framework
r := gin.Default()
r.Group("/doc/*w").GET("", gin.WrapH(s)) 
r.Run()

```

## Integrate with an artifact

1.  Create document configuration options including enable/disable switch,
    
    ```
    // Example
    type DocConfig struct {
        Enabled   bool
        Title     string
        Desc      string
        Version   string
        Servers   []string
        SpecPath  string
        
        UiEnabled bool
        UiPath    string
    }
    
    ```
    
2.  Create new file `doc.go` inside `internal/transport/http` package. (or create a new package)  
    
3.  Create `InitDocs` function
 
    ```
    
    func InitDocs(router) {
        // check configuration switch
        if !config.DocConf.Enabled {
            return
        }
    
        doc := qdoc.NewDoc(qdoc.Config{
    		Title:       config.DocConf.Title,
    		Description: config.DocConf.Desc,
    		Version:     config.DocConf.Version,
    		Servers: qdoc.Servers(config.DocConf.Servers...),
    		SpecPath: config.DocConf.SpecPath, // doc/json
      
    		UiConfig: qdoc.UiConfig{
    			Enabled:      config.DocConf.UiEnabled,
    			Path:         config.DocConf.UiPath, // doc/ui
    			DefaultTheme: ui.SWAGGER_UI,
    			ThemeByQuery: true,
    		},
        })
        
        // add enpoints
        
        doc.Get(...)
        
        doc.Post(...)
        
        // compile doc
        cd, err := doc.Compile() // returns *CompiledDoc object
        if err != nil {
    		log.Error("failed to compile Open API doc")
        }
        
        // serving
        s := cd.ServeMux()
        router.handle("/doc/*", s)
    }
    
    
    ```
    
4.  Call `InitDocs` while router initialization
    
    ```
     func InitRouter() {
         ...
         router := ...
         ...
         InitDocs(router)
         ...
     }
    ```
    

# Architecture

Quick Doc module consists of 3 packages,

1.  `schema` package
    
2.  `ui` package
    
3.  `qdoc` package
    

### Schema Package

Go `reflect` is a built-in package which expose runtime reflections. This reflection details can be used to extract information from Golang types such as field types, tags. [https://pkg.go.dev/reflect](https://pkg.go.dev/reflect)

This schema package gets use of Go `reflect` package to generate schema from Golang objects. It builds a `Property` object according to given Go object.

**Features**

-   Configurable object exploration algorithm
    
-   Respect json tag to generate key name
    
-   Support pointer variables
    
-   Extracting information from nil values
    
-   Support Integer, Float, Boolean, String, Array, Slice, Map, Struct types.
    
-   Support additional tags
    
-   [ToDo] Support constraint validations
    
-   [ToDo] Support recursive types
    

**Example**

```
type UserAccount struct {
	Name string   	`json:"name"`
	Age  int      	`json:"age,omitempty"`
	Logs []LogEntry `json:"logs"`
}

type LogEntry struct {
	Date  string `json:"date"`
	Msg   bool   `json:"msg"`
}


sb := schema.NewBuilderDefault()
prop, err := sb.GetSchema(UserAccount{
	Name: "Test User",
	Age:  22,
	Logs: []LogEntry{
		{
			Date:  "2022-01-21",
			Message: "account created",
        }
	},
})
```

### UI Package

There are plenty of different OpenAPI viewers can be found. The most famous OpenAPI viewer is Swagger UI. It supports both OpenAPI 2, OpenAPI 3 specifications. [Swagger UI](https://petstore.swagger.io/)

This ui package provides a couple of configurable and easy to use Open API web viewers.

**Features**
-   Configurable
-   [Swagger UI](https://petstore.swagger.io/)
-   [Rapi Doc UI](https://rapidocweb.com/)
    

### Quick doc Package

This is the primary package provides by quick doc module. Its features can be split into 3 primary functions.

-   Building API documentation config    
-   Compiling documentation config into OpenAPI specs 
-   Serve OpenAPI JSON and web UI over HTTP
    

**Features**

-   Easy and simple configuration    
-   Builtin OpenAPI web viewer    
-   Generate OpenAPI schema from Go objects  
-   Supported Auth methods: Basic, Bearer Token   
-   Multiple response support with schema   
-   JSON, Form and Multipart form request body support  
-   Query, Path parameter support  
-   Header parameter support    
-   Endpoint tag support    
-   and more…
    

# References

1.  [https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md](https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md)
    
2.  [https://rapidocweb.com/](https://rapidocweb.com/)
    
3.  [https://editor.swagger.io/](https://editor.swagger.io/)
    
4.  [https://petstore.swagger.io/](https://petstore.swagger.io/)
    
5.  [https://pkg.go.dev/reflect](https://pkg.go.dev/reflect)
