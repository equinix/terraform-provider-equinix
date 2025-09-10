# Workload Identity Federation (WIF) using Equinix STS

This guide walks you through setting up Workload Identity Federation (WIF) using Equinix STS. It enables your workloads to securely authenticate with Equinix services without relying on long-lived credentials.

## Prerequisites

- An Equinix account with an organization and project - [Sign up here](https://portal.equinix.com)
- Access to your identity provider (e.g., Azura AD, Terraform HCP)
- Equinix API credentials for an administrator user - See [Generating Client ID and Client Secret key](https://docs.equinix.com/equinix-api/api-authentication#generate-client-id-and-client-secret) for more details

## Step 1: Obtain Authentication Token

First, get an authentication token to make API calls:

```bash
export CLIENT_ID="your_client_id"
export CLIENT_SECRET="your_client_secret"

TOKEN=$(curl -s "https://api.equinix.com/oauth2/v1/token" \
  --json "{
    \"grant_type\": \"client_credentials\",
    \"client_id\": \"$CLIENT_ID\",
    \"client_secret\": \"$CLIENT_SECRET\"
  }" | jq -r '.access_token')
```

## Step 2: Establish Trust with Identity Provider

Create a trust relationship with your workload's identity provider:

```bash
ORG_ID="your_organization_id"

OIDCP=$(curl -s "https://sts.eqix.equinix.com/use/createOidcProvider" \
  -H "Authorization: Bearer $TOKEN" \
  --json '{
    "name": "Your Provider Name",
    "issuerLocation": "https://your-idp-issuer-url",
    "trustedClientIds": [
      "your-client-id"
    ],
    "idpPrefix": "your-prefix"
  }')

# Save the IdP ID for later use
IDP_ID=$(echo "$OIDCP" | jq -r '.result.idpId')
echo "Identity Provider ID: $IDP_ID"
```

## Step 3: Authorize Your Workloads

You can authorize workloads using either role assignments or access policies:

The subject in the principal name should match the sub claim of the JWT token issued by your identity provider. This ensures that the workload can be authenticated and authorized correctly.

### Option A: Using Role Assignments

```bash
# First get a JWT token
JWT=$(curl -s "https://api.equinix.com/oauth2/v1/userinfo" \
  -H "Authorization: Bearer $TOKEN" \
  | jq -r '.jwt_token')

# Create role assignment
curl -s "https://api.equinix.com/am/v3/assignments" \
  -H "Authorization: Bearer $JWT" \
  --json '{
    "principal": {
      "type": "FEDERATED",
      "name": "principal:'$ORG_ID':'${IDP_ID:4}':{subject}"
    },
    "roleName": "your-required-role",
    "resource": {
      "id": "'$ORG_ID'",
      "type": "ORGANIZATION"
    }
  }'
```

### Option B: Using Access Policies

```bash
ACCESS_URL="https://access.equinix.com"

curl -s "$ACCESS_URL/use/createAccessPolicy" \
  -H "Authorization: Bearer $TOKEN" \
  --json '{
    "accessPolicyId": "accesspolicy:your-policy-name",
    "grants": [
      "principal:'$ORG_ID':'${IDP_ID:4}':{subject}"
    ],
    "tags": {},
    "permissions": [{
      "serviceActions": [{
        "serviceId": "Equinix Service ID", 
        "actions": ["Action1", "Action2"]
      }],
      "resources": "all"
    }]
  }'
```

## Troubleshooting

If your workloads fail to authenticate:

1. Verify the trust relationship was established correctly
2. Check that the workload's identity matches exactly what's specified in your access policies or role assignments
3. Ensure the required permissions have been granted
4. Look for any errors in the token exchange process

For additional support, contact Equinix customer service.