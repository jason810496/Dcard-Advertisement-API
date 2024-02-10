// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/ad": {
            "get": {
                "description": "query ad by query parameters age, country",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "ad"
                ],
                "summary": "public API",
                "parameters": [
                    {
                        "maximum": 100,
                        "minimum": 1,
                        "type": "integer",
                        "description": "age",
                        "name": "age",
                        "in": "query"
                    },
                    {
                        "enum": [
                            "TW",
                            "HK",
                            "JP",
                            "US",
                            "KR"
                        ],
                        "type": "string",
                        "description": "country",
                        "name": "country",
                        "in": "query"
                    },
                    {
                        "enum": [
                            "ios",
                            "android",
                            "web"
                        ],
                        "type": "string",
                        "description": "platform",
                        "name": "platform",
                        "in": "query"
                    },
                    {
                        "enum": [
                            "F",
                            "M"
                        ],
                        "type": "string",
                        "description": "gender",
                        "name": "gender",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "limit",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "offset",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/github_com_jason810496_Dcard-Advertisement-API_pkg_schemas.PublicAdResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.HTTPError"
                        }
                    }
                }
            },
            "post": {
                "description": "create advertisement with ` + "`" + `startAt` + "`" + `, ` + "`" + `endAt` + "`" + ` and ` + "`" + `condition` + "`" + `",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json",
                    "application/json"
                ],
                "tags": [
                    "ad"
                ],
                "summary": "admin API",
                "parameters": [
                    {
                        "description": "advertisement request schema",
                        "name": "advertisement",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_jason810496_Dcard-Advertisement-API_pkg_schemas.CreateAdRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/github_com_jason810496_Dcard-Advertisement-API_pkg_schemas.CreateAdResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.HTTPError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "github_com_jason810496_Dcard-Advertisement-API_pkg_schemas.CreateAdRequest": {
            "type": "object",
            "required": [
                "conditions",
                "endAt",
                "startAt",
                "title"
            ],
            "properties": {
                "conditions": {
                    "type": "array",
                    "items": {
                        "type": "object",
                        "required": [
                            "ageEnd",
                            "ageStart",
                            "country",
                            "platform"
                        ],
                        "properties": {
                            "ageEnd": {
                                "type": "integer",
                                "example": 30
                            },
                            "ageStart": {
                                "type": "integer",
                                "example": 18
                            },
                            "country": {
                                "type": "array",
                                "items": {
                                    "type": "string"
                                },
                                "example": [
                                    "TW",
                                    "JP",
                                    "KR",
                                    "US"
                                ]
                            },
                            "platform": {
                                "type": "array",
                                "items": {
                                    "type": "string"
                                },
                                "example": [
                                    "ios",
                                    "web"
                                ]
                            }
                        }
                    }
                },
                "endAt": {
                    "type": "string",
                    "example": "2024-02-01 10:00:00"
                },
                "startAt": {
                    "type": "string",
                    "example": "2024-01-01 10:00:00"
                },
                "title": {
                    "type": "string",
                    "example": "AD 123"
                }
            }
        },
        "github_com_jason810496_Dcard-Advertisement-API_pkg_schemas.CreateAdResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "create success"
                }
            }
        },
        "github_com_jason810496_Dcard-Advertisement-API_pkg_schemas.PublicAdItem": {
            "type": "object",
            "properties": {
                "endAt": {
                    "type": "string",
                    "example": "2021-12-31 23:59:59"
                },
                "title": {
                    "type": "string",
                    "example": "This is an AD title"
                }
            }
        },
        "github_com_jason810496_Dcard-Advertisement-API_pkg_schemas.PublicAdResponse": {
            "type": "object",
            "properties": {
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/github_com_jason810496_Dcard-Advertisement-API_pkg_schemas.PublicAdItem"
                    }
                }
            }
        },
        "utils.HTTPError": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 400
                },
                "message": {
                    "type": "string",
                    "example": "status bad request"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
