# DotNet Source - filtered

FROM alpine AS dotnet-source
ARG CONFIGURATION=Release
WORKDIR /app
COPY Source/Backups/ ./Source/
RUN rm -rf Source/Cli

# DotNet Build
FROM mcr.microsoft.com/dotnet/sdk:5.0-alpine AS dotnet-build
ARG CONFIGURATION=Release
WORKDIR /app
COPY --from=dotnet-source /app/Source ./Source/

WORKDIR /app/Source/Head
RUN dotnet restore
RUN dotnet publish -c ${CONFIGURATION} -o out

# Image
FROM mcr.microsoft.com/dotnet/aspnet:5.0
ARG CONFIGURATION=Release

RUN echo Configuration = $CONFIGURATION

RUN if [ "$CONFIGURATION" = "Debug" ] ; then apt-get update && \
    apt-get install -y --no-install-recommends unzip procps && \
    rm -rf /var/lib/apt/lists/* \
    ; fi

RUN if [ "$CONFIGURATION" = "debug" ] ; then curl -sSL https://aka.ms/getvsdbgsh | bash /dev/stdin -v latest -l ~/vsdbg ; fi
RUN apt-get update && apt-get install -y libc-dev

WORKDIR /app
COPY --from=dotnet-build /app/Source/Head/out ./

EXPOSE 5000
EXPOSE 5001
ENV ASPNETCORE_URLS=http://*:5000/

ENTRYPOINT ["dotnet", "Head.dll"]
