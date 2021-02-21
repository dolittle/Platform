// Copyright (c) Dolittle. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

using Dolittle.SDK;
using Dolittle.Data.Backups.Events;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.OpenApi.Models;
using Dolittle.Data.Backups.Filters;

namespace Dolittle.Data.Backups.Head
{
    public class Startup
    {
        public Startup(IConfiguration configuration)
        {
            Configuration = configuration;
        }

        public IConfiguration Configuration { get; }

        public void ConfigureServices(IServiceCollection services)
        {
            var dolittleClient = BuildClient();
            dolittleClient.Start();
            services.AddSingleton(dolittleClient);
            services.AddControllers();
            services.AddSwaggerGen(c =>
            {
                c.SwaggerDoc("v1", new OpenApiInfo { Title = "Backup/Head", Version = "v1" });
            });
        }

        public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
        {
            if (env.IsDevelopment())
            {
                app.UseDeveloperExceptionPage();
                app.UseSwagger();
                app.UseSwaggerUI(c => c.SwaggerEndpoint("/swagger/v1/swagger.json", "Backup/Head v1"));
            }

            app.UseRouting();

            app.UseEndpoints(endpoints =>
            {
                endpoints.MapControllers();
            });
        }
        Client BuildClient()
            => new ClientBuilder("702967ce-9808-44f7-9399-d45cafa9be07")
                .WithRuntimeOn("localhost", 50053)
                .WithEventTypes(_ => _.RegisterAllFrom(typeof(EventTypeRegistry).Assembly))
                .WithFilters(filtersBuilder => 
                    filtersBuilder
                        .CreateBackupFilter()).Build();
    }
}
