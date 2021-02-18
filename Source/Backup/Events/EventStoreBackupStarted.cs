// Copyright (c) Dolittle. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

using System;
using Dolittle.SDK.Events;

namespace Dolittle.Platform.Backup.Events
{
    [EventType(EventTypeRegistry.EventStoreBackupStartedId, EventTypeRegistry.EventStoreBackupStartedGeneration)]
    public class EventStoreBackupStarted
    {
        public EventStoreBackupStarted(
            DateTimeOffset startTime,
            Guid application,
            string environment,
            string applicationName,
            string shareName,
            string backupFileName)
        {
            StartTime = startTime;
            Application = application;
            Environment = environment;
            ApplicationName = applicationName;
            ShareName = shareName;
            BackupFileName = backupFileName;
        }

        public DateTimeOffset StartTime { get; }
        public Guid Application { get; }
        public string Environment { get; }
        public string ApplicationName { get; }
        public string ShareName { get; }
        public string BackupFileName { get; }
    }
}
