// Copyright (c) Dolittle. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

using KubeOps.Operator.Webhooks;

namespace Controller.Customers;

public class Validation : IValidationWebhook<Customer>
{
    readonly ILogger<Validation> _logger;

    public Validation(ILogger<Validation> logger)
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