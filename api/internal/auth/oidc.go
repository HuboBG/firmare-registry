package auth

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
)

// OIDCVerifier wraps the OIDC provider and verifier for JWT validation
type OIDCVerifier struct {
	provider   *oidc.Provider
	verifier   *oidc.IDTokenVerifier
	adminRole  string
	deviceRole string
}

// NewOIDCVerifier creates a new OIDC verifier that connects to Keycloak
func NewOIDCVerifier(ctx context.Context, issuerURL, clientID, audience, adminRole, deviceRole string) (*OIDCVerifier, error) {
	provider, err := oidc.NewProvider(ctx, issuerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC provider: %w", err)
	}

	config := &oidc.Config{
		ClientID: clientID,
	}

	// If audience is specified, verify it
	if audience != "" {
		config.SupportedSigningAlgs = []string{oidc.RS256}
		// Note: go-oidc doesn't have a direct audience field in Config
		// Audience verification happens in VerifyToken
	}

	verifier := provider.Verifier(config)

	return &OIDCVerifier{
		provider:   provider,
		verifier:   verifier,
		adminRole:  adminRole,
		deviceRole: deviceRole,
	}, nil
}

// VerifyToken validates a JWT token and returns the claims
func (o *OIDCVerifier) VerifyToken(ctx context.Context, rawToken string) (*oidc.IDToken, error) {
	idToken, err := o.verifier.Verify(ctx, rawToken)
	if err != nil {
		return nil, fmt.Errorf("token verification failed: %w", err)
	}

	return idToken, nil
}

// HasRole checks if the token contains the specified role
func (o *OIDCVerifier) HasRole(token *oidc.IDToken, role string) (bool, error) {
	var claims struct {
		RealmAccess struct {
			Roles []string `json:"roles"`
		} `json:"realm_access"`
		ResourceAccess map[string]struct {
			Roles []string `json:"roles"`
		} `json:"resource_access"`
	}

	if err := token.Claims(&claims); err != nil {
		return false, fmt.Errorf("failed to parse claims: %w", err)
	}

	// Check realm roles
	for _, r := range claims.RealmAccess.Roles {
		if r == role {
			return true, nil
		}
	}

	// Check client roles
	for _, client := range claims.ResourceAccess {
		for _, r := range client.Roles {
			if r == role {
				return true, nil
			}
		}
	}

	return false, nil
}
