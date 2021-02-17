// Copyright (c) Dolittle. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

using Dolittle.SDK.Events.Filters;

namespace Dolittle.Platform.Backup.Filters
{
    public static class EventFiltersBuilderExtensions
    {
        /// <summary>
        /// 
        /// </summary>
        /// <param name="builder"></param>
        /// <returns></returns>
        public static EventFiltersBuilder BuildBackupStream(this EventFiltersBuilder builder)
            => builder.CreatePublicFilter(
                "584546c2-d3be-40ca-8321-d23dc7ed1d65",
                _ => _.Handle((@event, context) =>
                {
                    
                    return new PartitionedFilterResult(IsBackupEvent(), context.CommittedExecutionContext.Tenant.Value);
                }));


        bool IsBackupEvent()
            => 
    }
}
