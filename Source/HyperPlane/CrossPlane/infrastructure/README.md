

### Uninstall
Follow this https://crossplane.io/docs/v1.10/reference/uninstall.html.
Very important that providers and provider config is deleted correctly. You might end up in a state with
resources that waits for finalizers to be deleted (usually I have seen this on some CRDs and ProviderConfig). If that is the case
it is possible to just edit that resource by deleting the finalizers. 