// Package deprecations contains deprecation messages for various provider components.
package deprecations

// MetalDeprecationMessage is the standard deprecation message for all Equinix Metal resources and data sources.
// Equinix Metal platform will sunset on June 30, 2026, and all Metal functionality will be removed
// in provider version 5.0.0.
const MetalDeprecationMessage = "Equinix Metal will reach end of life on June 30, 2026. All Metal resources will be removed in version 5.0.0 of this provider. Use version 4.x of this provider for continued use through sunset. See https://docs.equinix.com/metal/ for more information."
