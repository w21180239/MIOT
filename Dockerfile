FROM scratch
EXPOSE 8080
ADD bin/main /
CMD ["/main"]
