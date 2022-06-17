package main

type ProviderStatusNew string

// List of ProviderStatus
const (
	AVAILABLE_ProviderStatus        ProviderStatusNew = "AVAILABLE"
	DEPROVISIONED_ProviderStatus    ProviderStatusNew = "DEPROVISIONED"
	DEPROVISIONING_ProviderStatus   ProviderStatusNew = "DEPROVISIONING"
	FAILED_ProviderStatus           ProviderStatusNew = "FAILED"
	NOT_AVAILABLE_ProviderStatus    ProviderStatusNew = "NOT_AVAILABLE"
	PENDING_APPROVAL_ProviderStatus ProviderStatusNew = "PENDING_APPROVAL"
	PROVISIONED_ProviderStatus      ProviderStatusNew = "PROVISIONED"
)
