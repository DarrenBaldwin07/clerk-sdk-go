package clerk

import (
	"net/http"
)

type InstanceService service

type UpdateInstanceParams struct {
	// TestMode can be used to toggle test mode for this instance.
	// Defaults to true for development instances.
	TestMode *bool `json:"test_mode,omitempty"`

	// HIBP is used to configure whether Clerk should use the
	// "Have I Been Pawned" service to check passwords against
	// known security breaches.
	// By default, this is enabled in all instances.
	HIBP *bool `json:"hibp,omitempty"`

	// EnhancedEmailDeliverability controls how Clerk delivers emails.
	// Specifically, when set to true, if the instance is a production
	// instance, OTP verification emails are sent by the Clerk's shared
	// domain via Postmark.
	EnhancedEmailDeliverability *bool `json:"enhanced_email_deliverability,omitempty"`

	// SupportEmail is the contact email address that will be displayed
	// on the frontend, in case your instance users need support.
	// If the empty string is provided, the support email that is currently
	// configured in the instance will be removed.
	SupportEmail *string `json:"support_email,omitempty"`

	// ClerkJSVersion allows you to request a specific Clerk JS version on the Clerk Hosted Account pages.
	// If an empty string is provided, the stored version will be removed.
	// If an explicit version is not set, the Clerk JS version will be automatically be resolved.
	ClerkJSVersion *string `json:"clerk_js_version,omitempty"`
}

func (s *InstanceService) Update(params UpdateInstanceParams) error {
	req, _ := s.client.NewRequest(http.MethodPatch, "instance", &params)

	_, err := s.client.Do(req, nil)
	return err
}