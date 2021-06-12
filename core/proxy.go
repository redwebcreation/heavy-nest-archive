package core

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/redwebcreation/hez2/globals"
	"go.uber.org/zap"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ProxiableContainer struct {
	Name        string
	Ipv4        string
	VirtualHost string
	VirtualPort string
	Container   *types.ContainerJSON
}

func GetWhitelistedDomains() []string {
	var domains []string

	proxiableContainers, _ := GetProxiableContainers()

	for _, proxiableContainer := range proxiableContainers {
		domains = append(domains, proxiableContainer.VirtualHost)
	}

	return domains
}

func GetProxiableContainers() ([]ProxiableContainer, error) {
	var proxiableContainers []ProxiableContainer

	for _, application := range globals.Config.Applications {
		container, _ := globals.GetContainer(application.Name())

		inspectedContainer, err := globals.Docker.ContainerInspect(context.Background(), container.ID)
		if err != nil {
			return nil, err
		}

		containerNetwork, err := globals.FindNetwork(application.Network)

		if err != nil {
			return nil, err
		}

		var networkConfiguration types.EndpointResource

		for _, c := range containerNetwork.Containers {
			if container.Names[0] == c.Name {
				networkConfiguration = c
				break
			}
		}

		virtualHost := ""
		virtualPort := "80"

		for _, envVariable := range inspectedContainer.Config.Env {
			if strings.HasPrefix(envVariable, "VIRTUAL_HOST=") {
				virtualHost = strings.SplitAfter(envVariable, "VIRTUAL_HOST=")[1]
			}

			if strings.HasPrefix(envVariable, "VIRTUAL_PORT=") {
				virtualPort = strings.SplitAfter(envVariable, "VIRTUAL_PORT=")[1]
			}
		}

		if virtualHost == "" {
			continue
		}

		proxiableContainers = append(proxiableContainers, ProxiableContainer{
			Name:        networkConfiguration.Name,
			Ipv4:        networkConfiguration.IPv4Address,
			VirtualHost: virtualHost,
			VirtualPort: virtualPort,
			Container:   &inspectedContainer,
		})
	}

	return proxiableContainers, nil
}

func HandleRequest(lastApplyExecution string, proxiables []ProxiableContainer) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		if lastApplyExecution != GetLastApplyExecution() {
			proxiables, _ = GetProxiableContainers()
		}

		ip, _, _ := net.SplitHostPort(request.RemoteAddr)
		request.Header.Set("X-Forwarded-For", ip)

		for _, proxiableContainer := range proxiables {
			if request.Host == proxiableContainer.VirtualHost {
				success := ForwardRequest(proxiableContainer, writer, request)
				if success {
					globals.Logger.Info(
						"request.success",
						zap.String("method", request.Method),
						zap.String("ip", ip),
						zap.String("vhost", request.Host),
					)
				}
				return
			}
		}

		globals.Logger.Info(
			"request.invalid",
			zap.String("method", request.Method),
			zap.String("ip", ip),
			zap.String("vhost", request.Host),
		)
		writer.WriteHeader(404)
		writer.Write([]byte("404. That’s an error. \nThe requested URL " + request.RequestURI + " was not found on this server. That’s all we know."))
	}
}

func ForwardRequest(container ProxiableContainer, writer http.ResponseWriter, request *http.Request) bool {
	containerUrl, err := url.Parse("http://" + container.Ipv4 + ":" + container.VirtualPort)

	if err != nil {
		zap.L().Error(
			"url.invalid",
			zap.String("error", err.Error()),
		)
		internalServerError(writer)
		return false
	}

	request.Host = containerUrl.Host + ":" + container.VirtualPort
	request.URL.Host = containerUrl.Host + ":" + container.VirtualPort
	request.URL.Scheme = containerUrl.Scheme
	request.RequestURI = ""

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		zap.L().Error(
			"forward.failed",
			zap.String("container_url", containerUrl.String()),
			zap.String("error", err.Error()),
		)
		internalServerError(writer)
		return false
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
	return true
}

func internalServerError(writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusInternalServerError)
	writer.Write([]byte("Internal Server Error"))
}
