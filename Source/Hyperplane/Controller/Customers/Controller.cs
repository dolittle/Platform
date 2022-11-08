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

            ns.EnsureMetadata().EnsureAnnotations().TryGetValue("customers.dolittle.io/owner-namespace", out var ownerNamespace);
            ns.EnsureMetadata().EnsureAnnotations().TryGetValue("customers.dolittle.io/owner-name", out var ownerName);

            if (ownerNamespace != customer.Namespace() || ownerName != customer.Name())
            {
                _logger.LogWarning("Namespace {Name} is not owned by the correct Customer", ns.Name());
                await _events.PublishAsync(customer, "NamespaceOwnedByOther", "Namespace exists but not owned by Customer", EventType.Warning).ConfigureAwait(false);
                return ResourceControllerResult.RequeueEvent(TimeSpan.FromSeconds(10));
            }

            // TODO: Check that we own the NS
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

        // TODO: Set owner reference (should manage deletion for us?)
        // Oh - maybe we cant since this is cross-namespace. But maybe we can use finalizers...
        _logger.LogInformation("Creating namespace for Customer {Id}", customer.Spec.Id);
        await _client.Create(ns).ConfigureAwait(false);
        await _events.PublishAsync(customer, "NamespaceCreated", "Created Namespace for Customer").ConfigureAwait(false);
        return null;
    }

    async Task SetNamespaceThings(V1Namespace ns, Customer customer)
    {
        ns.EnsureMetadata().EnsureAnnotations()["customers.dolittle.io/owner-namespace"] = customer.Metadata.Namespace();
        ns.EnsureMetadata().EnsureAnnotations()["customers.dolittle.io/owner-name"] = customer.Metadata.Name;

        ns.EnsureMetadata().EnsureAnnotations()["dolittle.io/customer-id"] = customer.Spec.Id;
        ns.EnsureMetadata().EnsureLabels()["customer"] = customer.Spec.Name;


        _logger.LogInformation("Adding finalizers");
        await _finalizers.RegisterAllFinalizersAsync(customer).ConfigureAwait(false);
    }

    public async Task<ResourceControllerResult?> ReconcileAsync(V1Namespace v1Namespace)
    {
        _logger.LogInformation("Reconciling on Namespace {}", v1Namespace.Name());
        return null;
    }
}