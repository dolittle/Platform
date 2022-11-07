// Copyright (c) Dolittle. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

using Controller.Entities;
using DotnetKubernetesClient;
using k8s.Models;
using KubeOps.Operator.Controller;
using KubeOps.Operator.Controller.Results;
using KubeOps.Operator.Rbac;

namespace Controller.Controllers;

[EntityRbac(typeof(Customer), Verbs = RbacVerb.Get | RbacVerb.List | RbacVerb.Watch)]
public class Customers : IResourceController<Customer>
{
    readonly IKubernetesClient _client;
    readonly ILogger _logger;

    public Customers(IKubernetesClient client, ILogger<Customers> logger)
    {
        _client = client;
        _logger = logger;
    }

    public async Task<ResourceControllerResult?> ReconcileAsync(Customer customer)
    {
        var nsName = $"customer-{customer.Spec.Id}";

        var ns = await _client.Get<V1Namespace>(nsName).ConfigureAwait(false);
        if (ns is not null)
        {
            _logger.LogInformation("Namespace for Customer {Id} already exists", customer.Spec.Id);
            // TODO: Check annotations etc.
            return null;
        }

        ns = new V1Namespace
        {
            Metadata = new V1ObjectMeta
            {
                Name = nsName,
                Annotations = new Dictionary<string, string>
                {
                    { "dolittle.io/customer-id", customer.Spec.Id }
                },
                Labels = new Dictionary<string, string>
                {
                    { "customer", customer.Spec.Name },
                },
            },
        };
        // TODO: Set owner reference (should manage deletion for us?)
        // Oh - maybe we cant since this is cross-namespace. But maybe we can use finalizers...
        _logger.LogInformation("Creating namespace for Customer {Id}", customer.Spec.Id);
        await _client.Create(ns).ConfigureAwait(false);
        return null;
    }
}