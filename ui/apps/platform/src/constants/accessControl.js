/* constants specific to Roles */
export const NO_ACCESS = 'NO_ACCESS';
export const READ_ACCESS = 'READ_ACCESS';
export const READ_WRITE_ACCESS = 'READ_WRITE_ACCESS';

export const defaultRoles = {
    Admin: true,
    Analyst: true,
    'Continuous Integration': true,
    None: true,
    'Sensor Creator': true,
};

/* constants specific to Auth Providers */
export const availableAuthProviders = [
    {
        label: 'Auth0',
        value: 'auth0',
    },
    {
        label: 'OpenID Connect',
        value: 'oidc',
    },
    {
        label: 'SAML 2.0',
        value: 'saml',
    },
    {
        label: 'User Certificates',
        value: 'userpki',
    },
    {
        label: 'Google IAP',
        value: 'iap',
    },
];

export const oidcCallbackValues = {
    auto: 'Auto-select (recommended)',
    post: 'HTTP POST',
    fragment: 'Fragment',
    query: 'Query',
};

export function getAuthProviderLabelByValue(value) {
    return availableAuthProviders.find((e) => e.value === value).label;
}

export const defaultMinimalReadAccessResources = [
    'Alert',
    'Cluster',
    'Config',
    'Deployment',
    'Image',
    'Namespace',
    'NetworkPolicy',
    'NetworkGraph',
    'Node',
    'Policy',
    'Secret',
];

// Default to giving new roles read access to specific resources.
export const defaultNewRolePermissions = defaultMinimalReadAccessResources.reduce(
    (map, resource) => {
        const newMap = map;
        newMap[resource] = READ_ACCESS;
        return newMap;
    },
    {}
);

export const defaultSelectedRole = {
    name: '',
    globalAccess: 'NO_ACCESS',
    resourceToAccess: defaultNewRolePermissions,
};
