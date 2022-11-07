// Copyright (c) Dolittle. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

using KubeOps.Operator;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddKubernetesOperator()
    //.AddWebhookLocaltunnel() // ONLY USE THIS WHEN NO SECRETS ARE INVOLVED (AND DEFINITELY NOT IN PRODUCTION - DATA GOES TO A 3RD PARTY SERVER!!!
;

var app = builder.Build();

app.UseKubernetesOperator();

await app.RunOperatorAsync(args).ConfigureAwait(false);
