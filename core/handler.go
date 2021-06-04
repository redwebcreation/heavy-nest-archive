package core

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func ForwardRequest(container ProxiableContainer, writer http.ResponseWriter, request *http.Request) {
	containerUrl, err := url.Parse("http://" + container.Ipv4 + ":" + container.VirtualPort)

	if err != nil {
		Logger().Error(
			"url.invalid",
			zap.String("error", err.Error()),
		)
		internalServerError(writer)
		return
	}

	request.Host = containerUrl.Host + ":" + container.VirtualPort
	request.URL.Host = containerUrl.Host + ":" + container.VirtualPort
	request.URL.Scheme = containerUrl.Scheme
	request.RequestURI = ""

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		Logger().Error(
			"forward.failed",
			zap.String("container_url", containerUrl.String()),
			zap.String("error", err.Error()),
		)
		internalServerError(writer)
		return
	}

	for key, values := range response.Header {
		for _, value := range values {
			writer.Header().Set(key, value)
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
}

func internalServerError(writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusInternalServerError)
	writer.Write([]byte("Internal Server Error"))
}

func Logger() *zap.Logger {
	config, _ := FindConfig(ConfigFile()).Resolve()

	var loggerConfig zap.Config

	loggerConfig.Level = zap.NewAtomicLevelAt(zapcore.Level(config.Proxy.Logs.Level))
	loggerConfig.OutputPaths = config.Proxy.Logs.Redirections
	loggerConfig.Encoding = "json"
	loggerConfig.EncoderConfig = zapcore.EncoderConfig{
		MessageKey:  "message",
		LevelKey:    "level",
		EncodeLevel: zapcore.LowercaseLevelEncoder,
	}

	logger, _ := loggerConfig.Build()

	defer logger.Sync()

	return logger
}
