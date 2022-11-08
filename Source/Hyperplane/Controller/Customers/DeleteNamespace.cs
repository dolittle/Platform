// Copyright (c) Dolittle. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

using DotnetKubernetesClient;
using k8s.Models;
using KubeOps.Operator.Finalizer;

namespace Controller.Customers;

public class DeleteNamespace : IResourceFinalizer<Customer>
{
    readonly IKubernetesClient _client;
    readonly ILogger _logger;
    
    public DeleteNamespace(IKubernetesClient client, ILogger<DeleteNamespace> logger)
    {
        _client = client;
        _logger = logger;
    }

    public async Task FinalizeAsync(Customer customer)
    {
        var nsName = $"customer-{customer.Spec.Id}";
        var ns = new V1Namespace
        {
            Metadata = new V1ObjectMeta
            {
                Name = nsName
            }
        };
        // TODO: Make sure we really want to delete this namespace?
        _logger.LogInformation("Deleting Namespace {Name}", ns.Name());
        await _client.Delete(ns).ConfigureAwait(false);
        _logger.LogInformation("Deleted Namespace {Name}", ns.Name());
    }

    async Task WaitForDeletion(V1Namespace ns)
    {
        var watcher = await _client.Watch<V1Namespace>(TimeSpan.FromSeconds(10), (x, y) => {});
        watcher.
    }
}