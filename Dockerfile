FROM gcr.io/distroless/static:latest
LABEL maintainers="Kubernetes COSI Authors"
LABEL description="Azure COSI driver"

COPY ./bin/azure-cosi-driver azure-cosi-driver
ENTRYPOINT ["/azure-cosi-driver"]
