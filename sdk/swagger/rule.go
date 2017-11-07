/*
 * Oathkeeper
 *
 * Oathkeeper
 *
 * OpenAPI spec version: Latest
 * Contact: hi@ory.am
 * Generated by: https://github.com/swagger-api/swagger-codegen.git
 */

package swagger

// A rule
type Rule struct {

	// AllowAnonymous sets if the endpoint is public, thus not needing any authorization at all.
	AllowAnonymous bool `json:"allowAnonymous,omitempty"`

	// BypassAccessControlPolicies if set true disables checking access control policies.
	BypassAccessControlPolicies bool `json:"bypassAccessControlPolicies,omitempty"`

	// BypassAuthorization if set true disables firewall capabilities.
	BypassAuthorization bool `json:"bypassAuthorization,omitempty"`

	// Description describes the rule.
	Description string `json:"description,omitempty"`

	// ID the a unique id of a rule.
	Id string `json:"id,omitempty"`

	// MatchesMethods is a list of HTTP methods that this rule matches.
	MatchesMethods []string `json:"matchesMethods,omitempty"`

	// MatchesPathCompiled is a regular expression of paths this rule matches.
	MatchesPath string `json:"matchesPath,omitempty"`

	// RequiredScopes is the action this rule requires.
	RequiredAction string `json:"requiredAction,omitempty"`

	// RequiredScopes is the resource this rule requires.
	RequiredResource string `json:"requiredResource,omitempty"`

	// RequiredScopes is a list of scopes that are required by this rule.
	RequiredScopes []string `json:"requiredScopes,omitempty"`
}
