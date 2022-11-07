// Copyright (c) Dolittle. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

using KubeOps.Operator;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddKubernetesOperator();

var app = builder.Build();

app.UseKubernetesOperator();

await app.RunOperatorAsync(args).ConfigureAwait(false);
