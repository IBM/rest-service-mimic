FROM scratch
COPY output/rest-service-mimic /
ENTRYPOINT ["/rest-service-mimic"]
