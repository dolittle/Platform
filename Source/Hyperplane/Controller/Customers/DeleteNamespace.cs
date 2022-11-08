// Copyright (c) Dolittle. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

using DotnetKubernetesClient;
using k8s;
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

    public string Identifier => "customers.dolittle.io/delete-namespace";

    public async Task FinalizeAsync(Customer customer)
    {
        _logger.LogInformation("Finalizing Customer {Name}", customer.Name());
        var nsName = $"customer-{customer.Spec.Id}";
        var ns = await _client.Get<V1Namespace>(nsName).ConfigureAwait(false);

        if (ns is null)
        {
            // NS does not exist - so no need to do anything
            return;
        }
        
        var ownerNamespace = ns.GetAnnotation("customers.dolittle.io/owner-namespace");
        var ownerName= ns.GetAnnotation("customers.dolittle.io/owner-name");
        if (ownerNamespace != customer.Namespace() || ownerName != customer.Name())
        {
            // NS is not owned by us - so no need to do anything
            return;
        }
        
        _logger.LogInformation("Deleting Namespace {Name}", ns.Name());
        await _client.Delete(ns).ConfigureAwait(false);
        await WaitForDeletion(ns).ConfigureAwait(false);
    }

    async Task WaitForDeletion(V1Namespace ns)
    {
        var tcs = new TaskCompletionSource(TaskCreationOptions.RunContinuationsAsynchronously);
        
        using var watcher = await _client.Watch<V1Namespace>(
            TimeSpan.FromMinutes(1),
            (evt, evtNS) =>
            {
                if (evtNS.Uid() != ns.Uid())
                {
                    return;
                }

                if (evt != WatchEventType.Deleted)
                {
                    return;
                }
                
                _logger.LogInformation("Was notified of NS deletion {Name}", evtNS.Name());
                tcs.SetResult();
            },
            (ex) =>
            {
                tcs.SetException(ex);
            },
            () =>
            {
                tcs.SetCanceled();
            });

        await tcs.Task;
    }
}