// Copyright (c) Dolittle. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

using System;
using Dolittle.SDK.Events;

namespace Dolittle.Platform.Backup.Events
{
    [EventType(EventTypeRegistry.EventStoreBackupStoredId, EventTypeRegistry.EventStoreBackupStoredGeneration)]
    public class EventStoreBackupStored
    {
        public EventStoreBackupStored(
            Guid application,
            string environment,
            string applicationName,
            string shareName,
            string backupFileName)
        {
            Application = application;
            Environment = environment;
            ApplicationName = applicationName;
            ShareName = shareName;
            BackupFileName = backupFileName;
        }

        public Guid Application { get; }
        public string Environment { get; }
        public string ApplicationName { get; }
        public string ShareName { get; }
        public string BackupFileName { get; }
    }
}
