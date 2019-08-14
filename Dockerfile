FROM scratch
EXPOSE 8080
ADD src/bin /
CMD ["/main"]
