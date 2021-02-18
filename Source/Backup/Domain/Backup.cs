// Copyright (c) Dolittle. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.


using Dolittle.SDK.Aggregates;
using Dolittle.SDK.Events;
using Dolittle.Platform.Backup.Events;
using System;

namespace Dolittle.Platform.Backup.Domain
{
    [AggregateRoot("9f34d1f8-6bc4-4146-9d47-bde3d1cf25e9")]
    public class Backup : AggregateRoot
    {
        public Backup(EventSourceId eventSource)
            : base(eventSource)
        {
        }

        public void StartBackup(DateTimeOffset startTime, string filepath, string environment, Guid application)
        {
            ApplyPublic(new BackupStarted(startTime, filepath, environment, application));
        }
        
        public void NotifyOfBackupStored(string filepath, string environment, Guid application)
        {
            ApplyPublic(new BackupStored(filepath, environment, application));
        }
    }
}
