FROM gcr.io/distroless/static-debian11

COPY  ./build/4stats /4stats

CMD ["/4stats"]
