package engine

import (
	"fmt"

	"github.com/go-logr/logr"
	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	kyvernov2alpha1 "github.com/kyverno/kyverno/api/kyverno/v2alpha1"
	"github.com/kyverno/kyverno/pkg/config"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	matched "github.com/kyverno/kyverno/pkg/utils/match"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

func findExceptions(
	selector engineapi.PolicyExceptionSelector,
	policy kyvernov1.PolicyInterface,
	rule string,
) ([]*kyvernov2alpha1.PolicyException, error) {
	if selector == nil {
		return nil, nil
	}
	polexs, err := selector.List(labels.Everything())
	if err != nil {
		return nil, err
	}
	var result []*kyvernov2alpha1.PolicyException
	policyName, err := cache.MetaNamespaceKeyFunc(policy)
	if err != nil {
		return nil, fmt.Errorf("failed to compute policy key: %w", err)
	}
	for _, polex := range polexs {
		if polex.Contains(policyName, rule) {
			result = append(result, polex)
		}
	}
	return result, nil
}

// matchesException checks if an exception applies to the resource being admitted
func matchesException(
	selector engineapi.PolicyExceptionSelector,
	policyContext engineapi.PolicyContext,
	rule kyvernov1.Rule,
	cfg config.Configuration,
) (*kyvernov2alpha1.PolicyException, error) {
	candidates, err := findExceptions(selector, policyContext.Policy(), rule.Name)
	if err != nil {
		return nil, err
	}
	gvk, subresource := policyContext.ResourceKind()
	for _, candidate := range candidates {
		err := matched.CheckMatchesResources(
			policyContext.NewResource(),
			candidate.Spec.Match,
			policyContext.NamespaceLabels(),
			policyContext.AdmissionInfo(),
			cfg.GetExcludedGroups(),
			gvk,
			subresource,
		)
		// if there's no error it means a match
		if err == nil {
			return candidate, nil
		}
	}
	return nil, nil
}

// hasPolicyExceptions returns nil when there are no matching exceptions.
// A rule response is returned when an exception is matched, or there is an error.
func (e *engine) hasPolicyExceptions(
	logger logr.Logger,
	ruleType engineapi.RuleType,
	ctx engineapi.PolicyContext,
	rule kyvernov1.Rule,
) *engineapi.RuleResponse {
	// if matches, check if there is a corresponding policy exception
	exception, err := matchesException(e.exceptionSelector, ctx, rule, e.configuration)
	if err != nil {
		logger.Error(err, "failed to match exceptions")
		return nil
	}
	if exception == nil {
		return nil
	}
	key, err := cache.MetaNamespaceKeyFunc(exception)
	if err != nil {
		logger.Error(err, "failed to compute policy exception key", "namespace", exception.GetNamespace(), "name", exception.GetName())
		return engineapi.RuleError(rule.Name, ruleType, "failed to compute exception key", err)
	} else {
		logger.V(3).Info("policy rule skipped due to policy exception", "exception", key)
		return engineapi.RuleSkip(rule.Name, ruleType, "rule skipped due to policy exception "+key).WithException(exception)
	}
}
