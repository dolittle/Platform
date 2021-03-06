﻿// Copyright (c) Dolittle. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

using System.Threading.Tasks;
using Dolittle.Data.Backups.Events;
using Dolittle.SDK.Events;
using Dolittle.SDK.Events.Filters;

namespace Dolittle.Data.Backups.Filters
{
    public static class EventFiltersBuilderExtensions
    {
        public const string BackupFilterId = "584546c2-d3be-40ca-8321-d23dc7ed1d65";

        public static EventFiltersBuilder CreateBackupFilter(this EventFiltersBuilder builder)
            => builder.CreatePublicFilter(
                BackupFilterId,
                _ => _.Handle((@event, context) => Task.FromResult(new PartitionedFilterResult(IsBackupEvent(context.Type), PartitionId.Unspecified))));


        static bool IsBackupEvent(EventType type)
            => type.Id.Value.ToString() switch
            {
                EventTypeRegistry.EventStoreAndReadModelsBackedUpId => true,
                _ => false
            };
    }
}
