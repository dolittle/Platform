// Copyright (c) Dolittle. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

using System;
using Dolittle.SDK.Events;

namespace Dolittle.Platform.Backup.Events
{
    [EventType("3a86e422-c958-40f8-91b3-383ea0f70d0a")]
    public class BackupStarted
    {
        public BackupStarted(string dumpFilename, string environment, Guid application)
        {
            DumpFilename = dumpFilename;
            Environment = environment;
            Application = application;
        }

        public string DumpFilename { get; }
        public string Environment { get; }
        public Guid Application { get; }
    }
}
