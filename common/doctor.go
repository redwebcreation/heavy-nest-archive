package common

import (
	"context"
	"fmt"
	"io/ioutil"
	"log/syslog"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/spf13/cobra"
	"github.com/wormable/nest/ansi"
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
			ValidateApplicationsConfigurations,
			ValidateLogPolicies,
			EnsureDnsRecordPointsToHost,
		},
	}
	for _, check := range diagnosis.Checks {
		check(diagnosis)
	}

	return diagnosis
}

func EnsureDnsRecordPointsToHost(diagnosis *Diagnosis) {
	response, err := http.Get("http://checkip.amazonaws.com")
	if err == nil {
		defer response.Body.Close()
	} else if strings.Contains(err.Error(), "no such host") {
		diagnosis.NewWarning(Warning{
			Title:  "It looks like you're not connected to internet, or AWS is down.",
			Advice: "If this is a production server, GL HF, you'll need it.",
		})
		return
	} else {
		ansi.Check(err)
	}

	rawPublicIp, err := ioutil.ReadAll(response.Body)
	ansi.Check(err)

	publicIp := strings.TrimSpace(string(rawPublicIp))

	for _, application := range Config.Applications {
		ips, err := net.LookupIP(application.Host)
		if err != nil {
			diagnosis.NewWarning(Warning{
				Title:  fmt.Sprintf("DNS record for %s does not exist", application.Host),
				Advice: fmt.Sprintf("Create a DNS record for %s pointing to %s", application.Host, publicIp),
			})
			continue
		}

		hasMatchingIp := false

		for _, ip := range ips {
			if ip.String() == publicIp {
				hasMatchingIp = true
			}
		}

		if !hasMatchingIp {
			diagnosis.NewWarning(Warning{
				Title:  fmt.Sprintf("DNS record for %s does not point to %s", application.Host, publicIp),
				Advice: fmt.Sprintf("Run `dig %s` to check if the DNS record is correct", application.Host),
			})
		}
	}
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

func ValidateApplicationsConfigurations(d *Diagnosis) {
	for _, application := range Config.Applications {
		if application.Registry != "" {
			validRegistry := false

			for _, registry := range Config.Registries {
				if application.Registry == registry.Name {
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
			for _, policy := range Config.LogPolicies {
				if application.LogPolicy == policy.Name {
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

		// -1 is the default value
		if application.Quotas.Memory != "-1" {
			_, err := ParseMemoryQuota(application.Quotas.Memory)

			d.NewError(Error{
				Title: fmt.Sprintf("invalid memory quota [%s]", application.Quotas.Memory),
				Error: err,
			})
		}
	}
}

func ValidateLogPolicies(d *Diagnosis) {
	for _, policy := range Config.LogPolicies {
		for _, rule := range policy.Rules {
			_, err := rule.ShouldLog(syslog.LOG_ERR)

			if err != nil {
				d.NewError(Error{
					Title: fmt.Sprintf("invalid log policy %s [%s]", policy.Name, rule),
					Error: err,
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
