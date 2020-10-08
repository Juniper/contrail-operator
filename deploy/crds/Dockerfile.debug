FROM gcr.io/distroless/base-debian10:debug
SHELL ["/busybox/sh", "-c"]
ADD https://storage.googleapis.com/kubernetes-release/release/v1.18.0/bin/linux/amd64/kubectl /app/
RUN chmod 700 /app/kubectl
ADD *_crd.yaml /crds/
ENTRYPOINT ["sh","-c","/app/kubectl apply -f /crds && /app/kubectl wait crds --for=condition=Established --timeout=2m managers.contrail.juniper.net"]
