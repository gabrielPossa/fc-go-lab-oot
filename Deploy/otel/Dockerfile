# Imagem do OTEL Collector + um busybox estatico (musl) usado apenas como
# probe do healthcheck. A imagem oficial e distroless (sem sh/wget/nc), entao
# sem isso o Docker nao consegue sondar a extensao health_check (porta 13133).
FROM otel/opentelemetry-collector:latest

COPY --from=busybox:musl /bin/busybox /busybox
