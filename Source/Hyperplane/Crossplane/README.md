
## ytt
Try using ytt https://carvel.dev/ytt/docs/v0.43.0/ for templating.

### Installation
https://carvel.dev/ytt/docs/v0.43.0/install/

### Applying compositions
ytt -f <composition> -f ytt | k apply -f- 

## Notes
- Do not modify resource definition while there are claims using it.
- When setting the name of a managed resource (at least for the kubernetes provider Object type) one should set the metadata.name field
- Installing jet-based providers installs A LOT of CRDs and that might give us performance issues on the kubernetes api-server. See discussion https://crossplane.slack.com/archives/C01718T2476/p1645793325962729
- 

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

### Deleting weirdness

When deleting an XR (Composite resource) it will delete the XR before deleting its managed resources.

## References
- Versioning of resources: https://crossplane.io/docs/v1.10/concepts/managed-resources.html#versioning
- 