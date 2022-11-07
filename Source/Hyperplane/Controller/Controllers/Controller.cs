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

    public Task<ResourceControllerResult?> ReconcileAsync(Customer customer)
    {
        var namespace = $"{customer.Spec.Id}";
        _client.Get<V1Namespace>()

        _client.Create<V1Namespace>(new V1Namespace{Spec = new V1NamespaceSpec())
        _logger.LogInformation("Reconciling Customer {Id} - {Name}", customer.Spec.Id, customer.Spec.Name);

        return Task.FromResult<ResourceControllerResult?>(null);
    }
}