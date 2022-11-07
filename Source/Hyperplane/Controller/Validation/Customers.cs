// Copyright (c) Dolittle. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

using Controller.Entities;
using KubeOps.Operator.Webhooks;

namespace Controller.Validation;

public class Customers : IValidationWebhook<Customer>
{
    readonly ILogger<Customers> _logger;

    public Customers(ILogger<Customers> logger)
    {
        _logger = logger;
    }

    public AdmissionOperations Operations => AdmissionOperations.Update;

    public Task<ValidationResult> UpdateAsync(Customer oldEntity, Customer newEntity, bool dryRun)
    {
        _logger.LogInformation("Validating Customer {OldId} - {NewId}", oldEntity.Spec.Id, newEntity.Spec.Id);
        
        if (oldEntity.Spec.Id != newEntity.Spec.Id)
        {
            return Task.FromResult(ValidationResult.Fail(1, "Cannot change customer id"));
        }

        return Task.FromResult(ValidationResult.Success());
    }
}