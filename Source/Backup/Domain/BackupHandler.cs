// Copyright (c) Dolittle. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

using System;
using System.ComponentModel.DataAnnotations;
using System.Threading.Tasks;
using Dolittle.SDK;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;

namespace Dolittle.Platform.Backup.Domain
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
        public async Task<IActionResult> Start(StartBackupRequest request)
        {
            _logger.LogInformation("Starting backup");
            await _client
                .AggregateOf<Backup>(request.EventSource, _ => _.ForTenant(request.Tenant))
                .Perform(_ => _.StartBackup(DateTimeOffset.UtcNow, request.DumpFilepath, request.Environment, request.Application));
            return Ok();
        }

        [HttpPost("notify")]
        public async Task<IActionResult> NotifyStored(NotifyBackupStoredRequest request)
        {
            _logger.LogInformation("Notifying that backup has been stored");
            await _client
                .AggregateOf<Backup>(request.EventSource, _ => _.ForTenant(request.Tenant))
                .Perform(_ => _.NotifyOfBackupStored(request.DumpFilepath, request.Environment, request.Application));
            return Ok();
        }
    }
    public record StartBackupRequest(
        [Required]string DumpFilepath,
        [Required]Guid Tenant,
        [Required]string Environment,
        [Required]Guid EventSource,
        [Required]Guid Application);
    public record NotifyBackupStoredRequest(
        [Required]string DumpFilepath,
        [Required]Guid Tenant,
        [Required]string Environment,
        [Required]Guid EventSource,
        [Required]Guid Application);
}
