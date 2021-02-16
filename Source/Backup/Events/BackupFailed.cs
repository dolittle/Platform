// Copyright (c) Dolittle. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

using Dolittle.SDK.Events;

namespace Dolittle.Platform.Backup.Events
{
    [EventType("02a06c14-51b2-4a2e-8e32-f5dbb9783d5e")]
    public class BackupFailed
    {
        public BackupFailed(string environment)
        {
            Environment = environment;
        }

        public string Environment { get; }
    }
}
