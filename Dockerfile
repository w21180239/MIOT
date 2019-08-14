FROM scratch
EXPOSE 8080
ADD bin /
CMD ["/main"]
