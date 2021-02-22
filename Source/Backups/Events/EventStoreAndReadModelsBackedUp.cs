// Copyright (c) Dolittle. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

using System;
using Dolittle.SDK.Events;

namespace Dolittle.Data.Backups.Events
{
    [EventType(EventTypeRegistry.DatabaseBackupStoredId, EventTypeRegistry.DatabaseBackupStoredGeneration)]
    public record EventStoreAndReadModelsBackedup(
        Guid Application,
        string Environment,
        string ShareName,
        string BackupFileName);
}
