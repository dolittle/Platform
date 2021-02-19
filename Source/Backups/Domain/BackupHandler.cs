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

        [HttpPost("start")]
        public async Task<IActionResult> Start(Request request)
        {
            _logger.LogInformation("Starting backup");
            await _client
                .AggregateOf<Backup>(request.EventSource, _ => _.ForTenant(request.Tenant))
                .Perform(_ => _.StartBackup(
                    DateTimeOffset.UtcNow,
                    request.Application,
                    request.Environment,
                    request.ApplicationName,
                    request.ShareName,
                    request.BackupFileName));
            return Ok();
        }

        [HttpPost("stored")]
        public async Task<IActionResult> NotifyStored(Request request)
        {
            _logger.LogInformation("Notifying that backup has been stored");
            await _client
                .AggregateOf<Backup>(request.EventSource, _ => _.ForTenant(request.Tenant))
                .Perform(_ => _.NotifyOfBackupStored(
                    request.Application,
                    request.Environment,
                    request.ApplicationName,
                    request.ShareName,
                    request.BackupFileName));
            return Ok();
        }
    }
    public record Request(
        string BackupFileName,
        Guid Tenant,
        string Environment,
        Guid EventSource,
        Guid Application,
        string ApplicationName,
        string ShareName);
}
