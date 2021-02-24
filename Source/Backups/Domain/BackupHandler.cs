// Copyright (c) Dolittle. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

using System;
using System.Threading.Tasks;
using Dolittle.Data.Backups.Events;
using Dolittle.SDK;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;

namespace Dolittle.Data.Backups.Domain
{
    [ApiController]
    [Route("api/backup")]
    public class BackupHandler : ControllerBase
    {
        readonly ILogger<BackupHandler> _logger;
        readonly Client _client;

        public BackupHandler(ILogger<BackupHandler> logger, Client client)
        {
            _logger = logger;
            _client = client;
        }

        [HttpPost("stored")]
        public async Task<IActionResult> Stored(NotifyStoredRequest request)
        {
            var eventSource = EventSources.From(request.Application, request.Environment);
            _logger.LogDebug("Received {Request} on event source {EventSource}", request, eventSource);
            await _client
                .EventStore.ForTenant(request.Tenant)
                .Commit(_ =>
                    _.CreatePublicEvent(
                        new EventStoreAndReadModelsBackedUp(
                            request.Application,
                            request.Environment,
                            request.ShareName,
                            request.BackupFileName,
                            request.DurationInSeconds))
                    .FromEventSource(eventSource));
            return Ok();
        }
    }
    public record Request
    {
        public string BackupFileName { get; init; }
        public Guid Tenant { get; init; }
        public string Environment { get; init; }
        public Guid Application { get; init; }
        public string ShareName { get; init; }
        public uint DurationInSeconds { get; init; }
    }
    public record NotifyStoredRequest : Request;
}
