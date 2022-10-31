

### Install
Apply namespace
Apply roles and role bindings
Apply deployments
Apply provider

NOTES:
Important that the Providers are installed in correct order. The ProviderConfig should not be made before the Provider

### Uninstall
Follow this https://crossplane.io/docs/v1.10/reference/uninstall.html.
Very important that providers and provider config is deleted correctly. You might end up in a state with
resources that waits for finalizers to be deleted (usually I have seen this on some CRDs and ProviderConfig). If that is the case
it is possible to just edit that resource by deleting the finalizers.

- Delete resources
- Delete compositions
- Delete resource definitions
- Delete providers (see above)
- Delete Crossplane
- Delete crossplane CRDs
- Delete dolittle CRDs

## References
- Versioning of resources: https://crossplane.io/docs/v1.10/concepts/managed-resources.html#versioning
- 