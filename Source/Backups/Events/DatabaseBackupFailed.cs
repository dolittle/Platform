// Copyright (c) Dolittle. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

using System;
using Dolittle.SDK.Events;

namespace Dolittle.Data.Backups.Events
{
    [EventType(EventTypeRegistry.DatabaseBackupFailedId, EventTypeRegistry.DatabaseBackupFailedGeneration)]
    public class DatabaseBackupFailed
    {
        public DatabaseBackupFailed(
            Guid application,
            string environment,
            string shareName,
            string backupFileName)
        {
            Application = application;
            Environment = environment;
            ShareName = shareName;
            BackupFileName = backupFileName;
        }

        public Guid Application { get; }
        public string Environment { get; }
        public string ShareName { get; }
        public string BackupFileName { get; }
    }
}
