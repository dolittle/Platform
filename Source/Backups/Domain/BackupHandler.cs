// Copyright (c) Dolittle. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

using System;
using System.Threading.Tasks;
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
        public async Task<IActionResult> Stored(Request request)
        {
            var eventSource = EventSources.From(request.Application, request.Environment);
            await _client
                .EventStore.ForTenant(request.Tenant)
                .Commit(_ =>
                    _.CreatePublicEvent(
                        new Events.DatabaseBackupStored(
                            request.Application,
                            request.Environment,
                            request.ShareName,
                            request.BackupFileName))
                    .FromEventSource(eventSource));
            return Ok();
        }

        [HttpPost("failed")]
        public async Task<IActionResult> Failed(Request request)
        {
            var eventSource = EventSources.From(request.Application, request.Environment);
            await _client
                .EventStore.ForTenant(request.Tenant)
                .Commit(_ =>
                    _.CreatePublicEvent(
                        new Events.DatabaseBackupFailed(
                            request.Application,
                            request.Environment,
                            request.ShareName,
                            request.BackupFileName))
                    .FromEventSource(eventSource));
            return Ok();
        }
    }
    public record Request(
        string BackupFileName,
        Guid Tenant,
        string Environment,
        Guid Application,
        string ShareName);
}
