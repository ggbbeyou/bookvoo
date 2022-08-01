// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
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
        "/api/v1/depth": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "交易相关"
                ],
                "summary": "深度信息",
                "parameters": [
                    {
                        "type": "string",
                        "description": "eg: ethusd",
                        "name": "symbol",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "默认100，最大5000",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/common.Response"
                        }
                    }
                }
            }
        },
        "/api/v1/market/klines": {
            "get": {
                "description": "行情k线数据接口",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "行情接口"
                ],
                "summary": "行情k线数据",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/common.Response"
                        }
                    }
                }
            }
        },
        "/api/v1/order/cancel": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "取消还未完成的订单",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "交易相关"
                ],
                "summary": "取消一个委托订单",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer 用户令牌",
                        "name": "Authorization",
                        "in": "header"
                    },
                    {
                        "description": "请求参数",
                        "name": "object",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.cancel_order_request"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/common.Response"
                        }
                    }
                }
            }
        },
        "/api/v1/order/new": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "新订单，支持限价单、市价单\n不同订单类型的参数要求：\n限价单: {\"symbol\": \"ethusd\", \"order_type\": \"limit\", \"side\": \"sell\", \"price\": \"1.00\", \"quantity\": \"100\"}\n市价-按数量: {\"symbol\": \"ethusd\", \"order_type\": \"market\", \"side\": \"sell\", \"quantity\": \"100\"}\n市价-按金额: {\"symbol\": \"ethusd\", \"order_type\": \"market\", \"side\": \"sell\", \"amount\": \"1000.00\"}",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "交易相关"
                ],
                "summary": "创建一个新委托订单",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer 用户令牌",
                        "name": "Authorization",
                        "in": "header"
                    },
                    {
                        "description": "请求参数",
                        "name": "object",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.new_order_request"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/common.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.cancel_order_request": {
            "type": "object",
            "properties": {
                "order_id": {
                    "type": "string"
                }
            }
        },
        "api.new_order_request": {
            "type": "object",
            "required": [
                "order_type",
                "side",
                "symbol"
            ],
            "properties": {
                "amount": {
                    "type": "string",
                    "example": "100.00"
                },
                "order_type": {
                    "type": "string",
                    "example": "limit/market"
                },
                "price": {
                    "type": "string",
                    "example": "1.00"
                },
                "quantity": {
                    "type": "string",
                    "example": "12"
                },
                "side": {
                    "type": "string",
                    "example": "sell/buy"
                },
                "symbol": {
                    "type": "string",
                    "example": "ethusd"
                }
            }
        },
        "common.Response": {
            "type": "object",
            "properties": {
                "data": {},
                "ok": {
                    "type": "integer"
                },
                "reason": {
                    "type": "string"
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
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
