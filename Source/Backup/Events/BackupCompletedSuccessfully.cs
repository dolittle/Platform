// Copyright (c) Dolittle. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

using Dolittle.SDK.Events;

namespace Dolittle.Platform.Backup.Events
{
    [EventType("4de270ec-3ee6-4482-8c0c-513f13bf755f")]
    public class BackupCompletedSuccessfully
    {
        public BackupCompletedSuccessfully(string environment, string fileName)
        {
            Environment = environment;
            FileName = fileName;
        }

        public string Environment { get; }
        public string FileName { get; } 
    }
}
