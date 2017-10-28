package common

import (
	"fmt"
	"time"
	"strings"
	"github.com/gin-gonic/gin"
)

const (
	OriginKey           = "Origin"
	MaxAgeKey           = "Access-Control-Max-Age"
	AllowOriginKey      = "Access-Control-Allow-Origin"
	AllowHeadersKey     = "Access-Control-Allow-Headers"
	AllowMethodsKey     = "Access-Control-Allow-Methods"
	AllowCredentialsKey = "Access-Control-Allow-Credentials"
	ExposeHeadersKey    = "Access-Control-Expose-Headers"
	RequestMethodKey    = "Access-Control-Request-Method"
	RequestHeadersKey   = "Access-Control-Request-Headers"
)

const (
	optionsMethod = "OPTIONS"
)

type CorsConfig struct {
	ValidateHeaders bool
	Origins         string
	origins         []string
	RequestHeaders  string
	requestHeaders  []string
	ExposedHeaders  string
	Methods         string
	methods         []string
	MaxAge          time.Duration
	maxAge          string
	Credentials     bool
	credentials     string
}

func (config *CorsConfig) prepare() {

	config.origins = strings.Split(config.Origins, ", ")
	config.methods = strings.Split(config.Methods, ", ")
	config.requestHeaders = strings.Split(config.RequestHeaders, ", ")
	config.maxAge = fmt.Sprintf("%.f", config.MaxAge.Seconds())
	config.credentials = fmt.Sprintf("%t", config.Credentials)

	for idx, header := range config.requestHeaders {
		config.requestHeaders[idx] = strings.ToLower(header)
	}
}

func CorsMiddleware(config CorsConfig) gin.HandlerFunc {

	forceOriginMatch := false

	if config.Origins == "" {
		panic("You must set at least a single valid origin. If you don't want CORS, to apply, simply remove the middleware.")
	}

	if config.Origins == "*" {
		forceOriginMatch = true
	}

	config.prepare()

	return func(context *gin.Context) {

		currentOrigin := context.Request.Header.Get(OriginKey)
		context.Writer.Header().Add("Vary", OriginKey)
		context.Writer.Header().Add("Content-Type", "application/json")

		if currentOrigin == "" {
			return
		}

		originMatch := false

		if forceOriginMatch == false {
			originMatch = matchOrigin(currentOrigin, config)
		}

		if forceOriginMatch || originMatch {

			valid := false
			preFlight := false

			if context.Request.Method == optionsMethod {

				requestMethod := context.Request.Header.Get(RequestMethodKey)

				if requestMethod != "" {
					preFlight = true
					valid = handlePreFlight(context, config, requestMethod)
				}
			}

			if preFlight == false {
				valid = handleRequest(context, config)
			}

			if valid {

				if config.Credentials {
					context.Writer.Header().Set(AllowCredentialsKey, config.credentials)
					context.Writer.Header().Set(AllowOriginKey, currentOrigin)
				} else if forceOriginMatch {
					context.Writer.Header().Set(AllowOriginKey, "*")
				} else {
					context.Writer.Header().Set(AllowOriginKey, currentOrigin)
				}

				if preFlight {
					context.AbortWithStatus(200)
				}

				return
			}
		}

		context.Abort()
	}
}

func handlePreFlight(context *gin.Context, config CorsConfig, requestMethod string) bool {

	if ok := validateRequestMethod(requestMethod, config); ok == false {
		return false
	}

	if ok := validateRequestHeaders(context.Request.Header.Get(RequestHeadersKey), config); ok {

		context.Writer.Header().Set(AllowMethodsKey, config.Methods)
		context.Writer.Header().Set(AllowHeadersKey, config.RequestHeaders)

		if config.maxAge != "0" {
			context.Writer.Header().Set(MaxAgeKey, config.maxAge)
		}

		return true
	}

	return false
}

func handleRequest(context *gin.Context, config CorsConfig) bool {

	if config.ExposedHeaders != "" {
		context.Writer.Header().Set(ExposeHeadersKey, config.ExposedHeaders)
	}

	return true
}

func matchOrigin(origin string, config CorsConfig) bool {

	for _, value := range config.origins {

		if value == origin {
			return true
		}
	}

	return false
}

func validateRequestMethod(requestMethod string, config CorsConfig) bool {

	if config.ValidateHeaders == false {
		return true
	}

	if requestMethod != "" {

		for _, value := range config.methods {

			if value == requestMethod {
				return true
			}
		}
	}

	return false
}

func validateRequestHeaders(requestHeaders string, config CorsConfig) bool {

	if config.ValidateHeaders == false {
		return true
	}

	headers := strings.Split(requestHeaders, ",")

	for _, header := range headers {

		match := false
		header = strings.ToLower(strings.Trim(header, " \t\r\n"))

		for _, value := range config.requestHeaders {

			if value == header {
				match = true
				break
			}
		}

		if match == false {
			return false
		}
	}

	return true
}