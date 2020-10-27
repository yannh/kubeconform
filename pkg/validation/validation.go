package validation

import (
	"fmt"

	"github.com/yannh/kubeconform/pkg/config"
	"github.com/yannh/kubeconform/pkg/resource"
	"github.com/yannh/kubeconform/pkg/schema"
)

// Validator TODO
type Validator interface {
	Validate(res resource.Resource) Result
}

type validator struct {
	repo *schema.Repository
	conf config.Config
}

// func (v validator) Validate(res resource.Resource) Result {
// 	return Result{}
// }

// New TODO
func New(repo *schema.Repository, c config.Config) Validator {
	return &validator{
		repo: repo,
		conf: c,
	}
}

// Single TODO
func (v validator) Validate(res resource.Resource) Result {
	sig, err := res.Signature()
	if err != nil {
		return Result{
			Resource: res,
			Err:      err,
			Status:   Error,
		}
	}

	if sig.Kind == "" {
		// We skip resoures that don't have a Kind defined
		return Result{
			Resource: res,
			Err:      nil,
			Status:   Skipped,
		}
	}
	if _, ok := v.conf.SkipKinds[sig.Kind]; ok {
		return Result{
			Resource: res,
			Err:      nil,
			Status:   Skipped,
		}
	}

	s, err := v.repo.Get(sig.Kind, sig.APIVersion, v.conf.KubernetesVersion)
	if err != nil {
		if v.conf.IgnoreMissingSchemas {
			return Result{
				Resource: res,
				Err:      nil,
				Status:   Skipped,
			}
		}

		return Result{
			Resource: res,
			Err:      fmt.Errorf("could not find schema for %s", sig.Kind),
			Status:   Error,
		}
	}

	re, err := res.AsMap()
	if err != nil {
		return Result{
			Resource: res,
			Err:      err,
			Status:   Error,
		}
	}

	err = s.Validate(re)
	status := Valid
	if err != nil {
		status = Error
	}
	return Result{
		Resource: res,
		Err:      err,
		Status:   status,
	}
}
