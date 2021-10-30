package common

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/wormable/nest/globals"
)

type Warning struct {
	Title  string
	Advice string
}

type Error struct {
	Title     string
	Error     error
	Solutions []string // An array of command names that can resolve the problem
}

type Diagnosis struct {
	Config       Configuration
	Checks       []func(*Diagnosis)
	WarningCount int
	ErrorCount   int
	Warnings     []Warning
	Errors       []Error
}

func (d *Diagnosis) NewWarning(w Warning) *Diagnosis {
	d.Warnings = append(d.Warnings, w)
	d.WarningCount++
	return d
}
func (d *Diagnosis) NewError(e Error) *Diagnosis {
	d.Errors = append(d.Errors, e)
	d.ErrorCount++
	return d
}

func AnalyseConfig() *Diagnosis {
	diagnosis := &Diagnosis{
		Config: Config,
		Checks: []func(*Diagnosis){
			ValidateRegistriesConfig,
			DefaultNetworkIsValid,
			PrivateCertificateAuthorityIsValid,
			BackendsAreConnected,
			ValidateApplicationsConfigurations,
			ValidateLogPolicies,
		},
	}
	for _, check := range diagnosis.Checks {
		check(diagnosis)
	}

	return diagnosis
}

func ValidateRegistriesConfig(d *Diagnosis) {
	for _, registry := range Config.Registries {
		_, err := globals.Docker.RegistryLogin(context.Background(), types.AuthConfig{
			Username:      registry.Username,
			Password:      registry.Password,
			ServerAddress: registry.Host,
		})
		if err != nil {
			d.NewError(Error{
				Title: fmt.Sprintf("could not login to registry [%s]", registry.Host),
				Error: err,
			})
		}
	}
}

func DefaultNetworkIsValid(d *Diagnosis) {
	err := networkIsValid(Config.DefaultNetwork)
	if err != nil {
		d.NewError(*err)
	}
}

func BackendsAreConnected(d *Diagnosis) {
	for _, backend := range Config.Backends {
		// TODO: Check if the backend joined the network
		response, err := http.Get("http://" + backend + "/version")

		if err != nil {
			d.NewError(Error{
				Title: fmt.Sprintf("Cound not connect to backend %s", backend),
				Error: err,
			})
			continue
		}

		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			d.NewError(Error{
				Title: fmt.Sprintf("could not read backend %s response", backend),
				Error: err,
			})
			continue
		}

		nameAndVersion := strings.Split(string(body), "@")

		if len(nameAndVersion) != 2 && nameAndVersion[0] == "nest" {
			d.NewError(Error{
				Title: fmt.Sprintf("backend %s returned an invalid response, are you sure it's a node?", backend),
			})
		}

		if nameAndVersion[1] != globals.Version {
			// TODO: Maybe it should only be a warning
			d.NewError(Error{
				Title: fmt.Sprintf("Mismatch between current version %s and backend's version %s", globals.Version, nameAndVersion[1]),
			})
		}
	}
}

func ValidateApplicationsConfigurations(d *Diagnosis) {
	for _, application := range Config.Applications {
		if application.Registry != "" {
			validRegistry := false

			for name := range Config.Registries {
				if application.Registry == name {
					validRegistry = true
				}
			}

			if !validRegistry {
				d.NewError(Error{
					Title: fmt.Sprintf("registry [%s] not found", application.Registry),
				})
			}
		}

		for _, envFile := range application.EnvFiles {
			_, err := os.Stat(envFile)

			if err != nil {
				d.NewError(Error{
					Title: fmt.Sprintf("can not retrieve env file [%s]", envFile),
					Error: err,
				})
			}
		}

		for _, volume := range application.Volumes {
			fromAndTo := strings.Split(volume, ":")
			_, err := os.Stat(fromAndTo[0])
			if err != nil {
				d.NewError(Error{
					Title: fmt.Sprintf("invalid volume source [%s]", fromAndTo[0]),
					Error: err,
				})
			}
		}

		if application.LogPolicy != "" {
			validLogPolicy := false
			for name := range Config.LogPolicies {
				if application.LogPolicy == name {
					validLogPolicy = true
				}
			}

			if !validLogPolicy {
				d.NewError(Error{
					Title: fmt.Sprintf("log policy [%s] not found", application.LogPolicy),
				})
			}
		}

		if application.Network != "" {
			err := networkIsValid(application.Network)
			if err != nil {
				d.NewError(*err)
			}
		}
	}
}

func ValidateLogPolicies(d *Diagnosis) {
	for name, rules := range Config.LogPolicies {
		for _, rule := range rules {
			_, err := rule.MustCompile(ErrorLevel)

			if err != nil {
				d.NewError(Error{
					Title: fmt.Sprintf("invalid log policy %s [%s]", name, rule),
					Error: err,
				})
			}
		}
	}
}

func PrivateCertificateAuthorityIsValid(d *Diagnosis) {
	now := time.Now()
	for _, cert := range []string{globals.CACertificate, globals.ServerCertificate} {
		bytes, err := os.ReadFile(cert)

		if err != nil {
			d.NewWarning(Warning{
				Title:     fmt.Sprintf("could not read certificate %s", cert),
				Advice: "run `nest certificates init` to generate needed certificates",
			})
			continue
		}

		block, _ := pem.Decode(bytes)
		certificate, err := x509.ParseCertificate(block.Bytes)

		if err != nil {
			d.NewWarning(Warning{
				Title:     fmt.Sprintf("cound not parse certificate %s", cert),
				Advice: "run `nest certificates init` to generate needed certificates",
			})
			continue
		}

		if now.AddDate(0, 6, 0).After(certificate.NotAfter) {
			expiresIn := certificate.NotAfter.Sub(now).Hours()
			var expiresInFormatted string

			if expiresIn > 24 {
				expiresInFormatted = fmt.Sprintf("%.1f days", expiresIn/24.0)
			} else {
				expiresInFormatted = fmt.Sprintf("%.1f hours", expiresIn)
			}

			if expiresIn < 0 {
				d.NewError(Error{
					Title:     fmt.Sprintf("certificate %s expired", cert),
					Error:     fmt.Errorf("certificate expired at %s", certificate.NotAfter.UTC().String()),
					Solutions: []string{"certificates init"},
				})
			} else {
				d.NewWarning(Warning{
					Title:  fmt.Sprintf("certificate %s expiring in %s", cert, expiresInFormatted),
					Advice: "<link to the docs>", // todo
				})
			}

		}
	}
}

func networkIsValid(name string) *Error {
	networks, err := globals.Docker.NetworkList(context.Background(), types.NetworkListOptions{
		Filters: filters.NewArgs(
			filters.KeyValuePair{
				Key:   "name",
				Value: name,
			},
		),
	})
	if err != nil {
		return &Error{
			Title: fmt.Sprintf("network [%s] not found", name),
			Error: err,
		}
	}

	if len(networks) == 0 {
		return &Error{
			Title: fmt.Sprintf("network [%s] not found", name),
		}
	}

	return nil
}

func (d *Diagnosis) CommandIsSolution(cmd *cobra.Command) bool {
	for _, err := range d.Errors {
		for _, name := range err.Solutions {
			cmp := cmd.Name()
			current := cmd.Parent()
			for {
				if current.Name() == "nest" {
					break
				}

				cmp = current.Name() + " " + cmp
				current = current.Parent()
			}

			if name == cmp {
				return true
			}
		}
	}

	return false
}
