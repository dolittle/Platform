﻿// Copyright (c) Dolittle. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

using System;
using Dolittle.SDK.Events;

namespace Dolittle.Platform.Backup.Events
{
    [EventType(EventTypeRegistry.BackupStartedId, EventTypeRegistry.BackupStartedGeneration)]
    public class BackupStarted
    {
        public BackupStarted(DateTimeOffset startTime, string dumpFilepath, string environment, Guid application)
        {
            StartTime = startTime;
            DumpFilepath = dumpFilepath;
            Environment = environment;
            Application = application;
        }

        public DateTimeOffset StartTime { get; }
        public string DumpFilepath { get; }
        public string Environment { get; }
        public Guid Application { get; }
    }
}
