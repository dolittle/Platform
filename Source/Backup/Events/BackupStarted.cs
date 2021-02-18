// Copyright (c) Dolittle. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

using System;
using Dolittle.SDK.Events;

namespace Dolittle.Platform.Backup.Events
{
    [EventType(EventTypeRegistry.BackupStartedId, EventTypeRegistry.BackupStartedGeneration)]
    public class BackupStarted
    {
        public BackupStarted(DateTimeOffset startTime, string dumpFilename, string environment, Guid application)
        {
            StartTime = startTime;
            DumpFilename = dumpFilename;
            Environment = environment;
            Application = application;
        }

        public DateTimeOffset StartTime { get; }
        public string DumpFilename { get; }
        public string Environment { get; }
        public Guid Application { get; }
    }
}
