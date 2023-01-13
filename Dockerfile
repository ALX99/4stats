FROM gcr.io/distroless/static-debian11

COPY  ./build/yonsuu /yonsuu

CMD ["/yonsuu"]
