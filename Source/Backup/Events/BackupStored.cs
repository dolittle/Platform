// Copyright (c) Dolittle. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

using System;
using Dolittle.SDK.Events;

namespace Dolittle.Platform.Backup.Events
{
    [EventType(EventTypeRegistry.BackupStoredId, EventTypeRegistry.BackupStoredGeneration)]
    public class BackupStored
    {
        public BackupStored(string dumpFilepath, string environment, Guid application)
        {
            DumpFilepath = dumpFilepath;
            Environment = environment;
            Application = application;
        }

        public string DumpFilepath { get; }
        public string Environment { get; }
        public Guid Application { get; }
    }
}
