ARG BASE_REGISTRY=registry.access.redhat.com
ARG BASE_IMAGE=ubi8-minimal
ARG BASE_TAG=8.7

FROM ${BASE_REGISTRY}/${BASE_IMAGE}:${BASE_TAG} AS extracted_bundle

ADD bundle.tar.gz /bundle

FROM ${BASE_REGISTRY}/${BASE_IMAGE}:${BASE_TAG}

LABEL name="central-db" \
      vendor="StackRox" \
      maintainer="support@stackrox.com" \
      summary="Stackrox Central Database for the StackRox Kubernetes Security Platform" \
      description="This image provides Database services to Stackrox Central in the StackRox Kubernetes Security Platform."

ENV PG_MAJOR=13
ENV PATH="$PATH:/usr/pgsql-$PG_MAJOR/bin/"
ENV LANG=en_US.utf8

COPY signatures/RPM-GPG-KEY-PGDG-13 /
COPY scripts/docker-entrypoint.sh /usr/local/bin/
COPY scripts/init-entrypoint.sh /usr/local/bin/
COPY --from=extracted_bundle /bundle/postgres.rpm /bundle/postgres-libs.rpm /bundle/postgres-server.rpm /bundle/postgres-contrib.rpm /tmp/

RUN microdnf upgrade -y --nobest && \
    # groupadd is in shadow-utils package that is not installed by default.
    microdnf --setopt=install_weak_deps=0 --setopt=tsflags=nodocs install -y shadow-utils && \
    groupadd -g 70 postgres && \
    adduser postgres -u 70 -g 70 -d /var/lib/postgresql -s /bin/sh && \
    rpm --import RPM-GPG-KEY-PGDG-13 && \
    microdnf --setopt=install_weak_deps=0 --setopt=tsflags=nodocs install -y \
        ca-certificates libicu systemd-sysv \
        glibc-locale-source glibc-langpack-en \
        perl-libs libxslt && \
    rpm -i /tmp/postgres-libs.rpm /tmp/postgres-server.rpm /tmp/postgres.rpm /tmp/postgres-contrib.rpm && \
    # Restore /usr/share/zoneinfo that's empty in ubi-minimal because postgres reads timezone data from it.
    # https://access.redhat.com/solutions/5616681
    microdnf --setopt=install_weak_deps=0 --setopt=tsflags=nodocs reinstall tzdata && \
    microdnf clean all && \
    rpm -e --nodeps $(rpm -qa curl '*rpm*' '*dnf*' '*libsolv*' '*hawkey*' 'yum*') && \
    rm -rf /var/cache/dnf /var/cache/yum && \
    localedef -f UTF-8 -i en_US en_US.UTF-8 && \
    chown postgres:postgres /usr/local/bin/docker-entrypoint.sh && \
    chmod +x /usr/local/bin/docker-entrypoint.sh && \
    mkdir /docker-entrypoint-initdb.d

# Use SIGINT to bring down with Fast Shutdown mode
STOPSIGNAL SIGINT

USER postgres:postgres

ENTRYPOINT ["docker-entrypoint.sh"]

EXPOSE 5432
CMD ["postgres", "-c", "config_file=/etc/stackrox.d/config/postgresql.conf"]