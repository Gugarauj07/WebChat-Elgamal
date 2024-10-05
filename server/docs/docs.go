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
        "/connect": {
            "post": {
                "description": "Registra um novo usuário com seu ID e chave pública",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Conecta um usuário ao servidor",
                "parameters": [
                    {
                        "description": "Informações do usuário",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/decrypt": {
            "post": {
                "description": "Decripta uma mensagem usando a chave privada fornecida",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "messages"
                ],
                "summary": "Decripta uma mensagem",
                "parameters": [
                    {
                        "description": "Mensagem encriptada e chave privada",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.DecryptRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/disconnect": {
            "post": {
                "description": "Remove um usuário do servidor",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Desconecta um usuário do servidor",
                "parameters": [
                    {
                        "description": "ID do usuário",
                        "name": "userId",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/encrypt": {
            "post": {
                "description": "Encripta uma mensagem usando a chave pública fornecida",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "messages"
                ],
                "summary": "Encripta uma mensagem",
                "parameters": [
                    {
                        "description": "Mensagem e chave pública",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.EncryptRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.EncryptedMessage"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/generate-keys": {
            "get": {
                "description": "Gera e retorna um novo par de chaves pública e privada",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "keys"
                ],
                "summary": "Gera um novo par de chaves",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.KeyPair"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/public-key/{userId}": {
            "get": {
                "description": "Retorna a chave pública de um usuário específico",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Obtém a chave pública de um usuário",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID do usuário",
                        "name": "userId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/receive-messages": {
            "post": {
                "description": "Retorna as mensagens encriptadas para um usuário específico",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "messages"
                ],
                "summary": "Recebe mensagens",
                "parameters": [
                    {
                        "description": "ID do usuário",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.ReceiveMessagesRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.ChatMessage"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/send-message": {
            "post": {
                "description": "Envia uma mensagem encriptada para um usuário específico",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "messages"
                ],
                "summary": "Envia uma mensagem",
                "parameters": [
                    {
                        "description": "Mensagem encriptada e IDs de remetente e destinatário",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.SendMessageRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/users": {
            "get": {
                "description": "Retorna uma lista de IDs de todos os usuários conectados",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Lista todos os usuários",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.DecryptRequest": {
            "type": "object",
            "properties": {
                "encryptedMessage": {
                    "$ref": "#/definitions/models.EncryptedMessage"
                },
                "privateKey": {
                    "type": "integer"
                }
            }
        },
        "api.EncryptRequest": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "publicKey": {
                    "type": "integer"
                }
            }
        },
        "api.ReceiveMessagesRequest": {
            "type": "object",
            "properties": {
                "userId": {
                    "type": "string"
                }
            }
        },
        "api.SendMessageRequest": {
            "type": "object",
            "properties": {
                "encryptedMessage": {
                    "$ref": "#/definitions/models.EncryptedMessage"
                },
                "receiverId": {
                    "type": "string"
                },
                "senderId": {
                    "type": "string"
                }
            }
        },
        "models.ChatMessage": {
            "type": "object",
            "properties": {
                "content": {
                    "$ref": "#/definitions/models.EncryptedMessage"
                },
                "senderId": {
                    "type": "string"
                }
            }
        },
        "models.EncryptedMessage": {
            "type": "object",
            "properties": {
                "c1": {
                    "type": "string"
                },
                "c2": {
                    "type": "string"
                }
            }
        },
        "models.KeyPair": {
            "type": "object",
            "properties": {
                "privateKey": {
                    "type": "integer"
                },
                "publicKey": {
                    "type": "integer"
                }
            }
        },
        "models.User": {
            "type": "object",
            "properties": {
                "publicKey": {
                    "type": "integer"
                },
                "userId": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:3000",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Chat API",
	Description:      "Esta é uma API de servidor de chat simples.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
