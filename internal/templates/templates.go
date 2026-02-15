package templates

import _ "embed"

//go:embed tenant/rbac.yaml
var TenantRBAC []byte

//go:embed tenant/quota.yaml
var TenantQuota []byte

//go:embed tenant/limitrange.yaml
var TenantLimitRange []byte

//go:embed tenant/networkpolicy.yaml
var TenantNetworkPolicy []byte
