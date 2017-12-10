FROM scratch
#ADD /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ADD main /
CMD ["/main"]
