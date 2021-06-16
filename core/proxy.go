package core

import (
	"go.uber.org/zap"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func GetWhitelistedDomains() []string {
	var domains = make([]string, len(Config.Applications))

	for _, application := range Config.Applications {
		domains = append(domains, application.Host)
	}

	return domains
}

func HandleRequest(writer http.ResponseWriter, request *http.Request) {
	ip, _, _ := net.SplitHostPort(request.RemoteAddr)
	request.Header.Set("X-Forwarded-For", ip)
	request.Header.Set("X-Forwarded-Proto", request.Proto)

	for _, application := range Config.Applications {
		if request.Host == application.Host {
			success := ForwardRequest(application, writer, request)
			if success {
				Logger.Info(
					"request.success",
					zap.String("method", request.Method),
					zap.String("ip", ip),
					zap.String("vhost", request.Host),
				)
			} else {
				internalServerError(writer)
			}
			return
		}
	}

	Logger.Info(
		"request.invalid",
		zap.String("method", request.Method),
		zap.String("ip", ip),
		zap.String("vhost", request.Host),
	)
	writer.WriteHeader(404)
	writer.Write([]byte("404. That’s an error. \nThe requested URL " + request.RequestURI + " was not found on this server. That’s all we know."))
}

func ForwardRequest(application Application, writer http.ResponseWriter, request *http.Request) bool {
	container, err := application.GetContainer(AnyType)
	if err != nil {
		Logger.Error(
			"container.missing",
			zap.String("vhost", request.Host),
			zap.String("container_name", application.Name(ApplicationContainer)),
		)
		return false
	}

	containerUrl, err := url.Parse("http://" + container.Ip + ":" + application.ContainerPort)
	if err != nil {
		Logger.Error(
			"url.invalid",
			zap.String("error", err.Error()),
		)
		return false
	}

	request.Host = application.Host
	request.URL.Host = containerUrl.Host
	request.URL.Scheme = containerUrl.Scheme
	request.RequestURI = ""

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		Logger.Error(
			"forward.failed",
			zap.String("container_url", containerUrl.String()),
			zap.String("error", err.Error()),
		)
		return false
	}

	for key, values := range response.Header {
		for _, value := range values {
			writer.Header().Add(key, value)
		}
	}


	done := make(chan bool)

	go func() {
		for {
			select {
			case <-time.Tick(10 * time.Millisecond):
				writer.(http.Flusher).Flush()
			case <-done:
				return
			}
		}
	}()

	var trailerKeys []string

	for key := range response.Trailer {
		trailerKeys = append(trailerKeys, key)
	}

	writer.Header().Set("Strict-Transport-Security", "max-age=15768000 ; includeSubDomains")
	writer.Header().Set("Trailer", strings.Join(trailerKeys, ","))

	writer.WriteHeader(response.StatusCode)
	io.Copy(writer, response.Body)

	for key, values := range response.Trailer {
		for _, value := range values {
			writer.Header().Set(key, value)
		}
	}

	close(done)
	return true
}

func internalServerError(writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusInternalServerError)
	writer.Write([]byte("Internal Server Error"))
}
