// Copyright (c) Dolittle. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

using k8s.Models;
using KubeOps.Operator.Entities;
using KubeOps.Operator.Entities.Annotations;

namespace Controller.Entities;

public class CustomerSpec
{
    /// <summary>
    /// The Name of the customer.
    /// </summary>
    [Required]
    public string Name { get; set;}

    /// <summary>
    /// The Id of the customer.
    /// </summary>
    [Required, Pattern("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}")]
    public string Id { get; set; }
}

[KubernetesEntity(Group = "dolittle.io", ApiVersion = "v1alpha1")]
public class Customer : CustomKubernetesEntity<CustomerSpec>
{

}