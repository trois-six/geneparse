FROM golang AS build

WORKDIR /geneparse
COPY . ./
RUN make build

FROM alpine

WORKDIR /
COPY --from=build /geneparse/geneparse /usr/local/bin/geneparse
USER nobody:nobody
ENTRYPOINT ["/usr/local/bin/geneparse"]
