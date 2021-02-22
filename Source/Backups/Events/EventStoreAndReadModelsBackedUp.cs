// Copyright (c) Dolittle. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

using System;
using Dolittle.SDK.Events;

namespace Dolittle.Data.Backups.Events
{
    [EventType(EventTypeRegistry.EventStoreAndReadModelsBackedUpId, EventTypeRegistry.EventStoreAndReadModelsBackedUpGeneration)]
    public record EventStoreAndReadModelsBackedUp(
        Guid Application,
        string Environment,
        string ShareName,
        string BackupFileName,
        uint BackupDuration);
}
