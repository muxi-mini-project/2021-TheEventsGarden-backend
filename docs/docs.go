// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terrms/",
        "contact": {
            "name": "TAODEI",
            "email": "864978550@qq.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/": {
            "post": {
                "description": "学号密码登录",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "登录",
                "parameters": [
                    {
                        "description": "登录的用户信息",
                        "name": "object",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "将student_id作为token保留",
                        "schema": {
                            "$ref": "#/definitions/handler.Token"
                        }
                    },
                    "400": {
                        "description": "输入有误，格式错误"
                    },
                    "401": {
                        "description": "用户名或密码错误"
                    }
                }
            }
        },
        "/crawler": {
            "post": {
                "description": "爬取用户云课堂作业",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "爬取作业",
                "parameters": [
                    {
                        "type": "string",
                        "description": "token",
                        "name": "token",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "password",
                        "name": "password",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "获取成功",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.homework"
                            }
                        }
                    },
                    "400": {
                        "description": "输入有误，格式错误"
                    },
                    "401": {
                        "description": "找不到该用户信息，请先登录"
                    }
                }
            }
        }
    },
    "definitions": {
        "handler.Token": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "model.User": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "model.homework": {
            "type": "object",
            "properties": {
                "class": {
                    "type": "string"
                },
                "content": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "teacher": {
                    "type": "string"
                },
                "time": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0.0",
	Host:        "localhost",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "TheEventsGarden API",
	Description: "记事园API",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
