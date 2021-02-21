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
        public async Task<IActionResult> Stored(BackupStoredRequest request)
        {
            await _client
                .EventStore.ForTenant(request.Tenant)
                .Commit(_ =>
                    _.CreatePublicEvent(
                        new DatabaseBackupStored(
                            request.Application,
                            request.Environment,
                            request.ShareName,
                            request.BackupFileName))
                    .FromEventSource(EventSources.From(request.Application, request.Environment)));
            return Ok();
        }

        [HttpPost("failed")]
        public async Task<IActionResult> Failed(BackupFailedRequest request)
        {
            await _client
                .EventStore.ForTenant(request.Tenant)
                .Commit(_ =>
                    _.CreatePublicEvent(
                        new DatabaseBackupFailed(
                            request.Application,
                            request.Environment,
                            request.ShareName,
                            request.BackupFileName,
                            request.FailureReason))
                    .FromEventSource(EventSources.From(request.Application, request.Environment)));
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
    }
    public record BackupStoredRequest : Request;
    public record BackupFailedRequest : Request
    {
        public string FailureReason { get; init; }
    }
}
