// Copyright (c) Dolittle. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

using DotnetKubernetesClient;
using k8s.Models;
using KubeOps.Operator.Controller;
using KubeOps.Operator.Controller.Results;
using KubeOps.Operator.Events;
using KubeOps.Operator.Finalizer;
using KubeOps.Operator.Rbac;

namespace Controller.Customers;

[EntityRbac(typeof(Customer), Verbs = RbacVerb.Get | RbacVerb.List | RbacVerb.Watch)]
public class Controller : IResourceController<Customer>
{
    readonly IKubernetesClient _client;
    readonly IEventManager _events;
    readonly IFinalizerManager<Customer> _finalizers;
    readonly ILogger _logger;

    public Controller(IKubernetesClient client, IEventManager events, IFinalizerManager<Customer> finalizers, ILogger<Controller> logger)
    {
        _client = client;
        _events = events;
        _finalizers = finalizers;
        _logger = logger;
    }

    public async Task<ResourceControllerResult?> ReconcileAsync(Customer customer)
    {
        var nsName = $"customer-{customer.Spec.Id}";
        var ns = await _client.Get<V1Namespace>(nsName).ConfigureAwait(false);
        if (ns is not null)
        {
            _logger.LogInformation("Namespace for Customer {Id} already exists", customer.Spec.Id);

            var ownerNamespace = ns.GetAnnotation("customers.dolittle.io/owner-namespace");
            var ownerName= ns.GetAnnotation("customers.dolittle.io/owner-name");

            if (ownerNamespace != customer.Namespace() || ownerName != customer.Name())
            {
                _logger.LogWarning("Namespace {Name} is not owned by the correct Customer", ns.Name());
                await _events.PublishAsync(customer, "NamespaceOwnedByOther", "Namespace exists but not owned by Customer", EventType.Warning).ConfigureAwait(false);
                return ResourceControllerResult.RequeueEvent(TimeSpan.FromSeconds(10));
            }

            SetNamespaceThings(ns, customer);
            await _client.Update(ns).ConfigureAwait(false);
            return null;
        }

        ns = new V1Namespace
        {
            Metadata = new V1ObjectMeta
            {
                Name = nsName,
            }
        };
        SetNamespaceThings(ns, customer);

        _logger.LogInformation("Creating namespace for Customer {Id}", customer.Spec.Id);
        await _client.Create(ns).ConfigureAwait(false);
        await _events.PublishAsync(customer, "NamespaceCreated", "Created Namespace for Customer").ConfigureAwait(false);
        return null;
    }

    async Task SetNamespaceThings(V1Namespace ns, Customer customer)
    {
        ns.SetAnnotation("customers.dolittle.io/owner-namespace", customer.Namespace());
        ns.SetAnnotation("customers.dolittle.io/owner-name", customer.Name());

        ns.SetAnnotation("dolittle.io/customer-id", customer.Spec.Id);
        ns.SetLabel("customer", customer.Spec.Name);

        _logger.LogInformation("Adding finalizers");
        await _finalizers.RegisterAllFinalizersAsync(customer).ConfigureAwait(false);
    }
}